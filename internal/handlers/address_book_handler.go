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

// GetAddressBooks lists address books with pagination, search, and member_id filter
func GetAddressBooks(c *fiber.Ctx) error {
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

	search := c.Query("search")
	memberID := c.Query("member_id")

	// Count query
	countQuery := database.DB.Table("address_books").
		Where("address_books.tenant_id = ? AND address_books.trashed_at IS NULL", tenantID)
	if search != "" {
		countQuery = countQuery.Where(
			"(address_books.recipient_name ILIKE ? OR address_books.phone ILIKE ? OR address_books.address ILIKE ?)",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}
	if memberID != "" {
		countQuery = countQuery.Where("address_books.member_id = ?", memberID)
	}
	var total int64
	countQuery.Count(&total)

	// Data query with LEFT JOIN on countries
	query := database.DB.Table("address_books").
		Select("address_books.*, countries.name as country_name").
		Joins("LEFT JOIN countries ON countries.id = address_books.country_id").
		Where("address_books.tenant_id = ? AND address_books.trashed_at IS NULL", tenantID)
	if search != "" {
		query = query.Where(
			"(address_books.recipient_name ILIKE ? OR address_books.phone ILIKE ? OR address_books.address ILIKE ?)",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}
	if memberID != "" {
		query = query.Where("address_books.member_id = ?", memberID)
	}

	var addressBooks []models.AddressBook
	if err := query.
		Order("address_books.is_default DESC, address_books.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&addressBooks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch address books"})
	}

	return c.JSON(fiber.Map{
		"data":  addressBooks,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetAddressBook gets a single address book entry by ID
func GetAddressBook(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid address book ID"})
	}

	var addressBook models.AddressBook
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&addressBook).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Address book entry not found"})
	}

	return c.JSON(addressBook)
}

// CreateAddressBook creates a new address book entry
func CreateAddressBook(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var addressBook models.AddressBook
	if err := c.BodyParser(&addressBook); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	addressBook.TenantID = &tenantID

	// If setting as default, clear other defaults for the same member
	if addressBook.IsDefault {
		database.DB.Model(&models.AddressBook{}).
			Where("tenant_id = ? AND member_id = ? AND trashed_at IS NULL", tenantID, addressBook.MemberID).
			Update("is_default", false)
	}

	if err := database.DB.Create(&addressBook).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create address book entry"})
	}

	return c.Status(201).JSON(addressBook)
}

// UpdateAddressBook updates an existing address book entry
func UpdateAddressBook(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid address book ID"})
	}

	var addressBook models.AddressBook
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&addressBook).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Address book entry not found"})
	}

	if err := c.BodyParser(&addressBook); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	addressBook.ID = id
	addressBook.TenantID = &tenantID

	// If setting as default, clear other defaults for the same member
	if addressBook.IsDefault {
		database.DB.Model(&models.AddressBook{}).
			Where("tenant_id = ? AND member_id = ? AND id != ? AND trashed_at IS NULL", tenantID, addressBook.MemberID, id).
			Update("is_default", false)
	}

	if err := database.DB.Save(&addressBook).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update address book entry"})
	}

	return c.JSON(addressBook)
}

// DeleteAddressBook soft-deletes an address book entry
func DeleteAddressBook(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid address book ID"})
	}

	var addressBook models.AddressBook
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&addressBook).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Address book entry not found"})
	}

	if err := database.DB.Model(&addressBook).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete address book entry"})
	}

	return c.JSON(fiber.Map{"message": "Address book entry deleted successfully"})
}
