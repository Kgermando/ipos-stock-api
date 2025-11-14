package stocks

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"

	"github.com/gofiber/fiber/v2"
)

// Synchronisation Send data to Local
func GetDataSynchronisationStockEndommage(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01")
	var data []models.StockEndommage

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All StockEndommages",
		"data":    data,
	})
}

// Get All data
func GetAllByUUIDStockEndommages(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	productUUID := c.Params("product_uuid")

	var data []models.StockEndommage
	db.
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("product_uuid = ?", productUUID).
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All stockEndommages",
		"data":    data,
	})
}

// Paginate
func GetPaginatedStockEndommage(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")

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

	var dataList []models.StockEndommage

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.StockEndommage{}).
		Where("product_uuid = ?", productUUID).
		Joins("JOIN products ON stock_endommages.product_uuid=products.uuid").
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("product_uuid = ?", productUUID).
		Joins("JOIN products ON stock_endommages.product_uuid=products.uuid").
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("stock_endommages.created_at DESC").
		Preload("Product").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch restitutions",
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
		"message":    "All stockEndommages paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get Total data
func GetTotalStockEndommage(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")

	var data []models.StockEndommage
	var totalQty float64

	db.Model(data).Where("product_uuid = ?", productUUID).Select("SUM(quantity)").Scan(&totalQty)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty stockEndommages",
		"data":    totalQty,
	})
}

// Get All data
func GetAllStockEndommages(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")
	var data []models.StockEndommage
	db.Where("product_uuid = ?", productUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All stockEndommages",
		"data":    data,
	})
}

// Get one data
func GetStockEndommage(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var stockEndommage models.StockEndommage
	db.Where("uuid = ?", uuid).First(&stockEndommage)
	if stockEndommage.UUID == "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No stockEndommage name found",
				"data":    nil,
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stockEndommage found",
			"data":    stockEndommage,
		},
	)
}

// Create data
func CreateStockEndommage(c *fiber.Ctx) error {
	p := &models.StockEndommage{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	// Vérifier si le stock endommagé existe déjà
	var existingStockEndommage models.StockEndommage
	database.DB.Where("uuid = ?", p.UUID).First(&existingStockEndommage)
	if existingStockEndommage.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "StockEndommage avec cet UUID existe déjà",
				"data":    nil,
			},
		)
	}

	p.Sync = true

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stockEndommage created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateStockEndommage(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		PosUUID        string  `json:"pos_uuid"`
		ProductUUID    string  `json:"product_uuid"`
		Quantity       float64 `json:"quantity"`
		PrixAchat      float64 `json:"prix_achat"`
		Raison         string  `json:"raison"` // Raison de l'endommagement
		Signature      string  `json:"signature"`
		EntrepriseUUID string  `json:"entreprise_uuid"`
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

	stockEndommage := new(models.StockEndommage)

	db.Where("uuid = ?", uuid).First(&stockEndommage)
	stockEndommage.PosUUID = updateData.PosUUID
	stockEndommage.ProductUUID = updateData.ProductUUID
	stockEndommage.Quantity = updateData.Quantity
	stockEndommage.PrixAchat = updateData.PrixAchat
	stockEndommage.Raison = updateData.Raison
	stockEndommage.Signature = updateData.Signature
	stockEndommage.EntrepriseUUID = updateData.EntrepriseUUID

	stockEndommage.Sync = true
	db.Save(&stockEndommage)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stockEndommage updated success",
			"data":    stockEndommage,
		},
	)

}

// Delete data
func DeleteStockEndommage(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var stockEndommage models.StockEndommage
	db.Where("uuid = ?", uuid).First(&stockEndommage)
	if stockEndommage.UUID == "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No stock name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&stockEndommage)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stockEndommage deleted success",
			"data":    nil,
		},
	)
}
