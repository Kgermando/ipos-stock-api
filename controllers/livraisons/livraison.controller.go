package livraisons

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
	var data []models.Livraison

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("livraisons.updated_at DESC").
			Preload("Pos").
			Preload("Client").
			Preload("Livreur").
			Preload("Zone").
			Preload("Commandes").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("livraisons.updated_at DESC").
			Preload("Pos").
			Preload("Client").
			Preload("Livreur").
			Preload("Zone").
			Preload("Commandes").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Livraisons",
		"data":    data,
	})
}

// Paginate
func GetPaginatedLivraison(c *fiber.Ctx) error {
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

	var dataList []models.Livraison
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Livraison{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("statut ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("statut ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("livraisons.updated_at DESC").
		Preload("Pos").
		Preload("Client").
		Preload("Livreur").
		Preload("Zone").
		Preload("Commandes").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch livraisons",
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
		"message":    "All livraisons paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllLivraisons(c *fiber.Ctx) error {
	entrepriseUUID := c.Params("entreprise_uuid")
	db := database.DB

	var data []models.Livraison
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Preload("Pos").
		Preload("Client").
		Preload("Livreur").
		Preload("Zone").
		Preload("Commandes").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All livraisons",
		"data":    data,
	})
}

// Get one data
func GetLivraison(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var livraison models.Livraison
	db.Where("uuid = ?", uuid).
		Preload("Pos").
		Preload("Client").
		Preload("Livreur").
		Preload("Zone").
		Preload("Commandes").
		First(&livraison)
	if livraison.UUID == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No livraison found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livraison found",
			"data":    livraison,
		},
	)
}

// Create data
func CreateLivraison(c *fiber.Ctx) error {
	p := &models.Livraison{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	// Vérifier si la livraison existe déjà
	var existingLivraison models.Livraison
	database.DB.Where("uuid = ?", p.UUID).First(&existingLivraison)
	if existingLivraison.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Livraison avec cet UUID existe déjà",
				"data":    nil,
			},
		)
	}

	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livraison created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateLivraison(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		ClientUUID     string `json:"client_uuid"`
		LivreurUUID    string `json:"livreur_uuid"`
		ZoneUUID       string `json:"zone_uuid"`
		Statut         string `json:"statut"`
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

	livraison := new(models.Livraison)

	db.Where("uuid = ?", uuid).First(&livraison)
	livraison.ClientUUID = updateData.ClientUUID
	livraison.LivreurUUID = updateData.LivreurUUID
	livraison.ZoneUUID = updateData.ZoneUUID
	livraison.Statut = updateData.Statut
	livraison.Signature = updateData.Signature
	livraison.EntrepriseUUID = updateData.EntrepriseUUID

	db.Save(&livraison)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livraison updated success",
			"data":    livraison,
		},
	)
}

// Delete data
func DeleteLivraison(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var livraison models.Livraison
	db.Where("uuid = ?", uuid).First(&livraison)
	if livraison.UUID == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No livraison found",
				"data":    nil,
			},
		)
	}

	db.Delete(&livraison)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livraison deleted success",
			"data":    nil,
		},
	)
}

func UploadCsvDataLivraison(c *fiber.Ctx) error {
	db := database.DB

	type UploadCSV struct {
		Data           []models.Livraison `json:"data"`
		EntrepriseUUID string             `json:"entreprise_uuid"`
		Signature      string             `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var lv models.Livraison

	for _, livraison := range dataUpload.Data {
		lv = models.Livraison{
			PosUUID:        livraison.PosUUID,
			ClientUUID:     livraison.ClientUUID,
			LivreurUUID:    livraison.LivreurUUID,
			ZoneUUID:       livraison.ZoneUUID,
			Statut:         livraison.Statut,
			Signature:      dataUpload.Signature,
			EntrepriseUUID: dataUpload.EntrepriseUUID,
		}
		if livraison.PosUUID == "" {
			continue
		}

		livraison.Sync = true
		db.Create(&lv)
	}

	fmt.Println("livraisons uploaded success")

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "livraisons uploaded success",
		},
	)
}
