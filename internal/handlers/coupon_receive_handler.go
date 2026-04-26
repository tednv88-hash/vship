package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
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
