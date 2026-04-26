package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// Public Pages (no layout)
// ============================================================================

// RenderHome renders the public home/landing page (standalone, no layout)
func RenderHome(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "vShip - 跨境集運管理系統",
	})
}

// RenderLogin renders the login page (standalone, no layout)
func RenderLogin(c *fiber.Ctx) error {
	return c.Render("pages/login", fiber.Map{
		"Title": "vShip - Login",
	})
}

// RenderResetPassword renders the password reset page (standalone, no layout)
func RenderResetPassword(c *fiber.Ctx) error {
	return c.Render("pages/reset_password", fiber.Map{
		"Title": "重設密碼 - vShip",
	})
}

// RenderSetupTenant renders the tenant setup page (standalone, no layout)
func RenderSetupTenant(c *fiber.Ctx) error {
	return c.Render("pages/setup_tenant", fiber.Map{
		"Title": "建立您的公司",
	})
}

// ============================================================================
// Dashboard
// ============================================================================

// RenderDashboard renders the dashboard page
func RenderDashboard(c *fiber.Ctx) error {
	return c.Render("pages/dashboard", fiber.Map{
		"Title":      "總覽",
		"ActiveMenu": "dashboard",
	}, "layouts/cms_layout")
}

// ============================================================================
// Package Management (ActiveMenu: "packages")
// ============================================================================

// RenderPackages renders the package list page
func RenderPackages(c *fiber.Ctx) error {
	return c.Render("pages/packages", fiber.Map{
		"Title":      "包裹管理",
		"ActiveMenu": "packages",
	}, "layouts/cms_layout")
}

// RenderPackageCreate renders the package create form
func RenderPackageCreate(c *fiber.Ctx) error {
	return c.Render("pages/package_create", fiber.Map{
		"Title":      "新增包裹",
		"ActiveMenu": "packages",
	}, "layouts/cms_layout")
}

// RenderPackageEdit renders the package edit form
func RenderPackageEdit(c *fiber.Ctx) error {
	return c.Render("pages/package_edit", fiber.Map{
		"Title":      "編輯包裹",
		"ActiveMenu": "packages",
	}, "layouts/cms_layout")
}

// RenderPackageScanIn renders the scan-in page
func RenderPackageScanIn(c *fiber.Ctx) error {
	return c.Render("pages/package_scan_in", fiber.Map{
		"Title":         "掃碼入庫",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "scan-in",
	}, "layouts/cms_layout")
}

// RenderPackageScanOut renders the scan-out page
func RenderPackageScanOut(c *fiber.Ctx) error {
	return c.Render("pages/package_scan_out", fiber.Map{
		"Title":         "掃碼出庫",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "scan-out",
	}, "layouts/cms_layout")
}

// RenderPackageForecast renders the package forecast page
func RenderPackageForecast(c *fiber.Ctx) error {
	return c.Render("pages/package_forecast", fiber.Map{
		"Title":         "預報包裹",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "forecast",
	}, "layouts/cms_layout")
}

// RenderPackageUnclaimed renders the unclaimed packages page
func RenderPackageUnclaimed(c *fiber.Ctx) error {
	return c.Render("pages/package_unclaimed", fiber.Map{
		"Title":         "待認領",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "unclaimed",
	}, "layouts/cms_layout")
}

// RenderPackageProblems renders the problem packages page
func RenderPackageProblems(c *fiber.Ctx) error {
	return c.Render("pages/package_problems", fiber.Map{
		"Title":         "問題件",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "problems",
	}, "layouts/cms_layout")
}

// RenderPackageEntryNew renders the new backend entry form
func RenderPackageEntryNew(c *fiber.Ctx) error {
	return c.Render("pages/package_entry_new", fiber.Map{
		"Title":         "新後台錄入",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "entry-new",
	}, "layouts/cms_layout")
}

// RenderPackageMultiEntry renders the multi-piece entry form
func RenderPackageMultiEntry(c *fiber.Ctx) error {
	return c.Render("pages/package_multi_entry", fiber.Map{
		"Title":         "一票多件錄入",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "multi-entry",
	}, "layouts/cms_layout")
}

// RenderPackageAdminForecast renders the admin forecast page
func RenderPackageAdminForecast(c *fiber.Ctx) error {
	return c.Render("pages/package_admin_forecast", fiber.Map{
		"Title":         "代客預報",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "admin-forecast",
	}, "layouts/cms_layout")
}

// RenderPackageClaim renders the customer claim page
func RenderPackageClaim(c *fiber.Ctx) error {
	return c.Render("pages/package_claim", fiber.Map{
		"Title":         "客戶認領",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "claim",
	}, "layouts/cms_layout")
}

// RenderPackagePendingPack renders the pending packing page
func RenderPackagePendingPack(c *fiber.Ctx) error {
	return c.Render("pages/package_pending_pack", fiber.Map{
		"Title":         "待打包",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "pending-pack",
	}, "layouts/cms_layout")
}

// RenderPackageAppointment renders the appointment packages page
func RenderPackageAppointment(c *fiber.Ctx) error {
	return c.Render("pages/package_appointment", fiber.Map{
		"Title":         "預約件",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "appointment",
	}, "layouts/cms_layout")
}

// RenderPackageReturns renders the returned packages page
func RenderPackageReturns(c *fiber.Ctx) error {
	return c.Render("pages/package_returns", fiber.Map{
		"Title":         "退貨件",
		"ActiveMenu":    "packages",
		"ActiveSubMenu": "returns",
	}, "layouts/cms_layout")
}

// ============================================================================
// Consolidation Orders (ActiveMenu: "orders")
// ============================================================================

// RenderOrders renders the order list page
func RenderOrders(c *fiber.Ctx) error {
	return c.Render("pages/orders", fiber.Map{
		"Title":      "集運訂單",
		"ActiveMenu": "orders",
	}, "layouts/cms_layout")
}

// RenderOrderCreate renders the order create form
func RenderOrderCreate(c *fiber.Ctx) error {
	return c.Render("pages/order_create", fiber.Map{
		"Title":      "新增訂單",
		"ActiveMenu": "orders",
	}, "layouts/cms_layout")
}

