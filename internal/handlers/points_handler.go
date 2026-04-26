package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetPointsLogs returns a paginated list of points logs
func GetPointsLogs(c *fiber.Ctx) error {
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
	var items []models.PointsLog

	query := database.DB.Where("tenant_id = ?", tenantID)

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if logType := c.Query("type"); logType != "" {
		query = query.Where("type = ?", logType)
	}

	query.Model(&models.PointsLog{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreatePointsLog creates a new points log entry
func CreatePointsLog(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.PointsLog
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create points log"})
	}

	return c.Status(201).JSON(item)
}

// GetPointsSetting returns the points setting for the current tenant
func GetPointsSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.Setting
	if err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, "points_setting").First(&item).Error; err != nil {
		return c.JSON(fiber.Map{"key": "points_setting", "value": ""})
	}

	return c.JSON(item)
}

// UpdatePointsSetting creates or updates the points setting
func UpdatePointsSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var body struct {
		Value string `json:"value"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var existing models.Setting
	err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, "points_setting").First(&existing).Error
	if err != nil {
		// Create new
		setting := models.Setting{
			TenantID: &tenantID,
			Key:      "points_setting",
			Value:    body.Value,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create points setting"})
		}
		return c.JSON(setting)
	}

	// Update existing
	if err := database.DB.Model(&existing).Update("value", body.Value).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update points setting"})
	}

	return c.JSON(existing)
}
