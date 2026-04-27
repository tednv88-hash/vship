package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"
	"vship/internal/services"
	"vship/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================
// Helper: consistent MP response format
// ============================================================

func mpOK(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"code": 200, "message": "success", "data": data})
}

func mpCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(201).JSON(fiber.Map{"code": 200, "message": "success", "data": data})
}

func mpErr(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"code": status, "message": msg, "data": nil})
}

func mpPaginate(c *fiber.Ctx) (page, limit, offset int) {
	page, _ = strconv.Atoi(c.Query("page", "1"))
	limit, _ = strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset = (page - 1) * limit
	return
}

// ============================================================
// Auto-migrate new MP tables
// ============================================================

func init() {
	// Will be called after database.DB is initialized.
	// We register a callback to auto-migrate tables used exclusively by MP.
	go func() {
		// Wait for DB to be ready (simple retry)
		for i := 0; i < 30; i++ {
			if database.DB != nil {
				_ = database.DB.AutoMigrate(
					&models.CartItem{},
					&models.Favorite{},
					&models.BrowsingHistory{},
					&models.Feedback{},
					&models.SigninLog{},
					&models.ContentPage{},
				)
				return
			}
			time.Sleep(time.Second)
		}
	}()
}

// ============================================================
// 1. Auth — Public (no auth required)
// ============================================================

// MpLogin handles miniprogram login (email/phone + password)
func MpLogin(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" && req.Phone == "" {
		return mpErr(c, 400, "Email or phone is required")
	}
	if req.Password == "" {
		return mpErr(c, 400, "Password is required")
	}

	var user models.User
	if req.Email != "" {
		if err := database.DB.Where("email = ? AND trashed_at IS NULL", req.Email).First(&user).Error; err != nil {
			return mpErr(c, 401, "Invalid credentials")
		}
	} else {
		if err := database.DB.Where("phone = ? AND trashed_at IS NULL", req.Phone).First(&user).Error; err != nil {
			return mpErr(c, 401, "Invalid credentials")
		}
	}

	if strings.TrimSpace(strings.ToLower(user.Status)) != "active" {
		return mpErr(c, 401, "Account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return mpErr(c, 401, "Invalid credentials")
	}

	tenantID := uuid.Nil
	if user.TenantID != nil {
		tenantID = *user.TenantID
	}

	token, err := utils.GenerateToken(user.ID, tenantID, user.Email, user.UserRole)
	if err != nil {
		return mpErr(c, 500, "Failed to generate token")
	}

	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", now)

	return mpOK(c, fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"name":        user.Name,
			"phone":       user.Phone,
			"user_role":   user.UserRole,
			"tenant_id":   tenantID,
			"profile_pic": user.ProfilePic,
		},
	})
}

// MpRegister handles miniprogram user registration
func MpRegister(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return mpErr(c, 400, "Email, password, and name are required")
	}
	if len(req.Password) < 6 {
		return mpErr(c, 400, "Password must be at least 6 characters")
	}

	var count int64
	database.DB.Model(&models.User{}).Where("email = ? AND trashed_at IS NULL", req.Email).Count(&count)
	if count > 0 {
		return mpErr(c, 409, "Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return mpErr(c, 500, "Failed to hash password")
	}

	var user models.User
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		subdomain := strings.Split(req.Email, "@")[0]
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

		user = models.User{
			TenantID:     &tenant.ID,
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			Name:         req.Name,
			Phone:        req.Phone,
			UserRole:     "user",
			Status:       "active",
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		_ = services.CreateDefaultPageDesigns(tx, tenant.ID)
		return nil
	})

	if err != nil {
		return mpErr(c, 500, "Failed to create account")
	}

	tenantID := uuid.Nil
	if user.TenantID != nil {
		tenantID = *user.TenantID
	}

	token, _ := utils.GenerateToken(user.ID, tenantID, user.Email, user.UserRole)

	return mpCreated(c, fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"phone":     user.Phone,
			"tenant_id": tenantID,
		},
	})
}

// MpWechatLogin is a placeholder for WeChat mini-program login
func MpWechatLogin(c *fiber.Ctx) error {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.Code == "" {
		return mpErr(c, 400, "Code is required")
	}
	return mpErr(c, 501, "WeChat login not configured")
}

// MpBindPhone is a placeholder for binding phone number
func MpBindPhone(c *fiber.Ctx) error {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.Phone == "" || req.Code == "" {
		return mpErr(c, 400, "Phone and code are required")
	}
	return mpErr(c, 501, "Phone binding not configured")
}