// RenderOrderDetail renders the order detail page
func RenderOrderDetail(c *fiber.Ctx) error {
	return c.Render("pages/order_detail", fiber.Map{
		"Title":      "訂單詳情",
		"ActiveMenu": "orders",
	}, "layouts/cms_layout")
}

// RenderPaymentAudits renders the payment audit management page
func RenderPaymentAudits(c *fiber.Ctx) error {
	return c.Render("pages/payment_audits", fiber.Map{
		"Title":         "付款審核",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "payment-audits",
	}, "layouts/cms_layout")
}

// RenderOrderReviews renders the order reviews page
func RenderOrderReviews(c *fiber.Ctx) error {
	return c.Render("pages/order_reviews", fiber.Map{
		"Title":         "訂單評價",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "order-reviews",
	}, "layouts/cms_layout")
}

// RenderOrdersPendingInspect renders orders pending inspection
func RenderOrdersPendingInspect(c *fiber.Ctx) error {
	return c.Render("pages/orders_pending_inspect", fiber.Map{
		"Title":         "待查驗",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "pending-inspect",
	}, "layouts/cms_layout")
}

// RenderOrdersPendingShip renders orders pending shipment
func RenderOrdersPendingShip(c *fiber.Ctx) error {
	return c.Render("pages/orders_pending_ship", fiber.Map{
		"Title":         "待發貨",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "pending-ship",
	}, "layouts/cms_layout")
}

// RenderOrdersShipped renders shipped orders
func RenderOrdersShipped(c *fiber.Ctx) error {
	return c.Render("pages/orders_shipped", fiber.Map{
		"Title":         "已發貨",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "shipped",
	}, "layouts/cms_layout")
}

// RenderOrdersArrived renders arrived orders
func RenderOrdersArrived(c *fiber.Ctx) error {
	return c.Render("pages/orders_arrived", fiber.Map{
		"Title":         "已到貨",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "arrived",
	}, "layouts/cms_layout")
}

// RenderOrdersCompleted renders completed orders
func RenderOrdersCompleted(c *fiber.Ctx) error {
	return c.Render("pages/orders_completed", fiber.Map{
		"Title":         "已完成",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "completed",
	}, "layouts/cms_layout")
}

// RenderOrdersUnpaid renders unpaid orders
func RenderOrdersUnpaid(c *fiber.Ctx) error {
	return c.Render("pages/orders_unpaid", fiber.Map{
		"Title":         "未支付",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "unpaid",
	}, "layouts/cms_layout")
}

// RenderOrdersProblems renders problem orders
func RenderOrdersProblems(c *fiber.Ctx) error {
	return c.Render("pages/orders_problems", fiber.Map{
		"Title":         "問題件",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "problems",
	}, "layouts/cms_layout")
}

// RenderOrdersOverdue renders overdue orders
func RenderOrdersOverdue(c *fiber.Ctx) error {
	return c.Render("pages/orders_overdue", fiber.Map{
		"Title":         "超時件",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "overdue",
	}, "layouts/cms_layout")
}

// RenderOrdersQuickPack renders the quick pack page
func RenderOrdersQuickPack(c *fiber.Ctx) error {
	return c.Render("pages/orders_quick_pack", fiber.Map{
		"Title":         "快速打包",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "quick-pack",
	}, "layouts/cms_layout")
}

// RenderOrdersMonthly renders monthly settlement orders
func RenderOrdersMonthly(c *fiber.Ctx) error {
	return c.Render("pages/orders_monthly", fiber.Map{
		"Title":         "月結訂單",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "monthly",
	}, "layouts/cms_layout")
}

// RenderOrdersCOD renders cash-on-delivery orders
func RenderOrdersCOD(c *fiber.Ctx) error {
	return c.Render("pages/orders_cod", fiber.Map{
		"Title":         "貨到付款",
		"ActiveMenu":    "orders",
		"ActiveSubMenu": "cod",
	}, "layouts/cms_layout")
}

// ============================================================================
// Goods Management (ActiveMenu: "goods")
// ============================================================================

// RenderGoods renders the goods list page
func RenderGoods(c *fiber.Ctx) error {
	return c.Render("pages/goods", fiber.Map{
		"Title":      "商品列表",
		"ActiveMenu": "goods",
	}, "layouts/cms_layout")
}

// RenderGoodsCategories renders the goods categories page
func RenderGoodsCategories(c *fiber.Ctx) error {
	return c.Render("pages/goods_categories", fiber.Map{
		"Title":         "商品分類",
		"ActiveMenu":    "goods",
		"ActiveSubMenu": "goods-categories",
	}, "layouts/cms_layout")
}

// RenderGoodsReviews renders the goods reviews page
func RenderGoodsReviews(c *fiber.Ctx) error {
	return c.Render("pages/goods_reviews", fiber.Map{
		"Title":         "商品評價",
		"ActiveMenu":    "goods",
		"ActiveSubMenu": "goods-reviews",
	}, "layouts/cms_layout")
}

// ============================================================================
// Shop Orders (ActiveMenu: "shop-orders")
// ============================================================================

// RenderShopOrders renders the shop orders list page
func RenderShopOrders(c *fiber.Ctx) error {
	return c.Render("pages/shop_orders", fiber.Map{
		"Title":      "商城訂單",
		"ActiveMenu": "shop-orders",
	}, "layouts/cms_layout")
}

// RenderShopOrderDetail renders the shop order detail page
func RenderShopOrderDetail(c *fiber.Ctx) error {
	return c.Render("pages/shop_order_detail", fiber.Map{
		"Title":      "商城訂單詳情",
		"ActiveMenu": "shop-orders",
	}, "layouts/cms_layout")
}

// RenderOrderRefunds renders the order refunds / after-sales page
func RenderOrderRefunds(c *fiber.Ctx) error {
	return c.Render("pages/order_refunds", fiber.Map{
		"Title":         "售後管理",
		"ActiveMenu":    "shop-orders",
		"ActiveSubMenu": "refunds",
	}, "layouts/cms_layout")
}

// RenderShopOrdersPendingPay renders shop orders pending payment
func RenderShopOrdersPendingPay(c *fiber.Ctx) error {
	return c.Render("pages/shop_orders_pending_pay", fiber.Map{
		"Title":         "待付款",
		"ActiveMenu":    "shop-orders",
		"ActiveSubMenu": "pending-pay",
	}, "layouts/cms_layout")
}

