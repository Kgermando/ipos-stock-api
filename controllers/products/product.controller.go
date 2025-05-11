package products

import (
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Paginate
func GetPaginatedProductEntreprise(c *fiber.Ctx) error {
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

	var dataList []models.Product

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Product{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		Preload("Stocks").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch products",
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
		"message":    "All products paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate by posUUID
func GetPaginatedProductByPosUUID(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posUUID := c.Params("pos_uuid")

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

	var dataList []models.Product

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Product{}).
		Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("updated_at DESC").
		// Preload("Stocks").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch products",
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
		"message":    "All products paginated by posUUID",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllProducts(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posUUID := c.Params("pos_uuid")

	var data []models.Product
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUID).
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All products",
		"data":    data,
	})
}

// Get All data by id
func GetAllProductBySearch(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posUUID := c.Params("pos_uuid")

	search := c.Query("search", "")

	var data []models.Product
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All products by search",
		"data":    data,
	})
}

// Get one data
func GetProduct(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var product models.Product
	db.Where("uuid = ?", uuid).First(&product)
	if product.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No product name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product found",
			"data":    product,
		},
	)
}

// Create data
func CreateProduct(c *fiber.Ctx) error {
	p := &models.Product{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateProduct(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		// Image          string  `json:"image"`
		Reference      string  `json:"reference"`
		Name           string  `json:"name"`
		Description    string  `json:"description"`
		UniteVente     string  `json:"unite_vente"`
		PrixVente      float64 `json:"prix_vente"`
		Tva            float64 `json:"tva"`
		Stock          float64 `json:"stock"` // stock disponible
		Signature      string  `json:"signature"`
		PosUUID        string  `json:"pos_uuid"`
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

	product := new(models.Product)

	db.Where("uuid = ?", uuid).First(&product)
	product.Reference = updateData.Reference
	product.Name = updateData.Name
	product.Description = updateData.Description
	product.UniteVente = updateData.UniteVente
	product.PrixVente = updateData.PrixVente
	product.Tva = updateData.Tva
	product.Stock = updateData.Stock
	// product.Image = updateData.Image
	product.Signature = updateData.Signature
	product.PosUUID = updateData.PosUUID
	product.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&product)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product updated success",
			"data":    product,
		},
	)

}

// Update data stock disponible
func UpdateProductStockDispo(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		Stock float64 `json:"stock"` // stock disponible
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

	product := new(models.Product)

	db.Where("uuid = ?", uuid).First(&product)
	product.Stock = updateData.Stock

	db.Save(&product)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product updated Stock success",
			"data":    product,
		},
	)
}

// Update data stock Endommage
func UpdateProductStockEndommage(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		StockEndommage float64 `json:"stock_endommage"` // stock endommage
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

	product := new(models.Product)

	db.Where("uuid = ?", uuid).First(&product)
	product.StockEndommage = updateData.StockEndommage

	db.Save(&product)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product updated Stock success",
			"data":    product,
		},
	)
}

// Update data Restitution
func UpdateProductRestitution(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		Restitution float64 `gorm:"default:0" json:"restitution"` // stock restitution
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

	product := new(models.Product)

	db.Where("uuid = ?", uuid).First(&product)
	product.Restitution = updateData.Restitution

	db.Save(&product)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product updated Stock success",
			"data":    product,
		},
	)
}

// Delete data
func DeleteProduct(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var product models.Product
	db.Where("uuid = ?", uuid).First(&product)
	if product.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No product name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&product)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "product deleted success",
			"data":    nil,
		},
	)
}
