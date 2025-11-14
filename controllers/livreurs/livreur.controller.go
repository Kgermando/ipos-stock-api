package livreurs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"

	"github.com/gofiber/fiber/v2"
)

// Synchronisation Send data to Local
func GetDataSynchronisation(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01")
	var data []models.Livreur

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("livreurs.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("livreurs.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Livreurs",
		"data":    data,
	})
}

// Paginate
func GetPaginatedLivreur(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

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

	var dataList []models.Livreur
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Livreur{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR telephone ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR telephone ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("livreurs.updated_at DESC").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch livreurs",
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
		"message":    "All livreurs paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllLivreurs(c *fiber.Ctx) error {
	entrepriseUUID := c.Params("entreprise_uuid")
	db := database.DB

	var data []models.Livreur
	db.Where("entreprise_uuid = ?", entrepriseUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All livreurs",
		"data":    data,
	})
}

// Get one data
func GetLivreur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var livreur models.Livreur
	db.Where("uuid = ?", uuid).First(&livreur)
	if livreur.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No livreur found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur found",
			"data":    livreur,
		},
	)
}

// Create data
func CreateLivreur(c *fiber.Ctx) error {
	p := &models.Livreur{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	// Vérifier si le livreur existe déjà
	var existingLivreur models.Livreur
	database.DB.Where("uuid = ?", p.UUID).First(&existingLivreur)
	if existingLivreur.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Livreur avec cet UUID existe déjà",
				"data":    nil,
			},
		)
	}

	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateLivreur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		TypeLivreur    string `json:"type_livreur"`
		Name           string `json:"name"`
		Telephone      string `json:"telephone"`
		Email          string `json:"email"`
		Adresse        string `json:"adresse"`
		Manager        string `json:"manager"`
		Signature      string `json:"signature"`
		EntrepriseUUID string `json:"entreprise_uuid"`
	}

	var updateData UpdateData

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    nil,
			},
		)
	}

	livreur := new(models.Livreur)

	db.Where("uuid = ?", uuid).First(&livreur)
	livreur.TypeLivreur = updateData.TypeLivreur
	livreur.Name = updateData.Name
	livreur.Telephone = updateData.Telephone
	livreur.Email = updateData.Email
	livreur.Adresse = updateData.Adresse
	livreur.Manager = updateData.Manager
	livreur.Signature = updateData.Signature
	livreur.EntrepriseUUID = updateData.EntrepriseUUID

	db.Save(&livreur)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur updated success",
			"data":    livreur,
		},
	)
}

// Delete data
func DeleteLivreur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var livreur models.Livreur
	db.Where("uuid = ?", uuid).First(&livreur)
	if livreur.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No livreur found",
				"data":    nil,
			},
		)
	}

	db.Delete(&livreur)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreur deleted success",
			"data":    nil,
		},
	)
}

func UploadCsvDataLivreur(c *fiber.Ctx) error {
	db := database.DB

	type UploadCSV struct {
		Data           []models.Livreur `json:"data"`
		EntrepriseUUID string           `json:"entreprise_uuid"`
		Signature      string           `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var lv models.Livreur

	for _, livreur := range dataUpload.Data {
		lv = models.Livreur{
			TypeLivreur:    livreur.TypeLivreur,
			Name:           livreur.Name,
			Telephone:      livreur.Telephone,
			Email:          livreur.Email,
			Adresse:        livreur.Adresse,
			Manager:        livreur.Manager,
			Signature:      dataUpload.Signature,
			EntrepriseUUID: dataUpload.EntrepriseUUID,
		}
		if livreur.Name == "" {
			continue
		}

		livreur.Sync = true
		db.Create(&lv)
	}

	fmt.Println("livreurs uploaded success")

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livreurs uploaded success",
		},
	)
}

// Get livreurs by type
func GetLivreursByType(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	typeLivreur := c.Params("type")

	var data []models.Livreur
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("type_livreur = ?", typeLivreur).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All livreurs by type",
		"data":    data,
	})
}
