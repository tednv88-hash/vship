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

// GetRechargeOrders returns a paginated list of recharge orders
func GetRechargeOrders(c *fiber.Ctx) error {
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

	var total int64
	var items []models.RechargeOrder

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.RechargeOrder{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetRechargeOrder returns a single recharge order by ID
func GetRechargeOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.RechargeOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Recharge order not found"})
	}

	return c.JSON(item)
}

// CreateRechargeOrder creates a new recharge order
func CreateRechargeOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.RechargeOrder
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create recharge order"})
	}

	return c.Status(201).JSON(item)
}

// GetRechargePlans returns a paginated list of recharge plans
func GetRechargePlans(c *fiber.Ctx) error {
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

	var total int64
	var items []models.RechargePlan

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.RechargePlan{}).Count(&total)
	query.Order("sort_order ASC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetRechargePlan returns a single recharge plan by ID
func GetRechargePlan(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.RechargePlan
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Recharge plan not found"})
	}

	return c.JSON(item)
}

// CreateRechargePlan creates a new recharge plan
func CreateRechargePlan(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.RechargePlan
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create recharge plan"})
	}

	return c.Status(201).JSON(item)
}

// UpdateRechargePlan updates an existing recharge plan
func UpdateRechargePlan(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.RechargePlan
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Recharge plan not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")
	delete(updates, "trashed_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update recharge plan"})
	}

	return c.JSON(item)
}

// DeleteRechargePlan soft-deletes a recharge plan
func DeleteRechargePlan(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.RechargePlan
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Recharge plan not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete recharge plan"})
	}

	return c.JSON(fiber.Map{"message": "Recharge plan deleted successfully"})
}

// GetRechargeSetting returns the recharge setting for the current tenant
func GetRechargeSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.Setting
	if err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, "recharge_setting").First(&item).Error; err != nil {
		return c.JSON(fiber.Map{"key": "recharge_setting", "value": ""})
	}

	return c.JSON(item)
}

// UpdateRechargeSetting creates or updates the recharge setting
func UpdateRechargeSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var body struct {
		Value string `json:"value"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var existing models.Setting
	err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, "recharge_setting").First(&existing).Error
	if err != nil {
		// Create new
		setting := models.Setting{
			TenantID: &tenantID,
			Key:      "recharge_setting",
			Value:    body.Value,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create recharge setting"})
		}
		return c.JSON(setting)
	}

	// Update existing
	if err := database.DB.Model(&existing).Update("value", body.Value).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update recharge setting"})
	}

	return c.JSON(existing)
}