// RenderShopOrdersPendingShip renders shop orders pending shipment
func RenderShopOrdersPendingShip(c *fiber.Ctx) error {
	return c.Render("pages/shop_orders_pending_ship", fiber.Map{
		"Title":         "待發貨",
		"ActiveMenu":    "shop-orders",
		"ActiveSubMenu": "pending-ship",
	}, "layouts/cms_layout")
}

// RenderShopOrdersPendingReceive renders shop orders pending receipt
func RenderShopOrdersPendingReceive(c *fiber.Ctx) error {
	return c.Render("pages/shop_orders_pending_receive", fiber.Map{
		"Title":         "待收貨",
		"ActiveMenu":    "shop-orders",
		"ActiveSubMenu": "pending-receive",
	}, "layouts/cms_layout")
}

// RenderShopOrdersCompleted renders completed/cancelled shop orders
func RenderShopOrdersCompleted(c *fiber.Ctx) error {
	return c.Render("pages/shop_orders_completed", fiber.Map{
		"Title":         "已完成",
		"ActiveMenu":    "shop-orders",
		"ActiveSubMenu": "shop-completed",
	}, "layouts/cms_layout")
}

// ============================================================================
// Batch Management (ActiveMenu: "batches")
// ============================================================================

// RenderBatches renders the batch list page
func RenderBatches(c *fiber.Ctx) error {
	return c.Render("pages/batches", fiber.Map{
		"Title":      "批次管理",
		"ActiveMenu": "batches",
	}, "layouts/cms_layout")
}

// RenderBatchCreate renders the batch create form
func RenderBatchCreate(c *fiber.Ctx) error {
	return c.Render("pages/batch_create", fiber.Map{
		"Title":      "新增批次",
		"ActiveMenu": "batches",
	}, "layouts/cms_layout")
}

// RenderBatchDetail renders the batch detail page
func RenderBatchDetail(c *fiber.Ctx) error {
	return c.Render("pages/batch_detail", fiber.Map{
		"Title":      "批次詳情",
		"ActiveMenu": "batches",
	}, "layouts/cms_layout")
}

// RenderBatchTemplates renders the batch templates page
func RenderBatchTemplates(c *fiber.Ctx) error {
	return c.Render("pages/batch_templates", fiber.Map{
		"Title":         "批次模板",
		"ActiveMenu":    "batches",
		"ActiveSubMenu": "batch-templates",
	}, "layouts/cms_layout")
}

// RenderBatchSettings renders the batch settings page
func RenderBatchSettings(c *fiber.Ctx) error {
	return c.Render("pages/batch_settings", fiber.Map{
		"Title":         "批次設置",
		"ActiveMenu":    "batches",
		"ActiveSubMenu": "batch-settings",
	}, "layouts/cms_layout")
}

// ============================================================================
// User Management (ActiveMenu: "users")
// ============================================================================

// RenderUsers renders the user list page
func RenderUsers(c *fiber.Ctx) error {
	return c.Render("pages/users", fiber.Map{
		"Title":      "用戶管理",
		"ActiveMenu": "users",
	}, "layouts/cms_layout")
}

// RenderMembers renders the member management page
func RenderMembers(c *fiber.Ctx) error {
	return c.Render("pages/members", fiber.Map{
		"Title":         "會員管理",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "members",
	}, "layouts/cms_layout")
}

// RenderUserAddresses renders the user addresses page
func RenderUserAddresses(c *fiber.Ctx) error {
	return c.Render("pages/user_addresses", fiber.Map{
		"Title":         "地址管理",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "user-addresses",
	}, "layouts/cms_layout")
}

// RenderUserMarks renders the user marks page
func RenderUserMarks(c *fiber.Ctx) error {
	return c.Render("pages/user_marks", fiber.Map{
		"Title":         "標記列表",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "user-marks",
	}, "layouts/cms_layout")
}

// RenderUserDiscounts renders the user discounts page
func RenderUserDiscounts(c *fiber.Ctx) error {
	return c.Render("pages/user_discounts", fiber.Map{
		"Title":         "折扣列表",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "user-discounts",
	}, "layouts/cms_layout")
}

// RenderBalanceLogs renders the balance logs page
func RenderBalanceLogs(c *fiber.Ctx) error {
	return c.Render("pages/balance_logs", fiber.Map{
		"Title":         "餘額明細",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "balance-logs",
	}, "layouts/cms_layout")
}

// RenderRechargeOrdersUser renders the recharge orders page under user management
func RenderRechargeOrdersUser(c *fiber.Ctx) error {
	return c.Render("pages/recharge_orders", fiber.Map{
		"Title":         "充值訂單",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "recharge-orders-user",
	}, "layouts/cms_layout")
}

// RenderUserBirthdays renders the user birthdays page
func RenderUserBirthdays(c *fiber.Ctx) error {
	return c.Render("pages/user_birthdays", fiber.Map{
		"Title":         "生日管理",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "user-birthdays",
	}, "layouts/cms_layout")
}

// RenderUserOrders renders the member orders page
func RenderUserOrders(c *fiber.Ctx) error {
	return c.Render("pages/user_orders", fiber.Map{
		"Title":         "會員訂單",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "user-orders",
	}, "layouts/cms_layout")
}

// RenderLevelSettings renders the level settings page
func RenderLevelSettings(c *fiber.Ctx) error {
	return c.Render("pages/level_settings", fiber.Map{
		"Title":         "等級設置",
		"ActiveMenu":    "users",
		"ActiveSubMenu": "level-settings",
	}, "layouts/cms_layout")
}

// ============================================================================
// Warehouse Management (ActiveMenu: "warehouses")
// ============================================================================

// RenderWarehouses renders the warehouse list page
func RenderWarehouses(c *fiber.Ctx) error {
	return c.Render("pages/warehouses", fiber.Map{
		"Title":      "倉庫管理",
		"ActiveMenu": "warehouses",
	}, "layouts/cms_layout")
}

// RenderWarehouseAddresses renders the warehouse addresses page
func RenderWarehouseAddresses(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_addresses", fiber.Map{
		"Title":         "倉庫地址",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-addresses",
	}, "layouts/cms_layout")
}

