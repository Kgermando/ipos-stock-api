package products

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"
	"github.com/xuri/excelize/v2"

	"github.com/gofiber/fiber/v2"
)

// Synchronisation Send data to Local
func GetDataSynchronisation(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01")
	var data []models.Product

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("products.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("products.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All products",
		"data":    data,
	})
}

// Paginate
func GetPaginatedProductEntreprise(c *fiber.Ctx) error {
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

	var dataList []models.Product

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Product{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
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

	var dataList []models.Product

	var totalRecords int64

	// Count total records matching the search query
	db.Model(&models.Product{}).
		Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("name ILIKE ? OR reference ILIKE ?", "%"+search+"%", "%"+search+"%").
		Count(&totalRecords)

	err = db.Where("entreprise_uuid = ?", entrepriseUUID).
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
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	var data []models.Product
	db.Where("entreprise_uuid = ?", entrepriseUUID).
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
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	search := c.Query("search", "")

	var data []models.Product
	db.Where("entreprise_uuid = ?", entrepriseUUID).
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

	// Vérifier si le produit existe déjà
	var existingProduct models.Product
	database.DB.Where("uuid = ?", p.UUID).First(&existingProduct)
	if existingProduct.UUID != "" {
		return c.Status(409).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Product avec cet UUID existe déjà",
				"data":    nil,
			},
		)
	}

	p.Sync = true

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
		Reference         string  `json:"reference"`
		Name              string  `json:"name"`
		Description       string  `json:"description"`
		UniteVente        string  `json:"unite_vente"`
		PrixVente         float64 `json:"prix_vente"`
		Tva               float64 `json:"tva"`
		PrixAchat         float64 `json:"prix_achat"`
		Remise            float64 `json:"remise"`                               // remise en pourcentage
		RemiseMinQuantity float64 `gorm:"default:0" json:"remise_min_quantity"` // la quantite minimale pour la remise
		Signature         string  `json:"signature"`
		PosUUID           string  `json:"pos_uuid"`
		EntrepriseUUID    string  `json:"entreprise_uuid"`
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
	product.PrixAchat = updateData.PrixAchat
	product.Remise = updateData.Remise
	product.RemiseMinQuantity = updateData.RemiseMinQuantity
	// product.Image = updateData.Image
	product.Signature = updateData.Signature
	product.PosUUID = updateData.PosUUID
	product.EntrepriseUUID = updateData.EntrepriseUUID

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

	product.Sync = true
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

