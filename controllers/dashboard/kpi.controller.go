package dashboard

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
)

// Fonction utilitaire pour calculer les variations
func calculateVariationKPI(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	return ((current - previous) / previous) * 100
}

// GlobalKpis retourne les KPI globaux pour le dashboard
func GlobalKpis(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut")
	dateFin := c.Query("date_fin")

	// Validation des paramètres requis
	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	// Si les dates ne sont pas fournies, utiliser la date d'aujourd'hui
	if dateDebut == "" || dateFin == "" {
		now := time.Now()
		dateDebut = now.Format("2006-01-02")
		dateFin = now.Format("2006-01-02")
	}

	var totalRevenue, totalRevenueYesterday float64
	var totalCommandes, totalCommandesYesterday int64
	var totalProduits int64
	var averageOrderValue, averageOrderValueYesterday float64

	// Calculs pour aujourd'hui/période actuelle
	baseQuery := `
		SELECT 
			COALESCE(SUM(total_ttc), 0) as total_revenue,
			COUNT(uuid) as total_commandes
		FROM commandes 
		WHERE entreprise_uuid = ? AND DATE(created_at) BETWEEN ? AND ? AND deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID, dateDebut, dateFin}
	if posUUID != "" && posUUID != "null" {
		baseQuery += " AND pos_uuid = ?"
		args = append(args, posUUID)
	}

	err := db.Raw(baseQuery, args...).Row().Scan(&totalRevenue, &totalCommandes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de commandes",
			"error":   err.Error(),
		})
	}

	// Calculs pour hier/période précédente (pour les variations)
	dateDebutYesterday := dateDebut
	dateFinYesterday := dateFin
	if dateDebut == dateFin { // Si c'est la même date (aujourd'hui)
		yesterday := time.Now().AddDate(0, 0, -1)
		dateDebutYesterday = yesterday.Format("2006-01-02")
		dateFinYesterday = yesterday.Format("2006-01-02")
	}

	argsYesterday := []interface{}{entrepriseUUID, dateDebutYesterday, dateFinYesterday}
	if posUUID != "" && posUUID != "null" {
		argsYesterday = append(argsYesterday, posUUID)
	}

	err = db.Raw(baseQuery, argsYesterday...).Row().Scan(&totalRevenueYesterday, &totalCommandesYesterday)
	if err != nil {
		// Si erreur, on continue avec des valeurs par défaut
		totalRevenueYesterday = 0
		totalCommandesYesterday = 0
	}

	// Calcul du panier moyen
	if totalCommandes > 0 {
		averageOrderValue = totalRevenue / float64(totalCommandes)
	}
	if totalCommandesYesterday > 0 {
		averageOrderValueYesterday = totalRevenueYesterday / float64(totalCommandesYesterday)
	}

	// Récupération du nombre total de produits
	productQuery := "SELECT COUNT(uuid) FROM products WHERE entreprise_uuid = ? AND deleted_at IS NULL"
	productArgs := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		productQuery += " AND pos_uuid = ?"
		productArgs = append(productArgs, posUUID)
	}

	err = db.Raw(productQuery, productArgs...).Row().Scan(&totalProduits)
	if err != nil {
		totalProduits = 0
	}

	// Calcul des variations en pourcentage
	revenueVariation := calculateVariationKPI(totalRevenue, totalRevenueYesterday)
	commandesVariation := calculateVariationKPI(float64(totalCommandes), float64(totalCommandesYesterday))
	averageOrderVariation := calculateVariationKPI(averageOrderValue, averageOrderValueYesterday)

	globalKpis := map[string]interface{}{
		"totalRevenue":               totalRevenue,
		"totalRevenueVariation":      revenueVariation,
		"totalCommandes":             totalCommandes,
		"totalCommandesVariation":    commandesVariation,
		"totalProduits":              totalProduits,
		"totalProduitsVariation":     0, // Peut être calculé si nécessaire
		"averageOrderValue":          averageOrderValue,
		"averageOrderValueVariation": averageOrderVariation,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "KPI globaux récupérés avec succès",
		"data":    globalKpis,
	})
}

// GetEvolutionVente retourne l'évolution des ventes
func GetEvolutionVente(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut")
	dateFin := c.Query("date_fin")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	// Si les dates ne sont pas fournies, utiliser la semaine dernière
	if dateDebut == "" || dateFin == "" {
		now := time.Now()
		dateFin = now.Format("2006-01-02")
		dateDebut = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -7).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	var ventesAujourdhui, ventesHier, ventesSemaine, ventesMois float64

	baseQuery := `
		SELECT COALESCE(SUM(total_ttc), 0) 
		FROM commandes 
		WHERE entreprise_uuid = ? AND DATE(created_at) = ? AND deleted_at IS NULL
	`
	args := []interface{}{entrepriseUUID}

	if posUUID != "" && posUUID != "null" {
		baseQuery += " AND pos_uuid = ?"
	}

	// Ventes d'aujourd'hui
	queryArgs := append(args, today)
	if posUUID != "" && posUUID != "null" {
		queryArgs = append(queryArgs, posUUID)
	}
	db.Raw(baseQuery, queryArgs...).Row().Scan(&ventesAujourdhui)

	// Ventes d'hier
	queryArgs = append(args, yesterday)
	if posUUID != "" && posUUID != "null" {
		queryArgs = append(queryArgs, posUUID)
	}
	db.Raw(baseQuery, queryArgs...).Row().Scan(&ventesHier)

	// Ventes de la semaine
	weekQuery := `
		SELECT COALESCE(SUM(total_ttc), 0) 
		FROM commandes 
		WHERE entreprise_uuid = ? AND DATE(created_at) >= ? AND deleted_at IS NULL
	`
	queryArgs = append(args, weekStart)
	if posUUID != "" && posUUID != "null" {
		weekQuery += " AND pos_uuid = ?"
		queryArgs = append(queryArgs, posUUID)
	}
	db.Raw(weekQuery, queryArgs...).Row().Scan(&ventesSemaine)

	// Ventes du mois
	queryArgs = append(args, monthStart)
	if posUUID != "" && posUUID != "null" {
		queryArgs = append(queryArgs, posUUID)
	}
	db.Raw(weekQuery, queryArgs...).Row().Scan(&ventesMois)

	// Calcul des variations
	variationJour := calculateVariationKPI(ventesAujourdhui, ventesHier)
	variationSemaine := 0.0 // Peut être calculé avec des données historiques
	variationMois := 0.0    // Peut être calculé avec des données historiques

	evolutionData := map[string]interface{}{
		"ventesAujourdhui": ventesAujourdhui,
		"ventesHier":       ventesHier,
		"ventesSemaine":    ventesSemaine,
		"ventesMois":       ventesMois,
		"variationJour":    variationJour,
		"variationSemaine": variationSemaine,
		"variationMois":    variationMois,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Évolution des ventes récupérée avec succès",
		"data":    evolutionData,
	})
}

// GetPerformanceVente retourne les performances de vente
func GetPerformanceVente(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut")
	dateFin := c.Query("date_fin")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	if dateDebut == "" || dateFin == "" {
		now := time.Now()
		dateFin = now.Format("2006-01-02")
		dateDebut = now.AddDate(0, 0, -30).Format("2006-01-02")
	}

	var totalRevenue, totalCost float64

	// Requête pour calculer le chiffre d'affaires et le coût
	query := `
		SELECT 
			COALESCE(SUM(cl.quantity * p.prix_vente), 0) as total_revenue,
			COALESCE(SUM(cl.quantity * p.prix_achat), 0) as total_cost
		FROM commandes c
		LEFT JOIN commande_lines cl ON cl.commande_uuid = c.uuid
		LEFT JOIN products p ON cl.product_uuid = p.uuid
		WHERE c.entreprise_uuid = ? AND c.created_at BETWEEN ? AND ? AND c.deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID, dateDebut, dateFin}
	if posUUID != "" && posUUID != "null" {
		query += " AND c.pos_uuid = ?"
		args = append(args, posUUID)
	}

	err := db.Raw(query, args...).Row().Scan(&totalRevenue, &totalCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de performance",
			"error":   err.Error(),
		})
	}

	// Calculs des marges
	margeBrute := totalRevenue - totalCost
	margeGlobale := 0.0
	if totalRevenue > 0 {
		margeGlobale = (margeBrute / totalRevenue) * 100
	}

	performanceData := map[string]interface{}{
		"totalRevenue": totalRevenue,
		"totalCost":    totalCost,
		"margeBrute":   margeBrute,
		"margeGlobale": margeGlobale,
		"performance":  margeGlobale, // Pourcentage de performance basé sur la marge
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Performance de vente récupérée avec succès",
		"data":    performanceData,
	})
}

