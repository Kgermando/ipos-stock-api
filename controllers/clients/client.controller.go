package clients

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
	var data []models.Client
	db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("created_at > ?", sync_created).
		Order("clients.updated_at DESC").
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Clients",
		"data":    data,
	})
}

// Paginate
func GetPaginatedClient(c *fiber.Ctx) error {
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

	var dataList []models.Client
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Client{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("fullname ILIKE ?", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("fullname ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("clients.updated_at DESC").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch provinces",
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
		"message":    "All clients paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllClients(c *fiber.Ctx) error {
	entrepriseUUID := c.Params("entreprise_uuid")
	db := database.DB

	var data []models.Client
	db.Where("entreprise_uuid = ?", entrepriseUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All clients",
		"data":    data,
	})
}

// Get one data
func GetClient(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var client models.Client
	db.Where("uuid = ?", uuid).First(&client)
	if client.Fullname == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No client found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "client found",
			"data":    client,
		},
	)
}

// Create data
func CreateClient(c *fiber.Ctx) error {
	p := &models.Client{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	p.UUID = utils.GenerateUUID()
	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "client created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateClient(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		Fullname   string `json:"fullname"`
		Telephone  string `json:"telephone"`
		Telephone2 string `json:"telephone2"`
		Email      string `json:"email"`
		Adress     string `json:"adress"`
		// Birthday       string `json:"birthday"`
		Organisation   string `json:"organisation"`
		WebSite        string `json:"website"`
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

	client := new(models.Client)

	db.Where("uuid = ?", uuid).First(&client)
	client.Fullname = updateData.Fullname
	client.Telephone = updateData.Telephone
	client.Telephone2 = updateData.Telephone2
	client.Email = updateData.Email
	client.Adress = updateData.Adress
	// client.Birthday = updateData.Birthday
	client.Organisation = updateData.Organisation
	client.WebSite = updateData.WebSite
	client.Signature = updateData.Signature
	client.EntrepriseUUID = updateData.EntrepriseUUID

	db.Save(&client)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "client updated success",
			"data":    client,
		},
	)

}

// Delete data
func DeleteClient(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var client models.Client
	db.Where("uuid = ?", uuid).First(&client)
	if client.Fullname == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No client found",
				"data":    nil,
			},
		)
	}

	db.Delete(&client)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "client deleted success",
			"data":    nil,
		},
	)
}

func UploadCsvDataClient(c *fiber.Ctx) error {
	db := database.DB

	type UploadCSV struct {
		Data           []models.Client `json:"data"`
		EntrepriseUUID string          `json:"entreprise_uuid"`
		Signature      string          `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var cl models.Client

	for _, client := range dataUpload.Data {
		cl = models.Client{
			Fullname:   client.Fullname,
			Telephone:  client.Telephone,
			Telephone2: client.Telephone2,
			Email:      client.Email,
			Adress:     client.Adress,
			// Birthday:       client.Birthday,
			Organisation:   client.Organisation,
			WebSite:        client.WebSite,
			Signature:      dataUpload.Signature,
			EntrepriseUUID: dataUpload.EntrepriseUUID,
		}
		if client.Fullname == "" {
			continue
		}

		client.Sync = true
		db.Create(&cl)
	}

	fmt.Println("clients uploaded success")

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "clients uploaded success",
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
