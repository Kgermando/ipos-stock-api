package stocks

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"

	"github.com/gofiber/fiber/v2"
)


// Synchronisation Send data to Local
func GetDataSynchronisationRestitution(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01") 

	var data []models.Restitution
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("created_at > ?", sync_created).
		Preload("Pos").
		Find(&data) 
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Restitutions",
		"data":    data,
	})
}

// Paginate
func GetPaginatedRestitution(c *fiber.Ctx) error {
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

	var dataList []models.Restitution

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Restitution{}).
		Where("product_uuid = ?", productUUID).
		Joins("JOIN products ON stocks.product_uuid=products.uuid").
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("product_uuid = ?", productUUID).
		Joins("JOIN products ON stocks.product_uuid=products.uuid").
		Where("products.name ILIKE ? OR products.reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("stocks.created_at DESC").
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
		"message":    "All restitutions paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get Total data
func GetTotalRestitution(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")

	var data []models.Restitution
	var totalQty int64

	db.Model(data).Where("product_uuid = ?", productUUID).Select("SUM(quantity)").Scan(&totalQty)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total qty restitutions",
		"data":    totalQty,
	})
}

// Get All data
func GetAllRestitutions(c *fiber.Ctx) error {
	db := database.DB
	productUUID := c.Params("product_uuid")
	var data []models.Restitution
	db.Where("product_uuid = ?", productUUID).Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All restitutions",
		"data":    data,
	})
}

// Get one data
func GetRestitution(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var restitution models.Restitution
	db.Where("uuid = ?", uuid).First(&restitution)
	if restitution.UUID == "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No restitution name found",
				"data":    nil,
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "restitution found",
			"data":    restitution,
		},
	)
}

// Create data
func CreateRestitution(c *fiber.Ctx) error {
	p := &models.Restitution{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}
	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "restitution created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateRestitution(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		PosUUID         string  `json:"pos_uuid"`
		ProductUUID     string  `json:"product_uuid"`
		Description     string  `json:"description"`
		Quantity        uint64  `json:"quantity"`
		PrixAchat       float64 `json:"prix_achat"`
		Motif           string  `json:"motif"`
		FournisseurUUID string  `json:"fournisseur_uuid"`
		Signature       string  `json:"signature"`
		EntrepriseUUID  string  `json:"entreprise_uuid"`
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

	restitution := new(models.Restitution)

	db.Where("uuid = ?", uuid).First(&restitution)
	restitution.PosUUID = updateData.PosUUID
	restitution.ProductUUID = updateData.ProductUUID
	restitution.Description = updateData.Description
	restitution.Quantity = updateData.Quantity
	restitution.PrixAchat = updateData.PrixAchat
	restitution.Motif = updateData.Motif
	restitution.FournisseurUUID = updateData.FournisseurUUID
	restitution.Signature = updateData.Signature
	restitution.EntrepriseUUID = updateData.EntrepriseUUID

	db.Save(&restitution)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "restitution updated success",
			"data":    restitution,
		},
	)

}

// Delete data
func DeleteRestitution(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var restitution models.Restitution
	db.Where("uuid = ?", uuid).First(&restitution)
	if restitution.UUID == "00000000-0000-0000-0000-000000000000" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No restitution name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&restitution)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "restitution deleted success",
			"data":    nil,
		},
	)
}
