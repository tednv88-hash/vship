package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetSmsLogs lists SMS logs with pagination and filters (log table, no trashed_at)
func GetSmsLogs(c *fiber.Ctx) error {
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

	if phone := c.Query("phone"); phone != "" {
		query = query.Where("phone = ?", phone)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if t := c.Query("type"); t != "" {
		query = query.Where("type = ?", t)
	}

	var total int64
	query.Model(&models.SmsLog{}).Count(&total)

	var items []models.SmsLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch SMS logs"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// SendSms creates an SMS log with status='pending'
func SendSms(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.SmsLog
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID
	item.Status = "pending"

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create SMS log"})
	}

	return c.Status(201).JSON(item)
}
