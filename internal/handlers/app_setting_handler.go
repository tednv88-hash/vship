package handlers

import (
	"encoding/json"
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAppSettings lists app settings with pagination and filters
func GetAppSettings(c *fiber.Ctx) error {
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

	if settingType := c.Query("setting_type"); settingType != "" {
		query = query.Where("setting_type = ?", settingType)
	}

	var total int64
	query.Model(&models.AppSetting{}).Count(&total)

	var items []models.AppSetting
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch app settings"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetAppSetting gets a single app setting by ID
func GetAppSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.AppSetting
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "App setting not found"})
	}

	return c.JSON(item)
}

// CreateAppSetting creates a new app setting
func CreateAppSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.AppSetting
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	mergeAppSettingFormFields(c, &item)

	item.TenantID = &tenantID
	if item.SettingType != "" {
		var existing models.AppSetting
		if err := database.DB.Where("tenant_id = ? AND setting_type = ?", tenantID, item.SettingType).First(&existing).Error; err == nil {
			item.ID = existing.ID
			if err := database.DB.Model(&existing).Updates(map[string]interface{}{
				"setting_type": item.SettingType,
				"config":       item.Config,
				"extra_fields": item.ExtraFields,
			}).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update app setting"})
			}
			return c.JSON(existing)
		}
	}

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create app setting"})
	}

	return c.Status(201).JSON(item)
}

// UpdateAppSetting updates an existing app setting
func UpdateAppSetting(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.AppSetting
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "App setting not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	normalizeAppSettingUpdates(c, updates)

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update app setting"})
	}

	return c.JSON(item)
}

func mergeAppSettingFormFields(c *fiber.Ctx, item *models.AppSetting) {
	var raw map[string]interface{}
	if err := json.Unmarshal(c.Body(), &raw); err != nil {
		return
	}
	if item.Config == nil {
		item.Config = models.JSONB{}
	}
	if item.ExtraFields == nil {
		item.ExtraFields = models.JSONB{}
	}
	for key, value := range raw {
		switch key {
		case "id", "tenant_id", "setting_type", "config", "created_at", "updated_at", "extra_fields":
			continue
		default:
			item.Config[key] = value
		}
	}
}

func normalizeAppSettingUpdates(c *fiber.Ctx, updates map[string]interface{}) {
	var item models.AppSetting
	if err := c.BodyParser(&item); err == nil {
		mergeAppSettingFormFields(c, &item)
		if len(item.Config) > 0 {
			updates["config"] = item.Config
		}
		if len(item.ExtraFields) > 0 {
			updates["extra_fields"] = item.ExtraFields
		}
	}
}
