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

// GetShippingMarks lists shipping marks with pagination
func GetShippingMarks(c *fiber.Ctx) error {
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
	database.DB.Model(&models.ShippingMark{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Count(&total)

	var marks []models.ShippingMark
	if err := database.DB.
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Order("code ASC").
		Offset(offset).Limit(limit).
		Find(&marks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch shipping marks"})
	}

	return c.JSON(fiber.Map{
		"data":  marks,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetShippingMark gets a single shipping mark by ID
func GetShippingMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid shipping mark ID"})
	}

	var mark models.ShippingMark
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&mark).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping mark not found"})
	}

	return c.JSON(mark)
}

// CreateShippingMark creates a new shipping mark
func CreateShippingMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var mark models.ShippingMark
	if err := c.BodyParser(&mark); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	mark.TenantID = &tenantID

	if err := database.DB.Create(&mark).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create shipping mark"})
	}

	return c.Status(201).JSON(mark)
}

// UpdateShippingMark updates an existing shipping mark
func UpdateShippingMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid shipping mark ID"})
	}

	var mark models.ShippingMark
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&mark).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping mark not found"})
	}

	if err := c.BodyParser(&mark); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	mark.ID = id
	mark.TenantID = &tenantID

	if err := database.DB.Save(&mark).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update shipping mark"})
	}

	return c.JSON(mark)
}

// DeleteShippingMark soft-deletes a shipping mark
func DeleteShippingMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid shipping mark ID"})
	}

	var mark models.ShippingMark
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&mark).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping mark not found"})
	}

	if err := database.DB.Model(&mark).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete shipping mark"})
	}

	return c.JSON(fiber.Map{"message": "Shipping mark deleted successfully"})
}
