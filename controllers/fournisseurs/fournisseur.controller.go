package fournisseurs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"

	"github.com/gofiber/fiber/v2"
)

// Synchronisation Send data to Local
func GetDataSynchronisation(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01")
	var data []models.Fournisseur

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("fournisseurs.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("fournisseurs.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Fournisseur",
		"data":    data,
	})
}

// Paginate
func GetPaginatedFournisseur(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

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

	var dataList []models.Fournisseur

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Fournisseur{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("entreprise_name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("entreprise_name ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("fournisseurs.updated_at DESC").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch fournisseurs",
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
		"message":    "All Fournisseurs paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllFournisseurs(c *fiber.Ctx) error {
	entrepriseUUID := c.Params("entreprise_uuid")
	db := database.DB

	var data []models.Fournisseur
	db.Where("entreprise_uuid = ?", entrepriseUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All fournisseurs",
		"data":    data,
	})
}

// Get one data
func GetFournisseur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var fournisseur models.Fournisseur
	db.Where("uuid = ?", uuid).First(&fournisseur)
	if fournisseur.EntrepriseName == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No fournisseur found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "fournisseur found",
			"data":    fournisseur,
		},
	)
}

// Create data
func CreateFournisseur(c *fiber.Ctx) error {
	p := &models.Fournisseur{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	p.UUID = utils.GenerateUUID()

	p.Sync = true

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "fournisseur created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateFournisseur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		EntrepriseName string `json:"entreprise_name"`
		Rccm           string `json:"rccm"`
		IdNat          string `json:"idnat"`
		NImpot         string `json:"nimpot"`
		Adresse        string `json:"adresse"`
		Email          string `json:"email"`     // Email officiel
		Telephone      string `json:"telephone"` // Telephone officiel
		Manager        string `json:"manager"`
		WebSite        string `json:"website"`
		TypeFourniture string `json:"type_fourniture"`
		Signature      string `json:"signature"`
		EntrepriseUUID string `json:"entreprise_uuid"`
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

	fournisseur := new(models.Fournisseur)

	db.Where("uuid = ?", uuid).First(&fournisseur)
	fournisseur.EntrepriseName = updateData.EntrepriseName
	fournisseur.Rccm = updateData.Rccm
	fournisseur.IdNat = updateData.IdNat
	fournisseur.NImpot = updateData.NImpot
	fournisseur.Adresse = updateData.Adresse
	fournisseur.Email = updateData.Email
	fournisseur.Telephone = updateData.Telephone
	fournisseur.Manager = updateData.Manager
	fournisseur.WebSite = updateData.WebSite
	fournisseur.Signature = updateData.Signature
	fournisseur.EntrepriseUUID = updateData.EntrepriseUUID

	fournisseur.Sync = true
	db.Save(&fournisseur)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "fournisseur updated success",
			"data":    fournisseur,
		},
	)

}

// Delete data
func DeleteFournisseur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var fournisseur models.Fournisseur
	err := db.Where("uuid = ?", uuid).First(&fournisseur)
	if err.Error != nil {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No fournisseur found",
				"data":    nil,
			},
		)
	}

	if fournisseur.UUID == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No fournisseur uuid found",
				"data":    nil,
			},
		)
	}

	db.Delete(&fournisseur)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "fournisseur deleted success",
			"data":    nil,
		},
	)
}

func UploadCsvDataFournisseur(c *fiber.Ctx) error {
	db := database.DB

	type UploadCSV struct {
		Data           []models.Fournisseur `json:"data"`
		EntrepriseUUID string               `json:"entreprise_uuid"`
		Signature      string               `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var cl models.Fournisseur

	for _, fournisseur := range dataUpload.Data {
		cl = models.Fournisseur{
			EntrepriseName: fournisseur.EntrepriseName,
			Rccm:           fournisseur.Rccm,
			IdNat:          fournisseur.IdNat,
			NImpot:         fournisseur.NImpot,
			Adresse:        fournisseur.Adresse,
			Email:          fournisseur.Email,
			Telephone:      fournisseur.Telephone,
			Manager:        fournisseur.Manager,
			WebSite:        fournisseur.WebSite,
			Signature:      dataUpload.Signature,
			EntrepriseUUID: dataUpload.EntrepriseUUID,
		}
		if fournisseur.EntrepriseName == "" {
			continue
		}
		db.Create(&cl)
	}

	fmt.Println("Fournisseura uploaded success")

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Fournisseura uploaded success",
			// "data":    dataUpload,
		},
	)
}

func GetDataUpload(data map[string]interface{}) ([]string, error) {
	var dataList []string

	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataList = append(dataList, string(dataStr))

	return dataList, nil
}
