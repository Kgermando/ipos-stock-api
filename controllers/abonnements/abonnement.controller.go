package abonnements

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedAbonnement(c *fiber.Ctx) error {
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

	var dataList []models.Abonnement

	var totalRecords int64

	db.Model(&models.Abonnement{}).
		Where("bouquet ILIKE ? OR moyen_payment ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("bouquet ILIKE ? OR moyen_payment ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("abonnement.updated_at DESC").
		Preload("Entreprise").
		Find(&dataList).Error

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

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All abonnements",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Query all data ID
func GetPaginatedAbonnementByID(c *fiber.Ctx) error {
	db := database.DB
	EntrepriseUUID := c.Params("entreprise_uuid")

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

	var dataList []models.Abonnement
	var totalRecords int64

	db.Model(&models.Abonnement{}).
		Where("entreprise_uuid = ?", EntrepriseUUID).
		Where("bouquet ILIKE ? OR moyen_payment ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

		err = db.Where("entreprise_uuid = ?", EntrepriseUUID).
		Where("bouquet ILIKE ? OR moyen_payment ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("abonnement.updated_at DESC").
		Preload("Entreprise").
		Find(&dataList).Error 

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

	// Return response
	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "abonnements retrieved successfully",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllAbonnements(c *fiber.Ctx) error {
	db := database.DB
	var data []models.Abonnement

	db.Preload("Entreprise").Find(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All abonnements",
		"data":    data,
	})
}

// Get All data
func GetAllAbonnementByUUId(c *fiber.Ctx) error {
	db := database.DB
	EntrepriseUUID := c.Params("entreprise_uuid")

	var data []models.Abonnement
	db.Where("entreprise_uuid = ?", EntrepriseUUID).Preload("Entreprise").Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All abonnements",
		"data":    data,
	})
}

// Get one data
func GetAbonnement(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var abonnement models.Abonnement
	db.Where("uuid = ?", uuid).Preload("Entreprise").First(&abonnement)
	if abonnement.Bouquet == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No abonnement name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "abonnement found",
			"data":    abonnement,
		},
	)
}

// Create data
func CreateAbonnement(c *fiber.Ctx) error {
	p := &models.Abonnement{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	p.UUID = utils.GenerateUUID()
	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "abonnement created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateAbonnement(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		EntrepriseUUID string  `json:"entreprise_uuid"`
		Bouquet        string  `json:"bouquet"` // Premium, Platinium, Fremium
		Montant        float64 `json:"montant"`
		MoyenPayment   string  `json:"moyen_payment"`
		Signature      string  `json:"signature"`
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

	abonnement := new(models.Abonnement)

	db.Where("uuid = ?", uuid).First(&abonnement)
	abonnement.EntrepriseUUID = updateData.EntrepriseUUID
	abonnement.Bouquet = updateData.Bouquet
	abonnement.Montant = updateData.Montant
	abonnement.MoyenPayment = updateData.MoyenPayment
	abonnement.Signature = updateData.Signature

	db.Save(&abonnement)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "abonnement updated success",
			"data":    abonnement,
		},
	)

}

// Delete data
func DeleteAbonnement(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var abonnement models.Abonnement
	db.Where("uuid = ?", uuid).First(&abonnement)
	if abonnement.Bouquet == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No abonnement Bouquet found",
				"data":    nil,
			},
		)
	}

	db.Delete(&abonnement)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "abonnement deleted success",
			"data":    nil,
		},
	)
}
