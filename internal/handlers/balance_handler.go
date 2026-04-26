package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetBalanceLogs returns a paginated list of user balance logs
func GetBalanceLogs(c *fiber.Ctx) error {
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
	var items []models.UserBalanceLog

	query := database.DB.Where("tenant_id = ?", tenantID)

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if logType := c.Query("type"); logType != "" {
		query = query.Where("type = ?", logType)
	}

	query.Model(&models.UserBalanceLog{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateBalanceLog creates a new balance log entry
func CreateBalanceLog(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.UserBalanceLog
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create balance log"})
	}

	return c.Status(201).JSON(item)
}
