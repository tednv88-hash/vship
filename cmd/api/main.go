package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vship/config"
	"vship/internal/database"
	"vship/internal/handlers"
	"vship/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Change working directory to executable location
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		if _, err := os.Stat(filepath.Join(execDir, ".env")); err == nil {
			os.Chdir(execDir)
		}
	}

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load config
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := handlers.RunMigrations(); err != nil {
		log.Printf("Warning: Migration error: %v", err)
	}

	// Setup template engine
	engine := html.New("./web/templates", ".html")
	engine.Reload(true) // Enable reload in development

	// Template functions
	engine.AddFunc("lower", strings.ToLower)
	engine.AddFunc("upper", strings.ToUpper)
	engine.AddFunc("title", strings.Title)
	engine.AddFunc("contains", strings.Contains)
	engine.AddFunc("hasPrefix", strings.HasPrefix)
	engine.AddFunc("hasSuffix", strings.HasSuffix)
	engine.AddFunc("replace", strings.ReplaceAll)
	engine.AddFunc("seq", func(n int) []int {
		s := make([]int, n)
		for i := range s {
			s[i] = i + 1
		}
		return s
	})
	engine.AddFunc("add", func(a, b int) int { return a + b })
	engine.AddFunc("sub", func(a, b int) int { return a - b })
	engine.AddFunc("mul", func(a, b float64) float64 { return a * b })
	engine.AddFunc("div", func(a, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	})
	engine.AddFunc("eq", func(a, b interface{}) bool { return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b) })
	engine.AddFunc("ne", func(a, b interface{}) bool { return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b) })
	engine.AddFunc("truncate", func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		return s[:n] + "..."
	})
	engine.AddFunc("statusBadge", func(status string) string {
		switch strings.ToLower(status) {
		case "active", "completed", "delivered", "arrived":
			return "success"
		case "pending", "processing", "in_transit":
			return "warning"
		case "inactive", "cancelled", "rejected":
			return "danger"
		case "draft":
			return "secondary"
		default:
			return "info"
		}
	})
	engine.AddFunc("formatWeight", func(w float64) string {
		return fmt.Sprintf("%.2f", w)
	})
	engine.AddFunc("formatMoney", func(amount float64) string {
		return fmt.Sprintf("%.2f", amount)
	})

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			accept := c.Get("Accept")
			if strings.Contains(accept, "application/json") || strings.HasPrefix(c.Path(), "/api") {
				return c.Status(code).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.Status(code).Render("pages/error", fiber.Map{
				"Title":   fmt.Sprintf("Error %d", code),
				"Code":    code,
				"Message": err.Error(),
			}, "layouts/cms_layout")
		},
	})

	// Middleware
	app.Use(recover.New())

	appEnv := strings.ToLower(os.Getenv("APP_ENV"))
	if appEnv != "production" && appEnv != "prod" {
		app.Use(logger.New())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Tenant-Subdomain",
		AllowCredentials: false,
	}))

	// Static files
	app.Static("/static", "./web/static", fiber.Static{
		Compress: true,
	})

	// ===== PUBLIC PAGE ROUTES =====
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/", handlers.RenderHome)
	app.Get("/login", handlers.RenderLogin)
	app.Get("/reset-password", handlers.RenderResetPassword)
	app.Get("/preview/page-designs/:id", handlers.RenderPageDesignPreview)

	// ===== PUBLIC API ROUTES =====
	app.Post("/api/v1/auth/login", handlers.Login)
	app.Post("/api/v1/auth/register", handlers.Register)
	app.Post("/api/v1/auth/logout", handlers.Logout)
	app.Post("/api/v1/auth/forgot-password", handlers.ForgotPassword)
	app.Post("/api/v1/auth/reset-password", handlers.ResetPassword)
	app.Get("/api/v1/auth/google/config", handlers.GetGoogleOAuthConfig)
	app.Post("/api/v1/auth/google", handlers.GoogleLogin)

	// ===== MINIPROGRAM PUBLIC ROUTES (no auth required) =====
	app.Post("/api/v1/mp/auth/login", handlers.MpLogin)
	app.Post("/api/v1/mp/auth/register", handlers.MpRegister)
	app.Post("/api/v1/mp/auth/wechat-login", handlers.MpWechatLogin)
	app.Post("/api/v1/mp/auth/send-code", handlers.MpSendCode)

	mpPub := app.Group("/api/v1/mp", middleware.MpPublicTenantMiddleware)

	// --- Goods (browsable without login) ---
	mpPub.Get("/goods", handlers.MpGetGoods)
	mpPub.Get("/goods/categories", handlers.MpGetCategories)
	mpPub.Get("/goods/categories/:id/goods", handlers.MpGetCategoryGoods)
	mpPub.Get("/goods/search", handlers.MpSearchGoods)
	mpPub.Get("/goods/:id", handlers.MpGetGoodsDetail)

	// --- Content (public) ---
	mpPub.Get("/help", handlers.MpGetHelpList)
	mpPub.Get("/help/:id", handlers.MpGetHelpDetail)
	mpPub.Get("/articles", handlers.MpGetArticles)
	mpPub.Get("/articles/:id", handlers.MpGetArticleDetail)
	mpPub.Get("/notices", handlers.MpGetNotices)
	mpPub.Get("/notices/:id", handlers.MpGetNoticeDetail)

	// --- Warehouses & Routes (public info) ---
	mpPub.Get("/warehouses", handlers.MpGetWarehouses)
	mpPub.Get("/warehouses/:id", handlers.MpGetWarehouseDetail)
	mpPub.Get("/routes", handlers.MpGetRoutes)
	mpPub.Get("/routes/:id", handlers.MpGetRouteDetail)
	mpPub.Post("/routes/estimate", handlers.MpCalculateEstimate)

	// --- Misc public ---
	mpPub.Get("/prohibited-items", handlers.MpGetProhibitedItems)
	mpPub.Get("/content/:slug", handlers.MpGetContent)
	mpPub.Get("/coupons", handlers.MpGetCoupons)
	mpPub.Get("/reviews", handlers.MpGetReviews)
	mpPub.Get("/value-added-services", handlers.MpGetValueAddedServices)
	mpPub.Get("/insurance-options", handlers.MpGetInsuranceOptions)
	mpPub.Get("/page-designs", handlers.MpGetPageDesigns)
	mpPub.Get("/app-settings", handlers.MpGetAppSettings)

	// ===== MINIPROGRAM AUTHENTICATED ROUTES =====
	mp := app.Group("/api/v1/mp", middleware.AuthMiddleware, middleware.TenantMiddleware)

	// --- User ---
	mp.Get("/user/info", handlers.MpGetUserInfo)
	mp.Put("/user/info", handlers.MpUpdateUserInfo)
	mp.Get("/user/balance", handlers.MpGetBalance)
	mp.Get("/user/points", handlers.MpGetPoints)
	mp.Post("/user/recharge", handlers.MpRecharge)
	mp.Post("/user/verify-identity", handlers.MpVerifyIdentity)
	mp.Get("/user/identity-status", handlers.MpGetIdentityStatus)
	mp.Post("/user/sign-in", handlers.MpSignIn)
	mp.Get("/user/sign-in/status", handlers.MpGetSignInStatus)
	mp.Post("/user/bind-phone", handlers.MpBindPhone)

	// --- Packages ---
	mp.Get("/packages", handlers.MpGetPackages)
	mp.Post("/packages/forecast", handlers.MpForecastPackage)
	mp.Get("/packages/forecasts", handlers.MpGetForecasts)
	mp.Post("/packages/merge", handlers.MpMergePackages)
	mp.Post("/packages/split", handlers.MpSplitPackage)
	mp.Get("/packages/:id", handlers.MpGetPackage)

	// --- Orders (Consolidation) ---
	mp.Get("/orders", handlers.MpGetOrders)
	mp.Get("/orders/:id", handlers.MpGetOrder)
	mp.Post("/orders", handlers.MpCreateOrder)
	mp.Put("/orders/:id/cancel", handlers.MpCancelOrder)
	mp.Put("/orders/:id/pay", handlers.MpPayOrder)

	// --- Shop/Mall Orders ---
	mp.Get("/shop-orders", handlers.MpGetShopOrders)
	mp.Get("/shop-orders/:id", handlers.MpGetShopOrder)
	mp.Post("/shop-orders/checkout", handlers.MpCheckoutCart)
	mp.Post("/shop-orders/:id/pay", handlers.MpPayShopOrder)
	mp.Put("/shop-orders/:id/cancel", handlers.MpCancelShopOrder)
	mp.Put("/shop-orders/:id/confirm", handlers.MpConfirmReceiveShopOrder)

	// --- Cart ---
	mp.Get("/cart", handlers.MpGetCart)
	mp.Post("/cart", handlers.MpAddToCart)
	mp.Put("/cart/:id", handlers.MpUpdateCartItem)
	mp.Delete("/cart/:id", handlers.MpDeleteCartItem)

	// --- Addresses ---
	mp.Get("/addresses", handlers.MpGetAddresses)
	mp.Get("/addresses/:id", handlers.MpGetAddress)
	mp.Post("/addresses", handlers.MpCreateAddress)
	mp.Put("/addresses/:id", handlers.MpUpdateAddress)
	mp.Delete("/addresses/:id", handlers.MpDeleteAddress)
	mp.Put("/addresses/:id/default", handlers.MpSetDefaultAddress)

	// --- Coupons (auth) ---
	mp.Get("/coupons/mine", handlers.MpGetMyCoupons)
	mp.Post("/coupons/:id/claim", handlers.MpClaimCoupon)

	// --- Tracking ---
	mp.Get("/tracking/:id", handlers.MpTrackPackage)

	// --- Messages ---
	mp.Get("/messages", handlers.MpGetMessages)
	mp.Put("/messages/:id/read", handlers.MpMarkRead)

	// --- Favorites ---
	mp.Get("/favorites", handlers.MpGetFavorites)
	mp.Post("/favorites", handlers.MpAddFavorite)
	mp.Delete("/favorites/:id", handlers.MpRemoveFavorite)

	// --- Browsing History ---
	mp.Get("/history", handlers.MpGetHistory)
	mp.Delete("/history", handlers.MpClearHistory)

	// --- Feedback ---
	mp.Post("/feedback", handlers.MpSubmitFeedback)

	// --- Reviews (auth) ---
	mp.Post("/reviews", handlers.MpCreateReview)

	// --- Refunds ---
	mp.Get("/refunds", handlers.MpGetRefunds)
	mp.Get("/refunds/:id", handlers.MpGetRefundDetail)
	mp.Post("/refunds", handlers.MpCreateRefund)

	// --- Dealer/Distribution ---
	mp.Get("/dealer/info", handlers.MpGetDealerInfo)
	mp.Post("/dealer/apply", handlers.MpApplyDealer)
	mp.Get("/dealer/orders", handlers.MpGetDealerOrders)
	mp.Get("/dealer/withdrawals", handlers.MpGetDealerWithdrawals)
	mp.Post("/dealer/withdraw", handlers.MpRequestWithdraw)
	mp.Get("/dealer/team", handlers.MpGetDealerTeam)

	// --- Misc (auth) ---
	mp.Get("/invite", handlers.MpGetInviteInfo)

	// ===== AUTHENTICATED ROUTES (no tenant required) =====
	userAPI := app.Group("/api/v1/user", middleware.AuthMiddleware)
	userAPI.Get("/me", handlers.GetCurrentUser)
	userAPI.Post("/setup-tenant", handlers.SetupTenant)

	// ===== AUTHENTICATED + TENANT ROUTES =====
	api := app.Group("/api/v1", middleware.AuthMiddleware, middleware.TenantMiddleware)

	// --- Dashboard ---
	api.Get("/dashboard/stats", handlers.GetDashboardStats)

	// --- Countries ---
	api.Get("/countries", handlers.GetCountries)
	api.Get("/countries/:id", handlers.GetCountry)
	api.Post("/countries", handlers.CreateCountry)
	api.Put("/countries/:id", handlers.UpdateCountry)
	api.Delete("/countries/:id", handlers.DeleteCountry)

	// --- Currencies ---
	api.Get("/currencies", handlers.GetCurrencies)
	api.Get("/currencies/:id", handlers.GetCurrency)
	api.Post("/currencies", handlers.CreateCurrency)
	api.Put("/currencies/:id", handlers.UpdateCurrency)
	api.Delete("/currencies/:id", handlers.DeleteCurrency)

	// --- Member Levels ---
	api.Get("/member-levels", handlers.GetMemberLevels)
	api.Get("/member-levels/:id", handlers.GetMemberLevel)
	api.Post("/member-levels", handlers.CreateMemberLevel)
	api.Put("/member-levels/:id", handlers.UpdateMemberLevel)
	api.Delete("/member-levels/:id", handlers.DeleteMemberLevel)

	// --- Package Categories ---
	api.Get("/package-categories", handlers.GetPackageCategories)
	api.Get("/package-categories/:id", handlers.GetPackageCategory)
	api.Post("/package-categories", handlers.CreatePackageCategory)
	api.Put("/package-categories/:id", handlers.UpdatePackageCategory)
	api.Delete("/package-categories/:id", handlers.DeletePackageCategory)

	// --- Route Categories ---
	api.Get("/route-categories", handlers.GetRouteCategories)
	api.Get("/route-categories/:id", handlers.GetRouteCategory)
	api.Post("/route-categories", handlers.CreateRouteCategory)
	api.Put("/route-categories/:id", handlers.UpdateRouteCategory)
	api.Delete("/route-categories/:id", handlers.DeleteRouteCategory)

	// --- Warehouses ---
	api.Get("/warehouses", handlers.GetWarehouses)
	api.Get("/warehouses/:id", handlers.GetWarehouse)
	api.Post("/warehouses", handlers.CreateWarehouse)
	api.Put("/warehouses/:id", handlers.UpdateWarehouse)
	api.Delete("/warehouses/:id", handlers.DeleteWarehouse)

	// --- Warehouse Shelves ---
	api.Get("/warehouses/:warehouseId/shelves", handlers.GetShelves)
	api.Post("/warehouses/:warehouseId/shelves", handlers.CreateShelf)
	api.Put("/shelves/:id", handlers.UpdateShelf)
	api.Delete("/shelves/:id", handlers.DeleteShelf)

	// --- Shipping Marks ---
	api.Get("/shipping-marks", handlers.GetShippingMarks)
	api.Get("/shipping-marks/:id", handlers.GetShippingMark)
	api.Post("/shipping-marks", handlers.CreateShippingMark)
	api.Put("/shipping-marks/:id", handlers.UpdateShippingMark)
	api.Delete("/shipping-marks/:id", handlers.DeleteShippingMark)

	// --- Packages ---
	api.Get("/packages", handlers.GetPackages)
	api.Get("/packages/:id", handlers.GetPackage)
	api.Post("/packages", handlers.CreatePackage)
	api.Put("/packages/:id", handlers.UpdatePackage)
	api.Delete("/packages/:id", handlers.DeletePackage)
	api.Post("/packages/forecast", handlers.ForecastPackage)
	api.Put("/packages/:id/receive", handlers.ReceivePackage)
	api.Put("/packages/:id/shelve", handlers.ShelvePackage)
	api.Put("/packages/:id/inspect", handlers.InspectPackage)
	api.Put("/packages/:id/ship-out", handlers.ShipOutPackage)

	// --- Shipping Routes ---
	api.Get("/shipping-routes", handlers.GetShippingRoutes)
	api.Get("/shipping-routes/:id", handlers.GetShippingRoute)
	api.Post("/shipping-routes", handlers.CreateShippingRoute)
	api.Put("/shipping-routes/:id", handlers.UpdateShippingRoute)
	api.Delete("/shipping-routes/:id", handlers.DeleteShippingRoute)

	// --- Route Pricing Tiers ---
	api.Get("/shipping-routes/:routeId/pricing", handlers.GetRoutePricingTiers)
	api.Post("/shipping-routes/:routeId/pricing", handlers.CreateRoutePricingTier)
	api.Put("/pricing-tiers/:id", handlers.UpdateRoutePricingTier)
	api.Delete("/pricing-tiers/:id", handlers.DeleteRoutePricingTier)

	// --- Pricing Calculator ---
	api.Post("/shipping-routes/calculate", handlers.CalculateShippingPrice)

	// --- Consolidation Orders ---
	api.Get("/orders", handlers.GetOrders)
	api.Get("/orders/:id", handlers.GetOrder)
	api.Post("/orders", handlers.CreateOrder)
	api.Put("/orders/:id", handlers.UpdateOrder)
	api.Delete("/orders/:id", handlers.DeleteOrder)
	api.Put("/orders/:id/pay", handlers.PayOrder)
	api.Put("/orders/:id/pack", handlers.PackOrder)
	api.Put("/orders/:id/ship", handlers.ShipOrder)
	api.Put("/orders/:id/complete", handlers.CompleteOrder)
	api.Put("/orders/:id/cancel", handlers.CancelOrder)

	// --- Order Packages (link packages to orders) ---
	api.Get("/orders/:orderId/packages", handlers.GetOrderPackages)
	api.Post("/orders/:orderId/packages", handlers.AddOrderPackage)
	api.Delete("/order-packages/:id", handlers.RemoveOrderPackage)

	// --- Payment Audits ---
	api.Get("/payment-audits", handlers.GetPaymentAudits)
	api.Get("/payment-audits/count", handlers.GetPaymentAuditCount)
	api.Get("/payment-audits/:id", handlers.GetPaymentAudit)
	api.Put("/payment-audits/:id/audit", handlers.AuditPaymentAudit)
	api.Put("/payment-audits/:id/approve", handlers.ApprovePaymentAudit)
	api.Put("/payment-audits/:id/reject", handlers.RejectPaymentAudit)

	// --- Shipping Batches ---
	api.Get("/batches", handlers.GetBatches)
	api.Get("/batches/:id", handlers.GetBatch)
	api.Post("/batches", handlers.CreateBatch)
	api.Put("/batches/:id", handlers.UpdateBatch)
	api.Delete("/batches/:id", handlers.DeleteBatch)
	api.Put("/batches/:id/depart", handlers.DepartBatch)
	api.Put("/batches/:id/arrive", handlers.ArriveBatch)

	// --- Batch Orders ---
	api.Get("/batches/:batchId/orders", handlers.GetBatchOrders)
	api.Post("/batches/:batchId/orders", handlers.AddBatchOrder)
	api.Delete("/batch-orders/:id", handlers.RemoveBatchOrder)

	// --- Batch Tracking ---
	api.Get("/batches/:batchId/tracking", handlers.GetBatchTracking)
	api.Post("/batches/:batchId/tracking", handlers.AddBatchTracking)

	// --- Insurance Products ---
	api.Get("/insurance-products", handlers.GetInsuranceProducts)
	api.Get("/insurance-products/:id", handlers.GetInsuranceProduct)
	api.Post("/insurance-products", handlers.CreateInsuranceProduct)
	api.Put("/insurance-products/:id", handlers.UpdateInsuranceProduct)
	api.Delete("/insurance-products/:id", handlers.DeleteInsuranceProduct)

	// --- Consumables ---
	api.Get("/consumables", handlers.GetConsumables)
	api.Get("/consumables/:id", handlers.GetConsumable)
	api.Post("/consumables", handlers.CreateConsumable)
	api.Put("/consumables/:id", handlers.UpdateConsumable)
	api.Delete("/consumables/:id", handlers.DeleteConsumable)

	// --- Value Added Services ---
	api.Get("/value-added-services", handlers.GetValueAddedServices)
	api.Get("/value-added-services/:id", handlers.GetValueAddedService)
	api.Post("/value-added-services", handlers.CreateValueAddedService)
	api.Put("/value-added-services/:id", handlers.UpdateValueAddedService)
	api.Delete("/value-added-services/:id", handlers.DeleteValueAddedService)

	// --- Logistics Companies ---
	api.Get("/logistics-companies", handlers.GetLogisticsCompanies)
	api.Get("/logistics-companies/:id", handlers.GetLogisticsCompany)
	api.Post("/logistics-companies", handlers.CreateLogisticsCompany)
	api.Put("/logistics-companies/:id", handlers.UpdateLogisticsCompany)
	api.Delete("/logistics-companies/:id", handlers.DeleteLogisticsCompany)

	// --- Tracking Templates ---
	api.Get("/tracking-templates", handlers.GetTrackingTemplates)
	api.Get("/tracking-templates/:id", handlers.GetTrackingTemplate)
	api.Post("/tracking-templates", handlers.CreateTrackingTemplate)
	api.Put("/tracking-templates/:id", handlers.UpdateTrackingTemplate)
	api.Delete("/tracking-templates/:id", handlers.DeleteTrackingTemplate)

	// --- Banners ---
	api.Get("/banners", handlers.GetBanners)
	api.Get("/banners/:id", handlers.GetBanner)
	api.Post("/banners", handlers.CreateBanner)
	api.Put("/banners/:id", handlers.UpdateBanner)
	api.Delete("/banners/:id", handlers.DeleteBanner)

	// --- Articles ---
	api.Get("/articles", handlers.GetArticles)
	api.Get("/articles/:id", handlers.GetArticle)
	api.Post("/articles", handlers.CreateArticle)
	api.Put("/articles/:id", handlers.UpdateArticle)
	api.Delete("/articles/:id", handlers.DeleteArticle)

	// --- Coupons ---
	api.Get("/coupons", handlers.GetCoupons)
	api.Get("/coupons/:id", handlers.GetCoupon)
	api.Post("/coupons", handlers.CreateCoupon)
	api.Put("/coupons/:id", handlers.UpdateCoupon)
	api.Delete("/coupons/:id", handlers.DeleteCoupon)

	// --- User Management ---
	api.Get("/users", handlers.GetUsers)
	api.Get("/users/:id", handlers.GetUser)
	api.Post("/users", handlers.CreateUser)
	api.Put("/users/:id", handlers.UpdateUser)
	api.Delete("/users/:id", handlers.DeleteUser)

	// --- Statistics ---
	api.Get("/statistics", handlers.GetStatisticsOverview)
	api.Get("/statistics/overview", handlers.GetStatisticsOverview)
	api.Get("/statistics/routes", handlers.GetRouteStatistics)
	api.Get("/statistics/countries", handlers.GetCountryStatistics)
	api.Get("/statistics/categories", handlers.GetCategoryStatistics)

	// --- Settings ---
	api.Get("/settings", handlers.GetSettings)
	api.Put("/settings", handlers.UpdateSettings)
	api.Get("/settings/:key", handlers.GetSettingByKey)
	api.Put("/settings/:key", handlers.UpdateSettingByKey)
	api.Post("/settings/clear-cache", handlers.ClearCache)

	// --- Members ---
	api.Get("/members", handlers.GetMembers)
	api.Get("/members/:id", handlers.GetMember)
	api.Post("/members", handlers.CreateMember)
	api.Put("/members/:id", handlers.UpdateMember)
	api.Delete("/members/:id", handlers.DeleteMember)

	// --- Address Book ---
	api.Get("/address-book", handlers.GetAddressBooks)
	api.Get("/address-book/:id", handlers.GetAddressBook)
	api.Post("/address-book", handlers.CreateAddressBook)
	api.Put("/address-book/:id", handlers.UpdateAddressBook)
	api.Delete("/address-book/:id", handlers.DeleteAddressBook)

	// --- Notifications ---
	api.Get("/notifications", handlers.GetNotifications)
	api.Get("/notifications/:id", handlers.GetNotification)
	api.Post("/notifications", handlers.CreateNotification)
	api.Put("/notifications/:id", handlers.UpdateNotification)
	api.Delete("/notifications/:id", handlers.DeleteNotification)
	api.Put("/notifications/:id/publish", handlers.PublishNotification)

	// --- Goods Management ---
	api.Get("/goods", handlers.GetGoodsList)
	api.Get("/goods/:id", handlers.GetGoods)
	api.Post("/goods", handlers.CreateGoods)
	api.Put("/goods/:id", handlers.UpdateGoods)
	api.Delete("/goods/:id", handlers.DeleteGoods)
	api.Get("/goods-categories", handlers.GetGoodsCategories)
	api.Get("/goods-categories/:id", handlers.GetGoodsCategory)
	api.Post("/goods-categories", handlers.CreateGoodsCategory)
	api.Put("/goods-categories/:id", handlers.UpdateGoodsCategory)
	api.Delete("/goods-categories/:id", handlers.DeleteGoodsCategory)
	api.Get("/goods-reviews", handlers.GetGoodsReviews)
	api.Get("/goods-reviews/:id", handlers.GetGoodsReview)
	api.Post("/goods-reviews", handlers.CreateGoodsReview)
	api.Put("/goods-reviews/:id", handlers.UpdateGoodsReview)
	api.Delete("/goods-reviews/:id", handlers.DeleteGoodsReview)

	// --- Shop/Mall Orders ---
	api.Get("/shop-orders", handlers.GetShopOrders)
	api.Get("/shop-orders/:id", handlers.GetShopOrder)
	api.Post("/shop-orders", handlers.CreateShopOrder)
	api.Put("/shop-orders/:id", handlers.UpdateShopOrder)
	api.Delete("/shop-orders/:id", handlers.DeleteShopOrder)
	api.Put("/shop-orders/:id/ship", handlers.ShipShopOrder)
	api.Put("/shop-orders/:id/complete", handlers.CompleteShopOrder)
	api.Put("/shop-orders/:id/cancel", handlers.CancelShopOrder)
	api.Get("/order-refunds", handlers.GetOrderRefunds)
	api.Get("/order-refunds/:id", handlers.GetOrderRefund)
	api.Post("/order-refunds", handlers.CreateOrderRefund)
	api.Put("/order-refunds/:id", handlers.UpdateOrderRefund)
	api.Put("/order-refunds/:id/approve", handlers.ApproveRefund)
	api.Put("/order-refunds/:id/reject", handlers.RejectRefund)

	// --- Content Management ---
	api.Get("/article-categories", handlers.GetArticleCategories)
	api.Get("/article-categories/:id", handlers.GetArticleCategory)
	api.Post("/article-categories", handlers.CreateArticleCategory)
	api.Put("/article-categories/:id", handlers.UpdateArticleCategory)
	api.Delete("/article-categories/:id", handlers.DeleteArticleCategory)
	api.Get("/files", handlers.GetFiles)
	api.Get("/files/deleted", handlers.GetDeletedFiles)
	api.Delete("/files", handlers.DeleteFile)
	api.Get("/files/:id", handlers.GetFile)
	api.Post("/files", handlers.CreateFile)
	api.Put("/files/:id", handlers.UpdateFile)
	api.Delete("/files/:id", handlers.DeleteFile)
	api.Put("/files/:id/restore", handlers.RestoreFile)
	api.Get("/file-groups", handlers.GetFileGroups)
	api.Get("/file-groups/:id", handlers.GetFileGroup)
	api.Post("/file-groups", handlers.CreateFileGroup)
	api.Put("/file-groups/:id", handlers.UpdateFileGroup)
	api.Delete("/file-groups/:id", handlers.DeleteFileGroup)

	// --- Marketing / Points / Recharge / Blind Box ---
	api.Get("/points/logs", handlers.GetPointsLogs)
	api.Post("/points/logs", handlers.CreatePointsLog)
	api.Get("/points/setting", handlers.GetPointsSetting)
	api.Put("/points/setting", handlers.UpdatePointsSetting)
	api.Get("/recharge-orders", handlers.GetRechargeOrders)
	api.Get("/recharge-orders/:id", handlers.GetRechargeOrder)
	api.Post("/recharge-orders", handlers.CreateRechargeOrder)
	api.Get("/recharge-plans", handlers.GetRechargePlans)
	api.Get("/recharge-plans/:id", handlers.GetRechargePlan)
	api.Post("/recharge-plans", handlers.CreateRechargePlan)
	api.Put("/recharge-plans/:id", handlers.UpdateRechargePlan)
	api.Delete("/recharge-plans/:id", handlers.DeleteRechargePlan)
	api.Get("/recharge/setting", handlers.GetRechargeSetting)
	api.Put("/recharge/setting", handlers.UpdateRechargeSetting)
	api.Get("/balance/logs", handlers.GetBalanceLogs)
	api.Post("/balance/logs", handlers.CreateBalanceLog)
	api.Get("/coupon-receives", handlers.GetCouponReceiveLogs)
	api.Get("/coupon-receives/:id", handlers.GetCouponReceiveLog)
	api.Post("/coupon-receives", handlers.CreateCouponReceiveLog)
	api.Put("/coupon-receives/:id", handlers.UpdateCouponReceiveLog)
	api.Delete("/coupon-receives/:id", handlers.DeleteCouponReceiveLog)
	api.Get("/blind-box/activities", handlers.GetBlindBoxActivities)
	api.Get("/blind-box/activities/:id", handlers.GetBlindBoxActivity)
	api.Post("/blind-box/activities", handlers.CreateBlindBoxActivity)
	api.Put("/blind-box/activities/:id", handlers.UpdateBlindBoxActivity)
	api.Delete("/blind-box/activities/:id", handlers.DeleteBlindBoxActivity)
	api.Get("/blind-box/draws", handlers.GetBlindBoxDraws)
	api.Post("/blind-box/draws", handlers.CreateBlindBoxDraw)
	api.Put("/blind-box/draws/:id", handlers.UpdateBlindBoxDraw)
	api.Delete("/blind-box/draws/:id", handlers.DeleteBlindBoxDraw)
	api.Get("/blind-box/setting", handlers.GetBlindBoxSetting)
	api.Put("/blind-box/setting", handlers.UpdateBlindBoxSetting)

	// --- User Management Expansion ---
	api.Get("/user-discounts", handlers.GetUserDiscounts)
	api.Get("/user-discounts/:id", handlers.GetUserDiscount)
	api.Post("/user-discounts", handlers.CreateUserDiscount)
	api.Put("/user-discounts/:id", handlers.UpdateUserDiscount)
	api.Delete("/user-discounts/:id", handlers.DeleteUserDiscount)
	api.Get("/user-marks", handlers.GetUserMarks)
	api.Get("/user-marks/:id", handlers.GetUserMark)
	api.Post("/user-marks", handlers.CreateUserMark)
	api.Put("/user-marks/:id", handlers.UpdateUserMark)
	api.Delete("/user-marks/:id", handlers.DeleteUserMark)
	api.Get("/user-birthdays", handlers.GetUserBirthdays)
	api.Get("/user-birthdays/:id", handlers.GetUserBirthday)
	api.Post("/user-birthdays", handlers.CreateUserBirthday)
	api.Put("/user-birthdays/:id", handlers.UpdateUserBirthday)
	api.Delete("/user-birthdays/:id", handlers.DeleteUserBirthday)

	// --- Warehouse Expansion ---
	api.Get("/warehouse-addresses", handlers.GetWarehouseAddresses)
	api.Get("/warehouse-addresses/:id", handlers.GetWarehouseAddress)
	api.Post("/warehouse-addresses", handlers.CreateWarehouseAddress)
	api.Put("/warehouse-addresses/:id", handlers.UpdateWarehouseAddress)
	api.Delete("/warehouse-addresses/:id", handlers.DeleteWarehouseAddress)
	api.Get("/warehouse-applications", handlers.GetWarehouseApplications)
	api.Get("/warehouse-applications/:id", handlers.GetWarehouseApplication)
	api.Put("/warehouse-applications/:id/approve", handlers.ApproveWarehouseApplication)
	api.Put("/warehouse-applications/:id/reject", handlers.RejectWarehouseApplication)
	api.Get("/warehouse-clerks", handlers.GetWarehouseClerks)
	api.Get("/warehouse-clerks/reviews", handlers.GetWarehouseClerkReviews)
	api.Get("/warehouse-clerks/:id", handlers.GetWarehouseClerk)
	api.Post("/warehouse-clerks", handlers.CreateWarehouseClerk)
	api.Put("/warehouse-clerks/:id", handlers.UpdateWarehouseClerk)
	api.Delete("/warehouse-clerks/:id", handlers.DeleteWarehouseClerk)
	api.Get("/warehouse/capital", handlers.GetWarehouseCapitalLogs)
	api.Get("/warehouse/bonuses", handlers.GetWarehouseBonuses)
	api.Get("/warehouse/bonuses/:id", handlers.GetWarehouseBonus)
	api.Post("/warehouse/bonuses", handlers.CreateWarehouseBonus)
	api.Put("/warehouse/bonuses/:id", handlers.UpdateWarehouseBonus)
	api.Put("/warehouse/bonuses/:id/pay", handlers.PayWarehouseBonus)
	api.Delete("/warehouse/bonuses/:id", handlers.DeleteWarehouseBonus)
	api.Get("/warehouse/withdrawals", handlers.GetWarehouseWithdrawals)
	api.Put("/warehouse/withdrawals/:id/approve", handlers.ApproveWarehouseWithdrawal)
	api.Put("/warehouse/withdrawals/:id/reject", handlers.RejectWarehouseWithdrawal)
	api.Put("/warehouse/withdrawals/:id/pay", handlers.PayWarehouseWithdrawal)

	// --- Batch Templates ---
	api.Get("/batch-templates", handlers.GetBatchTemplates)
	api.Get("/batch-templates/:id", handlers.GetBatchTemplate)
	api.Post("/batch-templates", handlers.CreateBatchTemplate)
	api.Put("/batch-templates/:id", handlers.UpdateBatchTemplate)
	api.Delete("/batch-templates/:id", handlers.DeleteBatchTemplate)

	// --- Order Reviews ---
	api.Get("/order-reviews", handlers.GetOrderReviews)
	api.Get("/order-reviews/:id", handlers.GetOrderReview)
	api.Post("/order-reviews", handlers.CreateOrderReview)
	api.Put("/order-reviews/:id/reply", handlers.ReplyOrderReview)
	api.Delete("/order-reviews/:id", handlers.DeleteOrderReview)

	// --- Dealer/Distribution ---
	api.Get("/dealer/applications", handlers.GetDealerApplications)
	api.Get("/dealer/applications/:id", handlers.GetDealerApplication)
	api.Put("/dealer/applications/:id/approve", handlers.ApproveDealerApplication)
	api.Put("/dealer/applications/:id/reject", handlers.RejectDealerApplication)
	api.Get("/dealers", handlers.GetDealers)
	api.Get("/dealers/:id", handlers.GetDealer)
	api.Put("/dealers/:id", handlers.UpdateDealer)
	api.Delete("/dealers/:id", handlers.DeleteDealer)
	api.Get("/dealer/orders", handlers.GetDealerOrders)
	api.Get("/dealer/withdrawals", handlers.GetDealerWithdrawals)
	api.Put("/dealer/withdrawals/:id/approve", handlers.ApproveDealerWithdrawal)
	api.Put("/dealer/withdrawals/:id/reject", handlers.RejectDealerWithdrawal)
	api.Put("/dealer/withdrawals/:id/pay", handlers.PayDealerWithdrawal)
	api.Get("/dealer/levels", handlers.GetDealerLevels)
	api.Get("/dealer/levels/:id", handlers.GetDealerLevel)
	api.Post("/dealer/levels", handlers.CreateDealerLevel)
	api.Put("/dealer/levels/:id", handlers.UpdateDealerLevel)
	api.Delete("/dealer/levels/:id", handlers.DeleteDealerLevel)
	api.Get("/dealer/posters", handlers.GetDealerPosters)
	api.Get("/dealer/posters/:id", handlers.GetDealerPoster)
	api.Post("/dealer/posters", handlers.CreateDealerPoster)
	api.Put("/dealer/posters/:id", handlers.UpdateDealerPoster)
	api.Delete("/dealer/posters/:id", handlers.DeleteDealerPoster)

	// --- Sharing ---
	api.Get("/sharing/orders", handlers.GetSharingOrders)
	api.Get("/sharing/verifications", handlers.GetSharingVerifications)
	api.Put("/sharing/verifications/:id", handlers.UpdateSharingVerification)
	api.Put("/sharing/verifications/:id/verify", handlers.VerifySharing)

	// --- Admin Roles ---
	api.Get("/admin-roles", handlers.GetAdminRoles)
	api.Get("/admin-roles/:id", handlers.GetAdminRole)
	api.Post("/admin-roles", handlers.CreateAdminRole)
	api.Put("/admin-roles/:id", handlers.UpdateAdminRole)
	api.Delete("/admin-roles/:id", handlers.DeleteAdminRole)

	// --- Communication ---
	api.Get("/sms", handlers.GetSmsLogs)
	api.Post("/sms", handlers.SendSms)
	api.Get("/emails", handlers.GetEmailLogs)
	api.Post("/emails", handlers.SendEmail)

	// --- Settings Expansion ---
	api.Get("/barcode-settings", handlers.GetBarcodeSettings)
	api.Get("/barcode-settings/:id", handlers.GetBarcodeSetting)
	api.Post("/barcode-settings", handlers.CreateBarcodeSetting)
	api.Put("/barcode-settings/:id", handlers.UpdateBarcodeSetting)
	api.Delete("/barcode-settings/:id", handlers.DeleteBarcodeSetting)
	api.Get("/bank-accounts", handlers.GetBankAccounts)
	api.Get("/bank-accounts/:id", handlers.GetBankAccount)
	api.Post("/bank-accounts", handlers.CreateBankAccount)
	api.Put("/bank-accounts/:id", handlers.UpdateBankAccount)
	api.Delete("/bank-accounts/:id", handlers.DeleteBankAccount)
	api.Get("/remittance-certificates", handlers.GetRemittanceCertificates)
	api.Get("/remittance-certificates/count", handlers.GetRemittanceCertificateCount)
	api.Get("/remittance-certificates/:id", handlers.GetRemittanceCertificate)
	api.Put("/remittance-certificates/:id/approve", handlers.ApproveRemittanceCertificate)
	api.Put("/remittance-certificates/:id/reject", handlers.RejectRemittanceCertificate)
	api.Get("/payment-flows", handlers.GetPaymentFlows)
	api.Get("/payment-flows/:id", handlers.GetPaymentFlow)
	api.Get("/channels", handlers.GetChannels)
	api.Get("/channels/:id", handlers.GetChannel)
	api.Post("/channels", handlers.CreateChannel)
	api.Put("/channels/:id", handlers.UpdateChannel)
	api.Delete("/channels/:id", handlers.DeleteChannel)
	api.Get("/navigation-items", handlers.GetNavigationItems)
	api.Get("/navigation-items/:id", handlers.GetNavigationItem)
	api.Post("/navigation-items", handlers.CreateNavigationItem)
	api.Put("/navigation-items/:id", handlers.UpdateNavigationItem)
	api.Delete("/navigation-items/:id", handlers.DeleteNavigationItem)
	api.Get("/help-articles", handlers.GetHelpArticles)
	api.Get("/help-articles/:id", handlers.GetHelpArticle)
	api.Post("/help-articles", handlers.CreateHelpArticle)
	api.Put("/help-articles/:id", handlers.UpdateHelpArticle)
	api.Delete("/help-articles/:id", handlers.DeleteHelpArticle)
	api.Get("/subscribe-messages", handlers.GetSubscribeMessages)
	api.Get("/subscribe-messages/:id", handlers.GetSubscribeMessage)
	api.Post("/subscribe-messages", handlers.CreateSubscribeMessage)
	api.Put("/subscribe-messages/:id", handlers.UpdateSubscribeMessage)
	api.Delete("/subscribe-messages/:id", handlers.DeleteSubscribeMessage)
	api.Get("/app-settings", handlers.GetAppSettings)
	api.Get("/app-settings/:id", handlers.GetAppSetting)
	api.Post("/app-settings", handlers.CreateAppSetting)
	api.Put("/app-settings/:id", handlers.UpdateAppSetting)
	api.Get("/page-designs", handlers.GetPageDesigns)
	api.Get("/page-designs/:id", handlers.GetPageDesign)
	api.Post("/page-designs", handlers.CreatePageDesign)
	api.Put("/page-designs/:id", handlers.UpdatePageDesign)
	api.Delete("/page-designs/:id", handlers.DeletePageDesign)
	api.Get("/web-menus", handlers.GetWebMenus)
	api.Get("/web-menus/:id", handlers.GetWebMenu)
	api.Post("/web-menus", handlers.CreateWebMenu)
	api.Put("/web-menus/:id", handlers.UpdateWebMenu)
	api.Delete("/web-menus/:id", handlers.DeleteWebMenu)
	api.Post("/web-menus/sort", handlers.SortWebMenus)
	api.Get("/web-links", handlers.GetWebLinks)
	api.Get("/web-links/:id", handlers.GetWebLink)
	api.Post("/web-links", handlers.CreateWebLink)
	api.Put("/web-links/:id", handlers.UpdateWebLink)
	api.Delete("/web-links/:id", handlers.DeleteWebLink)
	api.Get("/wechat-menus", handlers.GetWechatMenus)
	api.Get("/wechat-menus/:id", handlers.GetWechatMenu)
	api.Post("/wechat-menus", handlers.CreateWechatMenu)
	api.Post("/wechat-menus/publish", handlers.PublishWechatMenu)
	api.Put("/wechat-menus/:id", handlers.UpdateWechatMenu)
	api.Delete("/wechat-menus/:id", handlers.DeleteWechatMenu)
	api.Get("/languages", handlers.GetLanguages)
	api.Get("/languages/:id", handlers.GetLanguage)
	api.Post("/languages", handlers.CreateLanguage)
	api.Put("/languages/:id", handlers.UpdateLanguage)
	api.Delete("/languages/:id", handlers.DeleteLanguage)
	api.Get("/page-categories", handlers.GetPageCategory)
	api.Post("/page-categories", handlers.SavePageCategory)

	// ===== CMS PAGE ROUTES (authenticated) =====
	cms := app.Group("", middleware.AuthMiddleware)

	cms.Get("/dashboard", handlers.RenderDashboard)
	cms.Get("/setup-tenant", handlers.RenderSetupTenant)

	// --- Packages: specific string routes BEFORE parameterized routes ---
	cms.Get("/packages", handlers.RenderPackages)
	cms.Get("/packages/create", handlers.RenderPackageCreate)
	cms.Get("/packages/scan-in", handlers.RenderPackageScanIn)
	cms.Get("/packages/scan-out", handlers.RenderPackageScanOut)
	cms.Get("/packages/forecast", handlers.RenderPackageForecast)
	cms.Get("/packages/unclaimed", handlers.RenderPackageUnclaimed)
	cms.Get("/packages/problems", handlers.RenderPackageProblems)
	cms.Get("/packages/entry-new", handlers.RenderPackageEntryNew)
	cms.Get("/packages/multi-entry", handlers.RenderPackageMultiEntry)
	cms.Get("/packages/admin-forecast", handlers.RenderPackageAdminForecast)
	cms.Get("/packages/claim", handlers.RenderPackageClaim)
	cms.Get("/packages/pending-pack", handlers.RenderPackagePendingPack)
	cms.Get("/packages/appointment", handlers.RenderPackageAppointment)
	cms.Get("/packages/returns", handlers.RenderPackageReturns)
	cms.Get("/packages/:id/edit", handlers.RenderPackageEdit)

	// --- Orders: specific string routes BEFORE parameterized routes ---
	cms.Get("/orders", handlers.RenderOrders)
	cms.Get("/orders/create", handlers.RenderOrderCreate)
	cms.Get("/orders/reviews", handlers.RenderOrderReviews)
	cms.Get("/orders/pending-inspect", handlers.RenderOrdersPendingInspect)
	cms.Get("/orders/pending-ship", handlers.RenderOrdersPendingShip)
	cms.Get("/orders/shipped", handlers.RenderOrdersShipped)
	cms.Get("/orders/arrived", handlers.RenderOrdersArrived)
	cms.Get("/orders/completed", handlers.RenderOrdersCompleted)
	cms.Get("/orders/unpaid", handlers.RenderOrdersUnpaid)
	cms.Get("/orders/problems", handlers.RenderOrdersProblems)
	cms.Get("/orders/overdue", handlers.RenderOrdersOverdue)
	cms.Get("/orders/quick-pack", handlers.RenderOrdersQuickPack)
	cms.Get("/orders/monthly", handlers.RenderOrdersMonthly)
	cms.Get("/orders/cod", handlers.RenderOrdersCOD)
	cms.Get("/orders/:id", handlers.RenderOrderDetail)

	// --- Batches: specific string routes BEFORE parameterized routes ---
	cms.Get("/batches", handlers.RenderBatches)
	cms.Get("/batches/create", handlers.RenderBatchCreate)
	cms.Get("/batches/templates", handlers.RenderBatchTemplates)
	cms.Get("/batches/settings", handlers.RenderBatchSettings)
	cms.Get("/batches/:id", handlers.RenderBatchDetail)

	// --- Users ---
	cms.Get("/users", handlers.RenderUsers)
	cms.Get("/users/addresses", handlers.RenderUserAddresses)
	cms.Get("/users/marks", handlers.RenderUserMarks)
	cms.Get("/users/discounts", handlers.RenderUserDiscounts)
	cms.Get("/users/balance", handlers.RenderBalanceLogs)
	cms.Get("/users/recharge", handlers.RenderRechargeOrdersUser)
	cms.Get("/users/birthdays", handlers.RenderUserBirthdays)
	cms.Get("/users/orders", handlers.RenderUserOrders)
	cms.Get("/users/level-settings", handlers.RenderLevelSettings)

	// --- Warehouses ---
	cms.Get("/warehouses", handlers.RenderWarehouses)
	cms.Get("/warehouses/addresses", handlers.RenderWarehouseAddresses)
	cms.Get("/warehouses/applications", handlers.RenderWarehouseApplications)
	cms.Get("/warehouses/clerks", handlers.RenderWarehouseClerks)
	cms.Get("/warehouses/clerk-reviews", handlers.RenderWarehouseClerkReviews)
	cms.Get("/warehouses/capital", handlers.RenderWarehouseCapital)
	cms.Get("/warehouses/bonuses", handlers.RenderWarehouseBonuses)
	cms.Get("/warehouses/withdrawals", handlers.RenderWarehouseWithdrawals)
	cms.Get("/warehouses/settings", handlers.RenderWarehouseSettings)
	cms.Get("/warehouses/service-bonus", handlers.RenderWarehouseServiceBonus)
	cms.Get("/warehouses/shelf-data", handlers.RenderWarehouseShelfData)
	cms.Get("/warehouses/rack-data", handlers.RenderWarehouseRackData)

	// --- Goods ---
	cms.Get("/goods", handlers.RenderGoods)
	cms.Get("/goods/categories", handlers.RenderGoodsCategories)
	cms.Get("/goods/reviews", handlers.RenderGoodsReviews)

	// --- Shop Orders ---
	cms.Get("/shop-orders", handlers.RenderShopOrders)
	cms.Get("/shop-orders/refunds", handlers.RenderOrderRefunds)
	cms.Get("/shop-orders/pending-pay", handlers.RenderShopOrdersPendingPay)
	cms.Get("/shop-orders/pending-ship", handlers.RenderShopOrdersPendingShip)
	cms.Get("/shop-orders/pending-receive", handlers.RenderShopOrdersPendingReceive)
	cms.Get("/shop-orders/completed", handlers.RenderShopOrdersCompleted)
	cms.Get("/shop-orders/:id", handlers.RenderShopOrderDetail)

	// --- Content ---
	cms.Get("/content/articles", handlers.RenderArticles)
	cms.Get("/content/categories", handlers.RenderArticleCategories)
	cms.Get("/content/files", handlers.RenderFiles)
	cms.Get("/content/file-groups", handlers.RenderFileGroups)
	cms.Get("/content/file-recycle-bin", handlers.RenderFileRecycleBin)

	// --- Marketing ---
	cms.Get("/marketing/coupons", handlers.RenderCoupons)
	cms.Get("/marketing/coupon-receives", handlers.RenderCouponReceives)
	cms.Get("/marketing/points", handlers.RenderPointsLogs)
	cms.Get("/marketing/points/settings", handlers.RenderPointsSettings)
	cms.Get("/marketing/recharge-plans", handlers.RenderRechargePlans)
	cms.Get("/marketing/recharge-orders", handlers.RenderRechargeOrders)
	cms.Get("/marketing/sms", handlers.RenderSmsLogs)
	cms.Get("/marketing/email", handlers.RenderEmailLogs)
	cms.Get("/marketing/blind-box", handlers.RenderBlindBoxActivities)
	cms.Get("/marketing/blind-box/draws", handlers.RenderBlindBoxDraws)
	cms.Get("/marketing/blind-box/settings", handlers.RenderBlindBoxSettings)
	cms.Get("/marketing/blind-box/wall", handlers.RenderBlindBoxWall)
	cms.Get("/marketing/coupon-distribution", handlers.RenderCouponDistribution)
	cms.Get("/marketing/recharge-settings", handlers.RenderRechargeSettings)

	// --- Distribution/Dealers ---
	cms.Get("/distribution/applications", handlers.RenderDealerApplications)
	cms.Get("/distribution/dealers", handlers.RenderDealers)
	cms.Get("/distribution/orders", handlers.RenderDealerOrders)
	cms.Get("/distribution/withdrawals", handlers.RenderDealerWithdrawals)
	cms.Get("/distribution/levels", handlers.RenderDealerLevels)
	cms.Get("/distribution/settings", handlers.RenderDealerSettings)
	cms.Get("/distribution/sharing-orders", handlers.RenderSharingOrders)
	cms.Get("/distribution/sharing-settings", handlers.RenderSharingSettings)
	cms.Get("/distribution/sharing-verifications", handlers.RenderSharingVerifications)
	cms.Get("/distribution/poster", handlers.RenderDealerPoster)
	cms.Get("/distribution/group-leaders", handlers.RenderGroupLeaders)

	// --- Admin ---
	cms.Get("/admin/roles", handlers.RenderAdminRoles)
	cms.Get("/admin/users", handlers.RenderAdminUsers)

	// --- Client/App ---
	cms.Get("/client/settings", handlers.RenderAppSettings)
	cms.Get("/client/miniprogram", handlers.RenderClientMiniprogram)
	cms.Get("/client/h5", handlers.RenderClientH5)
	cms.Get("/client/web", handlers.RenderClientWeb)
	cms.Get("/client/web-menu", handlers.RenderClientWebMenu)
	cms.Get("/client/web-link", handlers.RenderClientWebLink)
	cms.Get("/client/wechat-menu", handlers.RenderClientWechatMenu)
	cms.Get("/client/wechat-reply", handlers.RenderClientWechatReply)
	cms.Get("/client/language", handlers.RenderClientLanguage)
	cms.Get("/client/page-links", handlers.RenderClientPageLinks)
	cms.Get("/client/page-category", handlers.RenderClientPageCategory)
	cms.Get("/client/page-designs", handlers.RenderPageDesigns)
	cms.Get("/client/page-designs/:id/edit", handlers.RenderPageDesignEdit)
	cms.Get("/client/help", handlers.RenderHelpArticles)
	cms.Get("/client/subscribe-messages", handlers.RenderSubscribeMessages)

	// --- Shipping Routes & Statistics ---
	cms.Get("/shipping-routes", handlers.RenderShippingRoutes)
	cms.Get("/statistics", handlers.RenderStatistics)
	cms.Get("/statistics/first-entry", handlers.RenderStatsFirstEntry)
	cms.Get("/statistics/country", handlers.RenderStatsCountry)
	cms.Get("/statistics/category", handlers.RenderStatsCategory)
	cms.Get("/statistics/channel", handlers.RenderStatsChannel)
	cms.Get("/statistics/order", handlers.RenderStatsOrder)
	cms.Get("/statistics/dashboard", handlers.RenderStatsDashboard)

	// --- Settings ---
	cms.Get("/settings", handlers.RenderSettings)
	cms.Get("/settings/countries", handlers.RenderCountries)
	cms.Get("/settings/currencies", handlers.RenderCurrencies)
	cms.Get("/settings/member-levels", handlers.RenderMemberLevels)
	cms.Get("/settings/package-categories", handlers.RenderPackageCategories)
	cms.Get("/settings/route-categories", handlers.RenderRouteCategories)
	cms.Get("/settings/logistics-companies", handlers.RenderLogisticsCompanies)
	cms.Get("/settings/insurance", handlers.RenderInsurance)
	cms.Get("/settings/consumables", handlers.RenderConsumables)
	cms.Get("/settings/value-added-services", handlers.RenderValueAddedServices)
	cms.Get("/settings/tracking-templates", handlers.RenderTrackingTemplates)
	cms.Get("/settings/banners", handlers.RenderBanners)
	cms.Get("/settings/articles", handlers.RenderArticles)
	cms.Get("/settings/coupons", handlers.RenderCoupons)
	cms.Get("/settings/payment-audits", handlers.RenderPaymentAudits)
	cms.Get("/settings/shipping-marks", handlers.RenderShippingMarks)
	cms.Get("/settings/address-book", handlers.RenderAddressBook)
	cms.Get("/settings/notifications", handlers.RenderNotifications)
	cms.Get("/settings/store-settings", handlers.RenderStoreSettings)
	cms.Get("/settings/notice", handlers.RenderNoticeSettings)
	cms.Get("/settings/payment", handlers.RenderPaymentSettings)
	cms.Get("/settings/sms", handlers.RenderSmsSettings)
	cms.Get("/settings/email", handlers.RenderEmailSettings)
	cms.Get("/settings/storage", handlers.RenderStorageSettings)
	cms.Get("/settings/printer", handlers.RenderPrinterSettings)
	cms.Get("/settings/template-messages", handlers.RenderTemplateMessages)
	cms.Get("/settings/admin-style", handlers.RenderAdminStyle)
	cms.Get("/settings/barcode", handlers.RenderBarcodeSettings)
	cms.Get("/settings/bank-accounts", handlers.RenderBankAccounts)
	cms.Get("/settings/remittance-certificates", handlers.RenderRemittanceCertificates)
	cms.Get("/settings/payment-flows", handlers.RenderPaymentFlows)
	cms.Get("/settings/channels", handlers.RenderChannels)
	cms.Get("/settings/navigation", handlers.RenderNavigationItems)
	cms.Get("/settings/menus", handlers.RenderMenuSettings)
	cms.Get("/settings/package-settings", handlers.RenderPackageSettingsPage)
	cms.Get("/settings/cache-clear", handlers.RenderCacheClear)
	cms.Get("/settings/user-client", handlers.RenderUserClientSettings)
	cms.Get("/settings/warehouse-keeper", handlers.RenderWarehouseKeeperSettings)
	cms.Get("/settings/ai-recognition", handlers.RenderAiRecognitionSettings)
	cms.Get("/settings/sms-prefix", handlers.RenderSmsPrefixSettings)

	// --- Members ---
	cms.Get("/members", handlers.RenderMembers)

	// --- Shipping Labels ---
	cms.Get("/shipping-labels", handlers.RenderShippingLabels)

	// --- Tools ---
	cms.Get("/tools", handlers.RenderTools)
	cms.Get("/tools/guide", handlers.RenderToolsGuide)
	cms.Get("/tools/api", handlers.RenderToolsApi)
	cms.Get("/tools/update-log", handlers.RenderToolsUpdateLog)
	cms.Get("/tools/freight-query", handlers.RenderToolsFreightQuery)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("vShip server starting on %s", addr)
	log.Fatal(app.Listen(addr))
}
