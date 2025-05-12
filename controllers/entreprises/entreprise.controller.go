package entreprises

import (
	"strconv"
	"time"

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
		Where("name ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("name ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("entreprises.updated_at DESC").
		Preload("Users").
		Preload("Pos").
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
			ID:              entreprise.ID,
			UUID:            entreprise.UUID,
			TypeEntreprise:  entreprise.TypeEntreprise,
			Name:            entreprise.Name,
			Code:            entreprise.Code,
			Rccm:            entreprise.Rccm,
			IdNat:           entreprise.IdNat,
			NImpot:          entreprise.NImpot,
			Adresse:         entreprise.Adresse,
			Email:           entreprise.Email,
			Telephone:       entreprise.Telephone,
			Manager:         entreprise.Manager,
			Status:          entreprise.Status,
			Signature:       entreprise.Signature,
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
		return err
	}

	p.UUID = utils.GenerateUUID()

	p.Sync = true

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "entreprise created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateEntreprise(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		TypeEntreprise string    `json:"type_entreprise"`
		Name           string    `json:"name"`
		Code           uint64    `json:"code"` // Code entreprise
		Rccm           string    `json:"rccm"`
		IdNat          string    `json:"idnat"`
		NImpot         string    `json:"nimpot"`
		Adresse        string    `json:"adresse"`
		Email          string    `json:"email"`     // Email officiel
		Telephone      string    `json:"telephone"` // Telephone officiel
		Manager        string    `json:"manager"`
		Status         bool      `json:"status"`
		TypeAbonnement string    `json:"type_abonnement"`
		Abonnement     time.Time `json:"abonnement"`
		Signature      string    `json:"signature"`
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
	entreprise.Code = updateData.Code
	entreprise.Rccm = updateData.Rccm
	entreprise.IdNat = updateData.IdNat
	entreprise.NImpot = updateData.NImpot
	entreprise.Adresse = updateData.Adresse
	entreprise.Email = updateData.Email
	entreprise.Telephone = updateData.Telephone
	entreprise.Manager = updateData.Manager
	entreprise.Status = updateData.Status
	entreprise.TypeAbonnement = updateData.TypeAbonnement
	entreprise.Signature = updateData.Signature

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
