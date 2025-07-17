package entreprises

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedEntreprise(c *fiber.Ctx) error {
	db := database.DB

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Default page number
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var dataList []models.Entreprise

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Entreprise{}).
		Where("name ILIKE ? OR rccm ILIKE ? OR id_nat ILIKE ? OR n_impot ILIKE ? OR email ILIKE ? OR telephone ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("name ILIKE ? OR rccm ILIKE ? OR id_nat ILIKE ? OR n_impot ILIKE ? OR email ILIKE ? OR telephone ILIKE ?",
		"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("entreprises.updated_at DESC").
		Preload("Users").
		Preload("Pos").
		Preload("Abonnement").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch entrepises",
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

	var entrepriseInfos []models.EntrepriseInfos
	for _, entreprise := range dataList {
		entrepriseInfos = append(entrepriseInfos, models.EntrepriseInfos{
			UUID:           entreprise.UUID,
			TypeEntreprise: entreprise.TypeEntreprise,
			Name:           entreprise.Name,
			Rccm:           entreprise.Rccm,
			IdNat:          entreprise.IdNat,
			NImpot:         entreprise.NImpot,
			Adresse:        entreprise.Adresse,
			Email:          entreprise.Email,
			Telephone:      entreprise.Telephone,
			Manager:        entreprise.Manager,
			Status:         entreprise.Status,
			Currency:       entreprise.Currency,
			Step:           entreprise.Step,
			TypeAbonnement: entreprise.TypeAbonnement,

			TotalUser:       len(entreprise.Users),
			TotalPos:        len(entreprise.Pos),
			TotalAbonnement: len(entreprise.Abonnement),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All entreprises",
		"data":       entrepriseInfos,
		"pagination": pagination,
	})
}

// Get All data
func GetAllEntreprises(c *fiber.Ctx) error {
	db := database.DB
	var data []models.Entreprise
	db.Preload("Users").Preload("Pos").Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All entreprises",
		"data":    data,
	})
}

// Get one data
func GetEntreprise(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var entreprise models.Entreprise

	db.Where("uuid = ?", uuid).
		Preload("Users").
		Preload("Pos").
		Preload("Abonnement").
		First(&entreprise)

	if entreprise.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No entreprise  name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "entreprise found",
			"data":    entreprise,
		},
	)
}

// Create data
func CreateEntreprise(c *fiber.Ctx) error {
	p := &models.Entreprise{}

	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON format",
			"data":    nil,
		})
	}

	db := database.DB

	// Vérifier si l'email existe déjà
	if p.Email != "" {
		var existingEntreprise models.Entreprise
		if err := db.Where("email = ?", p.Email).First(&existingEntreprise).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{
				"status":  "error",
				"message": "Une entreprise avec cet email existe déjà",
				"data":    existingEntreprise,
			})
		}
	}

	// Vérifier si le téléphone existe déjà
	if p.Telephone != "" {
		var existingEntreprise models.Entreprise
		if err := db.Where("telephone = ?", p.Telephone).First(&existingEntreprise).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{
				"status":  "error",
				"message": "Une entreprise avec ce numéro de téléphone existe déjà",
				"data":    existingEntreprise,
			})
		}
	}

	p.UUID = utils.GenerateUUID()

	if err := db.Create(p).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la création de l'entreprise",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Entreprise créée avec succès",
		"data":    p,
	})
}

// Update data
func UpdateEntreprise(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		TypeEntreprise string `json:"type_entreprise"` // PME, GE, Particulier
		Name           string `json:"name"`
		Rccm           string `json:"rccm"`
		IdNat          string `json:"idnat"`
		NImpot         string `json:"nimpot"`
		Adresse        string `json:"adresse"`
		Email          string `json:"email"`     // Email officiel
		Telephone      string `json:"telephone"` // Telephone officiel
		Manager        string `json:"manager"`
		Status         bool   `json:"status"`
		Currency       string `json:"currency"`        // Devise de l'entreprise, default CDF
		Step           int    `json:"step"`            // Etape de l'entreprise dans le processus d'inscription
		TypeAbonnement string `json:"type_abonnement"` // Pack starter, business, pro, entreprise
	}

	var updateData UpdateData

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your iunput",
				"data":    nil,
			},
		)
	}

	entreprise := new(models.Entreprise)

	db.Where("uuid = ?", uuid).First(&entreprise)
	entreprise.TypeEntreprise = updateData.TypeEntreprise
	entreprise.Name = updateData.Name
	entreprise.Rccm = updateData.Rccm
	entreprise.IdNat = updateData.IdNat
	entreprise.NImpot = updateData.NImpot
	entreprise.Adresse = updateData.Adresse
	entreprise.Email = updateData.Email
	entreprise.Telephone = updateData.Telephone
	entreprise.Manager = updateData.Manager
	entreprise.Status = updateData.Status
	entreprise.Currency = updateData.Currency
	entreprise.Step = updateData.Step
	entreprise.TypeAbonnement = updateData.TypeAbonnement

	db.Save(&entreprise)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "entreprise  updated success",
			"data":    entreprise,
		},
	)

}

// Delete data
func DeleteEntreprise(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var entreprise models.Entreprise
	db.Where("uuid = ?", uuid).First(&entreprise)
	if entreprise.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No Entreprise name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&entreprise)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Entreprise deleted success",
			"data":    nil,
		},
	)
}