// MpSendCode is a placeholder for sending verification code
func MpSendCode(c *fiber.Ctx) error {
	var req struct {
		Phone string `json:"phone"`
		Type  string `json:"type"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.Phone == "" {
		return mpErr(c, 400, "Phone is required")
	}
	return mpErr(c, 501, "SMS service not configured")
}

// ============================================================
// 2. User — Auth required
// ============================================================

// MpGetUserInfo returns current user info with balance, points, coupon count
func MpGetUserInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	if userID == uuid.Nil {
		return mpErr(c, 401, "Not authenticated")
	}

	var user models.User
	if err := database.DB.Where("id = ? AND trashed_at IS NULL", userID).First(&user).Error; err != nil {
		return mpErr(c, 404, "User not found")
	}

	// Get member balance
	var member models.Member
	var balance float64
	if err := database.DB.Where("email = ? AND tenant_id = ? AND trashed_at IS NULL", user.Email, tenantID).First(&member).Error; err == nil {
		balance = member.Balance
	}

	// Get points balance from latest points_log
	var latestPointsLog models.PointsLog
	var points int
	if err := database.DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Order("created_at DESC").First(&latestPointsLog).Error; err == nil {
		points = latestPointsLog.Balance
	}

	// Count unused coupons
	var couponCount int64
	database.DB.Model(&models.CouponReceiveLog{}).
		Where("user_id = ? AND tenant_id = ? AND status = ?", userID, tenantID, "unused").
		Count(&couponCount)

	return mpOK(c, fiber.Map{
		"id":           user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"phone":        user.Phone,
		"profile_pic":  user.ProfilePic,
		"user_role":    user.UserRole,
		"status":       user.Status,
		"tenant_id":    user.TenantID,
		"extra_fields": user.ExtraFields,
		"balance":      balance,
		"points":       points,
		"coupon_count": couponCount,
		"created_at":   user.CreatedAt,
	})
}

// MpUpdateUserInfo updates the current user's profile
func MpUpdateUserInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		return mpErr(c, 401, "Not authenticated")
	}

	var req struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		ProfilePic string `json:"profile_pic"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.Phone != "" {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}
	if req.ProfilePic != "" {
		updates["profile_pic"] = strings.TrimSpace(req.ProfilePic)
	}

	if len(updates) == 0 {
		return mpErr(c, 400, "No fields to update")
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return mpErr(c, 500, "Failed to update user info")
	}

	var user models.User
	database.DB.Where("id = ?", userID).First(&user)

	return mpOK(c, fiber.Map{
		"id":          user.ID,
		"email":       user.Email,
		"name":        user.Name,
		"phone":       user.Phone,
		"profile_pic": user.ProfilePic,
	})
}

// MpGetBalance returns user balance logs with pagination
func MpGetBalance(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.UserBalanceLog{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID)
	query.Count(&total)

	var logs []models.UserBalanceLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch balance logs")
	}

	// Get current balance
	var currentBalance float64
	if len(logs) > 0 {
		currentBalance = logs[0].Balance
	}

	return mpOK(c, fiber.Map{
		"balance": currentBalance,
		"logs":    logs,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// MpGetPoints returns user points logs with pagination
func MpGetPoints(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.PointsLog{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID)
	query.Count(&total)

	var logs []models.PointsLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch points logs")
	}

	var currentPoints int
	if len(logs) > 0 {
		currentPoints = logs[0].Balance
	}

	return mpOK(c, fiber.Map{
		"points": currentPoints,
		"logs":   logs,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// MpRecharge creates a recharge order for the current user
func MpRecharge(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		Amount    float64 `json:"amount"`
		PayMethod string  `json:"pay_method"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.Amount <= 0 {
		return mpErr(c, 400, "Amount must be greater than 0")
	}
	if req.PayMethod == "" {
		req.PayMethod = "balance"
	}

	order := models.RechargeOrder{
		TenantID:  &tenantID,
		UserID:    userID,
		OrderNo:   fmt.Sprintf("RC%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000)),
		Amount:    req.Amount,
		PayMethod: req.PayMethod,
		Status:    "pending",
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return mpErr(c, 500, "Failed to create recharge order")
	}

	return mpCreated(c, order)
}

// MpVerifyIdentity saves identity info to user extra_fields
func MpVerifyIdentity(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req struct {
		RealName   string `json:"real_name"`
		IDType     string `json:"id_type"`
		IDNumber   string `json:"id_number"`
		IDFrontImg string `json:"id_front_img"`
		IDBackImg  string `json:"id_back_img"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.RealName == "" || req.IDNumber == "" {
		return mpErr(c, 400, "Real name and ID number are required")
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return mpErr(c, 404, "User not found")
	}

	extra := user.ExtraFields
	if extra == nil {
		extra = make(models.JSONB)
	}
	extra["identity"] = map[string]interface{}{
		"real_name":    req.RealName,
		"id_type":      req.IDType,
		"id_number":    req.IDNumber,
		"id_front_img": req.IDFrontImg,
		"id_back_img":  req.IDBackImg,
		"status":       "pending",
		"submitted_at": time.Now().Format(time.RFC3339),
	}

	if err := database.DB.Model(&user).Update("extra_fields", extra).Error; err != nil {
		return mpErr(c, 500, "Failed to save identity info")
	}

	return mpOK(c, fiber.Map{"status": "pending"})
}

// MpGetIdentityStatus returns identity verification status
func MpGetIdentityStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return mpErr(c, 404, "User not found")
	}

	identity, ok := user.ExtraFields["identity"]
	if !ok {
		return mpOK(c, fiber.Map{"status": "none"})
	}

	return mpOK(c, identity)
}

// MpSignIn handles daily check-in
func MpSignIn(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	today := time.Now().Format("2006-01-02")

	// Check if already signed in today
	var existing models.SigninLog
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND signin_date = ?", userID, tenantID, today).
		First(&existing).Error; err == nil {
		return mpErr(c, 400, "Already signed in today")
	}

	// Calculate consecutive days
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var yesterdayLog models.SigninLog
	consecutiveDays := 1
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND signin_date = ?", userID, tenantID, yesterday).
		First(&yesterdayLog).Error; err == nil {
		consecutiveDays = yesterdayLog.ConsecutiveDays + 1
	}

	// Calculate bonus points (base 1, +1 every 7 consecutive days, max 10)
	pointsEarned := 1
	bonus := consecutiveDays / 7
	if bonus > 9 {
		bonus = 9
	}
	pointsEarned += bonus

	signinLog := models.SigninLog{
		TenantID:        &tenantID,
		UserID:          userID,
		SigninDate:      today,
		PointsEarned:    pointsEarned,
		ConsecutiveDays: consecutiveDays,
	}

	if err := database.DB.Create(&signinLog).Error; err != nil {
		return mpErr(c, 500, "Failed to sign in")
	}

	// Award points
	var lastPointsLog models.PointsLog
	currentPoints := 0
	if err := database.DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Order("created_at DESC").First(&lastPointsLog).Error; err == nil {
		currentPoints = lastPointsLog.Balance
	}

	pointsLog := models.PointsLog{
		TenantID:    &tenantID,
		UserID:      userID,
		Value:       pointsEarned,
		Balance:     currentPoints + pointsEarned,
		Type:        "signin",
		Description: fmt.Sprintf("Daily sign-in (Day %d)", consecutiveDays),
		RelatedID:   &signinLog.ID,
		RelatedType: "signin_log",
	}
	database.DB.Create(&pointsLog)

	return mpOK(c, fiber.Map{
		"signin_date":      today,
		"consecutive_days": consecutiveDays,
		"points_earned":    pointsEarned,
		"total_points":     pointsLog.Balance,
	})
}

// MpGetSignInStatus returns current month's signin records
func MpGetSignInStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	monthEnd := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location()).Format("2006-01-02")

	var logs []models.SigninLog
	database.DB.Where("user_id = ? AND tenant_id = ? AND signin_date >= ? AND signin_date <= ?",
		userID, tenantID, monthStart, monthEnd).
		Order("signin_date ASC").Find(&logs)

	// Today's status
	today := now.Format("2006-01-02")
	signedToday := false
	consecutiveDays := 0
	for _, l := range logs {
		if l.SigninDate == today {
			signedToday = true
			consecutiveDays = l.ConsecutiveDays
		}
	}

	dates := make([]string, 0, len(logs))
	for _, l := range logs {
		dates = append(dates, l.SigninDate)
	}

	return mpOK(c, fiber.Map{
		"signed_today":     signedToday,
		"consecutive_days": consecutiveDays,
		"signed_dates":     dates,
		"total_this_month": len(logs),
	})
}

// ============================================================
// 3. Packages
// ============================================================

// MpGetPackages lists packages for the current user
func MpGetPackages(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	// Inject user_id filter so GetPackages scopes to this user
	c.Request().URI().QueryArgs().Set("user_id", userID.String())
	return GetPackages(c)
}

// MpGetPackage returns a single package
func MpGetPackage(c *fiber.Ctx) error {
	return GetPackage(c)
}

// MpForecastPackage creates a forecast package for the current user
func MpForecastPackage(c *fiber.Ctx) error {
	return ForecastPackage(c)
}

// MpGetForecasts returns forecast packages for the current user
func MpGetForecasts(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	c.Request().URI().QueryArgs().Set("user_id", userID.String())
	c.Request().URI().QueryArgs().Set("status", "forecast")
	return GetPackages(c)
}

// MpMergePackages merges selected packages into a new consolidated package
func MpMergePackages(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		PackageIDs []uuid.UUID `json:"package_ids"`
		Remark     string      `json:"remark"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if len(req.PackageIDs) < 2 {
		return mpErr(c, 400, "At least 2 packages are required for merging")
	}

	// Verify all packages belong to this user and tenant
	var packages []models.Package
	if err := database.DB.Where("id IN ? AND tenant_id = ? AND user_id = ? AND trashed_at IS NULL",
		req.PackageIDs, tenantID, userID).Find(&packages).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch packages")
	}
	if len(packages) != len(req.PackageIDs) {
		return mpErr(c, 400, "Some packages not found or do not belong to you")
	}

	// Calculate totals
	var totalWeight, totalLength, totalWidth, totalHeight float64
	for _, p := range packages {
		totalWeight += p.Weight
		if p.Length > totalLength {
			totalLength = p.Length
		}
		if p.Width > totalWidth {
			totalWidth = p.Width
		}
		totalHeight += p.Height
	}

	// Create merged package
	mergedPkg := models.Package{
		TenantID:        &tenantID,
		TrackingNumber:  fmt.Sprintf("MRG%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000)),
		UserID:          &userID,
		Status:          "merged",
		Source:          "merge",
		Weight:          totalWeight,
		Length:          totalLength,
		Width:           totalWidth,
		Height:          totalHeight,
		ItemDescription: req.Remark,
		Remark:          fmt.Sprintf("Merged from %d packages", len(packages)),
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&mergedPkg).Error; err != nil {
			return err
		}

		// Mark original packages as merged, referencing the merged package
		for _, p := range packages {
			extra := p.ExtraFields
			if extra == nil {
				extra = make(models.JSONB)
			}
			extra["merged_into"] = mergedPkg.ID.String()
			tx.Model(&models.Package{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
				"status":       "merged",
				"extra_fields": extra,
			})
		}
		return nil
	})

	if err != nil {
		return mpErr(c, 500, "Failed to merge packages")
	}

	return mpCreated(c, mergedPkg)
}

// MpSplitPackage is a placeholder for splitting a package
func MpSplitPackage(c *fiber.Ctx) error {
	return mpErr(c, 501, "Package splitting is not available yet")
}

// ============================================================
// 4. Orders
// ============================================================

// MpGetOrders lists orders for the current user
func MpGetOrders(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	c.Request().URI().QueryArgs().Set("user_id", userID.String())
	return GetOrders(c)
}

// MpGetOrder returns a single order
func MpGetOrder(c *fiber.Ctx) error {
	return GetOrder(c)
}

// MpCreateOrder creates an order
func MpCreateOrder(c *fiber.Ctx) error {
	return CreateOrder(c)
}

// MpCancelOrder cancels an order
func MpCancelOrder(c *fiber.Ctx) error {
	return CancelOrder(c)
}

// MpPayOrder pays an order
func MpPayOrder(c *fiber.Ctx) error {
	return PayOrder(c)
}

// ============================================================
// 5. Shop Orders
// ============================================================

// MpGetShopOrders lists shop orders for the current user
func MpGetShopOrders(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	c.Request().URI().QueryArgs().Set("user_id", userID.String())
	return GetShopOrders(c)
}

// MpGetShopOrder returns a single shop order
func MpGetShopOrder(c *fiber.Ctx) error {
	return GetShopOrder(c)
}

// MpCheckoutCart creates a shop order from selected cart items
func MpCheckoutCart(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		CartIDs   []uuid.UUID `json:"cart_ids"`
		AddressID uuid.UUID   `json:"address_id"`
		Remark    string      `json:"remark"`
		PayMethod string      `json:"pay_method"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if len(req.CartIDs) == 0 {
		return mpErr(c, 400, "No cart items selected")
	}
	if req.AddressID == uuid.Nil {
		return mpErr(c, 400, "Address is required")
	}

	// Load address
	var addr models.AddressBook
	if err := database.DB.Where("id = ? AND member_id = ? AND tenant_id = ? AND trashed_at IS NULL",
		req.AddressID, userID, tenantID).First(&addr).Error; err != nil {
		return mpErr(c, 404, "Address not found")
	}

	// Load cart items with goods info
	type CartWithGoods struct {
		models.CartItem
		GoodsName     string  `json:"goods_name"`
		GoodsImageURL string  `json:"goods_image_url"`
		GoodsPrice    float64 `json:"goods_price"`
	}
	var cartRows []CartWithGoods
	if err := database.DB.Table("cart_items").
		Select("cart_items.*, goods.name AS goods_name, goods.image_url AS goods_image_url, goods.price AS goods_price").
		Joins("LEFT JOIN goods ON goods.id = cart_items.goods_id").
		Where("cart_items.id IN ? AND cart_items.user_id = ? AND cart_items.tenant_id = ?", req.CartIDs, userID, tenantID).
		Find(&cartRows).Error; err != nil {
		return mpErr(c, 500, "Failed to load cart")
	}
	if len(cartRows) == 0 {
		return mpErr(c, 400, "Cart items not found")
	}

	// Build order
	var total float64
	for _, ci := range cartRows {
		total += ci.GoodsPrice * float64(ci.Quantity)
	}

	receiverAddr := addr.Province + addr.City + addr.District + " " + addr.Address
	orderNo := time.Now().Format("20060102150405") + uuid.New().String()[0:6]

	order := models.ShopOrder{
		TenantID:        &tenantID,
		OrderNo:         orderNo,
		UserID:          &userID,
		Status:          "pending",
		TotalAmount:     total,
		PayAmount:       total,
		PayMethod:       req.PayMethod,
		Remark:          req.Remark,
		ReceiverName:    addr.RecipientName,
		ReceiverPhone:   addr.Phone,
		ReceiverAddress: receiverAddr,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return mpErr(c, 500, "Failed to create order")
	}

	for _, ci := range cartRows {
		oi := models.ShopOrderItem{
			TenantID:    &tenantID,
			OrderID:     order.ID,
			GoodsID:     ci.GoodsID,
			GoodsName:   ci.GoodsName,
			GoodsImage:  ci.GoodsImageURL,
			Price:       ci.GoodsPrice,
			Quantity:    ci.Quantity,
			TotalAmount: ci.GoodsPrice * float64(ci.Quantity),
		}
		if err := tx.Create(&oi).Error; err != nil {
			tx.Rollback()
			return mpErr(c, 500, "Failed to create order item")
		}
	}

	// Remove cart items
	if err := tx.Where("id IN ? AND user_id = ?", req.CartIDs, userID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return mpErr(c, 500, "Failed to clear cart")
	}

	tx.Commit()
	return mpCreated(c, order)
}

// ============================================================
// 6. Goods
// ============================================================

// MpGetGoods lists goods with status=active filter
func MpGetGoods(c *fiber.Ctx) error {
	c.Request().URI().QueryArgs().Set("status", "active")
	return GetGoodsList(c)
}

// MpGetGoodsDetail returns a single goods item
func MpGetGoodsDetail(c *fiber.Ctx) error {
	return GetGoods(c)
}

// MpGetCategories returns active goods categories
func MpGetCategories(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var categories []models.GoodsCategory
	if err := database.DB.Where("tenant_id = ? AND status = ? AND trashed_at IS NULL", tenantID, "active").
		Order("sort_order ASC, created_at ASC").Find(&categories).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch categories")
	}

	return mpOK(c, categories)
}

