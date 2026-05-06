package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCouponReceiveLogs returns a paginated list of coupon receive logs
func GetCouponReceiveLogs(c *fiber.Ctx) error {
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
	var items []models.CouponReceiveLog

	query := database.DB.Where("tenant_id = ?", tenantID)

	if couponID := c.Query("coupon_id"); couponID != "" {
		query = query.Where("coupon_id = ?", couponID)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.CouponReceiveLog{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetCouponReceiveLog returns a single coupon receive log by ID
func GetCouponReceiveLog(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.CouponReceiveLog
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Coupon receive log not found"})
	}

	return c.JSON(item)
}

// CreateCouponReceiveLog creates a new coupon receive log entry
func CreateCouponReceiveLog(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.CouponReceiveLog
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create coupon receive log"})
	}

	return c.Status(201).JSON(item)
}

// UpdateCouponReceiveLog updates a coupon receive log
func UpdateCouponReceiveLog(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.CouponReceiveLog
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Coupon receive log not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update coupon receive log"})
	}

	return c.JSON(item)
}

// DeleteCouponReceiveLog hard-deletes a coupon receive log
func DeleteCouponReceiveLog(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	result := database.DB.Delete(&models.CouponReceiveLog{}, "id = ? AND tenant_id = ?", id, tenantID)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete coupon receive log"})
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Coupon receive log not found"})
	}

	return c.JSON(fiber.Map{"message": "Coupon receive log deleted successfully"})
}
