package abonnements

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"
)

// GetAllAbonnements récupère tous les abonnements
func GetAllAbonnements(c *fiber.Ctx) error {
	db := database.DB

	var data []models.Subscription
	db.Order("created_at DESC").Find(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Abonnements",
		"data":    data,
	})
}

// GetPaginatedAbonnements récupère les abonnements avec pagination
func GetPaginatedAbonnements(c *fiber.Ctx) error {
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
		"message": "All Abonnements",
		"data":    data,
		"pagination": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetAbonnement récupère un abonnement par UUID
func GetAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var data models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&data).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// CreateAbonnement crée un nouvel abonnement
func CreateAbonnement(c *fiber.Ctx) error {
	db := database.DB

	var data models.Subscription
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Validation
	if err := utils.ValidateStruct(data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"error":   err,
		})
	}

	data.UUID = utils.GenerateUUID()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.Status = models.StatusPending

	if err := db.Create(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create abonnement",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Abonnement created successfully",
		"data":    data,
	})
}

// UpdateAbonnement met à jour un abonnement
func UpdateAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var abonnement models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	var updateData models.Subscription
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Update fields
	abonnement.Name = updateData.Name
	abonnement.Email = updateData.Email
	abonnement.Telephone = updateData.Telephone
	abonnement.Planid = updateData.Planid
	abonnement.PlanName = updateData.PlanName
	abonnement.Duration = updateData.Duration
	abonnement.Amount = updateData.Amount
	abonnement.PaymentMethod = updateData.PaymentMethod
	abonnement.Status = updateData.Status
	abonnement.UpdatedAt = time.Now()

	if err := db.Save(&abonnement).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update abonnement",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Abonnement updated successfully",
		"data":    abonnement,
	})
}

// DeleteAbonnement supprime un abonnement
func DeleteAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var abonnement models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	if err := db.Delete(&abonnement).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete abonnement",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Abonnement deleted successfully",
	})
}

// GetCurrentSubscription récupère l'abonnement actuel d'une entreprise
func GetCurrentSubscription(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Query("entreprise_uuid")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "entreprise_uuid is required",
		})
	}

	var subscription models.Subscription
	if err := db.Where("entreprise_uuid = ? AND status = ?", entrepriseUUID, models.StatusActive).
		Order("created_at DESC").
		First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No active subscription found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   subscription,
	})
}

// GetAvailablePlans récupère les plans disponibles
func GetAvailablePlans(c *fiber.Ctx) error {
	// Retourner les plans par défaut
	plans := []models.SubscriptionPlan{
		{
			UUID:        "basic",
			Name:        "Basic",
			Description: "Plan de base",
			Price:       29.99,
			Currency:    "USD",
			Duration:    1,
			MaxUsers:    5,
			MaxPOS:      2,
			StorageGB:   10,
			Features:    []string{"Gestion des stocks", "Ventes", "Support email"},
			Popular:     false,
		},
		{
			UUID:        "professional",
			Name:        "Professional",
			Description: "Plan professionnel",
			Price:       59.99,
			Currency:    "USD",
			Duration:    1,
			MaxUsers:    15,
			MaxPOS:      5,
			StorageGB:   50,
			Features:    []string{"Toutes fonctionnalités Basic", "Rapports avancés", "Support prioritaire"},
			Popular:     true,
		},
		{
			UUID:        "enterprise",
			Name:        "Enterprise",
			Description: "Plan entreprise",
			Price:       99.99,
			Currency:    "USD",
			Duration:    1,
			MaxUsers:    -1,
			MaxPOS:      -1,
			StorageGB:   200,
			Features:    []string{"Toutes fonctionnalités Professional", "Support dédié", "API complète"},
			Popular:     false,
		},
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   plans,
	})
}
