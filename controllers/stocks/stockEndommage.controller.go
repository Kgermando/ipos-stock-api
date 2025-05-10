package stocks

import (
	"fmt"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"strconv" 

	"github.com/gofiber/fiber/v2"
)

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

	var length int64
	db.Model(dataList).Where("product_uuid = ?", productUUID).Count(&length)
	db.Where("product_uuid = ?", productUUID).
		Joins("JOIN products ON stocks.product_uuid=products.uuid").
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("stocks.created_at DESC").
		Preload("Product").
		Preload("Fournisseur").
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
		"message":    "All stocks paginated",
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
		"message": "Total qty stocks",
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
		CodeEntreprise uint64  `json:"code_entreprise"`
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
	stockEndommage.CodeEntreprise = updateData.CodeEntreprise

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
