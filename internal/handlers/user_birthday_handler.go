package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUserBirthdays lists user birthdays with pagination and filters
func GetUserBirthdays(c *fiber.Ctx) error {
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

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Model(&models.UserBirthday{}).Count(&total)

	var items []models.UserBirthday
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user birthdays"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUserBirthday gets a single user birthday by ID
func GetUserBirthday(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.UserBirthday
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User birthday not found"})
	}

	return c.JSON(item)
}

// CreateUserBirthday creates a new user birthday
func CreateUserBirthday(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.UserBirthday
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user birthday"})
	}

	return c.Status(201).JSON(item)
}

// UpdateUserBirthday updates an existing user birthday
func UpdateUserBirthday(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.UserBirthday
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User birthday not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user birthday"})
	}

	return c.JSON(item)
}

// DeleteUserBirthday hard-deletes a user birthday reminder
func DeleteUserBirthday(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	result := database.DB.Delete(&models.UserBirthday{}, "id = ? AND tenant_id = ?", id, tenantID)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user birthday"})
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User birthday not found"})
	}

	return c.JSON(fiber.Map{"message": "User birthday deleted successfully"})
}
