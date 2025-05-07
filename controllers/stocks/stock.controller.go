package stocks

import (
	"fmt"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedStock(c *fiber.Ctx) error {
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

	var dataList []models.Stock

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

// Get data
func GetStockMargeBeneficiaire(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")

	var data models.Stock

	db.Model(data).Where("product_uuid = ?", productUUID).Preload("Product").Last(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty stocks",
		"data":    data,
	})
}

// Get Total data
func GetTotalStock(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")

	var data []models.Stock
	var totalQty int64

	db.Model(data).Where("product_uuid = ?", productUUID).Select("SUM(quantity)").Scan(&totalQty)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty stocks",
		"data":    totalQty,
	})
}

// Get All data
func GetAllStocks(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")
	var data []models.Stock
	db.Where("product_uuid = ?", productUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All stocks",
		"data":    data,
	})
}

// Get one data
func GetStock(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var stock models.Stock
	db.Where("uuid = ?", uuid).First(&stock)
	if stock.UUID == "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No stock name found",
				"data":    nil,
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stock found",
			"data":    stock,
		},
	)
}

// Create data
func CreateStock(c *fiber.Ctx) error {
	p := &models.Stock{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stock created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateStock(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		PosUUID         string    `json:"pos_uuid"`
		ProductUUID     string    `json:"product_uuid"`
		Description     string    `json:"description"`
		FournisseurUUID string    `json:"fournisseur_uuid"`
		Quantity        float64   `json:"quantity"`
		PrixAchat       float64   `json:"prix_achat"`
		DateExpiration  time.Time `json:"date_expiration"`
		Signature       string    `json:"signature"`
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

	stock := new(models.Stock)

	db.Where("uuid = ?", uuid).First(&stock)
	stock.PosUUID = updateData.PosUUID
	stock.ProductUUID = updateData.ProductUUID
	stock.Description = updateData.Description
	stock.FournisseurUUID = updateData.FournisseurUUID
	stock.Quantity = updateData.Quantity
	stock.PrixAchat = updateData.PrixAchat
	stock.DateExpiration = updateData.DateExpiration
	stock.Signature = updateData.Signature

	db.Save(&stock)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stock updated success",
			"data":    stock,
		},
	)

}

// Delete data
func DeleteStock(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var stock models.Stock
	db.Where("uuid = ?", uuid).First(&stock)
	if stock.UUID == "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No stock name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&stock)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "stock deleted success",
			"data":    nil,
		},
	)
}