// RenderWarehouseApplications renders the warehouse applications page
func RenderWarehouseApplications(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_applications", fiber.Map{
		"Title":         "入駐申請",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-applications",
	}, "layouts/cms_layout")
}

// RenderWarehouseClerks renders the warehouse clerks page
func RenderWarehouseClerks(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_clerks", fiber.Map{
		"Title":         "倉管員",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-clerks",
	}, "layouts/cms_layout")
}

// RenderWarehouseClerkReviews renders the warehouse clerk reviews page
func RenderWarehouseClerkReviews(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_clerk_reviews", fiber.Map{
		"Title":         "倉管評價",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-clerk-reviews",
	}, "layouts/cms_layout")
}

// RenderWarehouseCapital renders the warehouse capital details page
func RenderWarehouseCapital(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_capital", fiber.Map{
		"Title":         "資金明細",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-capital",
	}, "layouts/cms_layout")
}

// RenderWarehouseBonuses renders the warehouse bonuses page
func RenderWarehouseBonuses(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_bonuses", fiber.Map{
		"Title":         "獎金管理",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-bonuses",
	}, "layouts/cms_layout")
}

// RenderWarehouseWithdrawals renders the warehouse withdrawals page
func RenderWarehouseWithdrawals(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_withdrawals", fiber.Map{
		"Title":         "提現管理",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-withdrawals",
	}, "layouts/cms_layout")
}

// RenderWarehouseSettings renders the warehouse settings page
func RenderWarehouseSettings(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_settings", fiber.Map{
		"Title":         "倉庫設置",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "warehouse-settings",
	}, "layouts/cms_layout")
}

// RenderWarehouseServiceBonus renders the warehouse service bonus page
func RenderWarehouseServiceBonus(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_service_bonus", fiber.Map{
		"Title":         "服務分成",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "service-bonus",
	}, "layouts/cms_layout")
}

// RenderWarehouseShelfData renders the warehouse shelf data page
func RenderWarehouseShelfData(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_shelf_data", fiber.Map{
		"Title":         "貨位數據",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "shelf-data",
	}, "layouts/cms_layout")
}

// RenderWarehouseRackData renders the warehouse rack data page
func RenderWarehouseRackData(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_rack_data", fiber.Map{
		"Title":         "貨架數據",
		"ActiveMenu":    "warehouses",
		"ActiveSubMenu": "rack-data",
	}, "layouts/cms_layout")
}

// ============================================================================
// Shipping Routes (ActiveMenu: "shipping-routes")
// ============================================================================

// RenderShippingRoutes renders the shipping route list page
func RenderShippingRoutes(c *fiber.Ctx) error {
	return c.Render("pages/shipping_routes", fiber.Map{
		"Title":      "線路管理",
		"ActiveMenu": "shipping-routes",
	}, "layouts/cms_layout")
}

// ============================================================================
// Shipping Labels (ActiveMenu: "shipping-labels")
// ============================================================================

// RenderShippingLabels renders the shipping label management page
func RenderShippingLabels(c *fiber.Ctx) error {
	return c.Render("pages/shipping_labels", fiber.Map{
		"Title":      "運單標籤",
		"ActiveMenu": "shipping-labels",
	}, "layouts/cms_layout")
}

// ============================================================================
// Statistics (ActiveMenu: "statistics")
// ============================================================================

// RenderStatistics renders the statistics page
func RenderStatistics(c *fiber.Ctx) error {
	return c.Render("pages/statistics", fiber.Map{
		"Title":         "數據統計",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "overview",
	}, "layouts/cms_layout")
}

// RenderStatsFirstEntry renders the first entry statistics page
func RenderStatsFirstEntry(c *fiber.Ctx) error {
	return c.Render("pages/stats_first_entry", fiber.Map{
		"Title":         "首次入庫",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "first-entry",
	}, "layouts/cms_layout")
}

// RenderStatsCountry renders the country statistics page
func RenderStatsCountry(c *fiber.Ctx) error {
	return c.Render("pages/stats_country", fiber.Map{
		"Title":         "國家統計",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "country",
	}, "layouts/cms_layout")
}

// RenderStatsCategory renders the category statistics page
func RenderStatsCategory(c *fiber.Ctx) error {
	return c.Render("pages/stats_category", fiber.Map{
		"Title":         "類目統計",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "category",
	}, "layouts/cms_layout")
}

// RenderStatsChannel renders the channel statistics page
func RenderStatsChannel(c *fiber.Ctx) error {
	return c.Render("pages/stats_channel", fiber.Map{
		"Title":         "渠道統計",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "channel",
	}, "layouts/cms_layout")
}

// RenderStatsOrder renders the order statistics page
func RenderStatsOrder(c *fiber.Ctx) error {
	return c.Render("pages/stats_order", fiber.Map{
		"Title":         "訂單統計",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "order",
	}, "layouts/cms_layout")
}

// RenderStatsDashboard renders the data dashboard page
func RenderStatsDashboard(c *fiber.Ctx) error {
	return c.Render("pages/stats_dashboard", fiber.Map{
		"Title":         "數據大屏",
		"ActiveMenu":    "statistics",
		"ActiveSubMenu": "dashboard",
	}, "layouts/cms_layout")
}

// ============================================================================
// Content Management (ActiveMenu: "content")
// ============================================================================

// RenderArticles renders the article management page
func RenderArticles(c *fiber.Ctx) error {
	return c.Render("pages/articles", fiber.Map{
		"Title":         "文章管理",
		"ActiveMenu":    "content",
		"ActiveSubMenu": "articles",
	}, "layouts/cms_layout")
}

// RenderArticleCategories renders the article categories page
func RenderArticleCategories(c *fiber.Ctx) error {
	return c.Render("pages/article_categories", fiber.Map{
		"Title":         "文章分類",
		"ActiveMenu":    "content",
		"ActiveSubMenu": "article-categories",
	}, "layouts/cms_layout")
}

// RenderFiles renders the file management page
func RenderFiles(c *fiber.Ctx) error {
	return c.Render("pages/files", fiber.Map{
		"Title":         "文件管理",
		"ActiveMenu":    "content",
		"ActiveSubMenu": "files",
	}, "layouts/cms_layout")
}