// GetBestSellingProduct retourne les produits les plus vendus
func GetBestSellingProduct(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut")
	dateFin := c.Query("date_fin")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	if dateDebut == "" || dateFin == "" {
		now := time.Now()
		dateFin = now.Format("2006-01-02")
		dateDebut = now.AddDate(0, 0, -30).Format("2006-01-02")
	}

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	var products []models.TopProduct

	query := `
		SELECT 
			p.uuid,
			p.name,
			COALESCE(SUM(cl.quantity), 0) as quantite,
			COALESCE(SUM(cl.quantity * p.prix_vente), 0) as valeur,
			COALESCE(p.stock, 0) as stock,
			0 as variation
		FROM commandes c
		LEFT JOIN commande_lines cl ON cl.commande_uuid = c.uuid
		LEFT JOIN products p ON cl.product_uuid = p.uuid
		WHERE c.entreprise_uuid = ? AND c.created_at BETWEEN ? AND ? AND c.deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID, dateDebut, dateFin}
	if posUUID != "" && posUUID != "null" {
		query += " AND c.pos_uuid = ?"
		args = append(args, posUUID)
	}

	query += " GROUP BY p.uuid, p.name, p.stock ORDER BY quantite DESC LIMIT ?"
	args = append(args, limit)

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des produits les plus vendus",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var product models.TopProduct
		err := rows.Scan(&product.UUID, &product.Name, &product.Quantite, &product.Valeur, &product.Stock, &product.Variation)
		if err != nil {
			continue
		}
		products = append(products, product)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Produits les plus vendus récupérés avec succès",
		"data":    products,
	})
}

// GetStockKpis retourne les KPI liés au stock
func GetStockKpis(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	var totalStock, stockValeur, stockEndommage, stockRestitution float64
	var stockAlertes int64

	// Requête pour les données de stock
	query := `
		SELECT 
			COALESCE(SUM(stock), 0) as total_stock,
			COALESCE(SUM(stock * prix_vente), 0) as stock_valeur,
			COALESCE(SUM(stock_endommage * prix_achat), 0) as stock_endommage,
			COALESCE(SUM(restitution * prix_achat), 0) as stock_restitution,
			COUNT(CASE WHEN stock < 10 THEN 1 END) as stock_alertes
		FROM products 
		WHERE entreprise_uuid = ? AND deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		query += " AND pos_uuid = ?"
		args = append(args, posUUID)
	}

	err := db.Raw(query, args...).Row().Scan(&totalStock, &stockValeur, &stockEndommage, &stockRestitution, &stockAlertes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de stock",
			"error":   err.Error(),
		})
	}

	// Simulation d'un taux de rotation (à améliorer avec des données réelles)
	tauxRotation := 2.5 + (float64(stockAlertes) * 0.1)
	stockVariation := 0.0 // Peut être calculé avec des données historiques

	stockKpis := map[string]interface{}{
		"totalStock":       totalStock,
		"stockValeur":      stockValeur,
		"stockEndommage":   stockEndommage,
		"stockAlertes":     stockAlertes,
		"stockRestitution": stockRestitution,
		"tauxRotation":     tauxRotation,
		"stockVariation":   stockVariation,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "KPI stock récupérés avec succès",
		"data":    stockKpis,
	})
}

