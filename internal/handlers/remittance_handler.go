package handlers

import (
	"strconv"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetRemittanceCertificates lists remittance certificates with pagination and filters
func GetRemittanceCertificates(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("remark ILIKE ?", "%"+search+"%")
	}

	var total int64
	query.Model(&models.RemittanceCertificate{}).Count(&total)

	var items []models.RemittanceCertificate
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch remittance certificates"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetRemittanceCertificate gets a single remittance certificate by ID
func GetRemittanceCertificate(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.RemittanceCertificate
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Remittance certificate not found"})
	}

	return c.JSON(item)
}

// ApproveRemittanceCertificate approves a remittance certificate
func ApproveRemittanceCertificate(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.RemittanceCertificate
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Remittance certificate not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":     "approved",
		"auditor_id": userID,
		"audited_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve remittance certificate"})
	}

	item.Status = "approved"
	item.AuditorID = &userID
	item.AuditedAt = &now
	return c.JSON(item)
}

// RejectRemittanceCertificate rejects a remittance certificate
func RejectRemittanceCertificate(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		AuditRemark string `json:"audit_remark"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.RemittanceCertificate
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Remittance certificate not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":       "rejected",
		"audit_remark": body.AuditRemark,
		"auditor_id":   userID,
		"audited_at":   now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject remittance certificate"})
	}

	item.Status = "rejected"
	item.AuditRemark = body.AuditRemark
	item.AuditorID = &userID
	item.AuditedAt = &now
	return c.JSON(item)
}

// GetRemittanceCertificateCount returns count of pending remittance certificates
func GetRemittanceCertificateCount(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var count int64
	database.DB.Model(&models.RemittanceCertificate{}).
		Where("tenant_id = ? AND trashed_at IS NULL AND status = ?", tenantID, "pending").
		Count(&count)

	return c.JSON(fiber.Map{"count": count})
}
