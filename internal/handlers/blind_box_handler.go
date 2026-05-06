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

// GetBlindBoxActivities returns a paginated list of blind box activities
func GetBlindBoxActivities(c *fiber.Ctx) error {
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
	var items []models.BlindBoxActivity

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.BlindBoxActivity{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetBlindBoxActivity returns a single blind box activity by ID
func GetBlindBoxActivity(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.BlindBoxActivity
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blind box activity not found"})
	}

	return c.JSON(item)
}

// CreateBlindBoxActivity creates a new blind box activity
func CreateBlindBoxActivity(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.BlindBoxActivity
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create blind box activity"})
	}

	return c.Status(201).JSON(item)
}

// UpdateBlindBoxActivity updates an existing blind box activity
func UpdateBlindBoxActivity(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.BlindBoxActivity
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blind box activity not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update blind box activity"})
	}

	return c.JSON(item)
}

// DeleteBlindBoxActivity soft-deletes a blind box activity
func DeleteBlindBoxActivity(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.BlindBoxActivity
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blind box activity not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete blind box activity"})
	}

	return c.JSON(fiber.Map{"message": "Blind box activity deleted successfully"})
}

// GetBlindBoxDraws returns a paginated list of blind box draws
func GetBlindBoxDraws(c *fiber.Ctx) error {
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
	var items []models.BlindBoxDraw

	query := database.DB.Where("tenant_id = ?", tenantID)

	if activityID := c.Query("activity_id"); activityID != "" {
		query = query.Where("activity_id = ?", activityID)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Model(&models.BlindBoxDraw{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateBlindBoxDraw creates a new blind box draw record
func CreateBlindBoxDraw(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.BlindBoxDraw
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create blind box draw"})
	}

	return c.Status(201).JSON(item)
}

// UpdateBlindBoxDraw updates an existing blind box draw record
func UpdateBlindBoxDraw(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.BlindBoxDraw
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blind box draw not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update blind box draw"})
	}

	return c.JSON(item)
}

// DeleteBlindBoxDraw hard-deletes a blind box draw record
func DeleteBlindBoxDraw(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	result := database.DB.Delete(&models.BlindBoxDraw{}, "id = ? AND tenant_id = ?", id, tenantID)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete blind box draw"})
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Blind box draw not found"})
	}

	return c.JSON(fiber.Map{"message": "Blind box draw deleted successfully"})
}

// GetBlindBoxSetting returns the blind box setting for the current tenant
func GetBlindBoxSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.Setting
	if err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, "blindbox_setting").First(&item).Error; err != nil {
		return c.JSON(fiber.Map{"key": "blindbox_setting", "value": ""})
	}

	return c.JSON(item)
}

// UpdateBlindBoxSetting creates or updates the blind box setting
func UpdateBlindBoxSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var body struct {
		Value string `json:"value"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var existing models.Setting
	err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, "blindbox_setting").First(&existing).Error
	if err != nil {
		// Create new
		setting := models.Setting{
			TenantID: &tenantID,
			Key:      "blindbox_setting",
			Value:    body.Value,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create blind box setting"})
		}
		return c.JSON(setting)
	}

	// Update existing
	if err := database.DB.Model(&existing).Update("value", body.Value).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update blind box setting"})
	}

	return c.JSON(existing)
}
