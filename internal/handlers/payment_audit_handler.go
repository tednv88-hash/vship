package handlers

import (
	"encoding/json"
	"strconv"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetPaymentAudits lists payment audits with pagination and filters
func GetPaymentAudits(c *fiber.Ctx) error {
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

	query := database.DB.Model(&models.PaymentAudit{}).
		Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var audits []models.PaymentAudit
	if err := query.
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&audits).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch payment audits"})
	}

	return c.JSON(fiber.Map{
		"data":  audits,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetPaymentAudit gets a single payment audit by ID
func GetPaymentAudit(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment audit ID"})
	}

	var audit models.PaymentAudit
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&audit).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment audit not found"})
	}

	return c.JSON(audit)
}

// ApprovePaymentAudit approves a payment audit and updates the related order
func ApprovePaymentAudit(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment audit ID"})
	}

	var audit models.PaymentAudit
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&audit).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment audit not found"})
	}

	if audit.Status != "pending" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Payment audit must be in 'pending' status to be approved, current status: " + audit.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&audit).Updates(map[string]interface{}{
		"status":         "approved",
		"reviewed_by_id": userID,
		"reviewed_at":    now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve payment audit"})
	}

	// Update the related order's payment status
	database.DB.Model(&models.ConsolidationOrder{}).
		Where("id = ? AND tenant_id = ?", audit.OrderID, tenantID).
		Updates(map[string]interface{}{
			"payment_status": "paid",
			"paid_at":        now,
		})

	audit.Status = "approved"
	audit.ReviewedByID = &userID
	audit.ReviewedAt = &now
	return c.JSON(audit)
}

// AuditPaymentAudit handles the unified audit action (approve/reject) from frontend
func AuditPaymentAudit(c *fiber.Ctx) error {
	var body struct {
		Status string `json:"status"`
		Remark string `json:"remark"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.Status == "approved" {
		return ApprovePaymentAudit(c)
	} else if body.Status == "rejected" {
		// Store remark as reject_reason — use json.Marshal for safe encoding
		rejectBody, _ := json.Marshal(map[string]string{"reject_reason": body.Remark})
		c.Request().SetBody(rejectBody)
		return RejectPaymentAudit(c)
	}
	return c.Status(400).JSON(fiber.Map{"error": "Invalid audit status, must be 'approved' or 'rejected'"})
}

// GetPaymentAuditCount returns count of audits with optional status filter
func GetPaymentAuditCount(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	query := database.DB.Model(&models.PaymentAudit{}).Where("tenant_id = ?", tenantID)
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	var count int64
	query.Count(&count)
	return c.JSON(fiber.Map{"count": count})
}

// RejectPaymentAudit rejects a payment audit
func RejectPaymentAudit(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment audit ID"})
	}

	var body struct {
		RejectReason string `json:"reject_reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var audit models.PaymentAudit
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&audit).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment audit not found"})
	}

	if audit.Status != "pending" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Payment audit must be in 'pending' status to be rejected, current status: " + audit.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&audit).Updates(map[string]interface{}{
		"status":         "rejected",
		"reviewed_by_id": userID,
		"reviewed_at":    now,
		"reject_reason":  body.RejectReason,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject payment audit"})
	}

	audit.Status = "rejected"
	audit.ReviewedByID = &userID
	audit.ReviewedAt = &now
	audit.RejectReason = body.RejectReason
	return c.JSON(audit)
}
