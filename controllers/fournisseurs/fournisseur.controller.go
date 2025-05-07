package fournisseurs

import (
	"encoding/json"
	"fmt"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedFournisseur(c *fiber.Ctx) error {
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

	var dataList []models.Fournisseur

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("fullname ILIKE ?", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("fournisseurs.updated_at DESC").
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
		"message":    "All Fournisseurs paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllFournisseurs(c *fiber.Ctx) error {
	codeEntreprise := c.Params("code_entreprise")
	db := database.DB

	var data []models.Fournisseur
	db.Where("code_entreprise = ?", codeEntreprise).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All clients",
		"data":    data,
	})
}

// Get one data
func GetFournisseur(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var fournisseur models.Fournisseur
	db.Where("uuid = ?", uuid).First(&fournisseur)
	if fournisseur.Fullname == "" {
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
		Fullname   string `gorm:"not null" json:"fullname"`
		Telephone  string `gorm:"not null" json:"telephone"`
		Telephone2 string `json:"telephone2"`
		Email      string `json:"email"`
		Adress     string `json:"adress"`
		Entreprise string `json:"entreprise"`
		WebSite    string `json:"website"`

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

	fournisseur := new(models.Fournisseur)

	db.Where("uuid = ?", uuid).First(&fournisseur)
	fournisseur.Fullname = updateData.Fullname
	fournisseur.Telephone = updateData.Telephone
	fournisseur.Telephone2 = updateData.Telephone2
	fournisseur.Email = updateData.Email
	fournisseur.Adress = updateData.Adress
	fournisseur.Entreprise = updateData.Entreprise
	fournisseur.WebSite = updateData.WebSite
	fournisseur.Signature = updateData.Signature
	fournisseur.CodeEntreprise = updateData.CodeEntreprise

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
	db.Where("uuid = ?", uuid).First(&fournisseur)
	if fournisseur.Fullname == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No fournisseur found",
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
		CodeEntreprise uint64               `json:"code_entreprise"`
		Signature      string               `json:"signature"`
	}

	var dataUpload UploadCSV
	if err := json.Unmarshal(c.Body(), &dataUpload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var cl models.Fournisseur

	for _, Fournisseur := range dataUpload.Data {
		cl = models.Fournisseur{
			Fullname:       Fournisseur.Fullname,
			Telephone:      Fournisseur.Telephone,
			Telephone2:     Fournisseur.Telephone2,
			Email:          Fournisseur.Email,
			Adress:         Fournisseur.Adress,
			Entreprise:     Fournisseur.Entreprise,
			WebSite:        Fournisseur.WebSite,
			Signature:      dataUpload.Signature,
			CodeEntreprise: dataUpload.CodeEntreprise,
		}
		if Fournisseur.Fullname == "" {
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
