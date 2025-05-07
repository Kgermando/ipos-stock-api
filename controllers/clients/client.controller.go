package clients

import (
	"encoding/json"
	"fmt"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedClient(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

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

	var dataList []models.Client

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("fullname ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("clients.updated_at DESC").
		Find(&dataList)

	if err != nil {
		fmt.Println("error s'est produite: ", err)
		return c.Status(500).SendString(err.Error())
	}

	// Calculate total number of pages
	totalPages := len(dataList) / limit
	if remainder := len(dataList) % limit; remainder > 0 {
		totalPages++
	}
	pagination := map[string]interface{}{
		"total_pages": totalPages,
		"page":        page,
		"page_size":   limit,
		"length":      length,
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
	codeEntreprise := c.Params("code_entreprise")
	db := database.DB

	var data []models.Client
	db.Where("code_entreprise = ?", codeEntreprise).Find(&data)
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
		Fullname       string `json:"fullname"`
		Telephone      string `json:"telephone"`
		Telephone2     string `json:"telephone2"`
		Email          string `json:"email"`
		Adress         string `json:"adress"`
		// Birthday       string `json:"birthday"`
		Organisation   string `json:"organisation"`
		WebSite        string `json:"website"`
		Signature      string `json:"signature"`
		CodeEntreprise uint64 `json:"code_entreprise"`
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
	client.CodeEntreprise = updateData.CodeEntreprise

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
		CodeEntreprise uint64          `json:"code_entreprise"`
		Signature      string          `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var cl models.Client

	for _, client := range dataUpload.Data {
		cl = models.Client{
			Fullname:       client.Fullname,
			Telephone:      client.Telephone,
			Telephone2:     client.Telephone2,
			Email:          client.Email,
			Adress:         client.Adress,
			// Birthday:       client.Birthday,
			Organisation:   client.Organisation,
			WebSite:        client.WebSite,
			Signature:      dataUpload.Signature,
			CodeEntreprise: dataUpload.CodeEntreprise,
		}
		if client.Fullname == "" {
			continue
		}
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