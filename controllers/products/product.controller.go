package products

import (
	"fmt"
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

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("products.updated_at DESC").
		Preload("Stocks").
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
		"message":    "All products paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate
func GetPaginatedProduct(c *fiber.Ctx) error {
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

	var length int64
	db.Model(dataList).Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUID).Count(&length)
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("products.updated_at DESC").
		Preload("Stocks").
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
		"message":    "All products paginated",
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
	db.Where("uuid = ?", uuid).Find(&product)
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
		Reference      string  `json:"reference"`
		Name           string  `json:"name"`
		Description    string  `json:"description"`
		UniteVente     string  `json:"unite_vente"`
		PrixVente      float64 `json:"prix_vente"`
		Tva            float64 `json:"tva"`
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

	db.Where("uuid = ?", uuid).Find(&product)
	product.Reference = updateData.Reference
	product.Name = updateData.Name
	product.Description = updateData.Description
	product.UniteVente = updateData.UniteVente
	product.PrixVente = updateData.PrixVente
	product.Tva = updateData.Tva
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

// Delete data
func DeleteProduct(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var product models.Product
	db.Where("uuid = ?", uuid).Find(&product)
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
