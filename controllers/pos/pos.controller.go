package pos

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"

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

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Pos{}).
		Where("name ILIKE ? OR manager ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("name ILIKE ? OR manager ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("pos.updated_at DESC").
		Preload("Entreprise").
		Find(&dataList).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch pos",
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
		"message":    "All poss",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Query all data UUID
func GetPaginatedPosByUUID(c *fiber.Ctx) error {
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

	var dataList []models.Pos

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.CommandeLine{}).
		Where("entreprise_uuid = ?", EntrepriseUUID).
		Where("name ILIKE ? OR manager ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", EntrepriseUUID).
		Where("name ILIKE ? OR manager ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("pos.updated_at DESC").
		Preload("Entreprise").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch pos",
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
		"message":    "All pos by uuid",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllPoss(c *fiber.Ctx) error {
	db := database.DB

	var data []models.Pos
	db.Find(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All poss",
		"data":    data,
	})
}

// Get All data by UUID
func GetAllPosByUUId(c *fiber.Ctx) error {
	db := database.DB
	EntrepriseUUID := c.Params("entreprise_uuid")

	var data []models.Pos
	db.Where("entreprise_uuid = ?", EntrepriseUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All poss",
		"data":    data,
	})
}

// Get one data
func GetPos(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB
	var pos models.Pos

	db.Where("uuid = ?", uuid).
		Preload("Entreprise").
		Preload("Users").
		Preload("Products").
		Preload("Commandes").
		Preload("Clients").
		Preload("Fournisseurs").
		First(&pos)

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

	p.UUID = utils.GenerateUUID()
	p.Sync = true
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
	uuid := c.Params("uuid")
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

	db.Where("uuid = ?", uuid).First(&pos)
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
	uuid := c.Params("uuid")

	db := database.DB

	var pos models.Pos
	db.Where("uuid = ?", uuid).First(&pos)
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
