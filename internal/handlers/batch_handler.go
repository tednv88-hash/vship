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

// GetBatches lists shipping batches with pagination
func GetBatches(c *fiber.Ctx) error {
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
	database.DB.Model(&models.ShippingBatch{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Count(&total)

	var batches []models.ShippingBatch
	if err := database.DB.
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&batches).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch batches"})
	}

	return c.JSON(fiber.Map{
		"data":  batches,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetBatch gets a single batch by ID
func GetBatch(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var batch models.ShippingBatch
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&batch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Batch not found"})
	}

	return c.JSON(batch)
}

// CreateBatch creates a new batch with auto-generated batch number
func CreateBatch(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var batch models.ShippingBatch
	if err := c.BodyParser(&batch); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	batch.TenantID = &tenantID
	batch.BatchNumber = fmt.Sprintf("BAT%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000))

	if err := database.DB.Create(&batch).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create batch"})
	}

	return c.Status(201).JSON(batch)
}

// UpdateBatch updates an existing batch
func UpdateBatch(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var batch models.ShippingBatch
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&batch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Batch not found"})
	}

	originalBatchNumber := batch.BatchNumber

	if err := c.BodyParser(&batch); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	batch.ID = id
	batch.TenantID = &tenantID
	batch.BatchNumber = originalBatchNumber

	if err := database.DB.Save(&batch).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update batch"})
	}

	return c.JSON(batch)
}

// DeleteBatch soft-deletes a batch
func DeleteBatch(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var batch models.ShippingBatch
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&batch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Batch not found"})
	}

	if err := database.DB.Model(&batch).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete batch"})
	}

	return c.JSON(fiber.Map{"message": "Batch deleted successfully"})
}

// DepartBatch sets batch status to "in_transit"
func DepartBatch(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var batch models.ShippingBatch
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&batch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Batch not found"})
	}

	if batch.Status != "preparing" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Batch must be in 'preparing' status to depart, current status: " + batch.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&batch).Updates(map[string]interface{}{
		"status":      "in_transit",
		"departed_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to depart batch"})
	}

	batch.Status = "in_transit"
	batch.DepartedAt = &now
	return c.JSON(batch)
}

// ArriveBatch sets batch status to "arrived"
func ArriveBatch(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var batch models.ShippingBatch
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&batch).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Batch not found"})
	}

	if batch.Status != "in_transit" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Batch must be in 'in_transit' status to arrive, current status: " + batch.Status,
		})
	}

	now := time.Now()
	if err := database.DB.Model(&batch).Updates(map[string]interface{}{
		"status":     "arrived",
		"arrived_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to arrive batch"})
	}

	batch.Status = "arrived"
	batch.ArrivedAt = &now
	return c.JSON(batch)
}

// GetBatchOrders lists orders linked to a batch
func GetBatchOrders(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	batchID, err := uuid.Parse(c.Params("batchId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var batchOrders []models.BatchOrder
	if err := database.DB.
		Where("tenant_id = ? AND batch_id = ?", tenantID, batchID).
		Order("created_at ASC").
		Find(&batchOrders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch batch orders"})
	}

	return c.JSON(fiber.Map{"data": batchOrders})
}

// AddBatchOrder links an order to a batch
func AddBatchOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	batchID, err := uuid.Parse(c.Params("batchId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var body struct {
		OrderID uuid.UUID `json:"order_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	batchOrder := models.BatchOrder{
		TenantID: &tenantID,
		BatchID:  batchID,
		OrderID:  body.OrderID,
	}

	if err := database.DB.Create(&batchOrder).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add order to batch"})
	}

	return c.Status(201).JSON(batchOrder)
}

// RemoveBatchOrder removes an order link from a batch
func RemoveBatchOrder(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch order ID"})
	}

	var batchOrder models.BatchOrder
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&batchOrder).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Batch order link not found"})
	}

	if err := database.DB.Delete(&batchOrder).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove order from batch"})
	}

	return c.JSON(fiber.Map{"message": "Order removed from batch successfully"})
}

// GetBatchTracking lists tracking logs for a batch
func GetBatchTracking(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	batchID, err := uuid.Parse(c.Params("batchId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var trackingLogs []models.BatchTracking
	if err := database.DB.
		Where("tenant_id = ? AND batch_id = ?", tenantID, batchID).
		Order("created_at DESC").
		Find(&trackingLogs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch batch tracking logs"})
	}

	return c.JSON(fiber.Map{"data": trackingLogs})
}

// AddBatchTracking adds a tracking log entry for a batch
func AddBatchTracking(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	batchID, err := uuid.Parse(c.Params("batchId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid batch ID"})
	}

	var tracking models.BatchTracking
	if err := c.BodyParser(&tracking); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tracking.TenantID = &tenantID
	tracking.BatchID = batchID

	if err := database.DB.Create(&tracking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add tracking log"})
	}

	return c.Status(201).JSON(tracking)
}