// MpGetCategoryGoods returns goods for a specific category
func MpGetCategoryGoods(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid category ID")
	}
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.Goods{}).
		Where("tenant_id = ? AND category_id = ? AND status = ? AND trashed_at IS NULL", tenantID, categoryID, "active")
	query.Count(&total)

	var items []models.Goods
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch goods")
	}

	return mpOK(c, fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpSearchGoods searches goods by keyword
func MpSearchGoods(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		return mpErr(c, 400, "Keyword is required")
	}
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.Goods{}).
		Where("tenant_id = ? AND status = ? AND trashed_at IS NULL", tenantID, "active").
		Where("name ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	query.Count(&total)

	var items []models.Goods
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return mpErr(c, 500, "Failed to search goods")
	}

	return mpOK(c, fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// ============================================================
// 7. Cart
// ============================================================

// MpGetCart returns the current user's cart items with goods info
func MpGetCart(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	type CartItemWithGoods struct {
		models.CartItem
		GoodsName     string  `json:"goods_name"`
		GoodsImageURL string  `json:"goods_image_url"`
		GoodsPrice    float64 `json:"goods_price"`
		GoodsStock    int     `json:"goods_stock"`
		GoodsStatus   string  `json:"goods_status"`
	}

	var results []CartItemWithGoods
	if err := database.DB.Table("cart_items").
		Select("cart_items.*, goods.name AS goods_name, goods.image_url AS goods_image_url, goods.price AS goods_price, goods.stock AS goods_stock, goods.status AS goods_status").
		Joins("LEFT JOIN goods ON goods.id = cart_items.goods_id").
		Where("cart_items.user_id = ? AND cart_items.tenant_id = ?", userID, tenantID).
		Order("cart_items.created_at DESC").
		Find(&results).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch cart items")
	}

	return mpOK(c, results)
}