// RenderFileGroups renders the file groups page
func RenderFileGroups(c *fiber.Ctx) error {
	return c.Render("pages/file_groups", fiber.Map{
		"Title":         "文件分組",
		"ActiveMenu":    "content",
		"ActiveSubMenu": "file-groups",
	}, "layouts/cms_layout")
}

// RenderFileRecycleBin renders the file recycle bin page
func RenderFileRecycleBin(c *fiber.Ctx) error {
	return c.Render("pages/file_recycle_bin", fiber.Map{
		"Title":         "文件回收站",
		"ActiveMenu":    "content",
		"ActiveSubMenu": "file-recycle-bin",
	}, "layouts/cms_layout")
}

// ============================================================================
// Marketing (ActiveMenu: "marketing")
// ============================================================================

// RenderCoupons renders the coupon management page
func RenderCoupons(c *fiber.Ctx) error {
	return c.Render("pages/coupons", fiber.Map{
		"Title":         "優惠券",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "coupons",
	}, "layouts/cms_layout")
}

// RenderCouponReceives renders the coupon receive records page
func RenderCouponReceives(c *fiber.Ctx) error {
	return c.Render("pages/coupon_receives", fiber.Map{
		"Title":         "領取記錄",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "coupon-receives",
	}, "layouts/cms_layout")
}

// RenderPointsLogs renders the points logs page
func RenderPointsLogs(c *fiber.Ctx) error {
	return c.Render("pages/points_logs", fiber.Map{
		"Title":         "積分明細",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "points-logs",
	}, "layouts/cms_layout")
}

// RenderPointsSettings renders the points settings page
func RenderPointsSettings(c *fiber.Ctx) error {
	return c.Render("pages/points_settings", fiber.Map{
		"Title":         "積分設置",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "points-settings",
	}, "layouts/cms_layout")
}

// RenderRechargePlans renders the recharge plans page
func RenderRechargePlans(c *fiber.Ctx) error {
	return c.Render("pages/recharge_plans", fiber.Map{
		"Title":         "充值方案",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "recharge-plans",
	}, "layouts/cms_layout")
}

// RenderRechargeOrders renders the recharge orders page under marketing
func RenderRechargeOrders(c *fiber.Ctx) error {
	return c.Render("pages/recharge_orders", fiber.Map{
		"Title":         "充值訂單",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "recharge-orders",
	}, "layouts/cms_layout")
}

// RenderSmsLogs renders the SMS push logs page
func RenderSmsLogs(c *fiber.Ctx) error {
	return c.Render("pages/sms_logs", fiber.Map{
		"Title":         "短信推送",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "sms-push",
	}, "layouts/cms_layout")
}

// RenderEmailLogs renders the email push logs page
func RenderEmailLogs(c *fiber.Ctx) error {
	return c.Render("pages/email_logs", fiber.Map{
		"Title":         "郵件推送",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "email-push",
	}, "layouts/cms_layout")
}

// RenderBlindBoxActivities renders the blind box activities page
func RenderBlindBoxActivities(c *fiber.Ctx) error {
	return c.Render("pages/blind_box_activities", fiber.Map{
		"Title":         "盲盒活動",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "blind-box",
	}, "layouts/cms_layout")
}

// RenderBlindBoxDraws renders the blind box draw records page
func RenderBlindBoxDraws(c *fiber.Ctx) error {
	return c.Render("pages/blind_box_draws", fiber.Map{
		"Title":         "盲盒記錄",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "blind-box-draws",
	}, "layouts/cms_layout")
}

// RenderBlindBoxSettings renders the blind box settings page
func RenderBlindBoxSettings(c *fiber.Ctx) error {
	return c.Render("pages/blind_box_settings", fiber.Map{
		"Title":         "盲盒設置",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "blind-box-settings",
	}, "layouts/cms_layout")
}

// RenderCouponDistribution renders the coupon distribution settings page
func RenderCouponDistribution(c *fiber.Ctx) error {
	return c.Render("pages/coupon_distribution", fiber.Map{
		"Title":         "優惠券發放設置",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "coupon-distribution",
	}, "layouts/cms_layout")
}

// RenderRechargeSettings renders the recharge settings page
func RenderRechargeSettings(c *fiber.Ctx) error {
	return c.Render("pages/recharge_settings", fiber.Map{
		"Title":         "充值設置",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "recharge-settings",
	}, "layouts/cms_layout")
}

// RenderBlindBoxWall renders the blind box sharing wall page
func RenderBlindBoxWall(c *fiber.Ctx) error {
	return c.Render("pages/blind_box_wall", fiber.Map{
		"Title":         "盲盒分享牆",
		"ActiveMenu":    "marketing",
		"ActiveSubMenu": "blind-box-wall",
	}, "layouts/cms_layout")
}

// ============================================================================
// Client / App Settings (ActiveMenu: "client")
// ============================================================================

// RenderAppSettings renders the client app settings page (legacy, redirects to miniprogram)
func RenderAppSettings(c *fiber.Ctx) error {
	return c.Redirect("/client/miniprogram")
}

// RenderClientMiniprogram renders the miniprogram settings page
func RenderClientMiniprogram(c *fiber.Ctx) error {
	return c.Render("pages/client_miniprogram", fiber.Map{
		"Title":         "小程序設置",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "miniprogram",
	}, "layouts/cms_layout")
}

// RenderClientH5 renders the H5 settings page
func RenderClientH5(c *fiber.Ctx) error {
	return c.Render("pages/client_h5", fiber.Map{
		"Title":         "H5端設置",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "h5",
	}, "layouts/cms_layout")
}

// RenderClientWeb renders the PC web management page
func RenderClientWeb(c *fiber.Ctx) error {
	return c.Render("pages/client_web", fiber.Map{
		"Title":         "PC端管理",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "web",
	}, "layouts/cms_layout")
}

// RenderClientWebMenu renders the website menu management page
func RenderClientWebMenu(c *fiber.Ctx) error {
	return c.Render("pages/client_webmenu", fiber.Map{
		"Title":         "網站菜單",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "webmenu",
	}, "layouts/cms_layout")
}

// RenderClientWebLink renders the friendly links management page
func RenderClientWebLink(c *fiber.Ctx) error {
	return c.Render("pages/client_weblink", fiber.Map{
		"Title":         "友情鏈接",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "weblink",
	}, "layouts/cms_layout")
}

