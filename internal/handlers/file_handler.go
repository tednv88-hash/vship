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

// GetFiles returns a paginated list of files (non-trashed)
func GetFiles(c *fiber.Ctx) error {
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
	var items []models.File

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if search := c.Query("search"); search != "" {
		query = query.Where("file_name ILIKE ?", "%"+search+"%")
	}
	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if fileType := c.Query("file_type"); fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	query.Model(&models.File{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetFile returns a single file by ID
func GetFile(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.File
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	return c.JSON(item)
}

// CreateFile creates a new file metadata record (no actual upload)
func CreateFile(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.File
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create file"})
	}

	return c.Status(201).JSON(item)
}

// UpdateFile updates an existing file metadata
func UpdateFile(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.File
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update file"})
	}

	return c.JSON(item)
}

// DeleteFile soft-deletes a file (moves to recycle bin)
func DeleteFile(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	if c.Query("trashed") == "true" {
		if err := database.DB.Where("tenant_id = ? AND trashed_at IS NOT NULL", tenantID).Delete(&models.File{}).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to empty recycle bin"})
		}

		return c.JSON(fiber.Map{"message": "Recycle bin emptied successfully"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.File
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete file"})
	}

	return c.JSON(fiber.Map{"message": "File deleted successfully"})
}

// GetDeletedFiles returns a paginated list of soft-deleted files (recycle bin)
func GetDeletedFiles(c *fiber.Ctx) error {
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
	var items []models.File

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NOT NULL", tenantID)

	query.Model(&models.File{}).Count(&total)
	query.Order("trashed_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// RestoreFile restores a soft-deleted file from the recycle bin
func RestoreFile(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.File
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NOT NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Deleted file not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", nil).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to restore file"})
	}

	return c.JSON(fiber.Map{"message": "File restored successfully"})
}
