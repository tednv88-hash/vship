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

// GetMembers lists members with pagination, search, and status filter
func GetMembers(c *fiber.Ctx) error {
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
	status := c.Query("status")

	// Count query
	countQuery := database.DB.Table("members").
		Where("members.tenant_id = ? AND members.trashed_at IS NULL", tenantID)
	if search != "" {
		countQuery = countQuery.Where(
			"(members.name ILIKE ? OR members.email ILIKE ? OR members.phone ILIKE ?)",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}
	if status != "" {
		countQuery = countQuery.Where("members.status = ?", status)
	}
	var total int64
	countQuery.Count(&total)

	// Data query with LEFT JOIN on member_levels
	query := database.DB.Table("members").
		Select("members.*, member_levels.name as member_level_name").
		Joins("LEFT JOIN member_levels ON member_levels.id = members.member_level_id").
		Where("members.tenant_id = ? AND members.trashed_at IS NULL", tenantID)
	if search != "" {
		query = query.Where(
			"(members.name ILIKE ? OR members.email ILIKE ? OR members.phone ILIKE ?)",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}
	if status != "" {
		query = query.Where("members.status = ?", status)
	}

	var members []models.Member
	if err := query.
		Order("members.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&members).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch members"})
	}

	return c.JSON(fiber.Map{
		"data":  members,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetMember gets a single member by ID
func GetMember(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	var member models.Member
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&member).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}

	return c.JSON(member)
}

// CreateMember creates a new member
func CreateMember(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var member models.Member
	if err := c.BodyParser(&member); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	member.TenantID = &tenantID

	if err := database.DB.Create(&member).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create member"})
	}

	return c.Status(201).JSON(member)
}

// UpdateMember updates an existing member
func UpdateMember(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	var member models.Member
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&member).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}

	if err := c.BodyParser(&member); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	member.ID = id
	member.TenantID = &tenantID

	if err := database.DB.Save(&member).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update member"})
	}

	return c.JSON(member)
}

// DeleteMember soft-deletes a member
func DeleteMember(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	var member models.Member
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&member).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Member not found"})
	}

	if err := database.DB.Model(&member).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete member"})
	}

	return c.JSON(fiber.Map{"message": "Member deleted successfully"})
}
