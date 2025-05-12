package caisses

import (
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"

	"github.com/gofiber/fiber/v2"
)

// Get All data
func GetTotalAllCaisses(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	var dataList []models.CaisseItem
	db.Where("code_entreprise = ?", codeEntreprise).Find(&dataList)

	var total float64 = 0
	var totalEntree float64 = 0
	var totalSortie float64 = 0
	var solde float64 = 0
	var pourcent float64 = 0

	for _, item := range dataList {
		if item.TypeTransaction == "Entrée" {
			totalEntree += item.Montant
		}
		if item.TypeTransaction == "Sortie" {
			totalSortie += item.Montant
		}
	}

	total = totalEntree + totalSortie
	solde = totalEntree - totalSortie
	pourcent = solde * 100 / (totalEntree + totalSortie)

	response := map[string]interface{}{
		"total":       total,
		"totalentree": totalEntree,
		"totalsortie": totalSortie,
		"solde":       solde,
		"pourcent":    pourcent,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Total All caisses",
		"data":    response,
	})
}

// Get All data
func GetAllCaisses(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")

	var data []models.Caisse
	db.Where("code_entreprise = ?", codeEntreprise).
		Preload("Pos").
		Order("caisses.updated_at ASC").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses",
		"data":    data,
	})
}

// Get All data
func GetAllCaisseByPos(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posUUId := c.Params("pos_uuid")

	var data []models.Caisse
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUId).
		Preload("Pos").
		Order("caisses.updated_at ASC").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses",
		"data":    data,
	})
}

// Get All data by id
func GetAllCaisseBySearch(c *fiber.Ctx) error {
	db := database.DB
	codeEntreprise := c.Params("code_entreprise")
	posUUId := c.Params("pos_uuid")

	search := c.Query("search", "")

	var data []models.Caisse
	db.Where("code_entreprise = ?", codeEntreprise).
		Where("pos_uuid = ?", posUUId).
		Where("name ILIKE ?", "%"+search+"%").
		Preload("Pos").
		Order("caisses.updated_at ASC").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All caisses by id",
		"data":    data,
	})
}

// Get one data
func GetCaisse(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var caisse models.Caisse
	db.Where("uuid = ?", uuid).
		Preload("Pos").First(&caisse)

	if caisse.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No caisse name found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse found",
			"data":    caisse,
		},
	)
}

// Create data
func CreateCaisse(c *fiber.Ctx) error {
	p := &models.Caisse{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	p.Sync = true
	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse created success",
			"data":    p,
		},
	)
}

// Update data
func UpdateCaisse(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	type UpdateData struct {
		Name           string `gorm:"not null" json:"name"` // Nom de la caisse
		Signature      string `json:"signature"`            // Signature de la transaction
		PosUUID        string `json:"pos_uuid"`             // ID du point de vente
		CodeEntreprise uint64 `json:"code_entreprise"`      // ID de l'entreprise
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

	caisse := new(models.Caisse)

	db.Where("uuid = ?", uuid).First(&caisse)
	caisse.Name = updateData.Name
	caisse.Signature = updateData.Signature
	caisse.PosUUID = updateData.PosUUID
	caisse.CodeEntreprise = updateData.CodeEntreprise

	db.Save(&caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse updated success",
			"data":    caisse,
		},
	)

}

// Delete data
func DeleteCaisse(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	db := database.DB

	var caisse models.Caisse
	db.Where("uuid = ?", uuid).First(&caisse)
	if caisse.Name == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No caisse name found",
				"data":    nil,
			},
		)
	}

	db.Delete(&caisse)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "caisse deleted success",
			"data":    nil,
		},
	)
}
