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

// GetWarehouseCapitalLogs lists warehouse capital logs with pagination and filters (log table)
func GetWarehouseCapitalLogs(c *fiber.Ctx) error {
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

	query := database.DB.Where("tenant_id = ?", tenantID)

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	if t := c.Query("type"); t != "" {
		query = query.Where("type = ?", t)
	}

	var total int64
	query.Model(&models.WarehouseCapitalLog{}).Count(&total)

	var items []models.WarehouseCapitalLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch warehouse capital logs"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetWarehouseBonuses lists warehouse bonuses with pagination and filters
func GetWarehouseBonuses(c *fiber.Ctx) error {
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

	query := database.DB.Where("tenant_id = ?", tenantID)

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if month := c.Query("month"); month != "" {
		query = query.Where("month = ?", month)
	}

	var total int64
	query.Model(&models.WarehouseBonus{}).Count(&total)

	var items []models.WarehouseBonus
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch warehouse bonuses"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateWarehouseBonus creates a new warehouse bonus
func CreateWarehouseBonus(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.WarehouseBonus
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create warehouse bonus"})
	}

	return c.Status(201).JSON(item)
}

// UpdateWarehouseBonus updates an existing warehouse bonus
func UpdateWarehouseBonus(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WarehouseBonus
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse bonus not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update warehouse bonus"})
	}

	return c.JSON(item)
}

// PayWarehouseBonus marks a warehouse bonus as paid
func PayWarehouseBonus(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WarehouseBonus
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse bonus not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to pay warehouse bonus"})
	}

	item.Status = "paid"
	item.PaidAt = &now
	return c.JSON(item)
}

// GetWarehouseWithdrawals lists warehouse withdrawals with pagination and filters
func GetWarehouseWithdrawals(c *fiber.Ctx) error {
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

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.WarehouseWithdrawal{}).Count(&total)

	var items []models.WarehouseWithdrawal
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch warehouse withdrawals"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// ApproveWarehouseWithdrawal approves a warehouse withdrawal
func ApproveWarehouseWithdrawal(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WarehouseWithdrawal
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse withdrawal not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status": "approved",
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve warehouse withdrawal"})
	}

	item.Status = "approved"
	return c.JSON(item)
}

// RejectWarehouseWithdrawal rejects a warehouse withdrawal
func RejectWarehouseWithdrawal(c *fiber.Ctx) error {
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

	var item models.WarehouseWithdrawal
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse withdrawal not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":       "rejected",
		"audit_remark": body.AuditRemark,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject warehouse withdrawal"})
	}

	item.Status = "rejected"
	item.AuditRemark = body.AuditRemark
	return c.JSON(item)
}

// PayWarehouseWithdrawal marks a warehouse withdrawal as paid
func PayWarehouseWithdrawal(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WarehouseWithdrawal
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse withdrawal not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to pay warehouse withdrawal"})
	}

	item.Status = "paid"
	item.PaidAt = &now
	return c.JSON(item)
}
