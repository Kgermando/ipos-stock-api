package tablebox

import (
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
	var data []models.TableBox

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("table_boxes.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("table_boxes.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All table boxes sync data",
		"data":    data,
	})
}

// Paginate by entreprise
func GetPaginatedTableBoxEntreprise(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))
	search := c.Query("search", "")
	sort := c.Query("sort", "name")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var dataList []models.TableBox
	var totalRecords int64

	query := db.Model(&models.TableBox{}).Where("entreprise_uuid = ?", entrepriseUUID)

	if search != "" {
		query = query.Where("name LIKE ? OR category LIKE ?", "%"+search+"%", "%"+search+"%")
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
			"message": "Failed to fetch table boxes",
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
		"message":    "All table boxes",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate by posUUID
func GetPaginatedTableBoxByPosUUID(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))
	search := c.Query("search", "")
	sort := c.Query("sort", "name")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var dataList []models.TableBox
	var totalRecords int64

	query := db.Model(&models.TableBox{}).Where("entreprise_uuid = ?", entrepriseUUID).Where("pos_uuid = ?", posUUID)

	if search != "" {
		query = query.Where("name LIKE ? OR category LIKE ?", "%"+search+"%", "%"+search+"%")
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
			"message": "Failed to fetch table boxes",
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
		"message":    "All table boxes",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllTableBoxs(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	var data []models.TableBox
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All table boxes",
		"data":    data,
	})
}

// Get All data by search
func GetAllTableBoxBySearch(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	search := c.Query("search", "")

	var data []models.TableBox

	query := db.Where("entreprise_uuid = ?", entrepriseUUID).Where("pos_uuid = ?", posUUID)

	if search != "" {
		query = query.Where("name LIKE ? OR category LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Preload("Pos").Find(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All table boxes",
		"data":    data,
	})
}

// Get one data
func GetTableBox(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var tableBox models.TableBox
	db.Where("uuid = ?", uuid).Preload("Pos").First(&tableBox)
	if tableBox.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No table box found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Table box found",
			"data":    tableBox,
		},
	)
}

// Create data
func CreateTableBox(c *fiber.Ctx) error {
	p := &models.TableBox{}

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

	p.UUID = utils.GenerateUUID()
	p.Sync = true

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Table box created successfully",
			"data":    p,
		},
	)
}

// Update data
func UpdateTableBox(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var tableBox models.TableBox

	// Find the table box
	db.Where("uuid = ?", uuid).First(&tableBox)
	if tableBox.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No table box found",
				"data":    nil,
			},
		)
	}

	// Parse request body
	if err := c.BodyParser(&tableBox); err != nil {
		return err
	}

	tableBox.Sync = true

	// Save to database
	database.DB.Save(&tableBox)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Table box updated successfully",
			"data":    tableBox,
		},
	)
}

// Delete data
func DeleteTableBox(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var tableBox models.TableBox
	db.Where("uuid = ?", uuid).First(&tableBox)
	if tableBox.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No table box found",
				"data":    nil,
			},
		)
	}

	db.Delete(&tableBox)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Table box deleted successfully",
			"data":    nil,
		},
	)
}

// Get table boxes by category
func GetTableBoxsByCategory(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	category := c.Params("category")

	var data []models.TableBox
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("category = ?", category).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All table boxes by category",
		"data":    data,
	})
}
