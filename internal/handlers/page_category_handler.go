package handlers

import (
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetPageCategory gets the page category settings
func GetPageCategory(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.PageCategory
	if err := database.DB.Where("tenant_id = ?", tenantID).First(&item).Error; err != nil {
		// Return default settings if none exist
		return c.JSON(models.PageCategory{
			CategoryStyle: "20",
			ShareTitle:    "",
		})
	}
	return c.JSON(item)
}

// SavePageCategory creates or updates page category settings
func SavePageCategory(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var input struct {
		CategoryStyle string `json:"category_style"`
		ShareTitle    string `json:"share_title"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.PageCategory
	result := database.DB.Where("tenant_id = ?", tenantID).First(&item)
	if result.Error != nil {
		// Create new
		item.TenantID = &tenantID
		item.CategoryStyle = input.CategoryStyle
		item.ShareTitle = input.ShareTitle
		if err := database.DB.Create(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save page category settings"})
		}
	} else {
		// Update existing
		if err := database.DB.Model(&item).Updates(map[string]interface{}{
			"category_style": input.CategoryStyle,
			"share_title":    input.ShareTitle,
		}).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update page category settings"})
		}
	}

	return c.JSON(item)
}
