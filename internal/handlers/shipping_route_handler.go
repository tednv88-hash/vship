package handlers

import (
	"math"
	"strconv"
	"time"
	"vship/internal/database"
	"vship/internal/middleware"
	"vship/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetShippingRoutes lists shipping routes with pagination
func GetShippingRoutes(c *fiber.Ctx) error {
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
	database.DB.Model(&models.ShippingRoute{}).
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Count(&total)

	var routes []models.ShippingRoute
	if err := database.DB.
		Where("tenant_id = ? AND trashed_at IS NULL", tenantID).
		Order("sort_order ASC").
		Offset(offset).Limit(limit).
		Find(&routes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch shipping routes"})
	}

	return c.JSON(fiber.Map{
		"data":  routes,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetShippingRoute gets a single shipping route by ID
func GetShippingRoute(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid route ID"})
	}

	var route models.ShippingRoute
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&route).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping route not found"})
	}

	return c.JSON(route)
}

// CreateShippingRoute creates a new shipping route
func CreateShippingRoute(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var route models.ShippingRoute
	if err := c.BodyParser(&route); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	route.TenantID = &tenantID

	if err := database.DB.Create(&route).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create shipping route"})
	}

	return c.Status(201).JSON(route)
}

// UpdateShippingRoute updates an existing shipping route
func UpdateShippingRoute(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid route ID"})
	}

	var route models.ShippingRoute
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&route).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping route not found"})
	}

	if err := c.BodyParser(&route); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	route.ID = id
	route.TenantID = &tenantID

	if err := database.DB.Save(&route).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update shipping route"})
	}

	return c.JSON(route)
}

// DeleteShippingRoute soft-deletes a shipping route
func DeleteShippingRoute(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid route ID"})
	}

	var route models.ShippingRoute
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", id, tenantID).
		First(&route).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping route not found"})
	}

	if err := database.DB.Model(&route).Update("trashed_at", time.Now()).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete shipping route"})
	}

	return c.JSON(fiber.Map{"message": "Shipping route deleted successfully"})
}

// GetRoutePricingTiers lists pricing tiers for a route
func GetRoutePricingTiers(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	routeID, err := uuid.Parse(c.Params("routeId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid route ID"})
	}

	var tiers []models.RoutePricingTier
	if err := database.DB.
		Where("tenant_id = ? AND route_id = ?", tenantID, routeID).
		Order("weight_min ASC").
		Find(&tiers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pricing tiers"})
	}

	return c.JSON(fiber.Map{"data": tiers})
}

// CreateRoutePricingTier creates a new pricing tier for a route
func CreateRoutePricingTier(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	routeID, err := uuid.Parse(c.Params("routeId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid route ID"})
	}

	var tier models.RoutePricingTier
	if err := c.BodyParser(&tier); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tier.TenantID = &tenantID
	tier.RouteID = routeID

	if err := database.DB.Create(&tier).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create pricing tier"})
	}

	return c.Status(201).JSON(tier)
}

// UpdateRoutePricingTier updates an existing pricing tier
func UpdateRoutePricingTier(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid pricing tier ID"})
	}

	var tier models.RoutePricingTier
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&tier).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pricing tier not found"})
	}

	if err := c.BodyParser(&tier); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tier.ID = id
	tier.TenantID = &tenantID

	if err := database.DB.Save(&tier).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update pricing tier"})
	}

	return c.JSON(tier)
}

// DeleteRoutePricingTier deletes a pricing tier
func DeleteRoutePricingTier(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid pricing tier ID"})
	}

	var tier models.RoutePricingTier
	if err := database.DB.
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&tier).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pricing tier not found"})
	}

	if err := database.DB.Delete(&tier).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete pricing tier"})
	}

	return c.JSON(fiber.Map{"message": "Pricing tier deleted successfully"})
}

// CalculateShippingPrice calculates shipping price based on route, weight, and dimensions
func CalculateShippingPrice(c *fiber.Ctx) error {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		RouteID       uuid.UUID  `json:"route_id"`
		Weight        float64    `json:"weight"`
		Length        float64    `json:"length"`
		Width         float64    `json:"width"`
		Height        float64    `json:"height"`
		MemberLevelID *uuid.UUID `json:"member_level_id,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Find the route
	var route models.ShippingRoute
	if err := database.DB.
		Where("id = ? AND tenant_id = ? AND trashed_at IS NULL", req.RouteID, tenantID).
		First(&route).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Shipping route not found"})
	}

	// Calculate volume weight
	volumeWeightRatio := route.VolumeWeightRatio
	if volumeWeightRatio <= 0 {
		volumeWeightRatio = 5000
	}
	volumeWeight := (req.Length * req.Width * req.Height) / volumeWeightRatio

	// Chargeable weight = max(actual weight, volume weight)
	chargeableWeight := math.Max(req.Weight, volumeWeight)

	// Find matching pricing tier
	tierQuery := database.DB.
		Where("tenant_id = ? AND route_id = ? AND weight_min <= ? AND (weight_max >= ? OR weight_max = 0)",
			tenantID, req.RouteID, chargeableWeight, chargeableWeight)

	if req.MemberLevelID != nil && *req.MemberLevelID != uuid.Nil {
		tierQuery = tierQuery.Where("member_level_id = ?", *req.MemberLevelID)
	} else {
		tierQuery = tierQuery.Where("member_level_id IS NULL")
	}

	var tier models.RoutePricingTier
	if err := tierQuery.Order("weight_min ASC").First(&tier).Error; err != nil {
		// Fallback: try without member level filter
		if req.MemberLevelID != nil {
			if err := database.DB.
				Where("tenant_id = ? AND route_id = ? AND weight_min <= ? AND (weight_max >= ? OR weight_max = 0) AND member_level_id IS NULL",
					tenantID, req.RouteID, chargeableWeight, chargeableWeight).
				Order("weight_min ASC").
				First(&tier).Error; err != nil {
				return c.Status(404).JSON(fiber.Map{"error": "No matching pricing tier found"})
			}
		} else {
			return c.Status(404).JSON(fiber.Map{"error": "No matching pricing tier found"})
		}
	}

	// Calculate price based on billing mode
	var price float64
	billingMode := route.BillingMode

	switch billingMode {
	case "first_additional":
		price = tier.FirstWeightPrice + math.Max(0, (chargeableWeight-tier.FirstWeight))*tier.AdditionalWeightPrice
	default:
		// weight_interval or default
		price = chargeableWeight * tier.UnitPrice
	}

	return c.JSON(fiber.Map{
		"chargeable_weight": chargeableWeight,
		"volume_weight":     volumeWeight,
		"price":             math.Round(price*100) / 100,
		"billing_mode":      billingMode,
		"route_name":        route.Name,
	})
}
