package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetEmailLogs lists email logs with pagination and filters (log table, no trashed_at)
func GetEmailLogs(c *fiber.Ctx) error {
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

	if email := c.Query("email"); email != "" {
		query = query.Where("email = ?", email)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if t := c.Query("type"); t != "" {
		query = query.Where("type = ?", t)
	}

	var total int64
	query.Model(&models.EmailLog{}).Count(&total)

	var items []models.EmailLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch email logs"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// SendEmail creates an email log with status='pending'
func SendEmail(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.EmailLog
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID
	item.Status = "pending"

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create email log"})
	}

	return c.Status(201).JSON(item)
}
