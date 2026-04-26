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

// GetWarehouses lists warehouses with pagination
func GetWarehouses(c *fiber.Ctx) error {
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
	database.DB.Model(&models.Warehouse{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Count(&total)

	var warehouses []models.Warehouse
	if err := database.DB.
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Order("name ASC").
		Offset(offset).Limit(limit).
		Find(&warehouses).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch warehouses"})
	}

	return c.JSON(fiber.Map{
		"data":  warehouses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetWarehouse gets a single warehouse by ID
func GetWarehouse(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse ID"})
	}

	var warehouse models.Warehouse
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&warehouse).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse not found"})
	}

	return c.JSON(warehouse)
}

// CreateWarehouse creates a new warehouse
func CreateWarehouse(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var warehouse models.Warehouse
	if err := c.BodyParser(&warehouse); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	warehouse.TenantID = &tenantID

	if err := database.DB.Create(&warehouse).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create warehouse"})
	}

	return c.Status(201).JSON(warehouse)
}

// UpdateWarehouse updates an existing warehouse
func UpdateWarehouse(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse ID"})
	}

	var warehouse models.Warehouse
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&warehouse).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse not found"})
	}

	if err := c.BodyParser(&warehouse); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	warehouse.ID = id
	warehouse.TenantID = &tenantID

	if err := database.DB.Save(&warehouse).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update warehouse"})
	}

	return c.JSON(warehouse)
}

// DeleteWarehouse soft-deletes a warehouse
func DeleteWarehouse(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse ID"})
	}

	var warehouse models.Warehouse
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&warehouse).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Warehouse not found"})
	}

	if err := database.DB.Model(&warehouse).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete warehouse"})
	}

	return c.JSON(fiber.Map{"message": "Warehouse deleted successfully"})
}

// GetShelves lists shelves for a warehouse with pagination
func GetShelves(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	warehouseID, err := uuid.Parse(c.Params("warehouseId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse ID"})
	}

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
	database.DB.Model(&models.WarehouseShelf{}).
		Where("tenant_id = ? AND warehouse_id = ?", tenantID, warehouseID).
		Count(&total)

	var shelves []models.WarehouseShelf
	if err := database.DB.
		Where("tenant_id = ? AND warehouse_id = ?", tenantID, warehouseID).
		Order("shelf_code ASC").
		Offset(offset).Limit(limit).
		Find(&shelves).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch shelves"})
	}

	return c.JSON(fiber.Map{
		"data":  shelves,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateShelf creates a new shelf in a warehouse
func CreateShelf(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	warehouseID, err := uuid.Parse(c.Params("warehouseId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid warehouse ID"})
	}

	var shelf models.WarehouseShelf
	if err := c.BodyParser(&shelf); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	shelf.TenantID = &tenantID
	shelf.WarehouseID = warehouseID

	if err := database.DB.Create(&shelf).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create shelf"})
	}

	return c.Status(201).JSON(shelf)
}

// UpdateShelf updates an existing shelf
func UpdateShelf(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid shelf ID"})
	}

	var shelf models.WarehouseShelf
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&shelf).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shelf not found"})
	}

	if err := c.BodyParser(&shelf); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	shelf.ID = id
	shelf.TenantID = &tenantID

	if err := database.DB.Save(&shelf).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update shelf"})
	}

	return c.JSON(shelf)
}

// DeleteShelf hard-deletes a shelf
func DeleteShelf(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid shelf ID"})
	}

	var shelf models.WarehouseShelf
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&shelf).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shelf not found"})
	}

	if err := database.DB.Delete(&shelf).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete shelf"})
	}

	return c.JSON(fiber.Map{"message": "Shelf deleted successfully"})
}