// RenderClientWechatMenu renders the WeChat menu management page
func RenderClientWechatMenu(c *fiber.Ctx) error {
	return c.Render("pages/client_wechat_menu", fiber.Map{
		"Title":         "菜單管理",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "wechat-menu",
	}, "layouts/cms_layout")
}

// RenderClientWechatReply renders the WeChat auto-reply settings page
func RenderClientWechatReply(c *fiber.Ctx) error {
	return c.Render("pages/client_wechat_reply", fiber.Map{
		"Title":         "自動回覆",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "wechat-reply",
	}, "layouts/cms_layout")
}

// RenderClientLanguage renders the language settings page
func RenderClientLanguage(c *fiber.Ctx) error {
	return c.Render("pages/client_language", fiber.Map{
		"Title":         "語言設置",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "language",
	}, "layouts/cms_layout")
}

// RenderClientPageLinks renders the page links reference page
func RenderClientPageLinks(c *fiber.Ctx) error {
	return c.Render("pages/client_page_links", fiber.Map{
		"Title":         "頁面鏈接",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "page-links",
	}, "layouts/cms_layout")
}

// RenderClientPageCategory renders the category template settings page
func RenderClientPageCategory(c *fiber.Ctx) error {
	return c.Render("pages/client_page_category", fiber.Map{
		"Title":         "分類模板",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "page-category",
	}, "layouts/cms_layout")
}

// RenderPageDesigns renders the page designs page
func RenderPageDesigns(c *fiber.Ctx) error {
	return c.Render("pages/page_designs", fiber.Map{
		"Title":         "頁面設計",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "page-designs",
	}, "layouts/cms_layout")
}

// RenderPageDesignEdit renders the page design editor page
func RenderPageDesignEdit(c *fiber.Ctx) error {
	return c.Render("pages/page_design_edit", fiber.Map{
		"Title":         "頁面設計編輯",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "page-designs",
		"PageID":        c.Params("id"),
	}, "layouts/cms_layout")
}

// RenderPageDesignPreview renders the public H5 preview of a page design
func RenderPageDesignPreview(c *fiber.Ctx) error {
	return c.Render("pages/page_design_preview", fiber.Map{
		"Title":  "頁面預覽",
		"PageID": c.Params("id"),
	})
}

// RenderHelpArticles renders the help center page
func RenderHelpArticles(c *fiber.Ctx) error {
	return c.Render("pages/help_articles", fiber.Map{
		"Title":         "幫助中心",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "help-articles",
	}, "layouts/cms_layout")
}

// RenderSubscribeMessages renders the subscribe messages page
func RenderSubscribeMessages(c *fiber.Ctx) error {
	return c.Render("pages/subscribe_messages", fiber.Map{
		"Title":         "訂閱消息",
		"ActiveMenu":    "client",
		"ActiveSubMenu": "subscribe-messages",
	}, "layouts/cms_layout")
}

// ============================================================================
// Distribution / Dealers (ActiveMenu: "distribution")
// ============================================================================

// RenderDealerApplications renders the dealer applications page
func RenderDealerApplications(c *fiber.Ctx) error {
	return c.Render("pages/dealer_applications", fiber.Map{
		"Title":      "分銷商申請",
		"ActiveMenu": "distribution",
	}, "layouts/cms_layout")
}

// RenderDealers renders the dealers list page
func RenderDealers(c *fiber.Ctx) error {
	return c.Render("pages/dealers", fiber.Map{
		"Title":         "分銷商列表",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "dealer-list",
	}, "layouts/cms_layout")
}

// RenderDealerOrders renders the dealer orders page
func RenderDealerOrders(c *fiber.Ctx) error {
	return c.Render("pages/dealer_orders", fiber.Map{
		"Title":         "分銷訂單",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "dealer-orders",
	}, "layouts/cms_layout")
}

// RenderDealerWithdrawals renders the dealer withdrawals page
func RenderDealerWithdrawals(c *fiber.Ctx) error {
	return c.Render("pages/dealer_withdrawals", fiber.Map{
		"Title":         "分銷提現",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "dealer-withdrawals",
	}, "layouts/cms_layout")
}

// RenderDealerLevels renders the dealer levels page
func RenderDealerLevels(c *fiber.Ctx) error {
	return c.Render("pages/dealer_levels", fiber.Map{
		"Title":         "分銷等級",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "dealer-levels",
	}, "layouts/cms_layout")
}

// RenderDealerSettings renders the dealer settings page
func RenderDealerSettings(c *fiber.Ctx) error {
	return c.Render("pages/dealer_settings", fiber.Map{
		"Title":         "分銷設置",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "dealer-settings",
	}, "layouts/cms_layout")
}

// RenderSharingOrders renders the sharing orders page
func RenderSharingOrders(c *fiber.Ctx) error {
	return c.Render("pages/sharing_orders", fiber.Map{
		"Title":         "分享訂單",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "sharing-orders",
	}, "layouts/cms_layout")
}

// RenderSharingSettings renders the sharing settings page
func RenderSharingSettings(c *fiber.Ctx) error {
	return c.Render("pages/sharing_settings", fiber.Map{
		"Title":         "分享設置",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "sharing-settings",
	}, "layouts/cms_layout")
}

// RenderSharingVerifications renders the sharing verifications page
func RenderSharingVerifications(c *fiber.Ctx) error {
	return c.Render("pages/sharing_verifications", fiber.Map{
		"Title":         "分享核銷",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "sharing-verifications",
	}, "layouts/cms_layout")
}

// RenderDealerPoster renders the dealer poster page
func RenderDealerPoster(c *fiber.Ctx) error {
	return c.Render("pages/dealer_poster", fiber.Map{
		"Title":         "分銷海報",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "dealer-poster",
	}, "layouts/cms_layout")
}

// RenderGroupLeaders renders the group leaders page
func RenderGroupLeaders(c *fiber.Ctx) error {
	return c.Render("pages/group_leaders", fiber.Map{
		"Title":         "團長列表",
		"ActiveMenu":    "distribution",
		"ActiveSubMenu": "group-leaders",
	}, "layouts/cms_layout")
}