// MpAddToCart adds an item to cart (upsert: if exists, increment quantity)
func MpAddToCart(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		GoodsID  uuid.UUID  `json:"goods_id"`
		SkuID    *uuid.UUID `json:"sku_id"`
		Quantity int        `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.GoodsID == uuid.Nil {
		return mpErr(c, 400, "Goods ID is required")
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Check if item already in cart
	var existing models.CartItem
	query := database.DB.Where("user_id = ? AND tenant_id = ? AND goods_id = ?", userID, tenantID, req.GoodsID)
	if req.SkuID != nil {
		query = query.Where("sku_id = ?", *req.SkuID)
	} else {
		query = query.Where("sku_id IS NULL")
	}

	if err := query.First(&existing).Error; err == nil {
		// Update quantity
		newQty := existing.Quantity + req.Quantity
		if err := database.DB.Model(&existing).Update("quantity", newQty).Error; err != nil {
			return mpErr(c, 500, "Failed to update cart item")
		}
		existing.Quantity = newQty
		return mpOK(c, existing)
	}

	// Create new cart item
	item := models.CartItem{
		TenantID: &tenantID,
		UserID:   userID,
		GoodsID:  req.GoodsID,
		SkuID:    req.SkuID,
		Quantity: req.Quantity,
		Selected: true,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		return mpErr(c, 500, "Failed to add item to cart")
	}

	return mpCreated(c, item)
}

// MpUpdateCartItem updates a cart item's quantity
func MpUpdateCartItem(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid cart item ID")
	}

	var req struct {
		Quantity int   `json:"quantity"`
		Selected *bool `json:"selected"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}

	var item models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ? AND tenant_id = ?", id, userID, tenantID).
		First(&item).Error; err != nil {
		return mpErr(c, 404, "Cart item not found")
	}

	updates := map[string]interface{}{}
	if req.Quantity > 0 {
		updates["quantity"] = req.Quantity
	}
	if req.Selected != nil {
		updates["selected"] = *req.Selected
	}

	if len(updates) == 0 {
		return mpErr(c, 400, "No fields to update")
	}

	if err := database.DB.Model(&item).Updates(updates).Error; err != nil {
		return mpErr(c, 500, "Failed to update cart item")
	}

	database.DB.Where("id = ?", id).First(&item)
	return mpOK(c, item)
}

// MpDeleteCartItem deletes a cart item
func MpDeleteCartItem(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid cart item ID")
	}

	var item models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ? AND tenant_id = ?", id, userID, tenantID).
		First(&item).Error; err != nil {
		return mpErr(c, 404, "Cart item not found")
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		return mpErr(c, 500, "Failed to delete cart item")
	}

	return mpOK(c, fiber.Map{"message": "Cart item deleted"})
}

// ============================================================
// 8. Addresses
// ============================================================

// MpGetAddresses returns addresses for the current user
func MpGetAddresses(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var addresses []models.AddressBook
	if err := database.DB.Where("member_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		Order("is_default DESC, created_at DESC").Find(&addresses).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch addresses")
	}

	return mpOK(c, addresses)
}

// MpGetAddress returns a single address
func MpGetAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid address ID")
	}

	var addr models.AddressBook
	if err := database.DB.Where("id = ? AND member_id = ? AND tenant_id = ? AND trashed_at IS NULL",
		id, userID, tenantID).First(&addr).Error; err != nil {
		return mpErr(c, 404, "Address not found")
	}

	return mpOK(c, addr)
}

// MpCreateAddress creates a new address
func MpCreateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var addr models.AddressBook
	if err := c.BodyParser(&addr); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}

	addr.TenantID = &tenantID
	addr.MemberID = userID

	if addr.RecipientName == "" || addr.Phone == "" || addr.Address == "" {
		return mpErr(c, 400, "Recipient name, phone, and address are required")
	}

	// If this is the first address or marked default, ensure default uniqueness
	if addr.IsDefault {
		database.DB.Model(&models.AddressBook{}).
			Where("member_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
			Update("is_default", false)
	} else {
		// If no other addresses, set as default
		var count int64
		database.DB.Model(&models.AddressBook{}).
			Where("member_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
			Count(&count)
		if count == 0 {
			addr.IsDefault = true
		}
	}

	if err := database.DB.Create(&addr).Error; err != nil {
		return mpErr(c, 500, "Failed to create address")
	}

	return mpCreated(c, addr)
}

// MpUpdateAddress updates an existing address
func MpUpdateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid address ID")
	}

	var addr models.AddressBook
	if err := database.DB.Where("id = ? AND member_id = ? AND tenant_id = ? AND trashed_at IS NULL",
		id, userID, tenantID).First(&addr).Error; err != nil {
		return mpErr(c, 404, "Address not found")
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}

	delete(updates, "id")
	delete(updates, "tenant_id")
	delete(updates, "member_id")
	delete(updates, "created_at")
	delete(updates, "trashed_at")

	// Handle default flag
	if isDefault, ok := updates["is_default"]; ok {
		if val, ok := isDefault.(bool); ok && val {
			database.DB.Model(&models.AddressBook{}).
				Where("member_id = ? AND tenant_id = ? AND id != ? AND trashed_at IS NULL", userID, tenantID, id).
				Update("is_default", false)
		}
	}

	if err := database.DB.Model(&addr).Updates(updates).Error; err != nil {
		return mpErr(c, 500, "Failed to update address")
	}

	database.DB.Where("id = ?", id).First(&addr)
	return mpOK(c, addr)
}

// MpDeleteAddress soft-deletes an address
func MpDeleteAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid address ID")
	}

	var addr models.AddressBook
	if err := database.DB.Where("id = ? AND member_id = ? AND tenant_id = ? AND trashed_at IS NULL",
		id, userID, tenantID).First(&addr).Error; err != nil {
		return mpErr(c, 404, "Address not found")
	}

	if err := database.DB.Model(&addr).Update("trashed_at", time.Now()).Error; err != nil {
		return mpErr(c, 500, "Failed to delete address")
	}

	return mpOK(c, fiber.Map{"message": "Address deleted"})
}

// MpSetDefaultAddress sets an address as default
func MpSetDefaultAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid address ID")
	}

	var addr models.AddressBook
	if err := database.DB.Where("id = ? AND member_id = ? AND tenant_id = ? AND trashed_at IS NULL",
		id, userID, tenantID).First(&addr).Error; err != nil {
		return mpErr(c, 404, "Address not found")
	}

	// Clear all defaults
	database.DB.Model(&models.AddressBook{}).
		Where("member_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		Update("is_default", false)

	// Set this one as default
	database.DB.Model(&addr).Update("is_default", true)

	return mpOK(c, fiber.Map{"message": "Default address updated"})
}

// ============================================================
// 9. Coupons
// ============================================================

