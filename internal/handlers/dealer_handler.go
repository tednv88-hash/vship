package handlers

import (
	"encoding/json"
	"strconv"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const dealerPosterSettingKey = "dealer_posters"

// GetDealerApplications lists dealer applications with pagination and filters
func GetDealerApplications(c *fiber.Ctx) error {
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

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("(real_name ILIKE ? OR phone ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Model(&models.DealerApplication{}).Count(&total)

	var items []models.DealerApplication
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealer applications"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetDealerApplication gets a single dealer application by ID
func GetDealerApplication(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerApplication
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer application not found"})
	}

	return c.JSON(item)
}

// ApproveDealerApplication approves a dealer application
func ApproveDealerApplication(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerApplication
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer application not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status": "approved",
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve dealer application"})
	}

	item.Status = "approved"
	return c.JSON(item)
}

// RejectDealerApplication rejects a dealer application
func RejectDealerApplication(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		AuditRemark string `json:"audit_remark"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.DealerApplication
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer application not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":       "rejected",
		"audit_remark": body.AuditRemark,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject dealer application"})
	}

	item.Status = "rejected"
	item.AuditRemark = body.AuditRemark
	return c.JSON(item)
}

// GetDealers lists dealers with pagination and filters
func GetDealers(c *fiber.Ctx) error {
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

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("user_id::text ILIKE ?", "%"+search+"%")
	}

	var total int64
	query.Model(&models.Dealer{}).Count(&total)

	var items []models.Dealer
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealers"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetDealer gets a single dealer by ID
func GetDealer(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.Dealer
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer not found"})
	}

	return c.JSON(item)
}

// UpdateDealer updates an existing dealer
func UpdateDealer(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.Dealer
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update dealer"})
	}

	return c.JSON(item)
}

// DeleteDealer soft-deletes a dealer
func DeleteDealer(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.Dealer
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete dealer"})
	}

	return c.JSON(fiber.Map{"message": "Dealer deleted successfully"})
}

// GetDealerOrders lists dealer orders with pagination and filters
func GetDealerOrders(c *fiber.Ctx) error {
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

	if dealerID := c.Query("dealer_id"); dealerID != "" {
		query = query.Where("dealer_id = ?", dealerID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.DealerOrder{}).Count(&total)

	var items []models.DealerOrder
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealer orders"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetDealerWithdrawals lists dealer withdrawals with pagination and filters
func GetDealerWithdrawals(c *fiber.Ctx) error {
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

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if dealerID := c.Query("dealer_id"); dealerID != "" {
		query = query.Where("dealer_id = ?", dealerID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.DealerWithdrawal{}).Count(&total)

	var items []models.DealerWithdrawal
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealer withdrawals"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// ApproveDealerWithdrawal approves a dealer withdrawal
func ApproveDealerWithdrawal(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerWithdrawal
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer withdrawal not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status": "approved",
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to approve dealer withdrawal"})
	}

	item.Status = "approved"
	return c.JSON(item)
}

// RejectDealerWithdrawal rejects a dealer withdrawal
func RejectDealerWithdrawal(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		AuditRemark string `json:"audit_remark"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.DealerWithdrawal
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer withdrawal not found"})
	}

	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":       "rejected",
		"audit_remark": body.AuditRemark,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject dealer withdrawal"})
	}

	item.Status = "rejected"
	item.AuditRemark = body.AuditRemark
	return c.JSON(item)
}

// PayDealerWithdrawal marks a dealer withdrawal as paid
func PayDealerWithdrawal(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerWithdrawal
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer withdrawal not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&item).Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to pay dealer withdrawal"})
	}

	item.Status = "paid"
	item.PaidAt = &now
	return c.JSON(item)
}

// GetDealerLevels lists dealer levels with pagination
func GetDealerLevels(c *fiber.Ctx) error {
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

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	var total int64
	query.Model(&models.DealerLevel{}).Count(&total)

	var items []models.DealerLevel
	if err := query.Order("level_no ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealer levels"})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetDealerLevel gets a single dealer level by ID