// ============================================================================
// Admin (ActiveMenu: "admin")
// ============================================================================

// RenderAdminRoles renders the admin roles page
func RenderAdminRoles(c *fiber.Ctx) error {
	return c.Render("pages/admin_roles", fiber.Map{
		"Title":         "角色管理",
		"ActiveMenu":    "admin",
		"ActiveSubMenu": "admin-roles",
	}, "layouts/cms_layout")
}

// RenderAdminUsers renders the admin users list page
func RenderAdminUsers(c *fiber.Ctx) error {
	return c.Render("pages/admin_users", fiber.Map{
		"Title":         "管理員列表",
		"ActiveMenu":    "admin",
		"ActiveSubMenu": "admin-users",
	}, "layouts/cms_layout")
}

// ============================================================================
// Settings (ActiveMenu: "settings")
// ============================================================================

// RenderSettings renders the settings overview page
func RenderSettings(c *fiber.Ctx) error {
	return c.Render("pages/settings", fiber.Map{
		"Title":      "系統設置",
		"ActiveMenu": "settings",
	}, "layouts/cms_layout")
}

// RenderCountries renders the country management page
func RenderCountries(c *fiber.Ctx) error {
	return c.Render("pages/countries", fiber.Map{
		"Title":         "國家管理",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "countries",
	}, "layouts/cms_layout")
}

// RenderCurrencies renders the currency management page
func RenderCurrencies(c *fiber.Ctx) error {
	return c.Render("pages/currencies", fiber.Map{
		"Title":         "幣種管理",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "currencies",
	}, "layouts/cms_layout")
}

// RenderMemberLevels renders the member level management page
func RenderMemberLevels(c *fiber.Ctx) error {
	return c.Render("pages/member_levels", fiber.Map{
		"Title":         "會員等級",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "member-levels",
	}, "layouts/cms_layout")
}

// RenderPackageCategories renders the package category management page
func RenderPackageCategories(c *fiber.Ctx) error {
	return c.Render("pages/package_categories", fiber.Map{
		"Title":         "包裹分類",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "package-categories",
	}, "layouts/cms_layout")
}

// RenderRouteCategories renders the route category management page
func RenderRouteCategories(c *fiber.Ctx) error {
	return c.Render("pages/route_categories", fiber.Map{
		"Title":         "線路分類",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "route-categories",
	}, "layouts/cms_layout")
}

// RenderLogisticsCompanies renders the logistics company management page
func RenderLogisticsCompanies(c *fiber.Ctx) error {
	return c.Render("pages/logistics_companies", fiber.Map{
		"Title":         "物流公司",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "logistics-companies",
	}, "layouts/cms_layout")
}

// RenderInsurance renders the insurance product management page
func RenderInsurance(c *fiber.Ctx) error {
	return c.Render("pages/insurance", fiber.Map{
		"Title":         "保險產品",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "insurance",
	}, "layouts/cms_layout")
}

// RenderConsumables renders the consumable management page
func RenderConsumables(c *fiber.Ctx) error {
	return c.Render("pages/consumables", fiber.Map{
		"Title":         "耗材管理",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "consumables",
	}, "layouts/cms_layout")
}

// RenderValueAddedServices renders the value-added service management page
func RenderValueAddedServices(c *fiber.Ctx) error {
	return c.Render("pages/value_added_services", fiber.Map{
		"Title":         "增值服務",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "value-added-services",
	}, "layouts/cms_layout")
}

// RenderTrackingTemplates renders the tracking template management page
func RenderTrackingTemplates(c *fiber.Ctx) error {
	return c.Render("pages/tracking_templates", fiber.Map{
		"Title":         "軌跡模板",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "tracking-templates",
	}, "layouts/cms_layout")
}

// RenderBanners renders the banner management page
func RenderBanners(c *fiber.Ctx) error {
	return c.Render("pages/banners", fiber.Map{
		"Title":         "輪播圖",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "banners",
	}, "layouts/cms_layout")
}

// RenderShippingMarks renders the shipping mark management page
func RenderShippingMarks(c *fiber.Ctx) error {
	return c.Render("pages/shipping_marks", fiber.Map{
		"Title":         "貨運標記",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "shipping-marks",
	}, "layouts/cms_layout")
}

// RenderAddressBook renders the address book management page
func RenderAddressBook(c *fiber.Ctx) error {
	return c.Render("pages/address_book", fiber.Map{
		"Title":         "地址簿",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "address-book",
	}, "layouts/cms_layout")
}

// RenderNotifications renders the notification management page
func RenderNotifications(c *fiber.Ctx) error {
	return c.Render("pages/notifications", fiber.Map{
		"Title":         "通知管理",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "notifications",
	}, "layouts/cms_layout")
}

// RenderStoreSettings renders the store settings page
func RenderStoreSettings(c *fiber.Ctx) error {
	return c.Render("pages/store_settings", fiber.Map{
		"Title":         "商城設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "store-settings",
	}, "layouts/cms_layout")
}

// RenderNoticeSettings renders the notice settings page
func RenderNoticeSettings(c *fiber.Ctx) error {
	return c.Render("pages/notice_settings", fiber.Map{
		"Title":         "通知設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "notice-settings",
	}, "layouts/cms_layout")
}

// RenderPaymentSettings renders the payment settings page
func RenderPaymentSettings(c *fiber.Ctx) error {
	return c.Render("pages/payment_settings", fiber.Map{
		"Title":         "支付設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "payment-settings",
	}, "layouts/cms_layout")
}

// RenderSmsSettings renders the SMS settings page
func RenderSmsSettings(c *fiber.Ctx) error {
	return c.Render("pages/sms_settings", fiber.Map{
		"Title":         "短信設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "sms-settings",
	}, "layouts/cms_layout")
}

// RenderEmailSettings renders the email settings page
func RenderEmailSettings(c *fiber.Ctx) error {
	return c.Render("pages/email_settings", fiber.Map{
		"Title":         "郵件設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "email-settings",
	}, "layouts/cms_layout")
}

// RenderStorageSettings renders the storage settings page
func RenderStorageSettings(c *fiber.Ctx) error {
	return c.Render("pages/storage_settings", fiber.Map{
		"Title":         "存儲設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "storage-settings",
	}, "layouts/cms_layout")
}

