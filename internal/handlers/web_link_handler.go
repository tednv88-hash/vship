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

// GetWebLinks lists web links with pagination
func GetWebLinks(c *fiber.Ctx) error {
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

	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if linkType := c.Query("link_type"); linkType != "" {
		query = query.Where("link_type = ?", linkType)
	}

	var total int64
	query.Model(&models.WebLink{}).Count(&total)

	var items []models.WebLink
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch web links"})
	}
	return c.JSON(fiber.Map{"data": items, "total": total, "page": page, "limit": limit})
}

// GetWebLink gets a single web link by ID
func GetWebLink(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WebLink
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Web link not found"})
	}
	return c.JSON(item)
}

// CreateWebLink creates a new web link
func CreateWebLink(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	var item models.WebLink
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	item.TenantID = &tenantID
	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create web link"})
	}
	return c.Status(201).JSON(item)
}

// UpdateWebLink updates an existing web link
func UpdateWebLink(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WebLink
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Web link not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update web link"})
	}
	return c.JSON(item)
}

// DeleteWebLink soft-deletes a web link
func DeleteWebLink(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.WebLink
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Web link not found"})
	}
	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete web link"})
	}
	return c.JSON(fiber.Map{"message": "Web link deleted successfully"})
}