// Upload products from Excel file
func UploadProductsFromExcel(c *fiber.Ctx) error {
	// Récupération des paramètres d'URL
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	signature := c.FormValue("signature")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "entreprise_uuid et pos_uuid sont requis",
			"data":    nil,
		})
	}

	// Récupération du fichier uploadé
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Aucun fichier trouvé",
			"data":    nil,
		})
	}

	// Vérification de l'extension du fichier
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Seuls les fichiers Excel (.xlsx, .xls) sont acceptés",
			"data":    nil,
		})
	}

	// Ouverture du fichier uploadé
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de l'ouverture du fichier",
			"data":    nil,
		})
	}
	defer src.Close()

	// Lecture du fichier Excel avec excelize
	f, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la lecture du fichier Excel",
			"data":    nil,
		})
	}
	defer f.Close()

	// Récupération de la première feuille
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la lecture des lignes du fichier Excel",
			"data":    nil,
		})
	}

	if len(rows) <= 1 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le fichier Excel doit contenir au moins une ligne de données en plus des en-têtes",
			"data":    nil,
		})
	}

	db := database.DB
	var createdProducts []models.Product
	var errors []string
	successCount := 0
	errorCount := 0

	// Boucle pour traiter chaque ligne (en sautant la première ligne d'en-têtes)
	for i, row := range rows {
		if i == 0 {
			// Sauter la ligne d'en-têtes
			continue
		}

		// Vérifier qu'il y a assez de colonnes
		if len(row) < 6 {
			errors = append(errors, fmt.Sprintf("Ligne %d: Données insuffisantes (minimum 6 colonnes requises)", i+1))
			errorCount++
			continue
		}

		// Extraction des données de chaque colonne
		// Format attendu: Reference | Name | Description | UniteVente | PrixVente | PrixAchat | [Tva] | [Stock] | [Remise] | [RemiseMinQuantity]
		reference := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		description := strings.TrimSpace(row[2])
		uniteVente := strings.TrimSpace(row[3])

		// Validation des champs obligatoires
		if reference == "" || name == "" || uniteVente == "" {
			errors = append(errors, fmt.Sprintf("Ligne %d: Reference, Name et UniteVente sont obligatoires", i+1))
			errorCount++
			continue
		}

		// Conversion des prix
		prixVenteStr := strings.TrimSpace(row[4])
		prixVente, err := strconv.ParseFloat(prixVenteStr, 64)
		if err != nil || prixVente <= 0 {
			errors = append(errors, fmt.Sprintf("Ligne %d: Prix de vente invalide (%s)", i+1, prixVenteStr))
			errorCount++
			continue
		}

		prixAchatStr := strings.TrimSpace(row[5])
		prixAchat, err := strconv.ParseFloat(prixAchatStr, 64)
		if err != nil {
			prixAchat = 0 // Valeur par défaut
		}

		// Champs optionnels
		var tva, stock, remise, remiseMinQuantity float64

		if len(row) > 6 && strings.TrimSpace(row[6]) != "" {
			tva, _ = strconv.ParseFloat(strings.TrimSpace(row[6]), 64)
		}

		if len(row) > 7 && strings.TrimSpace(row[7]) != "" {
			stock, _ = strconv.ParseFloat(strings.TrimSpace(row[7]), 64)
		}

		if len(row) > 8 && strings.TrimSpace(row[8]) != "" {
			remise, _ = strconv.ParseFloat(strings.TrimSpace(row[8]), 64)
		}

		if len(row) > 9 && strings.TrimSpace(row[9]) != "" {
			remiseMinQuantity, _ = strconv.ParseFloat(strings.TrimSpace(row[9]), 64)
		}

		// Vérification de l'unicité de la référence pour cette entreprise
		var existingProduct models.Product
		if err := db.Where("reference = ? AND entreprise_uuid = ?", reference, entrepriseUUID).First(&existingProduct).Error; err == nil {
			errors = append(errors, fmt.Sprintf("Ligne %d: La référence '%s' existe déjà", i+1, reference))
			errorCount++
			continue
		}

		// Création du produit
		product := models.Product{
			UUID:              utils.GenerateUUID(),
			Reference:         reference,
			Name:              name,
			Description:       description,
			UniteVente:        uniteVente,
			PrixVente:         prixVente,
			PrixAchat:         prixAchat,
			Tva:               tva,
			Stock:             stock,
			Remise:            remise,
			RemiseMinQuantity: remiseMinQuantity,
			PosUUID:           posUUID,
			EntrepriseUUID:    entrepriseUUID,
			Signature:         signature,
			Sync:              true,
		}

		// Sauvegarde en base de données
		if err := db.Create(&product).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Ligne %d: Erreur lors de la sauvegarde (%s)", i+1, err.Error()))
			errorCount++
		} else {
			createdProducts = append(createdProducts, product)
			successCount++
		}
	}

	// Préparation de la réponse
	response := fiber.Map{
		"status":        "success",
		"message":       fmt.Sprintf("Import terminé: %d produits créés, %d erreurs", successCount, errorCount),
		"success_count": successCount,
		"error_count":   errorCount,
		"data":          createdProducts,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	if errorCount > 0 && successCount == 0 {
		response["status"] = "error"
		return c.Status(400).JSON(response)
	}

	return c.JSON(response)
}

// Structure pour valider les données Excel
type ExcelProductData struct {
	Reference         string  `excel:"Reference" validate:"required" example:"REF001"`
	Name              string  `excel:"Name" validate:"required" example:"Coca Cola 33cl"`
	Description       string  `excel:"Description" example:"Boisson gazeuse rafraîchissante"`
	UniteVente        string  `excel:"UniteVente" validate:"required" example:"pièce"`
	PrixVente         float64 `excel:"PrixVente" validate:"required,gt=0" example:"1.50"`
	PrixAchat         float64 `excel:"PrixAchat" example:"1.00"`
	Tva               float64 `excel:"Tva" example:"20"`
	Stock             float64 `excel:"Stock" example:"100"`
	Remise            float64 `excel:"Remise" example:"5"`
	RemiseMinQuantity float64 `excel:"RemiseMinQuantity" example:"10"`
}

