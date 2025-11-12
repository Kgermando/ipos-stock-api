package zones

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
	var data []models.Zone

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("zones.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("zones.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Zones",
		"data":    data,
	})
}

// Paginate
func GetPaginatedZone(c *fiber.Ctx) error {
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

	var dataList []models.Zone
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Zone{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("zones.updated_at DESC").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch zones",
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
		"message":    "All zones paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllZones(c *fiber.Ctx) error {
	entrepriseUUID := c.Params("entreprise_uuid")
	db := database.DB

	var data []models.Zone
	db.Where("entreprise_uuid = ?", entrepriseUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All zones",
		"data":    data,
	})
}

// Get one data
func GetZone(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var zone models.Zone
	db.Where("uuid = ?", uuid).First(&zone)
	if zone.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No zone found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "zone found",
			"data":    zone,
		},
	)
}

// Create data
func CreateZone(c *fiber.Ctx) error {
	p := &models.Zone{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	p.UUID = utils.GenerateUUID()
	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "zone created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateZone(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
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

	zone := new(models.Zone)

	db.Where("uuid = ?", uuid).First(&zone)
	zone.Name = updateData.Name
	zone.Description = updateData.Description
	zone.Signature = updateData.Signature
	zone.EntrepriseUUID = updateData.EntrepriseUUID

	db.Save(&zone)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "zone updated success",
			"data":    zone,
		},
	)
}

// Delete data
func DeleteZone(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var zone models.Zone
	db.Where("uuid = ?", uuid).First(&zone)
	if zone.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No zone found",
				"data":    nil,
			},
		)
	}

	db.Delete(&zone)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "zone deleted success",
			"data":    nil,
		},
	)
}

func UploadCsvDataZone(c *fiber.Ctx) error {
	db := database.DB

	type UploadCSV struct {
		Data           []models.Zone `json:"data"`
		EntrepriseUUID string        `json:"entreprise_uuid"`
		Signature      string        `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var zn models.Zone

	for _, zone := range dataUpload.Data {
		zn = models.Zone{
			Name:           zone.Name,
			Description:    zone.Description,
			Signature:      dataUpload.Signature,
			EntrepriseUUID: dataUpload.EntrepriseUUID,
		}
		if zone.Name == "" {
			continue
		}

		zone.Sync = true
		db.Create(&zn)
	}

	fmt.Println("zones uploaded success")

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "zones uploaded success",
		},
	)
}
