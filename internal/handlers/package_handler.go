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

// createStatusLog creates a PackageStatusLog record
func createStatusLog(tenantID, packageID uuid.UUID, fromStatus, toStatus, remark string, operatorID uuid.UUID) {
	log := models.PackageStatusLog{
		TenantID:   &tenantID,
		PackageID:  packageID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Remark:     remark,
		OperatorID: operatorID,
	}
	database.DB.Create(&log)
}

// GetPackages lists packages with pagination and filters
func GetPackages(c *fiber.Ctx) error {
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

	query := database.DB.Model(&models.Package{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if trackingNumber := c.Query("tracking_number"); trackingNumber != "" {
		query = query.Where("tracking_number = ?", trackingNumber)
	}
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if shippingMarkID := c.Query("shipping_mark_id"); shippingMarkID != "" {
		query = query.Where("shipping_mark_id = ?", shippingMarkID)
	}
	if c.Query("unclaimed") == "true" {
		query = query.Where("user_id IS NULL AND status = 'received'")
	}
	if c.Query("is_problem") == "true" {
		query = query.Where("is_problem = true")
	}

	var total int64
	query.Count(&total)

	var packages []models.Package
	if err := query.
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&packages).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch packages"})
	}

	return c.JSON(fiber.Map{
		"data":  packages,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetPackage gets a single package by ID
func GetPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	return c.JSON(pkg)
}

// CreatePackage creates a new package
func CreatePackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var pkg models.Package
	if err := c.BodyParser(&pkg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	pkg.TenantID = &tenantID

	if err := database.DB.Create(&pkg).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create package"})
	}

	return c.Status(201).JSON(pkg)
}

// UpdatePackage updates an existing package
func UpdatePackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	if err := c.BodyParser(&pkg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	pkg.ID = id
	pkg.TenantID = &tenantID

	if err := database.DB.Save(&pkg).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update package"})
	}

	return c.JSON(pkg)
}

// DeletePackage soft-deletes a package
func DeletePackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	if err := database.DB.Model(&pkg).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete package"})
	}

	return c.JSON(fiber.Map{"message": "Package deleted successfully"})
}

// ForecastPackage creates a new package with status "forecast"
func ForecastPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var pkg models.Package
	if err := c.BodyParser(&pkg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	pkg.TenantID = &tenantID
	pkg.Status = "forecast"

	if pkg.Source == "" {
		pkg.Source = "user"
	}

	if err := database.DB.Create(&pkg).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create forecast package"})
	}

	return c.Status(201).JSON(pkg)
}

// ReceivePackage updates a package status to "received"
func ReceivePackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	if pkg.Status != "forecast" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Package must be in 'forecast' status to be received, current status: " + pkg.Status,
		})
	}

	now := time.Now()
	fromStatus := pkg.Status

	if err := database.DB.Model(&pkg).Updates(map[string]interface{}{
		"status":      "received",
		"received_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to receive package"})
	}

	createStatusLog(tenantID, id, fromStatus, "received", "Package received", operatorID)

	pkg.Status = "received"
	pkg.ReceivedAt = &now
	return c.JSON(pkg)
}

// ShelvePackage updates a package status to "shelved"
func ShelvePackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var body struct {
		ShelfID uuid.UUID `json:"shelf_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	now := time.Now()
	fromStatus := pkg.Status

	if err := database.DB.Model(&pkg).Updates(map[string]interface{}{
		"status":     "shelved",
		"shelf_id":   body.ShelfID,
		"shelved_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to shelve package"})
	}

	createStatusLog(tenantID, id, fromStatus, "shelved", "Package shelved", operatorID)

	pkg.Status = "shelved"
	pkg.ShelfID = &body.ShelfID
	pkg.ShelvedAt = &now
	return c.JSON(pkg)
}

// InspectPackage updates a package status to "inspected"
func InspectPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var body struct {
		Weight float64 `json:"weight"`
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	now := time.Now()
	fromStatus := pkg.Status

	updates := map[string]interface{}{
		"status":       "inspected",
		"inspected_at": now,
	}
	if body.Weight > 0 {
		updates["weight"] = body.Weight
	}
	if body.Length > 0 {
		updates["length"] = body.Length
	}
	if body.Width > 0 {
		updates["width"] = body.Width
	}
	if body.Height > 0 {
		updates["height"] = body.Height
	}

	if err := database.DB.Model(&pkg).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to inspect package"})
	}

	createStatusLog(tenantID, id, fromStatus, "inspected", "Package inspected", operatorID)

	// Reload to get updated fields
	database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&pkg)
	return c.JSON(pkg)
}

// ShipOutPackage updates a package status to "shipped_out" (掃描出庫)
func ShipOutPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid package ID"})
	}

	var pkg models.Package
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&pkg).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Package not found"})
	}

	now := time.Now()
	fromStatus := pkg.Status

	if err := database.DB.Model(&pkg).Updates(map[string]interface{}{
		"status":     "shipped_out",
		"shipped_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to ship out package"})
	}

	createStatusLog(tenantID, id, fromStatus, "shipped_out", "Package shipped out", operatorID)

	database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&pkg)
	return c.JSON(pkg)
}