// GetStockFaible retourne les produits avec un stock faible
func GetStockFaible(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	var products []models.TopProduct

	query := `
		SELECT 
			uuid,
			name,
			COALESCE(stock, 0) as quantite,
			COALESCE(stock * prix_vente, 0) as valeur,
			COALESCE(stock, 0) as stock,
			0 as variation
		FROM products
		WHERE entreprise_uuid = ? AND stock < 10 AND deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		query += " AND pos_uuid = ?"
		args = append(args, posUUID)
	}

	query += " ORDER BY stock ASC LIMIT ?"
	args = append(args, limit)

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des produits à stock faible",
			"error":   err.Error(),
		})
	}
	defer rows.Close()
	for rows.Next() {
		var product models.TopProduct
		err := rows.Scan(&product.UUID, &product.Name, &product.Quantite, &product.Valeur, &product.Stock, &product.Variation)
		if err != nil {
			continue
		}
		products = append(products, product)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Produits à stock faible récupérés avec succès",
		"data":    products,
	})
}

// ===== FONCTIONS LEGACY POUR COMPATIBILITÉ =====

// GlobalKpiSummary - fonction legacy (redirige vers GlobalKpis)
func GlobalKpiSummary(c *fiber.Ctx) error {
	return GlobalKpis(c)
}

// EvolutionVente - fonction legacy
func EvolutionVente(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut")
	dateFin := c.Query("date_fin")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}
	if dateDebut == "" || dateFin == "" {
		now := time.Now()
		dateFin = now.Format("2006-01-02")
		dateDebut = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	var sales []models.Sale

	query := `
		SELECT
			DATE(created_at) AS date, 
			COUNT(uuid) AS commande,
			COALESCE(SUM(total_ttc), 0) AS vente,
			COALESCE(SUM(total_tva), 0) AS tva
		FROM commandes
		WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID, dateDebut, dateFin}
	if posUUID != "" && posUUID != "null" {
		query += " AND pos_uuid = ?"
		args = append(args, posUUID)
	}

	query += " GROUP BY DATE(created_at) ORDER BY DATE(created_at)"

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données d'évolution",
			"error":   err.Error(),
		})
	}
	defer rows.Close()
	for rows.Next() {
		var sale models.Sale
		err := rows.Scan(&sale.Date, &sale.Commande, &sale.Vente, &sale.Tva)
		if err != nil {
			continue
		}
		sales = append(sales, sale)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Evolution de la vente",
		"data":    sales,
	})
}

