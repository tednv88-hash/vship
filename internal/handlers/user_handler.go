package handlers

import (
	"strconv"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers returns a paginated list of users with optional status and search filters
func GetUsers(c *fiber.Ctx) error {
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
	var items []models.User

	query := database.DB.Where("tenant_id = ? AND trashed_at IS NULL", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("(name ILIKE ? OR email ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	query.Model(&models.User{}).Count(&total)
	query.Select("id, tenant_id, email, name, phone, user_role, status, profile_pic, last_login_at, logged_out_at, created_at, updated_at, extra_fields").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&items)

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUser returns a single user by ID (without password_hash)
func GetUser(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.User
	if err := database.DB.Select("id, tenant_id, email, name, phone, user_role, status, profile_pic, last_login_at, logged_out_at, created_at, updated_at, extra_fields").
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(item)
}

// createUserRequest represents the request body for creating a user
type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	UserRole string `json:"user_role"`
	Status   string `json:"status"`
}

// CreateUser creates a new user with hashed password
func CreateUser(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email, password, and name are required"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	item := models.User{
		TenantID:     &tenantID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Phone:        req.Phone,
		UserRole:     req.UserRole,
		Status:       req.Status,
	}

	if item.UserRole == "" {
		item.UserRole = "user"
	}
	if item.Status == "" {
		item.Status = "active"
	}

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Clear password hash before returning
	item.PasswordHash = ""

	return c.Status(201).JSON(item)
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	UserRole *string `json:"user_role"`
	Status   *string `json:"status"`
	Password *string `json:"password"`
}

// UpdateUser updates an existing user
func UpdateUser(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.User
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var req updateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.UserRole != nil {
		updates["user_role"] = *req.UserRole
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		updates["password_hash"] = string(hashedPassword)
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
		}
	}

	// Clear password hash before returning
	item.PasswordHash = ""

	return c.JSON(item)
}

// DeleteUser soft-deletes a user
func DeleteUser(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var item models.User
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := database.DB.Model(&item).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
