package commandes

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"

	"github.com/gofiber/fiber/v2"
)

// Synchronisation Send data to Local
func GetDataSynchronisationCommandeLine(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01")
	var data []models.CommandeLine

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All CommandeLines",
		"data":    data,
	})
}

// Query all data ID
func GetPaginatedCommandeLineByID(c *fiber.Ctx) error {
	db := database.DB
	commandeUUID := c.Params("commande_uuid")

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

	var dataList []models.CommandeLine

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.CommandeLine{}).
		Joins("JOIN products ON commande_lines.product_uuid=products.uuid").
		Where("commande_lines.commande_uuid = ?", commandeUUID).
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.
		Joins("JOIN products ON commande_lines.product_uuid=products.uuid").
		Where("commande_lines.commande_uuid = ?", commandeUUID).
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Select(`
			commande_lines.uuid AS uuid,
			products.reference AS reference,
			products.name AS name,
			products.description AS description,
			products.unite_vente AS unite_vente,
			commande_lines.quantity AS quantity,
			products.prix_vente AS prix_vente,
			products.tva AS tva,
			SUM(commande_lines.quantity::FLOAT * products.prix_vente::FLOAT)
		`).
		Offset(offset).
		Limit(limit).
		Order("commande_lines.updated_at DESC").
		Find(&dataList).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch commande_lines",
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
		"message":    "All commandeLine by commande",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data by UUID
func GetAllCommandeLineByUUId(c *fiber.Ctx) error {
	db := database.DB
	commandeUUID := c.Params("commande_uuid")

	var dataList []models.CommandeLine
	db.Where("commande_uuid = ?", commandeUUID).
		Order("updated_at DESC").
		Preload("Commande").
		Preload("Product").
		Find(&dataList)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All CommandeLine",
		"data":    dataList,
	})
}

// Get All data
func GetAllCommandeLines(c *fiber.Ctx) error {
	db := database.DB
	var data []models.CommandeLine
	db.Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All CommandeLine",
		"data":    data,
	})
}

// Get Total data
func GetTotalCommandeLine(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")

	var data []models.CommandeLine
	var totalQty int64

	if productUUID != "00000000-0000-0000-0000-000000000000" {
		db.Model(data).Where("product_uuid = ?", productUUID).Select("SUM(quantity)").Scan(&totalQty)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty stocks",
		"data":    totalQty,
	})
}

// Get one data
func GetCommandeLine(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB
	var commandeLine models.CommandeLine

	db.Where("uuid = ?", uuid).First(&commandeLine)
	if commandeLine.ProductUUID != "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No commandeLine found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "commandeLine found",
			"data":    commandeLine,
		},
	)
}

// Create data
func CreateCommandeLine(c *fiber.Ctx) error {
	p := &models.CommandeLine{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	// Vérifier si la ligne de commande existe déjà
	var existingCommandeLine models.CommandeLine
	database.DB.Where("uuid = ?", p.UUID).First(&existingCommandeLine)
	if existingCommandeLine.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "CommandeLine avec cet UUID existe déjà",
				"data":    nil,
			},
		)
	}

	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "commande created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateCommandeLine(c *fiber.Ctx) error {
	db := database.DB
	uuid := c.Params("uuid")

	type UpdateData struct {
		CommandeUUID   string `json:"commande_uuid"`
		ProductUUID    string `json:"product_uuid"`
		Quantity       uint64 `json:"quantity"`
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

	commandeLine := new(models.CommandeLine)

	db.Where("uuid = ?", uuid).First(&commandeLine)
	commandeLine.CommandeUUID = updateData.CommandeUUID
	commandeLine.ProductUUID = updateData.ProductUUID
	commandeLine.Quantity = updateData.Quantity
	commandeLine.EntrepriseUUID = updateData.EntrepriseUUID

	commandeLine.Sync = true
	db.Save(&commandeLine)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "commandeLine updated success",
			"data":    commandeLine,
		},
	)

}

// Delete data
func DeleteCommandeLine(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var commandeLine models.CommandeLine
	db.Where("uuid = ?", uuid).First(&commandeLine)
	if commandeLine.ProductUUID != "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No commandeLine found",
				"data":    nil,
			},
		)
	}

	db.Delete(&commandeLine)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "CommandeLine deleted success",
			"data":    nil,
		},
	)
}
