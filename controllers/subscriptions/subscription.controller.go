package subscriptions

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"
)

// CreateSubscription crée une nouvelle souscription
func CreateSubscription(c *fiber.Ctx) error {
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
	data.Step = 1

	if err := db.Create(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create subscription",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription created successfully",
		"data":    data,
	})
}

// GetSubscriptionByID récupère une souscription par ID
func GetSubscriptionByID(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var data models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&data).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// UpdateSubscriptionPlan met à jour le plan d'abonnement
func UpdateSubscriptionPlan(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	var updateData struct {
		PlanID   string `json:"planid"`
		PlanName string `json:"plan_name"`
		Duration int    `json:"duration"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Calculate amount based on plan
	amount := calculatePlanAmount(updateData.PlanID, updateData.Duration)

	subscription.Planid = updateData.PlanID
	subscription.PlanName = updateData.PlanName
	subscription.Duration = updateData.Duration
	subscription.Amount = amount
	subscription.Step = 2
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update subscription",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Plan updated successfully",
		"data":    subscription,
	})
}

// ProcessPayment traite le paiement
func ProcessPayment(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	var paymentData struct {
		PaymentMethod    string `json:"payment_method"`
		PaymentReference string `json:"payment_reference"`
		TransactionID    string `json:"transaction_id"`
	}

	if err := c.BodyParser(&paymentData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	subscription.PaymentMethod = paymentData.PaymentMethod
	subscription.PaymentReference = paymentData.PaymentReference
	subscription.TransactionID = paymentData.TransactionID
	subscription.PaymentStatus = "pending"
	subscription.PaymentDate = time.Now()
	subscription.Status = models.StatusPaymentPending
	subscription.Step = 3
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process payment",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Payment processed successfully",
		"data":    subscription,
	})
}

// ConfirmPayment confirme le paiement
func ConfirmPayment(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var subscription models.Subscription
	if err := db.Where("uuid = ?", uuid).First(&subscription).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Subscription not found",
		})
	}

	// Simulation de confirmation de paiement
	subscription.PaymentStatus = "completed"
	subscription.Status = models.StatusValidationPending
	subscription.StartDate = time.Now()
	subscription.EndDate = time.Now().AddDate(0, subscription.Duration, 0)
	subscription.UpdatedAt = time.Now()

	if err := db.Save(&subscription).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to confirm payment",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Payment confirmed successfully",
		"data":    subscription,
	})
}

// GetDefaultPlans récupère les plans par défaut
func GetDefaultPlans(c *fiber.Ctx) error {
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

// Fonction utilitaire pour calculer le montant
func calculatePlanAmount(planID string, duration int) float64 {
	basePrice := map[string]float64{
		"basic":        29.99,
		"professional": 59.99,
		"enterprise":   99.99,
	}

	price := basePrice[planID]
	if price == 0 {
		price = 29.99
	}

	return price * float64(duration)
}
