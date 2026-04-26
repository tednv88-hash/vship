package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetWarehouseApplications lists warehouse applications with pagination and filters
func GetWarehouseApplications(c *fiber.Ctx) error {
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
	if search := c.Query("search"); search != "" {
		query = query.Where("(applicant_name ILIKE ? OR warehouse_name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Model(&models.WarehouseApplication{}).Count(&total)

	var items []models.WarehouseApplication
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch warehouse applications"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetWarehouseApplication gets a single warehouse application by ID
func GetWarehouseApplication(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WarehouseApplication
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse application not found"})
	}

	return c.JSON(item)
}

// ApproveWarehouseApplication approves a warehouse application
func ApproveWarehouseApplication(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WarehouseApplication
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse application not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status": "approved",
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve warehouse application"})
	}

	item.Status = "approved"
	return c.JSON(item)
}

// RejectWarehouseApplication rejects a warehouse application
func RejectWarehouseApplication(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
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

	var item models.WarehouseApplication
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse application not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":       "rejected",
		"audit_remark": body.AuditRemark,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject warehouse application"})
	}

	item.Status = "rejected"
	item.AuditRemark = body.AuditRemark
	return c.JSON(item)
}
