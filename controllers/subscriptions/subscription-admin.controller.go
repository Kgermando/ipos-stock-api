package subscriptions

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"
)

// GetSubscriptionsForAdmin récupère toutes les souscriptions pour l'admin
func GetSubscriptionsForAdmin(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}

	status := c.Query("status")
	search := c.Query("search")

	// Calculate offset
	offset := (page - 1) * limit

	query := db.Model(&models.Subscription{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var data []models.Subscription
	var total int64

	// Get total count
	query.Count(&total)

	// Get paginated results
	query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&data)

	// Calculate pagination info
	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Subscriptions",
		"data":    data,
		"pagination": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetSubscriptionByIDAdmin récupère une souscription par ID (admin)
func GetSubscriptionByIDAdmin(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	// Get history
	var history []models.SubscriptionHistory
	db.Where("subscription_uuid = ?", uuid).Order("action_date DESC").Find(&history)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"subscription": subscription,
			"history":      history,
		},
	})
}

// ApproveSubscription approuve un abonnement
func ApproveSubscription(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	var requestData struct {
		AdminNotes string `json:"admin_notes"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	subscription.Status = models.StatusActive
	subscription.ValidationDate = time.Now()
	subscription.Notes = requestData.AdminNotes
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to approve subscription",
		})
	}

	// Create history record
	history := models.SubscriptionHistory{
		UUID:             utils.GenerateUUID(),
		SubscriptionUUID: subscription.UUID,
		Action:           "approved",
		Amount:           subscription.Amount,
		ActionDate:       time.Now(),
		UserUUID:         "admin", // Should be from context
		Notes:            "Subscription approved",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(&history)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription approved successfully",
		"data":    subscription,
	})
}

// RejectSubscription rejette un abonnement
func RejectSubscription(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	var requestData struct {
		Reason string `json:"reason"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	subscription.Status = models.StatusRejected
	subscription.ValidationDate = time.Now()
	subscription.RejectedReason = requestData.Reason
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to reject subscription",
		})
	}

	// Create history record
	history := models.SubscriptionHistory{
		UUID:             utils.GenerateUUID(),
		SubscriptionUUID: subscription.UUID,
		Action:           "rejected",
		Amount:           subscription.Amount,
		ActionDate:       time.Now(),
		UserUUID:         "admin", // Should be from context
		Notes:            "Subscription rejected: " + requestData.Reason,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(&history)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription rejected successfully",
	})
}

// SuspendSubscription suspend un abonnement
func SuspendSubscription(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	var requestData struct {
		Reason string `json:"reason"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	subscription.Status = models.StatusSuspended
	subscription.Notes = requestData.Reason
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to suspend subscription",
		})
	}

	// Create history record
	history := models.SubscriptionHistory{
		UUID:             utils.GenerateUUID(),
		SubscriptionUUID: subscription.UUID,
		Action:           "suspended",
		Amount:           subscription.Amount,
		ActionDate:       time.Now(),
		UserUUID:         "admin", // Should be from context
		Notes:            "Subscription suspended: " + requestData.Reason,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(&history)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription suspended successfully",
	})
}

// GetSubscriptionStatsAdmin récupère les statistiques pour l'admin
func GetSubscriptionStatsAdmin(c *fiber.Ctx) error {
	db := database.DB

	var stats struct {
		TotalSubscriptions     int64   `json:"total_subscriptions"`
		ActiveSubscriptions    int64   `json:"active_subscriptions"`
		PendingSubscriptions   int64   `json:"pending_subscriptions"`
		ExpiredSubscriptions   int64   `json:"expired_subscriptions"`
		CancelledSubscriptions int64   `json:"cancelled_subscriptions"`
		TotalRevenue           float64 `json:"total_revenue"`
		MonthlyRevenue         float64 `json:"monthly_revenue"`
	}

	// Count subscriptions by status
	db.Model(&models.Subscription{}).Count(&stats.TotalSubscriptions)
	db.Model(&models.Subscription{}).Where("status = ?", models.StatusActive).Count(&stats.ActiveSubscriptions)
	db.Model(&models.Subscription{}).Where("status = ?", models.StatusPending).Count(&stats.PendingSubscriptions)
	db.Model(&models.Subscription{}).Where("status = ?", models.StatusExpired).Count(&stats.ExpiredSubscriptions)
	db.Model(&models.Subscription{}).Where("status = ?", models.StatusCancelled).Count(&stats.CancelledSubscriptions)

	// Calculate revenue
	type RevenueSum struct {
		Total float64
	}
	var totalRevenue RevenueSum
	db.Model(&models.Subscription{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("status IN (?)", []string{string(models.StatusActive), string(models.StatusExpired)}).
		Scan(&totalRevenue)
	stats.TotalRevenue = totalRevenue.Total

	// Monthly revenue
	var monthlyRevenue RevenueSum
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	db.Model(&models.Subscription{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("status = ? AND created_at >= ?", models.StatusActive, startOfMonth).
		Scan(&monthlyRevenue)
	stats.MonthlyRevenue = monthlyRevenue.Total

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

// UpdateSubscriptionAdmin met à jour un abonnement (admin)
func UpdateSubscriptionAdmin(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	var updateData models.Subscription
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Update allowed fields
	subscription.Name = updateData.Name
	subscription.Email = updateData.Email
	subscription.Telephone = updateData.Telephone
	subscription.Status = updateData.Status
	subscription.Notes = updateData.Notes
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update subscription",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription updated successfully",
		"data":    subscription,
	})
}
