package middleware

import (
	"strings"
	"time"
	"vship/internal/database"
	"vship/internal/models"
	"vship/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware JWT authentication middleware
func AuthMiddleware(c *fiber.Ctx) error {
	path := c.Path()

	// Public paths that don't require authentication
	publicPaths := []string{
		"/", "/login",
		"/api/v1/auth/login", "/api/v1/auth/register", "/api/v1/auth/logout",
	}
	for _, publicPath := range publicPaths {
		if path == publicPath || strings.HasPrefix(path, publicPath+"/") {
			return c.Next()
		}
	}

	unauthorized := func(jsonMsg string) error {
		accept := c.Get("Accept")
		if strings.HasPrefix(path, "/api") {
			return c.Status(401).JSON(fiber.Map{
				"error": jsonMsg,
			})
		}
		if strings.Contains(accept, "text/html") {
			c.ClearCookie("auth_token")
			if path != "/login" {
				return c.Redirect("/login")
			}
			return c.Next()
		}
		return c.Status(401).JSON(fiber.Map{
			"error": jsonMsg,
		})
	}

	var tokenString string

	// Try Authorization header first (API requests)
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			tokenString = authHeader
		}
	} else {
		// Fall back to cookie (web requests)
		tokenString = c.Cookies("auth_token")
		if tokenString == "" {
			tokenString = c.Cookies("token")
		}
	}

	if tokenString == "" {
		return unauthorized("Authorization header or cookie is required")
	}

	// Validate JWT token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return unauthorized("Invalid or expired token")
	}

	if claims.UserID == uuid.Nil || claims.Email == "" {
		return unauthorized("Incomplete session data")
	}

	// Store user info in locals
	c.Locals("user_id", claims.UserID)
	c.Locals("user_email", claims.Email)
	c.Locals("user_role", claims.Role)

	if claims.TenantID != uuid.Nil {
		c.Locals("tenant_id", claims.TenantID)
	}

	// Verify user still exists and is active
	var user models.User
	if claims.TenantID == uuid.Nil {
		if err := database.DB.Select("id", "tenant_id", "status", "logged_out_at").
			Where("id = ?", claims.UserID).
			First(&user).Error; err != nil {
			return unauthorized("User not found")
		}
		if user.TenantID != nil && *user.TenantID != uuid.Nil {
			c.Locals("tenant_id", *user.TenantID)
		}
	} else {
		if err := database.DB.Select("id", "tenant_id", "status", "logged_out_at").
			Where("id = ? AND tenant_id = ?", claims.UserID, claims.TenantID).
			First(&user).Error; err != nil {
			return unauthorized("User not found")
		}
	}

	if strings.TrimSpace(strings.ToLower(user.Status)) != "active" {
		return unauthorized("Account is inactive")
	}

	// Check if user has logged out since token was issued
	if claims.IssuedAt != nil && user.LoggedOutAt != nil {
		if claims.IssuedAt.Time.Before(*user.LoggedOutAt) {
			c.ClearCookie("auth_token")
			return unauthorized("Session expired (logged out)")
		}
	}

	// Verify tenant is active
	if claims.TenantID != uuid.Nil {
		var tenant models.Tenant
		if err := database.DB.Select("id", "status").
			Where("id = ?", claims.TenantID).
			First(&tenant).Error; err != nil {
			return unauthorized("Tenant not found")
		}
		tenantStatus := strings.TrimSpace(strings.ToLower(tenant.Status))
		if tenantStatus != "" && tenantStatus != "active" {
			return unauthorized("Tenant is inactive")
		}
	}

	return c.Next()
}

// TenantMiddleware extracts tenant info from request
func TenantMiddleware(c *fiber.Ctx) error {
	tenantID := GetTenantID(c)
	if tenantID != uuid.Nil {
		var tenant models.Tenant
		if err := database.DB.Where("id = ? AND status = ?", tenantID, "active").First(&tenant).Error; err == nil {
			c.Locals("tenant", tenant)
		}
		return c.Next()
	}

	// Try to get tenant from user if already authenticated
	if userID, ok := c.Locals("user_id").(uuid.UUID); ok && userID != uuid.Nil {
		var user models.User
		if err := database.DB.Select("tenant_id").Where("id = ?", userID).First(&user).Error; err == nil {
			if user.TenantID != nil && *user.TenantID != uuid.Nil {
				var tenant models.Tenant
				if err := database.DB.Where("id = ? AND status = ?", *user.TenantID, "active").First(&tenant).Error; err == nil {
					c.Locals("tenant_id", tenant.ID)
					c.Locals("tenant", tenant)
					return c.Next()
				}
			}
		}
	}

	// For API requests without tenant, return error
	path := c.Path()
	accept := c.Get("Accept")
	if strings.Contains(accept, "application/json") || strings.HasPrefix(path, "/api") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Tenant information is required",
		})
	}

	return c.Next()
}

// RequireRole checks user role
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("user_role").(string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(403).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *fiber.Ctx) uuid.UUID {
	if userID, ok := c.Locals("user_id").(uuid.UUID); ok {
		return userID
	}
	return uuid.Nil
}

// GetTenantID gets tenant ID from context
func GetTenantID(c *fiber.Ctx) uuid.UUID {
	if tenantID, ok := c.Locals("tenant_id").(uuid.UUID); ok {
		return tenantID
	}
	return uuid.Nil
}

// _ is used to prevent "imported and not used" for time package
var _ = time.Now

// MpPublicTenantMiddleware resolves tenant for public miniprogram routes.
// It checks X-Tenant-ID header first, then falls back to the default tenant.
// Unlike TenantMiddleware, it does NOT require authentication.
func MpPublicTenantMiddleware(c *fiber.Ctx) error {
	// Try X-Tenant-ID header
	tenantHeader := c.Get("X-Tenant-ID")
	if tenantHeader != "" {
		if tid, err := uuid.Parse(tenantHeader); err == nil && tid != uuid.Nil {
			c.Locals("tenant_id", tid)
			var tenant models.Tenant
			if err := database.DB.Where("id = ? AND status = ?", tid, "active").First(&tenant).Error; err == nil {
				c.Locals("tenant", tenant)
			}
			return c.Next()
		}
	}

	// Fall back: use the first active tenant (single-tenant deployment)
	var tenant models.Tenant
	if err := database.DB.Where("status = ?", "active").Order("created_at ASC").First(&tenant).Error; err == nil {
		c.Locals("tenant_id", tenant.ID)
		c.Locals("tenant", tenant)
	}

	return c.Next()
}
