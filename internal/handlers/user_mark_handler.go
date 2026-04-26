package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUserMarks returns a paginated list of user marks
func GetUserMarks(c *fiber.Ctx) error {
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
	var items []models.UserMark

	query := database.DB.Where("tenant_id = ?", tenantID)

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Model(&models.UserMark{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateUserMark creates a new user mark
func CreateUserMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.UserMark
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user mark"})
	}

	return c.Status(201).JSON(item)
}

// DeleteUserMark hard-deletes a user mark
func DeleteUserMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	result := database.DB.Delete(&models.UserMark{}, "id = ? AND tenant_id = ?", id, tenantID)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user mark"})
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User mark not found"})
	}

	return c.JSON(fiber.Map{"message": "User mark deleted successfully"})
}

// UpdateUserMark updates an existing user mark
func UpdateUserMark(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.UserMark
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User mark not found"})
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(body, "id")
	delete(body, "tenant_id")
	delete(body, "created_at")

	if err := database.DB.Model(&item).Updates(body).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user mark"})
	}

	database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item)
	return c.JSON(item)
}
