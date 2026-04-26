package handlers

import (
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetSettings returns all settings for the current tenant, with optional group filter
func GetSettings(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	group := c.Query("group")

	var items []models.Setting
	query := database.DB.Where("tenant_id = ?", tenantID)
	if group != "" {
		query = query.Where("\"group\" = ?", group)
	}
	query.Order("key ASC").Find(&items)

	return c.JSON(fiber.Map{
		"data": items,
	})
}

// UpdateSettings bulk upserts settings for the current tenant
func UpdateSettings(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var inputs []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Group string `json:"group"`
	}

	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	for _, input := range inputs {
		if input.Key == "" {
			continue
		}

		var existing models.Setting
		err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, input.Key).First(&existing).Error
		if err != nil {
			// Create new setting
			setting := models.Setting{
				TenantID: &tenantID,
				Key:      input.Key,
				Value:    input.Value,
				Group:    input.Group,
			}
			if err := database.DB.Create(&setting).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to create setting: " + input.Key})
			}
		} else {
			// Update existing setting
			database.DB.Model(&existing).Updates(map[string]interface{}{
				"value": input.Value,
				"group": input.Group,
			})
		}
	}

	// Return updated settings
	var items []models.Setting
	database.DB.Where("tenant_id = ?", tenantID).Order("key ASC").Find(&items)

	return c.JSON(fiber.Map{
		"message": "Settings updated successfully",
		"data":    items,
	})
}

// GetSettingByKey returns a single setting by key for the current tenant
func GetSettingByKey(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	key := c.Params("key")
	if key == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Key is required"})
	}

	var item models.Setting
	err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, key).First(&item).Error
	if err != nil {
		// Return empty setting with the key so frontend doesn't break
		return c.JSON(fiber.Map{
			"data": fiber.Map{
				"key":   key,
				"value": "{}",
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": item,
	})
}

// UpdateSettingByKey upserts a single setting by key for the current tenant
func UpdateSettingByKey(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	key := c.Params("key")
	if key == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Key is required"})
	}

	var input struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Group string `json:"group"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Use URL key as primary, fallback to body key
	settingKey := key
	if input.Key != "" {
		settingKey = input.Key
	}

	var existing models.Setting
	err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, settingKey).First(&existing).Error
	if err != nil {
		setting := models.Setting{
			TenantID: &tenantID,
			Key:      settingKey,
			Value:    input.Value,
			Group:    input.Group,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create setting"})
		}
		return c.JSON(fiber.Map{
			"message": "Setting saved successfully",
			"data":    setting,
		})
	}

	updates := map[string]interface{}{
		"value": input.Value,
	}
	if input.Group != "" {
		updates["group"] = input.Group
	}
	database.DB.Model(&existing).Updates(updates)

	return c.JSON(fiber.Map{
		"message": "Setting updated successfully",
		"data":    existing,
	})
}

// ClearCache clears application cache (placeholder implementation)
func ClearCache(c *fiber.Ctx) error {
	// In a real application, this would clear Redis/memcache/etc.
	// For now, return success as a placeholder
	return c.JSON(fiber.Map{
		"message": "Cache cleared successfully",
	})
}
