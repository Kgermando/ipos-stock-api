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

	var abonnements []models.Abonnement
	if err := db.Preload("Entreprise").Order("created_at DESC").Find(&abonnements).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get abonnements",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Abonnements",
		"data":    abonnements,
	})
}

// GetPaginatedAbonnements récupère les abonnements avec pagination et filtres
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
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var abonnements []models.Abonnement
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Abonnement{}).
		Joins("LEFT JOIN entreprises ON abonnements.entreprise_uuid = entreprises.uuid").
		Where("entreprises.name ILIKE ? OR abonnements.moyen_payment ILIKE ? OR abonnements.statut ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Joins("LEFT JOIN entreprises ON abonnements.entreprise_uuid = entreprises.uuid").
		Where("entreprises.name ILIKE ? OR abonnements.moyen_payment ILIKE ? OR abonnements.statut ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("abonnements.created_at DESC").
		Preload("Entreprise").
		Find(&abonnements).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch abonnements",
			"error":   err.Error(),
		})
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	// Prepare pagination metadata
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All abonnements",
		"data":       abonnements,
		"pagination": pagination,
	})
}

// GetPaginatedAbonnements récupère les abonnements avec pagination et filtres
func GetPaginatedAbonnementsEntreprise(c *fiber.Ctx) error {
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

	// Si limit = 0, retourner tous les résultats sans pagination
	if c.Query("limit") == "0" {
		limit = 0
	}

	statut := c.Query("statut")
	entrepriseUUID := c.Query("entreprise_uuid")

	query := db.Model(&models.Abonnement{}).Preload("Entreprise")

	// Apply filters
	if statut != "" {
		query = query.Where("statut = ?", statut)
	}
	if entrepriseUUID != "" {
		query = query.Where("entreprise_uuid = ?", entrepriseUUID)
	}

	var abonnements []models.Abonnement
	var total int64

	// Get total count
	query.Count(&total)

	// Apply pagination only if limit > 0
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	// Get results
	query.Order("created_at DESC").Find(&abonnements)

	// Calculate pagination info
	response := fiber.Map{
		"status":  "success",
		"message": "Abonnements retrieved",
		"data":    abonnements,
	}

	if limit > 0 {
		totalPages := (int(total) + limit - 1) / limit
		response["pagination"] = fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		}
	} else {
		response["total"] = total
	}

	return c.JSON(response)
}

// GetAbonnement récupère un abonnement par UUID
func GetAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var abonnement models.Abonnement
	if err := db.Preload("Entreprise").Where("uuid = ?", uuid).First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   abonnement,
	})
}

// CreateAbonnement crée un nouvel abonnement
func CreateAbonnement(c *fiber.Ctx) error {
	db := database.DB

	var abonnement models.Abonnement
	if err := c.BodyParser(&abonnement); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Validation des champs requis
	if abonnement.EntrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "EntrepriseUUID is required",
		})
	}

	if abonnement.Montant <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Montant must be greater than 0",
		})
	}

	if abonnement.Duree <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Duree must be greater than 0",
		})
	}

	if abonnement.MoyenPayment == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "MoyenPayment is required",
		})
	}

	// Validation
	if err := utils.ValidateStruct(abonnement); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"error":   err,
		})
	}

	// Vérifier si l'entreprise existe
	var entreprise models.Entreprise
	if err := db.Where("uuid = ?", abonnement.EntrepriseUUID).First(&entreprise).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Entreprise not found",
		})
	}

	abonnement.UUID = utils.GenerateUUID()
	abonnement.CreatedAt = time.Now()
	abonnement.UpdatedAt = time.Now()

	// Définir le statut par défaut si non spécifié
	if abonnement.Statut == "" {
		abonnement.Statut = "pending"
	}

	if err := db.Create(&abonnement).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create abonnement",
		})
	}

	// Recharger avec les relations
	db.Preload("Entreprise").Where("uuid = ?", abonnement.UUID).First(&abonnement)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Abonnement created successfully",
		"data":    abonnement,
	})
}

// UpdateAbonnement met à jour un abonnement
func UpdateAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var abonnement models.Abonnement
	if err := db.Where("uuid = ?", uuid).First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	var updateData models.Abonnement
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Update fields
	if updateData.EntrepriseUUID != "" {
		// Vérifier si l'entreprise existe
		var entreprise models.Entreprise
		if err := db.Where("uuid = ?", updateData.EntrepriseUUID).First(&entreprise).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Entreprise not found",
			})
		}
		abonnement.EntrepriseUUID = updateData.EntrepriseUUID
	}

	if updateData.Montant > 0 {
		abonnement.Montant = updateData.Montant
	}

	if updateData.MoyenPayment != "" {
		abonnement.MoyenPayment = updateData.MoyenPayment
	}

	if updateData.Duree > 0 {
		abonnement.Duree = updateData.Duree
	}

	if updateData.Statut != "" {
		abonnement.Statut = updateData.Statut
	}

	if updateData.Signature != "" {
		abonnement.Signature = updateData.Signature
	}

	abonnement.UpdatedAt = time.Now()

	if err := db.Save(&abonnement).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update abonnement",
		})
	}

	// Recharger avec les relations
	db.Preload("Entreprise").Where("uuid = ?", abonnement.UUID).First(&abonnement)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Abonnement updated successfully",
		"data":    abonnement,
	})
}

// DeleteAbonnement supprime un abonnement (soft delete)
func DeleteAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var abonnement models.Abonnement
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