// Générer un fichier Excel modèle pour l'upload de produits
func GenerateProductExcelTemplate(c *fiber.Ctx) error {
	// Création d'un nouveau fichier Excel
	f := excelize.NewFile()
	defer f.Close()

	// Nom de la feuille principale
	sheetName := "Produits"
	f.SetSheetName("Sheet1", sheetName)

	// Définition des en-têtes avec descriptions
	headers := []string{
		"Reference", "Name", "Description", "UniteVente",
		"PrixVente", "PrixAchat", "Tva", "Stock",
		"Remise", "RemiseMinQuantity",
	}

	// Descriptions pour chaque colonne
	descriptions := []string{
		"Référence unique du produit (obligatoire)",
		"Nom du produit (obligatoire)",
		"Description détaillée du produit",
		"Unité de vente: pièce, kg, litre, etc. (obligatoire)",
		"Prix de vente unitaire en CDF (obligatoire, > 0)",
		"Prix d'achat unitaire en CDF (optionnel)",
		"TVA en pourcentage (optionnel, ex: 20)",
		"Stock initial disponible (optionnel)",
		"Remise en pourcentage (optionnel, ex: 5)",
		"Quantité minimale pour remise (optionnel)",
	}

	// Exemples de données
	examples := [][]interface{}{
		{"REF001", "Coca Cola 33cl", "Boisson gazeuse rafraîchissante", "pièce", 1.50, 1.00, 20, 100, 5, 10},
		{"REF002", "Pain complet", "Pain complet bio 500g", "pièce", 2.00, 1.20, 10, 50, 0, 0},
		{"REF003", "Lait frais", "Lait frais entier 1L", "litre", 1.80, 1.40, 10, 30, 10, 5},
		{"REF004", "Pommes", "Pommes Golden fraîches", "kg", 3.50, 2.50, 10, 25, 15, 3},
		{"REF005", "Café moulu", "Café arabica moulu 250g", "paquet", 4.99, 3.20, 20, 40, 0, 0},
	}

	// Style pour les en-têtes
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})

	// Style pour les descriptions
	descStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "#7F7F7F",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F2F2F2"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "top",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#CCCCCC", Style: 1},
			{Type: "top", Color: "#CCCCCC", Style: 1},
			{Type: "bottom", Color: "#CCCCCC", Style: 1},
			{Type: "right", Color: "#CCCCCC", Style: 1},
		},
	})

	// Style pour les exemples
	exampleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#CCCCCC", Style: 1},
			{Type: "top", Color: "#CCCCCC", Style: 1},
			{Type: "bottom", Color: "#CCCCCC", Style: 1},
			{Type: "right", Color: "#CCCCCC", Style: 1},
		},
	})

	// Écriture des en-têtes (ligne 1)
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Écriture des descriptions (ligne 2)
	for i, desc := range descriptions {
		cell := fmt.Sprintf("%c2", 'A'+i)
		f.SetCellValue(sheetName, cell, desc)
		f.SetCellStyle(sheetName, cell, cell, descStyle)
	}

	// Définir la hauteur des lignes
	f.SetRowHeight(sheetName, 1, 25)
	f.SetRowHeight(sheetName, 2, 60)

	// Écriture des exemples (à partir de la ligne 3)
	for rowIndex, example := range examples {
		for colIndex, value := range example {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, 3+rowIndex)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, exampleStyle)
		}
	}

	// Ajustement automatique de la largeur des colonnes
	columnWidths := []float64{12, 20, 35, 15, 12, 12, 8, 10, 10, 18}
	for i, width := range columnWidths {
		colName := fmt.Sprintf("%c", 'A'+i)
		f.SetColWidth(sheetName, colName, colName, width)
	}

	// Création d'une feuille d'instructions
	instructionSheet := "Instructions"
	f.NewSheet(instructionSheet)

	instructions := []string{
		"INSTRUCTIONS POUR L'IMPORT DE PRODUITS",
		"",
		"1. COLONNES OBLIGATOIRES :",
		"   • Reference : Référence unique du produit",
		"   • Name : Nom du produit",
		"   • UniteVente : Unité de vente (pièce, kg, litre, etc.)",
		"   • PrixVente : Prix de vente (doit être > 0)",
		"",
		"2. COLONNES OPTIONNELLES :",
		"   • Description : Description du produit",
		"   • PrixAchat : Prix d'achat",
		"   • Tva : TVA en pourcentage",
		"   • Stock : Stock initial",
		"   • Remise : Remise en pourcentage",
		"   • RemiseMinQuantity : Quantité min pour remise",
		"",
		"3. RÈGLES IMPORTANTES :",
		"   • La première ligne doit contenir les en-têtes",
		"   • La deuxième ligne contient les descriptions (à ignorer)",
		"   • Les données commencent à partir de la ligne 3",
		"   • Les références doivent être uniques",
		"   • Le prix de vente doit être supérieur à 0",
		"   • Les nombres décimaux utilisent le point (.)",
		"",
		"4. UNITÉS DE VENTE COURANTES :",
		"   • pièce, kg, litre, mètre, paquet, boîte, etc.",
		"",
		"5. EXEMPLES DE DONNÉES :",
		"   Voir la feuille 'Produits' pour des exemples complets",
		"",
		"6. APRÈS L'IMPORT :",
		"   • Un rapport détaillé sera généré",
		"   • Les erreurs seront listées ligne par ligne",
		"   • Les produits valides seront créés automatiquement",
	}

	// Style pour les instructions
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  16,
			Color: "#2F5597",
		},
	})

	headerInstStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#1F4788",
		},
	})

	normalStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
	})

	// Écriture des instructions
	for i, instruction := range instructions {
		cell := fmt.Sprintf("A%d", i+1)
		f.SetCellValue(instructionSheet, cell, instruction)

		if i == 0 {
			f.SetCellStyle(instructionSheet, cell, cell, titleStyle)
		} else if strings.Contains(instruction, ".") && strings.Contains(instruction, ":") && len(instruction) < 50 {
			f.SetCellStyle(instructionSheet, cell, cell, headerInstStyle)
		} else {
			f.SetCellStyle(instructionSheet, cell, cell, normalStyle)
		}
	}

	// Largeur de la colonne des instructions
	f.SetColWidth(instructionSheet, "A", "A", 80)

	// Définir la feuille active
	sheetIndex, _ := f.GetSheetIndex(sheetName)
	f.SetActiveSheet(sheetIndex)

	// Configuration de la réponse HTTP
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=modele-import-produits.xlsx")

	// Sauvegarde du fichier en mémoire et envoi
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la génération du fichier Excel",
			"data":    nil,
		})
	}

	return c.Send(buffer.Bytes())
}

