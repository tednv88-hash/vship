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

// GetShopOrders returns a paginated list of shop orders with optional filters
func GetShopOrders(c *fiber.Ctx) error {
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
	var items []models.ShopOrder

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if search := c.Query("search"); search != "" {
		query = query.Where("order_no ILIKE ?", "%"+search+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Model(&models.ShopOrder{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetShopOrder returns a single shop order by ID
func GetShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.ShopOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shop order not found"})
	}

	return c.JSON(item)
}

// CreateShopOrder creates a new shop order
func CreateShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.ShopOrder
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create shop order"})
	}

	return c.Status(201).JSON(item)
}

// UpdateShopOrder updates an existing shop order
func UpdateShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.ShopOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shop order not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")
	delete(updates, "trashed_at")

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update shop order"})
	}

	return c.JSON(item)
}

// DeleteShopOrder soft-deletes a shop order
func DeleteShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.ShopOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shop order not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete shop order"})
	}

	return c.JSON(fiber.Map{"message": "Shop order deleted successfully"})
}

// ShipShopOrder marks a shop order as shipped with express info
func ShipShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.ShopOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shop order not found"})
	}

	var body struct {
		ExpressCompany string `json:"express_company"`
		ExpressNo      string `json:"express_no"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":          "shipped",
		"express_company": body.ExpressCompany,
		"express_no":      body.ExpressNo,
		"delivery_time":   now,
	}

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to ship order"})
	}

	return c.JSON(item)
}

// CompleteShopOrder marks a shop order as completed
func CompleteShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.ShopOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shop order not found"})
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":       "completed",
		"receive_time": now,
	}

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to complete order"})
	}

	return c.JSON(item)
}

// CancelShopOrder marks a shop order as cancelled
func CancelShopOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.ShopOrder
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shop order not found"})
	}

	if err := database.DB.Model(&item).Update("status", "cancelled").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel order"})
	}

	return c.JSON(item)
}
