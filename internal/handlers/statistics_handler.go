package handlers

import (
	"fmt"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetStatisticsOverview returns a unified statistics response with trends, distributions, etc.
// Called by frontend as: GET /statistics?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func GetStatisticsOverview(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parse date range
	now := time.Now()
	end := now
	start := now.AddDate(0, 0, -30)

	if startDate != "" && endDate != "" {
		if s, err := time.Parse("2006-01-02", startDate); err == nil {
			start = s
		}
		if e, err := time.Parse("2006-01-02", endDate); err == nil {
			end = e
		}
	}
	endOfDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())

	// --- Period aggregates ---
	var periodOrders int64
	database.DB.Model(&models.ConsolidationOrder{}).
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Count(&periodOrders)

	var periodPackages int64
	database.DB.Model(&models.Package{}).
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Count(&periodPackages)

	var periodUsers int64
	database.DB.Model(&models.User{}).
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Count(&periodUsers)

	type SumResult struct {
		Total float64
	}
	var periodRevenueResult SumResult
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("COALESCE(SUM(total_amount), 0) as total").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Scan(&periodRevenueResult)

	// --- Trend data (daily) ---
	type DailyCount struct {
		Day   string
		Count int64
	}
	type DailySum struct {
		Day   string
		Total float64
	}

	// Calculate number of days in range
	days := int(endOfDay.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}
	if days > 365 {
		days = 365
	}

	trendDates := make([]string, days)
	trendOrders := make([]int64, days)
	trendPackages := make([]int64, days)
	trendRevenue := make([]float64, days)

	// Build date labels
	dateKeyMap := make(map[string]int) // "M/D" -> index
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		key := fmt.Sprintf("%d/%d", int(d.Month()), d.Day())
		trendDates[i] = key
		dateKeyMap[key] = i
	}

	// Daily orders
	var dailyOrders []DailyCount
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("TO_CHAR(created_at, 'FMMM/FMDD') as day, COUNT(*) as count").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Group("TO_CHAR(created_at, 'FMMM/FMDD'), DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&dailyOrders)
	for _, d := range dailyOrders {
		if idx, ok := dateKeyMap[d.Day]; ok {
			trendOrders[idx] = d.Count
		}
	}

	// Daily packages
	var dailyPackages []DailyCount
	database.DB.Model(&models.Package{}).
		Select("TO_CHAR(created_at, 'FMMM/FMDD') as day, COUNT(*) as count").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Group("TO_CHAR(created_at, 'FMMM/FMDD'), DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&dailyPackages)
	for _, d := range dailyPackages {
		if idx, ok := dateKeyMap[d.Day]; ok {
			trendPackages[idx] = d.Count
		}
	}

	// Daily revenue
	var dailyRevenue []DailySum
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("TO_CHAR(created_at, 'FMMM/FMDD') as day, COALESCE(SUM(total_amount), 0) as total").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Group("TO_CHAR(created_at, 'FMMM/FMDD'), DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&dailyRevenue)
	for _, d := range dailyRevenue {
		if idx, ok := dateKeyMap[d.Day]; ok {
			trendRevenue[idx] = d.Total
		}
	}

	// --- Route distribution ---
	type RouteDistItem struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var routeDistribution []RouteDistItem
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("COALESCE(shipping_routes.name, '未指定') as name, COUNT(*) as count").
		Joins("LEFT JOIN shipping_routes ON shipping_routes.id = consolidation_orders.shipping_route_id").
		Where("consolidation_orders.tenant_id = ? AND consolidation_orders.trashed_at IS NULL AND consolidation_orders.created_at >= ? AND consolidation_orders.created_at <= ?", tenantID, start, endOfDay).
		Group("shipping_routes.name").
		Order("count DESC").
		Limit(10).
		Scan(&routeDistribution)
	if routeDistribution == nil {
		routeDistribution = []RouteDistItem{}
	}

	// --- Order status distribution ---
	type StatusDistItem struct {
		Label string `json:"label"`
		Count int64  `json:"count"`
	}

	statusLabelMap := map[string]string{
		"draft":     "草稿",
		"pending":   "待處理",
		"paid":      "已付款",
		"packed":    "已打包",
		"shipped":   "已發貨",
		"completed": "已完成",
		"cancelled": "已取消",
	}

	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ? AND trashed_at IS NULL AND created_at >= ? AND created_at <= ?", tenantID, start, endOfDay).
		Group("status").
		Scan(&statusCounts)

	statusDistribution := make([]StatusDistItem, 0, len(statusCounts))
	for _, sc := range statusCounts {
		label := statusLabelMap[sc.Status]
		if label == "" {
			label = sc.Status
		}
		statusDistribution = append(statusDistribution, StatusDistItem{
			Label: label,
			Count: sc.Count,
		})
	}

	return c.JSON(fiber.Map{
		"period_orders":       periodOrders,
		"period_packages":     periodPackages,
		"period_users":        periodUsers,
		"period_revenue":      periodRevenueResult.Total,
		"trend_dates":         trendDates,
		"trend_orders":        trendOrders,
		"trend_packages":      trendPackages,
		"trend_revenue":       trendRevenue,
		"route_distribution":  routeDistribution,
		"status_distribution": statusDistribution,
	})
}

// GetRouteStatistics returns order statistics grouped by shipping route
func GetRouteStatistics(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	type RouteStats struct {
		ShippingRouteID *string `json:"shipping_route_id"`
		RouteName       string  `json:"route_name"`
		OrderCount      int64   `json:"order_count"`
		TotalAmount     float64 `json:"total_amount"`
		TotalWeight     float64 `json:"total_weight"`
	}

	var results []RouteStats
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("consolidation_orders.shipping_route_id, COALESCE(shipping_routes.name, '未指定') as route_name, COUNT(*) as order_count, COALESCE(SUM(consolidation_orders.total_amount), 0) as total_amount, COALESCE(SUM(consolidation_orders.chargeable_weight), 0) as total_weight").
		Joins("LEFT JOIN shipping_routes ON shipping_routes.id = consolidation_orders.shipping_route_id").
		Where("consolidation_orders.tenant_id = ? AND consolidation_orders.trashed_at IS NULL", tenantID).
		Group("consolidation_orders.shipping_route_id, shipping_routes.name").
		Scan(&results)

	return c.JSON(fiber.Map{
		"data": results,
	})
}

// GetCountryStatistics returns order statistics grouped by recipient country
func GetCountryStatistics(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	type CountryStats struct {
		RecipientCountry string  `json:"recipient_country"`
		OrderCount       int64   `json:"order_count"`
		TotalAmount      float64 `json:"total_amount"`
		TotalWeight      float64 `json:"total_weight"`
	}

	var results []CountryStats
	database.DB.Model(&models.ConsolidationOrder{}).
		Select("recipient_country, COUNT(*) as order_count, COALESCE(SUM(total_amount), 0) as total_amount, COALESCE(SUM(chargeable_weight), 0) as total_weight").
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Group("recipient_country").
		Scan(&results)

	return c.JSON(fiber.Map{
		"data": results,
	})
}

// GetCategoryStatistics returns package statistics grouped by category
func GetCategoryStatistics(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	type CategoryStats struct {
		CategoryID   *string `json:"category_id"`
		CategoryName string  `json:"category_name"`
		PackageCount int64   `json:"package_count"`
		TotalWeight  float64 `json:"total_weight"`
	}

	var results []CategoryStats
	database.DB.Model(&models.Package{}).
		Select("packages.category_id, COALESCE(package_categories.name, '未分類') as category_name, COUNT(*) as package_count, COALESCE(SUM(packages.weight), 0) as total_weight").
		Joins("LEFT JOIN package_categories ON package_categories.id = packages.category_id").
		Where("packages.tenant_id = ? AND packages.trashed_at IS NULL", tenantID).
		Group("packages.category_id, package_categories.name").
		Scan(&results)

	return c.JSON(fiber.Map{
		"data": results,
	})
}
