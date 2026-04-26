package handlers

import (
	"fmt"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetDashboardStats returns summary statistics for the dashboard
func GetDashboardStats(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	// --- Counts ---
	var totalPackages, totalOrders, totalUsers, totalBatches int64
	database.DB.Model(&models.Package{}).Where("tenant_id = ? AND trashed_at IS NULL", tenantID).Count(&totalPackages)
	database.DB.Model(&models.ConsolidationOrder{}).Where("tenant_id = ? AND trashed_at IS NULL", tenantID).Count(&totalOrders)
	database.DB.Model(&models.User{}).Where("tenant_id = ? AND trashed_at IS NULL", tenantID).Count(&totalUsers)
	database.DB.Model(&models.ShippingBatch{}).Where("tenant_id = ? AND trashed_at IS NULL", tenantID).Count(&totalBatches)

	var pendingOrders int64
	database.DB.Model(&models.ConsolidationOrder{}).Where("tenant_id = ? AND trashed_at IS NULL AND status IN (?,?)", tenantID, "draft", "pending").Count(&pendingOrders)

	var shippingBatches int64
	database.DB.Model(&models.ShippingBatch{}).Where("tenant_id = ? AND trashed_at IS NULL AND status IN (?,?)", tenantID, "in_transit", "shipping").Count(&shippingBatches)

	// --- Today's stats ---
	var todayPackages, todayOrders int64
	database.DB.Model(&models.Package{}).Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, today).Count(&todayPackages)
	database.DB.Model(&models.ConsolidationOrder{}).Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, today).Count(&todayOrders)

	// --- Revenue stats ---
	type SumResult struct {
		Total float64
	}

	var totalRevenue SumResult
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("COALESCE(SUM(total_amount), 0) as total").
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Scan(&totalRevenue)

	var todayRevenue SumResult
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("COALESCE(SUM(total_amount), 0) as total").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, today).
		Scan(&todayRevenue)

	// Week revenue (since Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := today.AddDate(0, 0, -(weekday - 1))
	var weekRevenue SumResult
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("COALESCE(SUM(total_amount), 0) as total").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, weekStart).
		Scan(&weekRevenue)

	// Month revenue (since 1st of month)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var monthRevenue SumResult
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("COALESCE(SUM(total_amount), 0) as total").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, monthStart).
		Scan(&monthRevenue)

	// --- 30-day trend ---
	type DailyCount struct {
		Day   string
		Count int64
	}

	trendDates := make([]string, 30)
	trendPackages := make([]int64, 30)
	trendOrders := make([]int64, 30)

	// Generate date labels
	for i := 29; i >= 0; i-- {
		d := today.AddDate(0, 0, -i)
		trendDates[29-i] = fmt.Sprintf("%d/%d", int(d.Month()), d.Day())
	}

	// Query package counts per day for last 30 days
	thirtyDaysAgo := today.AddDate(0, 0, -29)
	var dailyPackages []DailyCount
	database.DB.Model(&models.Package{}).
		Select("TO_CHAR(created_at, 'FMMM/FMDD') as day, COUNT(*) as count").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, thirtyDaysAgo).
		Group("TO_CHAR(created_at, 'FMMM/FMDD'), DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&dailyPackages)

	// Map to trend array
	packageMap := make(map[string]int64)
	for _, dp := range dailyPackages {
		packageMap[dp.Day] = dp.Count
	}
	for i, dateStr := range trendDates {
		if v, ok := packageMap[dateStr]; ok {
			trendPackages[i] = v
		}
	}

	// Query order counts per day for last 30 days
	var dailyOrders []DailyCount
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("TO_CHAR(created_at, 'FMMM/FMDD') as day, COUNT(*) as count").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ?", tenantID, thirtyDaysAgo).
		Group("TO_CHAR(created_at, 'FMMM/FMDD'), DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&dailyOrders)

	orderMap := make(map[string]int64)
	for _, do := range dailyOrders {
		orderMap[do.Day] = do.Count
	}
	for i, dateStr := range trendDates {
		if v, ok := orderMap[dateStr]; ok {
			trendOrders[i] = v
		}
	}

	// --- Recent orders (last 10) ---
	type RecentOrder struct {
		ID           string    `json:"id"`
		OrderNo      string    `json:"order_no"`
		UserName     string    `json:"user_name"`
		PackageCount int       `json:"package_count"`
		TotalWeight  float64   `json:"total_weight"`
		TotalAmount  float64   `json:"total_amount"`
		Status       string    `json:"status"`
		CreatedAt    time.Time `json:"created_at"`
	}

	var recentOrders []RecentOrder
	database.DB.Model(&models.ConsolidationOrder{}).
		Select(`consolidation_orders.id, consolidation_orders.order_number as order_no,
			COALESCE(users.name, users.email, '') as user_name,
			consolidation_orders.package_count, consolidation_orders.chargeable_weight as total_weight,
			consolidation_orders.total_amount, consolidation_orders.status, consolidation_orders.created_at`).
		Joins("LEFT JOIN users ON users.id = consolidation_orders.user_id").
		Where("consolidation_orders.tenant_id = ? AND consolidation_orders.trashed_at IS NULL", tenantID).
		Order("consolidation_orders.created_at DESC").
		Limit(10).
		Scan(&recentOrders)

	if recentOrders == nil {
		recentOrders = []RecentOrder{}
	}

	return c.JSON(fiber.Map{
		"total_packages":   totalPackages,
		"total_orders":     totalOrders,
		"total_users":      totalUsers,
		"total_batches":    totalBatches,
		"pending_orders":   pendingOrders,
		"shipping_batches": shippingBatches,
		"today_packages":   todayPackages,
		"today_orders":     todayOrders,
		"today_revenue":    todayRevenue.Total,
		"week_revenue":     weekRevenue.Total,
		"month_revenue":    monthRevenue.Total,
		"total_revenue":    totalRevenue.Total,
		"trend_dates":      trendDates,
		"trend_packages":   trendPackages,
		"trend_orders":     trendOrders,
		"recent_orders":    recentOrders,
	})
}
