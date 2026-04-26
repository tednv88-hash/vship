package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"vship/config"
	"vship/internal/database"
	"vship/internal/email"
	"vship/internal/middleware"
	"vship/internal/models"
	"vship/internal/services"
	"vship/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login handles user login
func Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ? AND trashed_at IS NULL", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Check status
	if strings.TrimSpace(strings.ToLower(user.Status)) != "active" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Account is inactive",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Get tenant ID
	tenantID := uuid.Nil
	if user.TenantID != nil {
		tenantID = *user.TenantID
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, tenantID, user.Email, user.UserRole)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Update last login
	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", now)

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"user_role": user.UserRole,
			"tenant_id": tenantID,
		},
	})
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email, password, and name are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	// Check if email already exists
	var count int64
	database.DB.Model(&models.User{}).Where("email = ? AND trashed_at IS NULL", req.Email).Count(&count)
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create tenant and user in a transaction
	var user models.User
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Create a new tenant for this user
		subdomain := strings.Split(req.Email, "@")[0]
		// Ensure subdomain uniqueness by appending random suffix if needed
		var subdomainCount int64
		tx.Model(&models.Tenant{}).Where("subdomain = ?", subdomain).Count(&subdomainCount)
		if subdomainCount > 0 {
			subdomain = fmt.Sprintf("%s-%s", subdomain, uuid.New().String()[:8])
		}

		tenant := models.Tenant{
			Name:      req.Name + "的商城",
			Subdomain: subdomain,
			Plan:      "free",
			Status:    "active",
		}
		if err := tx.Create(&tenant).Error; err != nil {
			return fmt.Errorf("failed to create tenant: %w", err)
		}

		// Create user with tenant association
		user = models.User{
			TenantID:     &tenant.ID,
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			Name:         req.Name,
			Phone:        req.Phone,
			UserRole:     "admin",
			Status:       "active",
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Create default page designs for the new tenant
		if err := services.CreateDefaultPageDesigns(tx, tenant.ID); err != nil {
			log.Printf("Warning: failed to create default pages for tenant %s: %v", tenant.ID, err)
			// Non-fatal: don't rollback the entire registration
		}

		return nil
	})

	if err != nil {
		log.Printf("Error during registration: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Registration successful",
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"tenant_id": user.TenantID,
		},
	})
}

// Logout handles user logout
func Logout(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID != uuid.Nil {
		now := time.Now()
		database.DB.Model(&models.User{}).Where("id = ?", userID).Update("logged_out_at", now)
	}

	c.ClearCookie("auth_token")

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// GetCurrentUser returns current user info
func GetCurrentUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	var user models.User
	if err := database.DB.Where("id = ? AND trashed_at IS NULL", userID).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":        user.ID,
		"email":     user.Email,
		"name":      user.Name,
		"phone":     user.Phone,
		"user_role": user.UserRole,
		"status":    user.Status,
		"tenant_id": user.TenantID,
	})
}

// ===== Forgot Password =====

// ForgotPassword sends a password reset email
func ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email is required"})
	}

	// Always return success to prevent email enumeration
	safeMsg := fiber.Map{"message": "如果該電子郵件已註冊，我們已發送密碼重設連結"}

	// Find user
	var user models.User
	if err := database.DB.Where("LOWER(email) = LOWER(?) AND trashed_at IS NULL", req.Email).First(&user).Error; err != nil {
		return c.JSON(safeMsg)
	}

	// Generate opaque token
	token, tokenHash, err := newOpaqueToken()
	if err != nil {
		log.Printf("Error generating reset token: %v", err)
		return c.JSON(safeMsg)
	}

	// Save token to DB
	resetToken := models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&resetToken).Error; err != nil {
		log.Printf("Error saving reset token: %v", err)
		return c.JSON(safeMsg)
	}

	// Build reset URL
	cfg := config.Load()
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", cfg.BaseURL, url.QueryEscape(token))

	// Send email
	if err := email.SendPasswordResetEmail(cfg, user.Email, user.Name, resetURL); err != nil {
		log.Printf("Error sending password reset email to %s: %v", user.Email, err)
		// Still return success to prevent enumeration
	}

	return c.JSON(safeMsg)
}

// ResetPassword resets the user password using a token
func ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Token and new password are required"})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
	}

	// URL-decode the token
	tokenStr, _ := url.QueryUnescape(req.Token)
	if tokenStr == "" {
		tokenStr = req.Token
	}

	// Hash the token for lookup
	hash := sha256.Sum256([]byte(tokenStr))

	// Find valid token
	var t models.PasswordResetToken
	if err := database.DB.Where("token_hash = ? AND used_at IS NULL AND expires_at > NOW()", hash[:]).First(&t).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "無效或已過期的重設連結"})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Update user password
	database.DB.Model(&models.User{}).Where("id = ?", t.UserID).Update("password_hash", string(hashedPassword))

	// Mark token as used
	now := time.Now()
	database.DB.Model(&models.PasswordResetToken{}).Where("id = ?", t.ID).Update("used_at", now)

	return c.JSON(fiber.Map{"message": "密碼重設成功"})
}

