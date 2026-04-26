package handlers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetOrders lists orders with pagination and filters
func GetOrders(c *fiber.Ctx) error {
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

	query := database.DB.Model(&models.ConsolidationOrder{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status == "overdue" {
		// "overdue" is not a real DB status — it means orders in non-terminal
		// status (draft/paid/packed) created more than N days ago.
		days := 7
		if d, err := strconv.Atoi(c.Query("overdue_days")); err == nil && d > 0 {
			days = d
		}
		cutoff := time.Now().AddDate(0, 0, -days)
		query = query.Where("status NOT IN ('completed','cancelled','shipped') AND created_at < ?", cutoff)
	} else if status != "" {
		query = query.Where("status = ?", status)
	}
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if orderNumber := c.Query("order_number"); orderNumber != "" {
		query = query.Where("order_number = ?", orderNumber)
	}

	var total int64
	query.Count(&total)

	var orders []models.ConsolidationOrder
	if err := query.
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.JSON(fiber.Map{
		"data":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetOrder gets a single order by ID
func GetOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(order)
}

// CreateOrder creates a new order with auto-generated order number
func CreateOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var order models.ConsolidationOrder
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	order.TenantID = &tenantID
	order.Status = "draft"
	order.OrderNumber = fmt.Sprintf("ORD%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000))

	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}

	return c.Status(201).JSON(order)
}

// UpdateOrder updates an existing order
func UpdateOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	originalOrderNumber := order.OrderNumber
	originalStatus := order.Status

	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	order.ID = id
	order.TenantID = &tenantID
	order.OrderNumber = originalOrderNumber
	order.Status = originalStatus

	if err := database.DB.Save(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update order"})
	}

	return c.JSON(order)
}

// DeleteOrder soft-deletes an order
func DeleteOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if err := database.DB.Model(&order).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete order"})
	}

	return c.JSON(fiber.Map{"message": "Order deleted successfully"})
}

// PayOrder sets order status to "paid"
func PayOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Order must be in 'draft' status to be paid, current status: " + order.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&order).Updates(map[string]interface{}{
		"status":         "paid",
		"payment_status": "paid",
		"paid_at":        now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to pay order"})
	}

	order.Status = "paid"
	order.PaymentStatus = "paid"
	order.PaidAt = &now
	return c.JSON(order)
}

// PackOrder sets order status to "packed"
func PackOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.Status != "paid" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Order must be in 'paid' status to be packed, current status: " + order.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&order).Updates(map[string]interface{}{
		"status":    "packed",
		"packed_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to pack order"})
	}

	order.Status = "packed"
	order.PackedAt = &now
	return c.JSON(order)
}

// ShipOrder sets order status to "shipped"
func ShipOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.Status != "packed" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Order must be in 'packed' status to be shipped, current status: " + order.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&order).Updates(map[string]interface{}{
		"status":     "shipped",
		"shipped_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to ship order"})
	}

	order.Status = "shipped"
	order.ShippedAt = &now
	return c.JSON(order)
}

// CompleteOrder sets order status to "completed"
func CompleteOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.Status != "shipped" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Order must be in 'shipped' status to be completed, current status: " + order.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&order).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to complete order"})
	}

	order.Status = "completed"
	order.CompletedAt = &now
	return c.JSON(order)
}

// CancelOrder sets order status to "cancelled"
func CancelOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var order models.ConsolidationOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.Status == "completed" || order.Status == "cancelled" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Order cannot be cancelled from '" + order.Status + "' status",
		})
	}

	now := time.Now()
	if err := database.DB.Model(&order).Updates(map[string]interface{}{
		"status":       "cancelled",
		"cancelled_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel order"})
	}

	order.Status = "cancelled"
	order.CancelledAt = &now
	return c.JSON(order)
}

// GetOrderPackages lists packages linked to an order
func GetOrderPackages(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	orderID, err := uuid.Parse(c.Params("orderId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var orderPackages []models.OrderPackage
	if err := database.DB.
		Where("tenant_id = ? AND order_id = ?", tenantID, orderID).
		Order("created_at ASC").
		Find(&orderPackages).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch order packages"})
	}

	return c.JSON(fiber.Map{"data": orderPackages})
}

// AddOrderPackage links a package to an order
func AddOrderPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	orderID, err := uuid.Parse(c.Params("orderId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var body struct {
		PackageID uuid.UUID `json:"package_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	orderPackage := models.OrderPackage{
		TenantID:  &tenantID,
		OrderID:   orderID,
		PackageID: body.PackageID,
	}

	if err := database.DB.Create(&orderPackage).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add package to order"})
	}

	return c.Status(201).JSON(orderPackage)
}

// RemoveOrderPackage removes a package link from an order
func RemoveOrderPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order package ID"})
	}

	var orderPackage models.OrderPackage
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&orderPackage).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order package link not found"})
	}

	if err := database.DB.Delete(&orderPackage).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove package from order"})
	}

	return c.JSON(fiber.Map{"message": "Package removed from order successfully"})
}
