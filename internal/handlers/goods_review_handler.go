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

// GetGoodsReviews returns a paginated list of goods reviews with optional filters
func GetGoodsReviews(c *fiber.Ctx) error {
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
	var items []models.GoodsReview

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if goodsID := c.Query("goods_id"); goodsID != "" {
		query = query.Where("goods_id = ?", goodsID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.GoodsReview{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetGoodsReview returns a single goods review by ID
func GetGoodsReview(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.GoodsReview
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Goods review not found"})
	}

	return c.JSON(item)
}

// CreateGoodsReview creates a new goods review
func CreateGoodsReview(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.GoodsReview
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create goods review"})
	}

	return c.Status(201).JSON(item)
}

// UpdateGoodsReview updates an existing goods review (e.g. setting reply)
func UpdateGoodsReview(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.GoodsReview
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Goods review not found"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")
	delete(updates, "trashed_at")

	// If reply is being set, automatically set replied_at
	if _, ok := updates["reply"]; ok {
		now := time.Now()
		updates["replied_at"] = now
	}

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update goods review"})
	}

	return c.JSON(item)
}

// DeleteGoodsReview soft-deletes a goods review
func DeleteGoodsReview(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.GoodsReview
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Goods review not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete goods review"})
	}

	return c.JSON(fiber.Map{"message": "Goods review deleted successfully"})
}