func GetDealerLevel(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerLevel
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer level not found"})
	}

	return c.JSON(item)
}

// CreateDealerLevel creates a new dealer level
func CreateDealerLevel(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var item models.DealerLevel
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item.TenantID = &tenantID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create dealer level"})
	}

	return c.Status(201).JSON(item)
}

// UpdateDealerLevel updates an existing dealer level
func UpdateDealerLevel(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerLevel
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer level not found"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update dealer level"})
	}

	return c.JSON(item)
}

// DeleteDealerLevel soft-deletes a dealer level
func DeleteDealerLevel(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.DealerLevel
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dealer level not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete dealer level"})
	}

	return c.JSON(fiber.Map{"message": "Dealer level deleted successfully"})
}

func getDealerPosters(tenantID uuid.UUID) ([]map[string]interface{}, error) {
	var setting models.Setting
	err := database.DB.Where("tenant_id = ? AND key = ?", tenantID, dealerPosterSettingKey).First(&setting).Error
	if err != nil {
		return []map[string]interface{}{}, nil
	}

	var posters []map[string]interface{}
	if setting.Value == "" {
		return []map[string]interface{}{}, nil
	}
	if err := json.Unmarshal([]byte(setting.Value), &posters); err != nil {
		return nil, err
	}

	return posters, nil
}

func saveDealerPosters(tenantID uuid.UUID, posters []map[string]interface{}) error {
	data, err := json.Marshal(posters)
	if err != nil {
		return err
	}

	var setting models.Setting
	err = database.DB.Where("tenant_id = ? AND key = ?", tenantID, dealerPosterSettingKey).First(&setting).Error
	if err != nil {
		setting = models.Setting{TenantID: &tenantID, Key: dealerPosterSettingKey, Value: string(data), Group: "dealer"}
		return database.DB.Create(&setting).Error
	}

	return database.DB.Model(&setting).Updates(map[string]interface{}{"value": string(data), "group": "dealer"}).Error
}

// GetDealerPosters lists dealer poster settings stored as tenant JSON config.
func GetDealerPosters(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	posters, err := getDealerPosters(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealer posters"})
	}

	return c.JSON(fiber.Map{"data": posters, "total": len(posters)})
}

// GetDealerPoster returns a single dealer poster by ID.
func GetDealerPoster(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	posters, err := getDealerPosters(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch dealer poster"})
	}

	id := c.Params("id")
	for _, poster := range posters {
		if poster["id"] == id {
			return c.JSON(poster)
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "Dealer poster not found"})
}

// CreateDealerPoster creates a dealer poster config entry.
func CreateDealerPoster(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	posters, err := getDealerPosters(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create dealer poster"})
	}

	var poster map[string]interface{}
	if err := c.BodyParser(&poster); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	poster["id"] = uuid.New().String()
	posters = append(posters, poster)

	if err := saveDealerPosters(tenantID, posters); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create dealer poster"})
	}

	return c.Status(201).JSON(poster)
}

// UpdateDealerPoster updates a dealer poster config entry.
func UpdateDealerPoster(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	posters, err := getDealerPosters(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update dealer poster"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	id := c.Params("id")
	for i, poster := range posters {
		if poster["id"] == id {
			for key, value := range updates {
				if key != "id" {
					poster[key] = value
				}
			}
			posters[i] = poster
			if err := saveDealerPosters(tenantID, posters); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update dealer poster"})
			}
			return c.JSON(poster)
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "Dealer poster not found"})
}

// DeleteDealerPoster deletes a dealer poster config entry.
func DeleteDealerPoster(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	posters, err := getDealerPosters(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete dealer poster"})
	}

	id := c.Params("id")
	for i, poster := range posters {
		if poster["id"] == id {
			posters = append(posters[:i], posters[i+1:]...)
			if err := saveDealerPosters(tenantID, posters); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to delete dealer poster"})
			}
			return c.JSON(fiber.Map{"message": "Dealer poster deleted successfully"})
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "Dealer poster not found"})
}
