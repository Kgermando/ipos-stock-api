package caisses

import (
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

	var data []models.Caisse

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("caisses.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("caisses.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All Caisses",
		"data":    data,
	})
}

// Get All data
func GetTotalAllCaisses(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")

	var dataList []models.CaisseItem
	db.Where("entreprise_uuid = ?", entrepriseUUID).Find(&dataList)

	var total float64 = 0
	var totalEntree float64 = 0
	var totalSortie float64 = 0
	var totalMontantDebut float64 = 0
	var solde float64 = 0
	var pourcent float64 = 0

	for _, item := range dataList {
		if item.TypeTransaction == "Entrée" {
			totalEntree += item.Montant
		}
		if item.TypeTransaction == "Sortie" {
			totalSortie += item.Montant
		}
		if item.TypeTransaction == "MontantDebut" {
			totalMontantDebut += item.Montant
		}
	}

	total = totalEntree + totalSortie
	solde = totalEntree - totalSortie
	pourcent = solde * 100 / (totalEntree + totalSortie)

	response := map[string]interface{}{
		"total":         total,
		"totalentree":   totalEntree,
		"totalsortie":   totalSortie,
		"solde":         solde,
		"pourcent":      pourcent,
		"montant_debut": totalMontantDebut,
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
	entrepriseUUID := c.Params("entreprise_uuid")

	var data []models.Caisse
	db.Where("entreprise_uuid = ?", entrepriseUUID).
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
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUId := c.Params("pos_uuid")

	var data []models.Caisse
	db.Where("entreprise_uuid = ?", entrepriseUUID).
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
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUId := c.Params("pos_uuid")

	search := c.Query("search", "")

	var data []models.Caisse
	db.Where("entreprise_uuid = ?", entrepriseUUID).
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

	// Vérifier si la caisse existe déjà
	var existingCaisse models.Caisse
	database.DB.Where("uuid = ?", p.UUID).First(&existingCaisse)
	if existingCaisse.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Caisse avec cet UUID existe déjà",
				"data":    nil,
			},
		)
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
		Name string `json:"name"` // Nom de la caisse
		// Entree         float64 `json:"entree"`          // Montant d'entrée
		// Sortie         float64 `json:"sortie"`          // Montant de sortie
		Signature      string `json:"signature"`       // Signature de la transaction
		PosUUID        string `json:"pos_uuid"`        // ID du point de vente
		EntrepriseUUID string `json:"entreprise_uuid"` // ID de l'entreprise
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
	// caisse.Entree = updateData.Entree
	// caisse.Sortie = updateData.Sortie
	caisse.Signature = updateData.Signature
	caisse.PosUUID = updateData.PosUUID
	caisse.EntrepriseUUID = updateData.EntrepriseUUID

	caisse.Sync = true
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