// PerformanceVente - fonction legacy (redirige vers GetPerformanceVente)
func PerformanceVente(c *fiber.Ctx) error {
	return GetPerformanceVente(c)
}

// BestSellingProduct - fonction legacy (redirige vers GetBestSellingProduct)
func BestSellingProduct(c *fiber.Ctx) error {
	return GetBestSellingProduct(c)
}

// StockFaible - fonction legacy (redirige vers GetStockFaible)
func StockFaible(c *fiber.Ctx) error {
	return GetStockFaible(c)
}

// SetupStockChart - fonction pour les données de graphique de stock
func SetupStockChart(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	var stockDisponible, stockEndommage, stockRestitution float64

	query := `
		SELECT 
			COALESCE(SUM(stock), 0) as stock_disponible,
			COALESCE(SUM(stock_endommage), 0) as stock_endommage,
			COALESCE(SUM(restitution), 0) as stock_restitution
		FROM products 
		WHERE entreprise_uuid = ? AND deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		query += " AND pos_uuid = ?"
		args = append(args, posUUID)
	}

	err := db.Raw(query, args...).Row().Scan(&stockDisponible, &stockEndommage, &stockRestitution)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de stock pour le graphique",
			"error":   err.Error(),
		})
	}

	chartData := map[string]interface{}{
		"series": []float64{stockDisponible, stockEndommage, stockRestitution},
		"labels": []string{"Stock Disponible", "Stock Endommagé", "Stock Restitution"},
		"colors": []string{"#28a745", "#dc3545", "#ffc107"},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Données du graphique de stock récupérées avec succès", "data": chartData,
	})
}