// Obtenir des informations sur le format Excel attendu
func GetExcelFormatInfo(c *fiber.Ctx) error {
	formatInfo := map[string]interface{}{
		"status":  "success",
		"message": "Format Excel pour l'import de produits",
		"data": map[string]interface{}{
			"required_columns": []map[string]string{
				{"column": "A", "name": "Reference", "type": "string", "description": "Référence unique du produit"},
				{"column": "B", "name": "Name", "type": "string", "description": "Nom du produit"},
				{"column": "D", "name": "UniteVente", "type": "string", "description": "Unité de vente (pièce, kg, litre, etc.)"},
				{"column": "E", "name": "PrixVente", "type": "number", "description": "Prix de vente unitaire (> 0)"},
			},
			"optional_columns": []map[string]string{
				{"column": "C", "name": "Description", "type": "string", "description": "Description du produit"},
				{"column": "F", "name": "PrixAchat", "type": "number", "description": "Prix d'achat unitaire"},
				{"column": "G", "name": "Tva", "type": "number", "description": "TVA en pourcentage"},
				{"column": "H", "name": "Stock", "type": "number", "description": "Stock initial"},
				{"column": "I", "name": "Remise", "type": "number", "description": "Remise en pourcentage"},
				{"column": "J", "name": "RemiseMinQuantity", "type": "number", "description": "Quantité minimale pour remise"},
			},
			"rules": []string{
				"La première ligne doit contenir les en-têtes",
				"Les données commencent à partir de la ligne 2 ou 3",
				"Les références doivent être uniques pour l'entreprise",
				"Le prix de vente doit être supérieur à 0",
				"Formats de fichier acceptés: .xlsx, .xls",
				"Les nombres décimaux utilisent le point (.)",
			},
			"example_data": map[string]interface{}{
				"Reference":         "REF001",
				"Name":              "Coca Cola 33cl",
				"Description":       "Boisson gazeuse rafraîchissante",
				"UniteVente":        "pièce",
				"PrixVente":         1.50,
				"PrixAchat":         1.00,
				"Tva":               20,
				"Stock":             100,
				"Remise":            5,
				"RemiseMinQuantity": 10,
			},
			"unite_vente_examples": []string{
				"pièce", "kg", "litre", "mètre", "paquet", "boîte",
				"carton", "palette", "tonne", "ml", "cl", "gramme",
			},
		},
	}

	return c.JSON(formatInfo)
}
