package handlers

import (
	"strconv"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetPaymentFlows lists payment flows with pagination and filters (log table, read-only)
func GetPaymentFlows(c *fiber.Ctx) error {
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

	if orderType := c.Query("order_type"); orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}
	if orderNo := c.Query("order_no"); orderNo != "" {
		query = query.Where("order_no = ?", orderNo)
	}
	if payMethod := c.Query("pay_method"); payMethod != "" {
		query = query.Where("pay_method = ?", payMethod)
	}

	var total int64
	query.Model(&models.PaymentFlow{}).Count(&total)

	var items []models.PaymentFlow
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch payment flows"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetPaymentFlow gets a single payment flow by ID
func GetPaymentFlow(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.PaymentFlow
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment flow not found"})
	}

	return c.JSON(item)
}