// MpGetCoupons returns available coupons for claiming
func MpGetCoupons(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	now := time.Now()
	var total int64
	query := database.DB.Model(&models.Coupon{}).
		Where("tenant_id = ? AND status = ? AND trashed_at IS NULL", tenantID, "active").
		Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", now, now).
		Where("total_count = 0 OR used_count < total_count")
	query.Count(&total)

	var coupons []models.Coupon
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&coupons).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch coupons")
	}

	return mpOK(c, fiber.Map{
		"data":  coupons,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpGetMyCoupons returns user's claimed coupons
func MpGetMyCoupons(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)
	status := c.Query("status") // unused, used, expired

	var total int64
	query := database.DB.Model(&models.CouponReceiveLog{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	type CouponWithInfo struct {
		models.CouponReceiveLog
		CouponName     string  `json:"coupon_name"`
		CouponType     string  `json:"coupon_type"`
		CouponValue    float64 `json:"coupon_value"`
		MinOrderAmount float64 `json:"min_order_amount"`
		MaxDiscount    float64 `json:"max_discount"`
	}

	var results []CouponWithInfo
	q := database.DB.Table("coupon_receive_logs").
		Select("coupon_receive_logs.*, coupons.name AS coupon_name, coupons.type AS coupon_type, coupons.value AS coupon_value, coupons.min_order_amount, coupons.max_discount").
		Joins("LEFT JOIN coupons ON coupons.id = coupon_receive_logs.coupon_id").
		Where("coupon_receive_logs.user_id = ? AND coupon_receive_logs.tenant_id = ?", userID, tenantID)
	if status != "" {
		q = q.Where("coupon_receive_logs.status = ?", status)
	}
	if err := q.Order("coupon_receive_logs.created_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch user coupons")
	}

	return mpOK(c, fiber.Map{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpClaimCoupon claims a coupon for the current user
func MpClaimCoupon(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	couponID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid coupon ID")
	}

	// Check coupon exists and is available
	var coupon models.Coupon
	if err := database.DB.Where("id = ? AND tenant_id = ? AND status = ? AND trashed_at IS NULL",
		couponID, tenantID, "active").First(&coupon).Error; err != nil {
		return mpErr(c, 404, "Coupon not found or unavailable")
	}

	// Check availability
	now := time.Now()
	if coupon.StartAt != nil && now.Before(*coupon.StartAt) {
		return mpErr(c, 400, "Coupon is not yet available")
	}
	if coupon.EndAt != nil && now.After(*coupon.EndAt) {
		return mpErr(c, 400, "Coupon has expired")
	}
	if coupon.TotalCount > 0 && coupon.UsedCount >= coupon.TotalCount {
		return mpErr(c, 400, "Coupon is fully claimed")
	}

	// Check if user already claimed this coupon
	var existingCount int64
	database.DB.Model(&models.CouponReceiveLog{}).
		Where("coupon_id = ? AND user_id = ? AND tenant_id = ?", couponID, userID, tenantID).
		Count(&existingCount)
	if existingCount > 0 {
		return mpErr(c, 400, "You have already claimed this coupon")
	}

	// Claim it
	log := models.CouponReceiveLog{
		TenantID: &tenantID,
		CouponID: couponID,
		UserID:   userID,
		Status:   "unused",
	}
	if err := database.DB.Create(&log).Error; err != nil {
		return mpErr(c, 500, "Failed to claim coupon")
	}

	// Increment used count
	database.DB.Model(&coupon).Update("used_count", gorm.Expr("used_count + 1"))

	return mpOK(c, log)
}

// ============================================================
// 10. Content
// ============================================================

// MpGetHelpList returns help articles
func MpGetHelpList(c *fiber.Ctx) error {
	return GetHelpArticles(c)
}

// MpGetHelpDetail returns a single help article
func MpGetHelpDetail(c *fiber.Ctx) error {
	return GetHelpArticle(c)
}

// MpGetArticles returns articles
func MpGetArticles(c *fiber.Ctx) error {
	return GetArticles(c)
}

// MpGetArticleDetail returns a single article
func MpGetArticleDetail(c *fiber.Ctx) error {
	return GetArticle(c)
}

// MpGetNotices returns published notifications of type 'notice'
func MpGetNotices(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.Notification{}).
		Where("tenant_id = ? AND type = ? AND status = ? AND trashed_at IS NULL", tenantID, "notice", "published")
	query.Count(&total)

	var notices []models.Notification
	if err := query.Order("published_at DESC, created_at DESC").Offset(offset).Limit(limit).Find(&notices).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch notices")
	}

	return mpOK(c, fiber.Map{
		"data":  notices,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpGetNoticeDetail returns a single notice
func MpGetNoticeDetail(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid notice ID")
	}

	var notice models.Notification
	if err := database.DB.Where("id = ? AND tenant_id = ? AND type = ? AND trashed_at IS NULL",
		id, tenantID, "notice").First(&notice).Error; err != nil {
		return mpErr(c, 404, "Notice not found")
	}

	return mpOK(c, notice)
}

// ============================================================
// 11. Tracking
// ============================================================

// MpTrackPackage tracks a package by tracking number
func MpTrackPackage(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	trackingNo := strings.TrimSpace(c.Params("trackingNo"))
	if trackingNo == "" {
		return mpErr(c, 400, "Tracking number is required")
	}

	var pkg models.Package
	if err := database.DB.Where("tracking_number = ? AND tenant_id = ? AND trashed_at IS NULL",
		trackingNo, tenantID).First(&pkg).Error; err != nil {
		return mpErr(c, 404, "Package not found")
	}

	// Get status logs
	var statusLogs []models.PackageStatusLog
	database.DB.Where("package_id = ? AND tenant_id = ?", pkg.ID, tenantID).
		Order("created_at DESC").Find(&statusLogs)

	return mpOK(c, fiber.Map{
		"package":     pkg,
		"status_logs": statusLogs,
	})
}

// ============================================================
// 12. Warehouses & Routes
// ============================================================

// MpGetWarehouses lists warehouses
func MpGetWarehouses(c *fiber.Ctx) error {
	return GetWarehouses(c)
}

// MpGetWarehouseDetail returns a single warehouse
func MpGetWarehouseDetail(c *fiber.Ctx) error {
	return GetWarehouse(c)
}

// MpGetRoutes returns active shipping routes
func MpGetRoutes(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.ShippingRoute{}).
		Where("tenant_id = ? AND status = ? AND trashed_at IS NULL", tenantID, "active")
	query.Count(&total)

	var routes []models.ShippingRoute
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(limit).Find(&routes).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch routes")
	}

	return mpOK(c, fiber.Map{
		"data":  routes,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpGetRouteDetail returns a single route with pricing tiers
func MpGetRouteDetail(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid route ID")
	}

	var route models.ShippingRoute
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&route).Error; err != nil {
		return mpErr(c, 404, "Route not found")
	}

	var tiers []models.RoutePricingTier
	database.DB.Where("route_id = ? AND tenant_id = ?", id, tenantID).
		Order("weight_min ASC").Find(&tiers)

	return mpOK(c, fiber.Map{
		"route":         route,
		"pricing_tiers": tiers,
	})
}

// MpCalculateEstimate calculates shipping cost estimate
func MpCalculateEstimate(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		RouteID       uuid.UUID  `json:"route_id"`
		Weight        float64    `json:"weight"`
		Length        float64    `json:"length"`
		Width         float64    `json:"width"`
		Height        float64    `json:"height"`
		MemberLevelID *uuid.UUID `json:"member_level_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.RouteID == uuid.Nil {
		return mpErr(c, 400, "Route ID is required")
	}
	if req.Weight <= 0 {
		return mpErr(c, 400, "Weight must be greater than 0")
	}

	// Get route
	var route models.ShippingRoute
	if err := database.DB.Where("id = ? AND tenant_id = ? AND status = ? AND trashed_at IS NULL",
		req.RouteID, tenantID, "active").First(&route).Error; err != nil {
		return mpErr(c, 404, "Route not found")
	}

	// Calculate volume weight
	var volumeWeight float64
	if req.Length > 0 && req.Width > 0 && req.Height > 0 && route.VolumeWeightRatio > 0 {
		volumeWeight = (req.Length * req.Width * req.Height) / route.VolumeWeightRatio
	}

	chargeableWeight := math.Max(req.Weight, volumeWeight)

	// Find matching pricing tier
	tierQuery := database.DB.Where("route_id = ? AND tenant_id = ?", req.RouteID, tenantID).
		Where("weight_min <= ? AND (weight_max >= ? OR weight_max = 0)", chargeableWeight, chargeableWeight)
	if req.MemberLevelID != nil {
		tierQuery = tierQuery.Where("member_level_id = ? OR member_level_id IS NULL", *req.MemberLevelID).
			Order("member_level_id DESC NULLS LAST")
	} else {
		tierQuery = tierQuery.Where("member_level_id IS NULL")
	}

	var tier models.RoutePricingTier
	var shippingFee float64

	if err := tierQuery.First(&tier).Error; err == nil {
		if tier.FirstWeight > 0 && chargeableWeight > 0 {
			// First weight pricing
			if chargeableWeight <= tier.FirstWeight {
				shippingFee = tier.FirstWeightPrice
			} else {
				additionalWeight := chargeableWeight - tier.FirstWeight
				shippingFee = tier.FirstWeightPrice + (additionalWeight * tier.AdditionalWeightPrice)
			}
		} else {
			// Unit price
			shippingFee = chargeableWeight * tier.UnitPrice
		}
	} else {
		// Fallback: try any tier for this route
		var fallbackTier models.RoutePricingTier
		if err := database.DB.Where("route_id = ? AND tenant_id = ? AND member_level_id IS NULL",
			req.RouteID, tenantID).Order("weight_min ASC").First(&fallbackTier).Error; err == nil {
			shippingFee = chargeableWeight * fallbackTier.UnitPrice
		}
	}

	// Round to 2 decimal places
	shippingFee = math.Round(shippingFee*100) / 100

	return mpOK(c, fiber.Map{
		"route_name":        route.Name,
		"actual_weight":     req.Weight,
		"volume_weight":     math.Round(volumeWeight*1000) / 1000,
		"chargeable_weight": math.Round(chargeableWeight*1000) / 1000,
		"shipping_fee":      shippingFee,
		"currency":          "TWD",
		"estimated_days":    route.EstimatedDays,
		"transport_mode":    route.TransportMode,
	})
}

// ============================================================
// 13. Messages
// ============================================================

// MpGetMessages returns notifications for the current user
func MpGetMessages(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.Notification{}).
		Where("tenant_id = ? AND trashed_at IS NULL AND status = ?", tenantID, "published").
		Where("(target_type = 'all' OR (target_type = 'user' AND target_id = ?))", userID)
	query.Count(&total)

	var messages []models.Notification
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch messages")
	}

	return mpOK(c, fiber.Map{
		"data":  messages,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpMarkRead marks a notification as read (stored in extra_fields)
func MpMarkRead(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid message ID")
	}

	var notification models.Notification
	if err := database.DB.Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&notification).Error; err != nil {
		return mpErr(c, 404, "Message not found")
	}

	// Track read status in extra_fields
	extra := notification.ExtraFields
	if extra == nil {
		extra = make(models.JSONB)
	}

	readByRaw, ok := extra["read_by"]
	var readBy []string
	if ok {
		// Parse existing read_by list
		if rawBytes, err := json.Marshal(readByRaw); err == nil {
			_ = json.Unmarshal(rawBytes, &readBy)
		}
	}

	// Check if already read
	uid := userID.String()
	alreadyRead := false
	for _, r := range readBy {
		if r == uid {
			alreadyRead = true
			break
		}
	}

	if !alreadyRead {
		readBy = append(readBy, uid)
		extra["read_by"] = readBy
		database.DB.Model(&notification).Update("extra_fields", extra)
	}

	return mpOK(c, fiber.Map{"message": "Marked as read"})
}

// ============================================================
// 14. Favorites
// ============================================================

// MpGetFavorites returns user's favorited goods with goods info
func MpGetFavorites(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).Count(&total)

	type FavoriteWithGoods struct {
		models.Favorite
		GoodsName     string  `json:"goods_name"`
		GoodsImageURL string  `json:"goods_image_url"`
		GoodsPrice    float64 `json:"goods_price"`
		GoodsStatus   string  `json:"goods_status"`
	}

	var results []FavoriteWithGoods
	if err := database.DB.Table("favorites").
		Select("favorites.*, goods.name AS goods_name, goods.image_url AS goods_image_url, goods.price AS goods_price, goods.status AS goods_status").
		Joins("LEFT JOIN goods ON goods.id = favorites.goods_id").
		Where("favorites.user_id = ? AND favorites.tenant_id = ?", userID, tenantID).
		Order("favorites.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&results).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch favorites")
	}

	return mpOK(c, fiber.Map{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpAddFavorite adds a goods item to favorites
func MpAddFavorite(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		GoodsID uuid.UUID `json:"goods_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.GoodsID == uuid.Nil {
		return mpErr(c, 400, "Goods ID is required")
	}

	// Check if already favorited
	var count int64
	database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND tenant_id = ? AND goods_id = ?", userID, tenantID, req.GoodsID).
		Count(&count)
	if count > 0 {
		return mpErr(c, 400, "Already in favorites")
	}

	fav := models.Favorite{
		TenantID: &tenantID,
		UserID:   userID,
		GoodsID:  req.GoodsID,
	}

	if err := database.DB.Create(&fav).Error; err != nil {
		return mpErr(c, 500, "Failed to add to favorites")
	}

	return mpCreated(c, fav)
}

// MpRemoveFavorite removes a goods item from favorites
func MpRemoveFavorite(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid favorite ID")
	}

	var fav models.Favorite
	if err := database.DB.Where("id = ? AND user_id = ? AND tenant_id = ?", id, userID, tenantID).
		First(&fav).Error; err != nil {
		return mpErr(c, 404, "Favorite not found")
	}

	if err := database.DB.Delete(&fav).Error; err != nil {
		return mpErr(c, 500, "Failed to remove from favorites")
	}

	return mpOK(c, fiber.Map{"message": "Removed from favorites"})
}

// ============================================================
// 15. History
// ============================================================

// MpGetHistory returns user's browsing history with goods info
func MpGetHistory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	database.DB.Model(&models.BrowsingHistory{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).Count(&total)

	type HistoryWithGoods struct {
		models.BrowsingHistory
		GoodsName     string  `json:"goods_name"`
		GoodsImageURL string  `json:"goods_image_url"`
		GoodsPrice    float64 `json:"goods_price"`
		GoodsStatus   string  `json:"goods_status"`
	}

	var results []HistoryWithGoods
	if err := database.DB.Table("browsing_histories").
		Select("browsing_histories.*, goods.name AS goods_name, goods.image_url AS goods_image_url, goods.price AS goods_price, goods.status AS goods_status").
		Joins("LEFT JOIN goods ON goods.id = browsing_histories.goods_id").
		Where("browsing_histories.user_id = ? AND browsing_histories.tenant_id = ?", userID, tenantID).
		Order("browsing_histories.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&results).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch browsing history")
	}

	return mpOK(c, fiber.Map{
		"data":  results,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpClearHistory clears all browsing history for the current user
func MpClearHistory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	if err := database.DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Delete(&models.BrowsingHistory{}).Error; err != nil {
		return mpErr(c, 500, "Failed to clear browsing history")
	}

	return mpOK(c, fiber.Map{"message": "Browsing history cleared"})
}

// ============================================================
// 16. Feedback
// ============================================================

// MpSubmitFeedback creates a feedback record
func MpSubmitFeedback(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		Type        string   `json:"type"`
		Content     string   `json:"content"`
		Images      []string `json:"images"`
		ContactInfo string   `json:"contact_info"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.Content == "" {
		return mpErr(c, 400, "Content is required")
	}
	if req.Type == "" {
		req.Type = "general"
	}

	imagesJSON := models.JSONB{}
	if len(req.Images) > 0 {
		imagesJSON["list"] = req.Images
	}

	feedback := models.Feedback{
		TenantID:    &tenantID,
		UserID:      userID,
		Type:        req.Type,
		Content:     req.Content,
		Images:      imagesJSON,
		ContactInfo: req.ContactInfo,
		Status:      "pending",
	}

	if err := database.DB.Create(&feedback).Error; err != nil {
		return mpErr(c, 500, "Failed to submit feedback")
	}

	return mpCreated(c, feedback)
}

// ============================================================
// 17. Reviews
// ============================================================

// MpGetReviews returns reviews for a goods item
func MpGetReviews(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	goodsID := c.Query("goods_id")
	if goodsID == "" {
		return mpErr(c, 400, "goods_id is required")
	}
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.GoodsReview{}).
		Where("tenant_id = ? AND goods_id = ? AND status = ? AND trashed_at IS NULL", tenantID, goodsID, "approved")
	query.Count(&total)

	var reviews []models.GoodsReview
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch reviews")
	}

	return mpOK(c, fiber.Map{
		"data":  reviews,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpCreateReview creates a new review
func MpCreateReview(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		GoodsID uuid.UUID  `json:"goods_id"`
		OrderID *uuid.UUID `json:"order_id"`
		Content string     `json:"content"`
		Rating  int        `json:"rating"`
		Images  []string   `json:"images"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.GoodsID == uuid.Nil {
		return mpErr(c, 400, "Goods ID is required")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return mpErr(c, 400, "Rating must be between 1 and 5")
	}
	if req.Content == "" {
		return mpErr(c, 400, "Content is required")
	}

	imagesJSON := models.JSONB{}
	if len(req.Images) > 0 {
		imagesJSON["list"] = req.Images
	}

	review := models.GoodsReview{
		TenantID: &tenantID,
		GoodsID:  req.GoodsID,
		UserID:   userID,
		OrderID:  req.OrderID,
		Content:  req.Content,
		Rating:   req.Rating,
		Images:   imagesJSON,
		Status:   "pending",
	}

	if err := database.DB.Create(&review).Error; err != nil {
		return mpErr(c, 500, "Failed to create review")
	}

	return mpCreated(c, review)
}

// ============================================================
// 18. Refunds
// ============================================================

// MpGetRefunds returns refunds for the current user
func MpGetRefunds(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var total int64
	query := database.DB.Model(&models.OrderRefund{}).
		Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID)
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	var refunds []models.OrderRefund
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&refunds).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch refunds")
	}

	return mpOK(c, fiber.Map{
		"data":  refunds,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpGetRefundDetail returns a single refund
func MpGetRefundDetail(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return mpErr(c, 400, "Invalid refund ID")
	}

	var refund models.OrderRefund
	if err := database.DB.Where("id = ? AND user_id = ? AND tenant_id = ? AND trashed_at IS NULL",
		id, userID, tenantID).First(&refund).Error; err != nil {
		return mpErr(c, 404, "Refund not found")
	}

	return mpOK(c, refund)
}

// MpCreateRefund creates a refund request
func MpCreateRefund(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		OrderID     uuid.UUID  `json:"order_id"`
		OrderItemID *uuid.UUID `json:"order_item_id"`
		Type        string     `json:"type"`
		Reason      string     `json:"reason"`
		Description string     `json:"description"`
		Amount      float64    `json:"amount"`
		Images      []string   `json:"images"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.OrderID == uuid.Nil {
		return mpErr(c, 400, "Order ID is required")
	}
	if req.Reason == "" {
		return mpErr(c, 400, "Reason is required")
	}
	if req.Type == "" {
		req.Type = "refund"
	}

	imagesJSON := models.JSONB{}
	if len(req.Images) > 0 {
		imagesJSON["list"] = req.Images
	}

	refund := models.OrderRefund{
		TenantID:    &tenantID,
		OrderID:     req.OrderID,
		OrderItemID: req.OrderItemID,
		UserID:      userID,
		RefundNo:    fmt.Sprintf("RF%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000)),
		Type:        req.Type,
		Reason:      req.Reason,
		Description: req.Description,
		Amount:      req.Amount,
		Images:      imagesJSON,
		Status:      "pending",
	}

	if err := database.DB.Create(&refund).Error; err != nil {
		return mpErr(c, 500, "Failed to create refund request")
	}

	return mpCreated(c, refund)
}

// ============================================================
// 19. Dealer
// ============================================================

// MpGetDealerInfo returns the current user's dealer status
func MpGetDealerInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var dealer models.Dealer
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		First(&dealer).Error; err != nil {
		// Check if there's a pending application
		var app models.DealerApplication
		if err := database.DB.Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
			Order("created_at DESC").First(&app).Error; err == nil {
			return mpOK(c, fiber.Map{
				"is_dealer":          false,
				"application_status": app.Status,
				"application":        app,
			})
		}
		return mpOK(c, fiber.Map{
			"is_dealer":          false,
			"application_status": "none",
		})
	}

	// Get level info
	var level models.DealerLevel
	if dealer.LevelID != nil {
		database.DB.Where("id = ?", *dealer.LevelID).First(&level)
	}

	// Count team members
	var teamCount int64
	database.DB.Model(&models.Dealer{}).
		Where("parent_id = ? AND tenant_id = ? AND trashed_at IS NULL", dealer.ID, tenantID).
		Count(&teamCount)

	availableCommission := dealer.TotalCommission - dealer.WithdrawnCommission

	return mpOK(c, fiber.Map{
		"is_dealer":            true,
		"dealer":               dealer,
		"level":                level,
		"team_count":           teamCount,
		"available_commission": availableCommission,
	})
}

// MpApplyDealer submits a dealer application
func MpApplyDealer(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	// Check if already a dealer
	var dealerCount int64
	database.DB.Model(&models.Dealer{}).
		Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		Count(&dealerCount)
	if dealerCount > 0 {
		return mpErr(c, 400, "You are already a dealer")
	}

	// Check for pending application
	var pendingCount int64
	database.DB.Model(&models.DealerApplication{}).
		Where("user_id = ? AND tenant_id = ? AND status = ? AND trashed_at IS NULL", userID, tenantID, "pending").
		Count(&pendingCount)
	if pendingCount > 0 {
		return mpErr(c, 400, "You already have a pending application")
	}

	var req struct {
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Reason   string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.RealName == "" || req.Phone == "" {
		return mpErr(c, 400, "Real name and phone are required")
	}

	app := models.DealerApplication{
		TenantID: &tenantID,
		UserID:   userID,
		RealName: req.RealName,
		Phone:    req.Phone,
		Reason:   req.Reason,
		Status:   "pending",
	}

	if err := database.DB.Create(&app).Error; err != nil {
		return mpErr(c, 500, "Failed to submit application")
	}

	return mpCreated(c, app)
}

// MpGetDealerOrders returns dealer's commission orders
func MpGetDealerOrders(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	// Find dealer record
	var dealer models.Dealer
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		First(&dealer).Error; err != nil {
		return mpErr(c, 404, "You are not a dealer")
	}

	var total int64
	query := database.DB.Model(&models.DealerOrder{}).
		Where("dealer_id = ? AND tenant_id = ?", dealer.ID, tenantID)
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	var orders []models.DealerOrder
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch dealer orders")
	}

	return mpOK(c, fiber.Map{
		"data":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpGetDealerWithdrawals returns dealer's withdrawal records
func MpGetDealerWithdrawals(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var dealer models.Dealer
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		First(&dealer).Error; err != nil {
		return mpErr(c, 404, "You are not a dealer")
	}

	var total int64
	query := database.DB.Model(&models.DealerWithdrawal{}).
		Where("dealer_id = ? AND tenant_id = ? AND trashed_at IS NULL", dealer.ID, tenantID)
	query.Count(&total)

	var withdrawals []models.DealerWithdrawal
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&withdrawals).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch withdrawals")
	}

	return mpOK(c, fiber.Map{
		"data":  withdrawals,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MpRequestWithdraw requests a withdrawal of dealer commission
func MpRequestWithdraw(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var dealer models.Dealer
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		First(&dealer).Error; err != nil {
		return mpErr(c, 404, "You are not a dealer")
	}

	var req struct {
		Amount      float64                `json:"amount"`
		Method      string                 `json:"method"`
		AccountInfo map[string]interface{} `json:"account_info"`
	}
	if err := c.BodyParser(&req); err != nil {
		return mpErr(c, 400, "Invalid request body")
	}
	if req.Amount <= 0 {
		return mpErr(c, 400, "Amount must be greater than 0")
	}
	if req.Method == "" {
		return mpErr(c, 400, "Withdrawal method is required")
	}

	availableCommission := dealer.TotalCommission - dealer.WithdrawnCommission
	if req.Amount > availableCommission {
		return mpErr(c, 400, fmt.Sprintf("Insufficient balance. Available: %.2f", availableCommission))
	}

	accountInfo := models.JSONB{}
	for k, v := range req.AccountInfo {
		accountInfo[k] = v
	}

	withdrawal := models.DealerWithdrawal{
		TenantID:    &tenantID,
		DealerID:    dealer.ID,
		Amount:      req.Amount,
		Method:      req.Method,
		AccountInfo: accountInfo,
		Status:      "pending",
	}

	if err := database.DB.Create(&withdrawal).Error; err != nil {
		return mpErr(c, 500, "Failed to create withdrawal request")
	}

	return mpCreated(c, withdrawal)
}

// MpGetDealerTeam returns the dealer's team members
func MpGetDealerTeam(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	page, limit, offset := mpPaginate(c)

	var dealer models.Dealer
	if err := database.DB.Where("user_id = ? AND tenant_id = ? AND trashed_at IS NULL", userID, tenantID).
		First(&dealer).Error; err != nil {
		return mpErr(c, 404, "You are not a dealer")
	}

	var total int64
	query := database.DB.Model(&models.Dealer{}).
		Where("parent_id = ? AND tenant_id = ? AND trashed_at IS NULL", dealer.ID, tenantID)
	query.Count(&total)

	type TeamMember struct {
		models.Dealer
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
		UserPhone string `json:"user_phone"`
	}

	var members []TeamMember
	if err := database.DB.Table("dealers").
		Select("dealers.*, users.name AS user_name, users.email AS user_email, users.phone AS user_phone").
		Joins("LEFT JOIN users ON users.id = dealers.user_id").
		Where("dealers.parent_id = ? AND dealers.tenant_id = ? AND dealers.trashed_at IS NULL", dealer.ID, tenantID).
		Order("dealers.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&members).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch team members")
	}

	return mpOK(c, fiber.Map{
		"data":  members,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// ============================================================
// 20. Misc
// ============================================================

// MpGetProhibitedItems returns a static list of prohibited items
func MpGetProhibitedItems(c *fiber.Ctx) error {
	items := []fiber.Map{
		{"category": "Dangerous Goods", "items": []string{
			"Explosives and fireworks", "Flammable liquids and gases",
			"Toxic and infectious substances", "Radioactive materials",
			"Corrosive substances", "Compressed gas cylinders",
		}},
		{"category": "Restricted Items", "items": []string{
			"Drugs and narcotics", "Weapons and ammunition",
			"Counterfeit goods", "Pornographic materials",
			"Endangered wildlife products", "Currency and negotiable instruments",
		}},
		{"category": "Battery & Electronics", "items": []string{
			"Loose lithium batteries (>100Wh)", "Damaged or recalled electronics",
			"Hoverboards with non-certified batteries",
		}},
		{"category": "Food & Liquids", "items": []string{
			"Perishable food without proper packaging", "Alcoholic beverages (certain routes)",
			"Live animals", "Plants and seeds (without permit)",
		}},
	}

	return mpOK(c, items)
}

// MpGetContent returns a content page by slug
func MpGetContent(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" {
		return mpErr(c, 400, "Slug is required")
	}

	var page models.ContentPage
	if err := database.DB.Where("slug = ? AND tenant_id = ? AND status = ?", slug, tenantID, "active").
		First(&page).Error; err != nil {
		return mpErr(c, 404, "Content page not found")
	}

	return mpOK(c, page)
}

// MpGetInviteInfo returns the user's invite code and invite stats
func MpGetInviteInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return mpErr(c, 404, "User not found")
	}

	// Generate invite code from user ID (deterministic)
	inviteCode := strings.ToUpper(userID.String()[:8])
	if code, ok := user.ExtraFields["invite_code"].(string); ok && code != "" {
		inviteCode = code
	} else {
		// Save the invite code
		extra := user.ExtraFields
		if extra == nil {
			extra = make(models.JSONB)
		}
		extra["invite_code"] = inviteCode
		database.DB.Model(&user).Update("extra_fields", extra)
	}

	// Count invited users (users whose extra_fields has invited_by = this user's invite code)
	var invitedCount int64
	database.DB.Model(&models.User{}).
		Where("tenant_id = ? AND trashed_at IS NULL AND extra_fields->>'invited_by' = ?", tenantID, inviteCode).
		Count(&invitedCount)

	return mpOK(c, fiber.Map{
		"invite_code":   inviteCode,
		"invited_count": invitedCount,
	})
}

// MpGetValueAddedServices returns active value-added services
func MpGetValueAddedServices(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var services []models.ValueAddedService
	if err := database.DB.Where("tenant_id = ? AND status = ? AND trashed_at IS NULL", tenantID, "active").
		Order("sort_order ASC, created_at ASC").Find(&services).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch value-added services")
	}

	return mpOK(c, services)
}

// MpGetInsuranceOptions returns active insurance products
func MpGetInsuranceOptions(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var products []models.InsuranceProduct
	if err := database.DB.Where("tenant_id = ? AND status = ? AND trashed_at IS NULL", tenantID, "active").
		Order("sort_order ASC, created_at ASC").Find(&products).Error; err != nil {
		return mpErr(c, 500, "Failed to fetch insurance products")
	}

	return mpOK(c, products)
}

// MpGetPageDesigns returns page designs (proxy to existing handler)
func MpGetPageDesigns(c *fiber.Ctx) error {
	return GetPageDesigns(c)
}

// MpGetAppSettings returns app settings (proxy to existing handler)
func MpGetAppSettings(c *fiber.Ctx) error {
	return GetAppSettings(c)
}
