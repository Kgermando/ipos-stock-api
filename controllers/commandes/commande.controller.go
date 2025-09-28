package commandes

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
	var data []models.Commande

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
		"message": "All Commandes",
		"data":    data,
	})
}

// Paginate
func GetPaginatedCommandeEntreprise(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")

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

	var dataList []models.Commande
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Commande{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("ncommande::TEXT ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("ncommande::TEXT ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("commandes.updated_at DESC").
		Preload("CommandeLines").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch commande",
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
		"message":    "All commandes paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate
func GetPaginatedCommandePOS(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
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

	var dataList []models.Commande
	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Commande{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("ncommande::TEXT ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("ncommande::TEXT ILIKE ? OR status ILIKE ?", "%"+search+"%", "%"+search+"%").
		Offset(offset).
		Limit(limit).
		Order("commandes.updated_at DESC").
		Preload("CommandeLines").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch commande",
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
		"message":    "All commandes paginated",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllCommandes(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	var data []models.Commande
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Preload("CommandeLines").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All commandes",
		"data":    data,
	})
}

// Get one data
func GetCommande(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var commande models.Commande
	db.Where("uuid = ?", uuid).
		Preload("CommandeLines").
		First(&commande)
	if commande.Ncommande == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No commande name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "commande found",
			"data":    commande,
		},
	)
}

// Create data
func CreateCommande(c *fiber.Ctx) error {
	p := &models.Commande{}

	if err := c.BodyParser(&p); err != nil {
		return err
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
func UpdateCommande(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		PosUUID        string `json:"pos_uuid"`
		Ncommande      string `json:"ncommande"` // Number Random
		Status         string `json:"status"`    // Ouverte et Fermée
		ClientUUID     string `json:"client_uuid"`
		Signature      string `json:"signature"`
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

	commande := new(models.Commande)

	db.Where("uuid = ?", uuid).First(&commande)
	commande.PosUUID = updateData.PosUUID
	commande.Ncommande = updateData.Ncommande
	commande.Status = updateData.Status
	commande.ClientUUID = updateData.ClientUUID
	commande.Signature = updateData.Signature
	commande.EntrepriseUUID = updateData.EntrepriseUUID

	commande.Sync = true
	db.Save(&commande)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "commande updated success",
			"data":    commande,
		},
	)

}

// Delete data
func DeleteCommande(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var commande models.Commande
	db.Where("uuid = ?", uuid).First(&commande)
	if commande.Ncommande == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No commande name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&commande)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "commande deleted success",
			"data":    nil,
		},
	)
}
