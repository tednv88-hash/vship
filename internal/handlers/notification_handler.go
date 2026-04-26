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

// GetNotifications lists notifications with pagination, search, type filter, and status filter
func GetNotifications(c *fiber.Ctx) error {
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

	search := c.Query("search")
	notifType := c.Query("type")
	status := c.Query("status")

	// Count query
	countQuery := database.DB.Model(&models.Notification{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID)
	if search != "" {
		countQuery = countQuery.Where("title ILIKE ?", "%"+search+"%")
	}
	if notifType != "" {
		countQuery = countQuery.Where("type = ?", notifType)
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	var total int64
	countQuery.Count(&total)

	// Data query
	query := database.DB.
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID)
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if notifType != "" {
		query = query.Where("type = ?", notifType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var notifications []models.Notification
	if err := query.
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	return c.JSON(fiber.Map{
		"data":  notifications,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetNotification gets a single notification by ID
func GetNotification(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	var notification models.Notification
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&notification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	return c.JSON(notification)
}

// CreateNotification creates a new notification
func CreateNotification(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var notification models.Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	notification.TenantID = &tenantID

	if err := database.DB.Create(&notification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create notification"})
	}

	return c.Status(201).JSON(notification)
}

// UpdateNotification updates an existing notification
func UpdateNotification(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	var notification models.Notification
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&notification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	notification.ID = id
	notification.TenantID = &tenantID

	if err := database.DB.Save(&notification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update notification"})
	}

	return c.JSON(notification)
}

// DeleteNotification soft-deletes a notification
func DeleteNotification(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	var notification models.Notification
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&notification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	if err := database.DB.Model(&notification).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification deleted successfully"})
}

// PublishNotification sets a notification status to 'published' and records the publish time
func PublishNotification(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	var notification models.Notification
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&notification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&notification).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to publish notification"})
	}

	notification.Status = "published"
	notification.PublishedAt = &now

	return c.JSON(notification)
}
