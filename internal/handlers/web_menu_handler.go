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

// GetWebMenus lists web menus with pagination
func GetWebMenus(c *fiber.Ctx) error {
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
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	var total int64
	query.Model(&models.WebMenu{}).Count(&total)

	var items []models.WebMenu
	if err := query.Order("sort_order ASC, created_at ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch web menus"})
	}

	return c.JSON(fiber.Map{"data": items, "total": total, "page": page, "limit": limit})
}

// GetWebMenu gets a single web menu by ID
func GetWebMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WebMenu
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Web menu not found"})
	}
	return c.JSON(item)
}

// CreateWebMenu creates a new web menu
func CreateWebMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	var item models.WebMenu
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	item.TenantID = &tenantID
	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create web menu"})
	}
	return c.Status(201).JSON(item)
}

// UpdateWebMenu updates an existing web menu
func UpdateWebMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WebMenu
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Web menu not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update web menu"})
	}
	return c.JSON(item)
}

// DeleteWebMenu soft-deletes a web menu
func DeleteWebMenu(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WebMenu
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Web menu not found"})
	}
	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete web menu"})
	}
	return c.JSON(fiber.Map{"message": "Web menu deleted successfully"})
}

// SortWebMenus updates sort order for multiple menus
func SortWebMenus(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	var sortData []struct {
		ID        string `json:"id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&sortData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	for _, s := range sortData {
		database.DB.Model(&models.WebMenu{}).Where("id = ? AND tenant_id = ?", s.ID, tenantID).Update("sort_order", s.SortOrder)
	}
	return c.JSON(fiber.Map{"message": "Sort order updated successfully"})
}
