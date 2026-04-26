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

// GetWechatMenus lists wechat menus
func GetWechatMenus(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if parentID := c.Query("parent_id"); parentID != "" {
		if parentID == "null" || parentID == "0" {
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", parentID)
		}
	}

	var total int64
	query.Model(&models.WechatMenu{}).Count(&total)

	var items []models.WechatMenu
	if err := query.Order("sort_order ASC, created_at ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch wechat menus"})
	}
	return c.JSON(fiber.Map{"data": items, "total": total, "page": page, "limit": limit})
}

// GetWechatMenu gets a single wechat menu by ID
func GetWechatMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WechatMenu
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Wechat menu not found"})
	}
	return c.JSON(item)
}

// CreateWechatMenu creates a new wechat menu
func CreateWechatMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	var item models.WechatMenu
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	item.TenantID = &tenantID
	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create wechat menu"})
	}
	return c.Status(201).JSON(item)
}

// UpdateWechatMenu updates an existing wechat menu
func UpdateWechatMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WechatMenu
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Wechat menu not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update wechat menu"})
	}
	return c.JSON(item)
}

// DeleteWechatMenu soft-deletes a wechat menu
func DeleteWechatMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WechatMenu
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Wechat menu not found"})
	}
	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete wechat menu"})
	}
	return c.JSON(fiber.Map{"message": "Wechat menu deleted successfully"})
}

// PublishWechatMenu publishes the current wechat menu configuration
func PublishWechatMenu(c *fiber.Ctx) error {
	// In a real application, this would call the WeChat API to publish the menu
	// For now, return success as a placeholder
	return c.JSON(fiber.Map{
		"message": "Wechat menu published successfully",
	})
}