// RenderPrinterSettings renders the printer settings page
func RenderPrinterSettings(c *fiber.Ctx) error {
	return c.Render("pages/printer_settings", fiber.Map{
		"Title":         "打印機設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "printer-settings",
	}, "layouts/cms_layout")
}

// RenderTemplateMessages renders the template messages page
func RenderTemplateMessages(c *fiber.Ctx) error {
	return c.Render("pages/template_messages", fiber.Map{
		"Title":         "模板消息",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "template-messages",
	}, "layouts/cms_layout")
}

// RenderAdminStyle renders the admin style page
func RenderAdminStyle(c *fiber.Ctx) error {
	return c.Render("pages/admin_style", fiber.Map{
		"Title":         "後台樣式",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "admin-style",
	}, "layouts/cms_layout")
}

// RenderBarcodeSettings renders the barcode settings page
func RenderBarcodeSettings(c *fiber.Ctx) error {
	return c.Render("pages/barcode_settings", fiber.Map{
		"Title":         "條碼設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "barcode-settings",
	}, "layouts/cms_layout")
}

// RenderBankAccounts renders the bank accounts page
func RenderBankAccounts(c *fiber.Ctx) error {
	return c.Render("pages/bank_accounts", fiber.Map{
		"Title":         "銀行設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "bank-accounts",
	}, "layouts/cms_layout")
}

// RenderRemittanceCertificates renders the remittance certificates page
func RenderRemittanceCertificates(c *fiber.Ctx) error {
	return c.Render("pages/remittance_certificates", fiber.Map{
		"Title":         "匯款憑證審核",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "remittance-certificates",
	}, "layouts/cms_layout")
}

// RenderPaymentFlows renders the payment flows page
func RenderPaymentFlows(c *fiber.Ctx) error {
	return c.Render("pages/payment_flows", fiber.Map{
		"Title":         "支付流水",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "payment-flows",
	}, "layouts/cms_layout")
}

// RenderChannels renders the channels page
func RenderChannels(c *fiber.Ctx) error {
	return c.Render("pages/channels", fiber.Map{
		"Title":         "渠道設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "channels",
	}, "layouts/cms_layout")
}

// RenderNavigationItems renders the navigation items page
func RenderNavigationItems(c *fiber.Ctx) error {
	return c.Render("pages/navigation_items", fiber.Map{
		"Title":         "導航管理",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "navigation-items",
	}, "layouts/cms_layout")
}

// RenderMenuSettings renders the menu settings page
func RenderMenuSettings(c *fiber.Ctx) error {
	return c.Render("pages/menu_settings", fiber.Map{
		"Title":         "菜單管理",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "menu-settings",
	}, "layouts/cms_layout")
}

// RenderPackageSettingsPage renders the package settings page under settings
func RenderPackageSettingsPage(c *fiber.Ctx) error {
	return c.Render("pages/package_settings", fiber.Map{
		"Title":         "包裹設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "package-settings",
	}, "layouts/cms_layout")
}

// RenderCacheClear renders the cache clear page
func RenderCacheClear(c *fiber.Ctx) error {
	return c.Render("pages/cache_clear", fiber.Map{
		"Title":         "清除緩存",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "cache-clear",
	}, "layouts/cms_layout")
}

// RenderUserClientSettings renders the user client settings page
func RenderUserClientSettings(c *fiber.Ctx) error {
	return c.Render("pages/user_client_settings", fiber.Map{
		"Title":         "用戶端設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "user-client",
	}, "layouts/cms_layout")
}

// RenderWarehouseKeeperSettings renders the warehouse keeper settings page
func RenderWarehouseKeeperSettings(c *fiber.Ctx) error {
	return c.Render("pages/warehouse_keeper_settings", fiber.Map{
		"Title":         "倉管端設置",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "warehouse-keeper",
	}, "layouts/cms_layout")
}

// RenderAiRecognitionSettings renders the AI recognition settings page
func RenderAiRecognitionSettings(c *fiber.Ctx) error {
	return c.Render("pages/ai_recognition_settings", fiber.Map{
		"Title":         "智能AI識別",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "ai-recognition",
	}, "layouts/cms_layout")
}

// RenderSmsPrefixSettings renders the SMS prefix settings page
func RenderSmsPrefixSettings(c *fiber.Ctx) error {
	return c.Render("pages/sms_prefix_settings", fiber.Map{
		"Title":         "短信前綴",
		"ActiveMenu":    "settings",
		"ActiveSubMenu": "sms-prefix",
	}, "layouts/cms_layout")
}

// ============================================================================
// Tools (ActiveMenu: "tools")
// ============================================================================

// RenderTools renders the tools page
func RenderTools(c *fiber.Ctx) error {
	return c.Render("pages/tools", fiber.Map{
		"Title":      "工具",
		"ActiveMenu": "tools",
	}, "layouts/cms_layout")
}

// RenderToolsGuide renders the tools guide page
func RenderToolsGuide(c *fiber.Ctx) error {
	return c.Render("pages/tools_guide", fiber.Map{
		"Title":         "使用指南",
		"ActiveMenu":    "tools",
		"ActiveSubMenu": "tools-guide",
	}, "layouts/cms_layout")
}

// RenderToolsApi renders the tools API page
func RenderToolsApi(c *fiber.Ctx) error {
	return c.Render("pages/tools_api", fiber.Map{
		"Title":         "API接口",
		"ActiveMenu":    "tools",
		"ActiveSubMenu": "tools-api",
	}, "layouts/cms_layout")
}

// RenderToolsUpdateLog renders the tools update log page
func RenderToolsUpdateLog(c *fiber.Ctx) error {
	return c.Render("pages/tools_update_log", fiber.Map{
		"Title":         "更新日誌",
		"ActiveMenu":    "tools",
		"ActiveSubMenu": "tools-update-log",
	}, "layouts/cms_layout")
}

// RenderToolsFreightQuery renders the freight query tool page
func RenderToolsFreightQuery(c *fiber.Ctx) error {
	return c.Render("pages/tools_freight_query", fiber.Map{
		"Title":         "運費查詢",
		"ActiveMenu":    "tools",
		"ActiveSubMenu": "tools-freight-query",
	}, "layouts/cms_layout")
}