// UpdateStatutAbonnement met à jour le statut d'un abonnement
func UpdateStatutAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var request struct {
		Statut string `json:"statut"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	if request.Statut == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "statut is required",
		})
	}

	var abonnement models.Abonnement
	if err := db.Where("uuid = ?", uuid).First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	abonnement.Statut = request.Statut
	abonnement.UpdatedAt = time.Now()

	if err := db.Save(&abonnement).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update statut",
		})
	}

	// Recharger avec les relations
	db.Preload("Entreprise").Where("uuid = ?", abonnement.UUID).First(&abonnement)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Statut updated successfully",
		"data":    abonnement,
	})
}

// GetAbonnementActuel récupère l'abonnement actuel valide d'une entreprise
// basé sur la date de création, la durée et le statut
func GetAbonnementActuel(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Query("entreprise_uuid")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "entreprise_uuid is required",
		})
	}

	var abonnement models.Abonnement

	// Rechercher l'abonnement actif le plus récent
	if err := db.Preload("Entreprise").
		Where("entreprise_uuid = ? AND statut = ?", entrepriseUUID, "active").
		Order("created_at DESC").
		First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No active subscription found",
		})
	}

	// Calculer la date d'expiration
	dateExpiration := abonnement.CreatedAt.AddDate(0, abonnement.Duree, 0)

	// Vérifier si l'abonnement est encore valide
	maintenant := time.Now()
	estValide := maintenant.Before(dateExpiration)

	// Calculer les jours restants
	joursRestants := int(dateExpiration.Sub(maintenant).Hours() / 24)
	if joursRestants < 0 {
		joursRestants = 0
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"abonnement":      abonnement,
			"date_expiration": dateExpiration,
			"est_valide":      estValide,
			"jours_restants":  joursRestants,
			"statut":          abonnement.Statut,
		},
	})
}

// VerifierValiditeAbonnement vérifie si un abonnement est encore valide
func VerifierValiditeAbonnement(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	var abonnement models.Abonnement
	if err := db.Preload("Entreprise").Where("uuid = ?", uuid).First(&abonnement).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Abonnement not found",
		})
	}

	// Calculer la date d'expiration
	dateExpiration := abonnement.CreatedAt.AddDate(0, abonnement.Duree, 0)

	// Vérifier si l'abonnement est encore valide
	maintenant := time.Now()
	estValide := maintenant.Before(dateExpiration) && abonnement.Statut == "active"

	// Calculer les jours restants
	joursRestants := int(dateExpiration.Sub(maintenant).Hours() / 24)
	if joursRestants < 0 {
		joursRestants = 0
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"abonnement":      abonnement,
			"date_expiration": dateExpiration,
			"est_valide":      estValide,
			"jours_restants":  joursRestants,
			"statut":          abonnement.Statut,
		},
	})
}

// GetAbonnementsExpirant récupère les abonnements qui vont expirer dans X jours
func GetAbonnementsExpirant(c *fiber.Ctx) error {
	db := database.DB

	// Paramètre pour le nombre de jours (par défaut 30 jours)
	jours, err := strconv.Atoi(c.Query("jours", "30"))
	if err != nil || jours < 0 {
		jours = 30
	}

	var abonnements []models.Abonnement
	if err := db.Preload("Entreprise").
		Where("statut = ?", "active").
		Find(&abonnements).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get abonnements",
		})
	}

	var abonnementsExpirant []fiber.Map
	maintenant := time.Now()
	dateLimite := maintenant.AddDate(0, 0, jours)

	for _, abonnement := range abonnements {
		dateExpiration := abonnement.CreatedAt.AddDate(0, abonnement.Duree, 0)

		// Vérifier si l'abonnement expire dans les X jours
		if dateExpiration.After(maintenant) && dateExpiration.Before(dateLimite) {
			joursRestants := int(dateExpiration.Sub(maintenant).Hours() / 24)

			abonnementsExpirant = append(abonnementsExpirant, fiber.Map{
				"abonnement":      abonnement,
				"date_expiration": dateExpiration,
				"jours_restants":  joursRestants,
			})
		}
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   abonnementsExpirant,
		"count":  len(abonnementsExpirant),
	})
}

// GetStatistiquesAbonnements récupère les statistiques des abonnements
func GetStatistiquesAbonnements(c *fiber.Ctx) error {
	db := database.DB

	var stats struct {
		TotalAbonnements   int64   `json:"total_abonnements"`
		AbonnementsActifs  int64   `json:"abonnements_actifs"`
		AbonnementsPending int64   `json:"abonnements_pending"`
		AbonnementsSuspend int64   `json:"abonnements_suspended"`
		AbonnementsCancel  int64   `json:"abonnements_cancelled"`
		RevenuTotal        float64 `json:"revenu_total"`
		RevenuMensuel      float64 `json:"revenu_mensuel"`
	}

	// Total abonnements
	db.Model(&models.Abonnement{}).Count(&stats.TotalAbonnements)

	// Abonnements par statut
	db.Model(&models.Abonnement{}).Where("statut = ?", "active").Count(&stats.AbonnementsActifs)
	db.Model(&models.Abonnement{}).Where("statut = ?", "pending").Count(&stats.AbonnementsPending)
	db.Model(&models.Abonnement{}).Where("statut = ?", "suspended").Count(&stats.AbonnementsSuspend)
	db.Model(&models.Abonnement{}).Where("statut = ?", "cancelled").Count(&stats.AbonnementsCancel)

	// Revenu total
	db.Model(&models.Abonnement{}).Select("COALESCE(SUM(montant), 0)").Scan(&stats.RevenuTotal)

	// Revenu mensuel (ce mois)
	debutMois := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	finMois := debutMois.AddDate(0, 1, -1)
	db.Model(&models.Abonnement{}).
		Where("created_at >= ? AND created_at <= ?", debutMois, finMois).
		Select("COALESCE(SUM(montant), 0)").
		Scan(&stats.RevenuMensuel)

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}
