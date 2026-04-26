package handlers

import (
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"
	"vship/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SetupTenant handles initial tenant creation
func SetupTenant(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return c.Status(401).JSON(fiber.Map{"error": "Not authenticated"})
	}

	// Check if user already has a tenant
	var existingUser models.User
	if err := database.DB.First(&existingUser, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	if existingUser.TenantID != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User already has a tenant"})
	}

	var req struct {
		Name      string `json:"name"`
		Subdomain string `json:"subdomain"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.Subdomain == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and subdomain are required"})
	}

	// Check if subdomain is taken
	var count int64
	database.DB.Model(&models.Tenant{}).Where("subdomain = ? AND trashed_at IS NULL", req.Subdomain).Count(&count)
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{"error": "Subdomain already taken"})
	}

	tenant := models.Tenant{
		Name:      req.Name,
		Subdomain: req.Subdomain,
		Plan:      "free",
		Status:    "active",
	}

	tx := database.DB.Begin()

	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create tenant"})
	}

	// Associate user with tenant
	if err := tx.Model(&existingUser).Updates(map[string]interface{}{
		"tenant_id": tenant.ID,
	}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user with tenant"})
	}

	// Generate a new token with the tenant ID
	token, err := utils.GenerateToken(existingUser.ID, tenant.ID, existingUser.Email, existingUser.UserRole)
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate new token"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Transaction failed"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Tenant created successfully",
		"tenant":  tenant,
		"token":   token,
	})
}
