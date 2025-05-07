package pos

import (
	"fmt"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedPos(c *fiber.Ctx) error {
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

	var dataList []models.Pos

	var length int64
	db.Model(dataList).Count(&length)
	db.Where("name ILIKE ? OR manager ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("pos.updated_at DESC").
		Preload("Entreprise").
		Preload("Stocks").
		Preload("BonCommades").
		Preload("Commandes").
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
		"message":    "All poss",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Query all data ID
func GetPaginatedPosByID(c *fiber.Ctx) error {
	db := database.DB
	EntrepriseID := c.Params("entreprise_id")

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

	var dataList []models.Pos

	var length int64
	var data []models.Pos
	db.Model(data).Where("entreprise_id = ?", EntrepriseID).Count(&length)
	db.Where("entreprise_id = ?", EntrepriseID).
		Where("name ILIKE ? OR manager ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("pos.updated_at DESC").
		Preload("Entreprise").
		Preload("Stocks").
		Preload("BonCommades").
		Preload("Commandes").
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
		"message":    "All bonCommandeLine by BonCommande",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllPoss(c *fiber.Ctx) error {
	db := database.DB
	var data []models.Pos
	db.Preload("Stocks").
		// Preload("BonCommades").
		Preload("Commandes").Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All poss",
		"data":    data,
	})
}

// Get All data
func GetAllPosById(c *fiber.Ctx) error {
	db := database.DB
	EntrepriseID := c.Params("entreprise_id")

	var data []models.Pos
	db.Where("entreprise_id = ?", EntrepriseID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All poss",
		"data":    data,
	})
}

// Get one data
func GetPos(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var pos models.Pos
	db.Preload("Stocks").
		// Preload("BonCommades").
		Preload("Commandes").Find(&pos, id)
	if pos.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No pos name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "pos found",
			"data":    pos,
		},
	)
}

// Create data
func CreatePos(c *fiber.Ctx) error {
	p := &models.Pos{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "pos created success",
			"data":    p,
		},
	)
}

// Update data
func UpdatePos(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	type UpdateData struct {
		EntrepriseUUID string `json:"entreprise_uuid"`
		Name           string `json:"name"`
		Adresse        string `json:"adresse"`
		Email          string `json:"email"`
		Telephone      string `json:"telephone"`
		Manager        string `json:"manager"`
		Status         bool   `json:"status"` // Actif ou Inactif
		Signature      string `json:"signature"`
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

	pos := new(models.Pos)

	db.First(&pos, id)
	pos.EntrepriseUUID = updateData.EntrepriseUUID
	pos.Name = updateData.Name
	pos.Email = updateData.Email
	pos.Telephone = updateData.Telephone
	pos.Manager = updateData.Manager
	pos.Status = updateData.Status
	pos.Signature = updateData.Signature

	db.Save(&pos)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "pos updated success",
			"data":    pos,
		},
	)

}

// Delete data
func DeletePos(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	var pos models.Pos
	db.First(&pos, id)
	if pos.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No pos name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&pos)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "pos deleted success",
			"data":    nil,
		},
	)
}
