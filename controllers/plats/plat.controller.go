package plats

import (
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
	var data []models.Plat

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("plats.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("plats.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All plats sync data",
		"data":    data,
	})
}

// Paginate by entreprise
func GetPaginatedPlatEntreprise(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))
	search := c.Query("search", "")
	sort := c.Query("sort", "name")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var dataList []models.Plat
	var totalRecords int64

	query := db.Model(&models.Plat{}).Where("entreprise_uuid = ?", entrepriseUUID)

	if search != "" {
		query = query.Where("name LIKE ? OR reference LIKE ? OR categorie LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	query.Count(&totalRecords)

	// Get paginated results
	err := query.Order(sort + " " + order).
		Offset(offset).
		Limit(limit).
		Preload("Pos").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch plats",
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
		"message":    "All plats",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate by posUUID
func GetPaginatedPlatByPosUUID(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))
	search := c.Query("search", "")
	sort := c.Query("sort", "name")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var dataList []models.Plat
	var totalRecords int64

	query := db.Model(&models.Plat{}).Where("entreprise_uuid = ?", entrepriseUUID).Where("pos_uuid = ?", posUUID)

	if search != "" {
		query = query.Where("name LIKE ? OR reference LIKE ? OR categorie LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	query.Count(&totalRecords)

	// Get paginated results
	err := query.Order(sort + " " + order).
		Offset(offset).
		Limit(limit).
		Preload("Pos").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch plats",
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
		"message":    "All plats",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllPlats(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	var data []models.Plat
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All plats",
		"data":    data,
	})
}

// Get All data by search
func GetAllPlatBySearch(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	search := c.Query("search", "")

	var data []models.Plat

	query := db.Where("entreprise_uuid = ?", entrepriseUUID).Where("pos_uuid = ?", posUUID)

	if search != "" {
		query = query.Where("name LIKE ? OR reference LIKE ? OR categorie LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Preload("Pos").Find(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All plats",
		"data":    data,
	})
}

// Get one data
func GetPlat(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var plat models.Plat
	db.Where("uuid = ?", uuid).Preload("Pos").First(&plat)
	if plat.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No plat found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat found",
			"data":    plat,
		},
	)
}

// Create data
func CreatePlat(c *fiber.Ctx) error {
	p := &models.Plat{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	if p.Name == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Form not complete",
				"data":    nil,
			},
		)
	}

	// Vérifier si le plat existe déjà
	var existingPlat models.Plat
	database.DB.Where("uuid = ?", p.UUID).First(&existingPlat)
	if existingPlat.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Plat avec cet UUID existe déjà",
				"data":    nil,
			},
		)
	}

	p.Sync = true

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat created successfully",
			"data":    p,
		},
	)
}

// Update data
func UpdatePlat(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var plat models.Plat

	// Find the plat
	db.Where("uuid = ?", uuid).First(&plat)
	if plat.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No plat found",
				"data":    nil,
			},
		)
	}

	// Parse request body
	if err := c.BodyParser(&plat); err != nil {
		return err
	}

	plat.Sync = true

	// Save to database
	database.DB.Save(&plat)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat updated successfully",
			"data":    plat,
		},
	)
}

// Delete data
func DeletePlat(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var plat models.Plat
	db.Where("uuid = ?", uuid).First(&plat)
	if plat.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No plat found",
				"data":    nil,
			},
		)
	}

	db.Delete(&plat)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat deleted successfully",
			"data":    nil,
		},
	)
}

// Update availability
func UpdatePlatAvailability(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var plat models.Plat
	db.Where("uuid = ?", uuid).First(&plat)
	if plat.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No plat found",
				"data":    nil,
			},
		)
	}

	// Parse request body for availability
	var updateData struct {
		IsAvailable bool `json:"is_available"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return err
	}

	plat.IsAvailable = updateData.IsAvailable
	plat.Sync = true

	db.Save(&plat)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "plat availability updated successfully",
			"data":    plat,
		},
	)
}

// Get available plats only
func GetAvailablePlats(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	var data []models.Plat
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("is_available = ?", true).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All available plats",
		"data":    data,
	})
}