// newOpaqueToken generates a random token and its SHA-256 hash
func newOpaqueToken() (token string, tokenHash []byte, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", nil, err
	}
	token = base64.RawURLEncoding.EncodeToString(buf)
	sum := sha256.Sum256([]byte(token))
	return token, sum[:], nil
}

// ===== Google OAuth =====

// GoogleTokenInfo represents Google's token info response
type GoogleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Aud           string `json:"aud"`
	Exp           int64  `json:"exp"`
}

// GetGoogleOAuthConfig returns Google OAuth client configuration
func GetGoogleOAuthConfig(c *fiber.Ctx) error {
	cfg := config.Load()
	return c.JSON(fiber.Map{
		"client_id": cfg.GoogleOAuth.ClientID,
		"enabled":   cfg.GoogleOAuth.Enabled,
	})
}

// GoogleLogin handles Google OAuth login/register
func GoogleLogin(c *fiber.Ctx) error {
	cfg := config.Load()

	if !cfg.GoogleOAuth.Enabled {
		return c.Status(503).JSON(fiber.Map{"error": "Google login is not enabled"})
	}
	if cfg.GoogleOAuth.ClientID == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Google OAuth not configured"})
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Token is required"})
	}

	// Verify Google ID token
	tokenInfo, err := verifyGoogleToken(req.Token, cfg.GoogleOAuth.ClientID)
	if err != nil {
		log.Printf("Google token verification failed: %v", err)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Google token"})
	}

	if !tokenInfo.EmailVerified {
		return c.Status(400).JSON(fiber.Map{"error": "Email not verified by Google"})
	}
	if tokenInfo.Aud != cfg.GoogleOAuth.ClientID {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token audience"})
	}

	// Find or create user
	var user models.User
	now := time.Now()

	if err := database.DB.Where("LOWER(email) = LOWER(?) AND trashed_at IS NULL", tokenInfo.Email).First(&user).Error; err != nil {
		// Create new user with tenant in a transaction
		userName := tokenInfo.Name
		if len(userName) > 50 {
			userName = userName[:50]
		}
		if userName == "" {
			userName = strings.Split(tokenInfo.Email, "@")[0]
		}

		txErr := database.DB.Transaction(func(tx *gorm.DB) error {
			// Create tenant
			subdomain := strings.Split(strings.ToLower(tokenInfo.Email), "@")[0]
			var subdomainCount int64
			tx.Model(&models.Tenant{}).Where("subdomain = ?", subdomain).Count(&subdomainCount)
			if subdomainCount > 0 {
				subdomain = fmt.Sprintf("%s-%s", subdomain, uuid.New().String()[:8])
			}

			tenant := models.Tenant{
				Name:      userName + "的商城",
				Subdomain: subdomain,
				Plan:      "free",
				Status:    "active",
			}
			if err := tx.Create(&tenant).Error; err != nil {
				return fmt.Errorf("failed to create tenant: %w", err)
			}

			user = models.User{
				TenantID:     &tenant.ID,
				Email:        strings.ToLower(tokenInfo.Email),
				PasswordHash: "", // Google OAuth users have no password
				Name:         userName,
				UserRole:     "admin",
				Status:       "active",
				ProfilePic:   tokenInfo.Picture,
				LastLoginAt:  &now,
			}
			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			// Create default page designs
			if err := services.CreateDefaultPageDesigns(tx, tenant.ID); err != nil {
				log.Printf("Warning: failed to create default pages for Google OAuth tenant %s: %v", tenant.ID, err)
			}

			return nil
		})

		if txErr != nil {
			log.Printf("Error creating Google OAuth user: %v", txErr)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
		}
	} else {
		// Existing user - update last login
		updates := map[string]interface{}{
			"last_login_at": now,
		}
		if tokenInfo.Picture != "" && user.ProfilePic == "" {
			updates["profile_pic"] = tokenInfo.Picture
		}
		database.DB.Model(&user).Updates(updates)

		if strings.TrimSpace(strings.ToLower(user.Status)) != "active" {
			return c.Status(401).JSON(fiber.Map{"error": "Account is inactive"})
		}
	}

	// Get tenant ID
	tenantID := uuid.Nil
	requiresSetup := false
	if user.TenantID != nil {
		tenantID = *user.TenantID
	} else {
		requiresSetup = true
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, tenantID, user.Email, user.UserRole)
	if err != nil {
		log.Printf("Error generating token for Google OAuth user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"token":          token,
		"requires_setup": requiresSetup,
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"user_role": user.UserRole,
			"tenant_id": tenantID,
		},
	})
}

// verifyGoogleToken verifies a Google ID token
func verifyGoogleToken(idToken, clientID string) (*GoogleTokenInfo, error) {
	tokenInfoURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", url.QueryEscape(idToken))
	resp, err := http.Get(tokenInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token verification failed: %s", string(body))
	}

	var tokenInfo GoogleTokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse token info: %w", err)
	}

	if tokenInfo.Email == "" {
		return nil, fmt.Errorf("no email in token info")
	}

	// Check expiry
	if tokenInfo.Exp > 0 && time.Now().Unix() > tokenInfo.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &tokenInfo, nil
}
