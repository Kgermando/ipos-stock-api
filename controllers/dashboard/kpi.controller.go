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

// Fonction utilitaire pour calculer la période précédente
func calculatePreviousPeriod(dateDebut, dateFin string) (string, string) {
	debut, err := time.Parse("2006-01-02", dateDebut)
	if err != nil {
		return dateDebut, dateFin
	}

	fin, err := time.Parse("2006-01-02", dateFin)
	if err != nil {
		return dateDebut, dateFin
	}

	// Calculer la durée de la période
	duree := fin.Sub(debut)

	// Calculer la période précédente
	debutPrevious := debut.Add(-duree - 24*time.Hour)
	finPrevious := debut.Add(-24 * time.Hour)

	return debutPrevious.Format("2006-01-02"), finPrevious.Format("2006-01-02")
}

// Fonction utilitaire pour calculer la satisfaction client basée sur des métriques réelles
func calculateSatisfactionClient(margeGlobale, efficaciteStock, performancePOS float64) float64 {
	// Calcul d'un score de satisfaction basé sur les performances réelles
	// Score sur 100 basé sur plusieurs facteurs mesurables

	score := 0.0

	// Score basé sur la marge (0-40 points)
	if margeGlobale >= 30 {
		score += 40
	} else if margeGlobale >= 20 {
		score += 30
	} else if margeGlobale >= 10 {
		score += 20
	} else if margeGlobale >= 5 {
		score += 10
	}

	// Score basé sur l'efficacité du stock (0-30 points)
	if efficaciteStock >= 95 {
		score += 30
	} else if efficaciteStock >= 85 {
		score += 25
	} else if efficaciteStock >= 70 {
		score += 20
	} else if efficaciteStock >= 50 {
		score += 15
	} else if efficaciteStock >= 30 {
		score += 10
	}

	// Score basé sur la performance POS (0-30 points)
	if performancePOS >= 1000 {
		score += 30
	} else if performancePOS >= 500 {
		score += 25
	} else if performancePOS >= 200 {
		score += 20
	} else if performancePOS >= 100 {
		score += 15
	} else if performancePOS >= 50 {
		score += 10
	}

	return score
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
	// Requête séparée pour le revenu total (éviter les doublons)
	revenueQuery := `
		SELECT 
			COALESCE(SUM(total_ttc), 0) as total_revenue,
			COUNT(uuid) as total_commandes
		FROM commandes 
		WHERE entreprise_uuid = ? AND DATE(created_at) BETWEEN ? AND ? AND deleted_at IS NULL
	`
	// Requête pour la marge bénéficiaire basée sur les lignes de commande
	margeQuery := `
		SELECT 
			COALESCE(SUM(cl.quantity * (p.prix_vente - p.prix_achat)), 0) as margeBeneficiaire		
		FROM commandes c
		INNER JOIN commande_lines cl ON cl.commande_uuid = c.uuid
		INNER JOIN products p ON cl.product_uuid = p.uuid
		WHERE c.entreprise_uuid = ? AND DATE(c.created_at) BETWEEN ? AND ? AND c.deleted_at IS NULL
	`
	args := []interface{}{entrepriseUUID, dateDebut, dateFin}
	if posUUID != "" && posUUID != "null" {
		revenueQuery += " AND pos_uuid = ?"
		margeQuery += " AND c.pos_uuid = ?"
		args = append(args, posUUID)
	}

	// Exécuter la requête pour le revenu et le nombre de commandes
	err := db.Raw(revenueQuery, args...).Row().Scan(&totalRevenue, &totalCommandes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de revenu",
			"error":   err.Error(),
		})
	}
	// Exécuter la requête pour la marge bénéficiaire
	var margeBeneficiaireCalculee float64
	err = db.Raw(margeQuery, args...).Row().Scan(&margeBeneficiaireCalculee)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de marge",
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
	// Requêtes pour les données d'hier
	err = db.Raw(revenueQuery, argsYesterday...).Row().Scan(&totalRevenueYesterday, &totalCommandesYesterday)
	if err != nil {
		// Si erreur, on continue avec des valeurs par défaut
		totalRevenueYesterday = 0
		totalCommandesYesterday = 0
	}

	var margeBeneficiaireYesterday float64
	err = db.Raw(margeQuery, argsYesterday...).Row().Scan(&margeBeneficiaireYesterday)
	if err != nil {
		// Si erreur, on continue avec des valeurs par défaut
		margeBeneficiaireYesterday = 0
	}

	// Calcul du panier moyen
	if totalCommandes > 0 {
		averageOrderValue = totalRevenue / float64(totalCommandes)
	}
	if totalCommandesYesterday > 0 {
		averageOrderValueYesterday = totalRevenueYesterday / float64(totalCommandesYesterday)
	}
	// Calcul de la marge bénéficiaire - utilisation des valeurs calculées par les requêtes
	margeBeneficiaire := margeBeneficiaireCalculee
	margeBeneficiaireVariation := calculateVariationKPI(margeBeneficiaire, margeBeneficiaireYesterday)

	// Calcul du pourcentage de marge
	pourcentageMarge := 0.0
	if totalRevenue > 0 {
		pourcentageMarge = (margeBeneficiaire / totalRevenue) * 100
	}
	// Récupération du nombre total de produits (période actuelle)
	productQuery := "SELECT COUNT(uuid) FROM products WHERE entreprise_uuid = ? AND DATE(created_at) BETWEEN ? AND ? AND deleted_at IS NULL"
	productArgs := []interface{}{entrepriseUUID, dateDebut, dateFin}
	if posUUID != "" && posUUID != "null" {
		productQuery += " AND pos_uuid = ?"
		productArgs = append(productArgs, posUUID)
	}

	err = db.Raw(productQuery, productArgs...).Row().Scan(&totalProduits)
	if err != nil {
		totalProduits = 0
	}

	// Récupération du nombre total de produits (période précédente)
	var totalProduitsYesterday int64
	productArgsYesterday := []interface{}{entrepriseUUID, dateDebutYesterday, dateFinYesterday}
	if posUUID != "" && posUUID != "null" {
		productArgsYesterday = append(productArgsYesterday, posUUID)
	}

	err = db.Raw(productQuery, productArgsYesterday...).Row().Scan(&totalProduitsYesterday)
	if err != nil {
		totalProduitsYesterday = 0
	} // Calcul des variations en pourcentage
	revenueVariation := calculateVariationKPI(totalRevenue, totalRevenueYesterday)
	commandesVariation := calculateVariationKPI(float64(totalCommandes), float64(totalCommandesYesterday))
	averageOrderVariation := calculateVariationKPI(averageOrderValue, averageOrderValueYesterday)
	totalProduitsVariation := calculateVariationKPI(float64(totalProduits), float64(totalProduitsYesterday))

	// Calcul de la valeur totale du stock et de la marge bénéficiaire totale potentielle (sans filtre de dates)
	var valeurTotaleStock, margeBeneficiaireTotale float64

	stockValueQuery := `
		SELECT 
			COALESCE(SUM(stocks.quantity - stock_endommages.quantity - commande_lines.quantity * prix_vente), 0) as valeur_totale_stock,
			COALESCE(SUM((stocks.quantity - stock_endommages.quantity - commande_lines.quantity) * (prix_vente - prix_achat)), 0) as marge_beneficiaire_totale
		FROM products 
		LEFT JOIN stocks ON products.uuid = stocks.product_uuid
		LEFT JOIN stock_endommages ON products.uuid = stock_endommages.product_uuid
		LEFT JOIN commande_lines ON products.uuid = commande_lines.product_uuid
		WHERE products.entreprise_uuid = ? AND products.deleted_at IS NULL
	`

	stockValueArgs := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		stockValueQuery += " AND products.pos_uuid = ?"
		stockValueArgs = append(stockValueArgs, posUUID)
	}

	err = db.Raw(stockValueQuery, stockValueArgs...).Row().Scan(&valeurTotaleStock, &margeBeneficiaireTotale)
	if err != nil {
		// Si erreur, on continue avec des valeurs par défaut
		valeurTotaleStock = 0
		margeBeneficiaireTotale = 0
	}

	globalKpis := map[string]interface{}{
		"totalRevenue":               totalRevenue,
		"totalRevenueVariation":      revenueVariation,
		"totalCommandes":             totalCommandes,
		"totalCommandesVariation":    commandesVariation,
		"totalProduits":              totalProduits,
		"totalProduitsVariation":     totalProduitsVariation,
		"averageOrderValue":          averageOrderValue,
		"averageOrderValueVariation": averageOrderVariation,
		"margeBeneficiaire":          margeBeneficiaire,
		"margeBeneficiaireVariation": margeBeneficiaireVariation,
		"pourcentageMarge":           pourcentageMarge,
		"valeurTotaleStock":          valeurTotaleStock,
		"margeBeneficiaireTotale":    margeBeneficiaireTotale,
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
	// Calcul des variations avec des données historiques réelles

	// Variation jour (déjà calculée)
	variationJour := calculateVariationKPI(ventesAujourdhui, ventesHier)

	// Variation semaine - comparer avec la semaine précédente
	var ventesSemainePrecedente float64
	weekStartPrevious := now.AddDate(0, 0, -14).Format("2006-01-02")
	weekEndPrevious := now.AddDate(0, 0, -8).Format("2006-01-02")

	weekPreviousQuery := `
		SELECT COALESCE(SUM(total_ttc), 0) 
		FROM commandes 
		WHERE entreprise_uuid = ? AND DATE(created_at) BETWEEN ? AND ? AND deleted_at IS NULL
	`
	queryArgsPrevious := append(args, weekStartPrevious, weekEndPrevious)
	if posUUID != "" && posUUID != "null" {
		weekPreviousQuery += " AND pos_uuid = ?"
		queryArgsPrevious = append(queryArgsPrevious, posUUID)
	}
	db.Raw(weekPreviousQuery, queryArgsPrevious...).Row().Scan(&ventesSemainePrecedente)

	variationSemaine := calculateVariationKPI(ventesSemaine, ventesSemainePrecedente)

	// Variation mois - comparer avec le mois précédent
	var ventesMoisPrecedent float64
	lastMonth := now.AddDate(0, -1, 0)
	monthStartPrevious := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location()).Format("2006-01-02")
	monthEndPrevious := time.Date(lastMonth.Year(), lastMonth.Month()+1, 0, 0, 0, 0, 0, lastMonth.Location()).Format("2006-01-02")

	queryArgsMonthPrevious := append(args, monthStartPrevious, monthEndPrevious)
	if posUUID != "" && posUUID != "null" {
		queryArgsMonthPrevious = append(queryArgsMonthPrevious, posUUID)
	}
	db.Raw(weekPreviousQuery, queryArgsMonthPrevious...).Row().Scan(&ventesMoisPrecedent)

	variationMois := calculateVariationKPI(ventesMois, ventesMoisPrecedent)

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

	var totalRevenue, totalCost, totalTva float64
	var totalCommandes int64

	// Requête pour calculer le chiffre d'affaires, le coût et la TVA
	query := `
		SELECT 
			COALESCE(SUM(cl.quantity * p.prix_vente), 0) as total_revenue,
			COALESCE(SUM(cl.quantity * p.prix_achat), 0) as total_cost,
			COALESCE(SUM(c.total_tva), 0) as total_tva,
			COUNT(DISTINCT c.uuid) as total_commandes
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

	err := db.Raw(query, args...).Row().Scan(&totalRevenue, &totalCost, &totalTva, &totalCommandes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Erreur lors de la récupération des données de performance",
			"error":   err.Error(),
		})
	}

	// Calcul des données de la période précédente pour la croissance
	dateDebutPrevious, dateFinPrevious := calculatePreviousPeriod(dateDebut, dateFin)
	var totalRevenuePrevious float64

	argsPrevious := []interface{}{entrepriseUUID, dateDebutPrevious, dateFinPrevious}
	if posUUID != "" && posUUID != "null" {
		argsPrevious = append(argsPrevious, posUUID)
	}

	queryPrevious := `
		SELECT COALESCE(SUM(cl.quantity * p.prix_vente), 0) as total_revenue
		FROM commandes c
		LEFT JOIN commande_lines cl ON cl.commande_uuid = c.uuid
		LEFT JOIN products p ON cl.product_uuid = p.uuid
		WHERE c.entreprise_uuid = ? AND c.created_at BETWEEN ? AND ? AND c.deleted_at IS NULL
	`

	if posUUID != "" && posUUID != "null" {
		queryPrevious += " AND c.pos_uuid = ?"
	}

	db.Raw(queryPrevious, argsPrevious...).Row().Scan(&totalRevenuePrevious)
	// Calcul des données de stock pour l'efficacité
	var totalStock, stockEndommage, stockRestitution float64
	stockQuery := `
		SELECT 
			COALESCE(SUM(stocks.quantity - stock_endommages.quantity - commande_lines.quantity), 0) as total_stock,
			COALESCE(SUM(stock_endommages.quantity), 0) as stock_endommage,
			COALESCE(SUM(restitutions.quantity), 0) as stock_restitution
		FROM products 
		LEFT JOIN stocks ON products.uuid = stocks.product_uuid
		LEFT JOIN stock_endommages ON products.uuid = stock_endommages.product_uuid
		LEFT JOIN restitutions ON products.uuid = restitutions.product_uuid
		LEFT JOIN commande_lines ON products.uuid = commande_lines.product_uuid
		WHERE products.entreprise_uuid = ? AND products.deleted_at IS NULL
	`
	stockArgs := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		stockQuery += " AND products.pos_uuid = ?"
		stockArgs = append(stockArgs, posUUID)
	}

	db.Raw(stockQuery, stockArgs...).Row().Scan(&totalStock, &stockEndommage, &stockRestitution)

	// Calculs des métriques
	margeBrute := totalRevenue - totalCost
	margeGlobale := 0.0
	if totalRevenue > 0 {
		margeGlobale = (margeBrute / totalRevenue) * 100
	}

	// Croissance des ventes
	croissanceVentes := calculateVariationKPI(totalRevenue, totalRevenuePrevious)

	// Performance POS (basée sur le chiffre d'affaires par commande)
	performancePOS := 0.0
	if totalCommandes > 0 {
		performancePOS = totalRevenue / float64(totalCommandes)
	}
	// Efficacité stock (pourcentage de stock utilisable)
	efficaciteStock := 100.0
	if totalStock > 0 {
		stockUtilisable := totalStock - stockEndommage - stockRestitution
		efficaciteStock = (stockUtilisable / totalStock) * 100
		// S'assurer que l'efficacité ne soit pas négative
		if efficaciteStock < 0 {
			efficaciteStock = 0
		}
	}
	// Satisfaction client (calcul basé sur les métriques de performance réelles)
	satisfactionClient := calculateSatisfactionClient(margeGlobale, efficaciteStock, performancePOS)

	performanceData := map[string]interface{}{
		"totalRevenue":       totalRevenue,
		"totalCost":          totalCost,
		"margeGlobale":       margeGlobale,
		"margeBrute":         margeBrute,
		"croissanceVentes":   croissanceVentes,
		"performancePOS":     performancePOS,
		"efficaciteStock":    efficaciteStock,
		"satisfactionClient": satisfactionClient,
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
			COALESCE(SUM(stocks.quantity - stock_endommages.quantity - commande_lines.quantity), 0) as total_stock,
			COALESCE(SUM(stocks.quantity * prix_vente), 0) as stock_valeur,
			COALESCE(SUM(stock_endommages.quantity * prix_achat), 0) as stock_endommage,
			COALESCE(SUM(restitutions.quantity * prix_achat), 0) as stock_restitution,
			COUNT(CASE WHEN stocks.quantity < 10 THEN 1 END) as stock_alertes
		FROM products 
		LEFT JOIN stocks ON products.uuid = stocks.product_uuid
		LEFT JOIN stock_endommages ON products.uuid = stock_endommages.product_uuid
		LEFT JOIN restitutions ON products.uuid = restitutions.product_uuid
		LEFT JOIN commande_lines ON products.uuid = commande_lines.product_uuid
		WHERE products.entreprise_uuid = ? AND products.deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		query += " AND products.pos_uuid = ?"
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
	// Calcul du taux de rotation du stock (basé sur les données réelles)
	// Formule: Taux de rotation = Coût des marchandises vendues / Stock moyen
	// Période de calcul: 30 derniers jours
	now := time.Now()
	dateDebut30j := now.AddDate(0, 0, -30).Format("2006-01-02")
	dateFin30j := now.Format("2006-01-02")

	var coutMarchandisesVendues float64
	rotationQuery := `
		SELECT COALESCE(SUM(cl.quantity * p.prix_achat), 0) as cout_marchandises
		FROM commandes c
		LEFT JOIN commande_lines cl ON cl.commande_uuid = c.uuid
		LEFT JOIN products p ON cl.product_uuid = p.uuid
		WHERE c.entreprise_uuid = ? AND c.created_at BETWEEN ? AND ? AND c.deleted_at IS NULL
	`

	rotationArgs := []interface{}{entrepriseUUID, dateDebut30j, dateFin30j}
	if posUUID != "" && posUUID != "null" {
		rotationQuery += " AND c.pos_uuid = ?"
		rotationArgs = append(rotationArgs, posUUID)
	}

	db.Raw(rotationQuery, rotationArgs...).Row().Scan(&coutMarchandisesVendues)

	// Calcul du stock moyen (stock actuel comme approximation)
	stockMoyen := totalStock

	// Calcul du taux de rotation
	tauxRotation := 0.0
	if stockMoyen > 0 {
		// Taux de rotation pour 30 jours, on l'annualise en multipliant par 12
		tauxRotation = (coutMarchandisesVendues / stockMoyen) * 12
	}

	// Calcul de la variation du stock par rapport à la semaine précédente
	var totalStockPrevious float64
	weekAgoPrevious := now.AddDate(0, 0, -7).Format("2006-01-02")

	stockPreviousQuery := `
		SELECT COALESCE(SUM(stocks.quantity - stock_endommages.quantity - commande_lines.quantity), 0) as total_stock
		FROM products 
		LEFT JOIN stocks ON products.uuid = stocks.product_uuid
		LEFT JOIN stock_endommages ON products.uuid = stock_endommages.product_uuid
		LEFT JOIN commande_lines ON products.uuid = commande_lines.product_uuid
		WHERE products.entreprise_uuid = ? AND products.deleted_at IS NULL AND products.created_at <= ?
	`
	stockPreviousArgs := []interface{}{entrepriseUUID, weekAgoPrevious}
	if posUUID != "" && posUUID != "null" {
		stockPreviousQuery += " AND products.pos_uuid = ?"
		stockPreviousArgs = append(stockPreviousArgs, posUUID)
	}

	db.Raw(stockPreviousQuery, stockPreviousArgs...).Row().Scan(&totalStockPrevious)
	stockVariation := calculateVariationKPI(totalStock, totalStockPrevious)

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

func GetAlertsKpis(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Le paramètre entreprise_uuid est requis",
		})
	}

	// Calcul des alertes de stock (produits avec stock < 10)
	var alertesStock int64
	stockQuery := `
		SELECT COUNT(uuid) 
		FROM products 
		WHERE entreprise_uuid = ? AND stock < 10 AND deleted_at IS NULL
	`
	stockArgs := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		stockQuery += " AND pos_uuid = ?"
		stockArgs = append(stockArgs, posUUID)
	}
	db.Raw(stockQuery, stockArgs...).Row().Scan(&alertesStock)

	// Calcul des alertes d'expiration (produits avec stock > 0 mais pas de vente récente)
	var alertesExpiration int64
	now := time.Now()
	dateLimit := now.AddDate(0, 0, -30).Format("2006-01-02") // 30 jours sans vente

	expirationQuery := `
		SELECT COUNT(DISTINCT p.uuid)
		FROM products p
		LEFT JOIN commande_lines cl ON cl.product_uuid = p.uuid
		LEFT JOIN commandes c ON c.uuid = cl.commande_uuid AND c.created_at >= ? AND c.deleted_at IS NULL
		WHERE p.entreprise_uuid = ? AND p.stock > 0 AND p.deleted_at IS NULL
		AND c.uuid IS NULL
	`
	expirationArgs := []interface{}{dateLimit, entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		expirationQuery += " AND p.pos_uuid = ?"
		expirationArgs = append(expirationArgs, posUUID)
	}
	db.Raw(expirationQuery, expirationArgs...).Row().Scan(&alertesExpiration)
	// Calcul des alertes de ventes (baisse significative par rapport à la période précédente)
	var alertesVentes int64
	today := now.Format("2006-01-02")
	weekAgo := now.AddDate(0, 0, -7).Format("2006-01-02")

	var ventesToday, ventesWeekAgo float64

	ventesQuery := `
		SELECT COALESCE(SUM(total_ttc), 0)
		FROM commandes 
		WHERE entreprise_uuid = ? AND DATE(created_at) = ? AND deleted_at IS NULL
	`

	ventesTodayArgs := []interface{}{entrepriseUUID, today}
	ventesWeekArgs := []interface{}{entrepriseUUID, weekAgo}

	if posUUID != "" && posUUID != "null" {
		ventesQuery += " AND pos_uuid = ?"
		ventesTodayArgs = append(ventesTodayArgs, posUUID)
		ventesWeekArgs = append(ventesWeekArgs, posUUID)
	}

	db.Raw(ventesQuery, ventesTodayArgs...).Row().Scan(&ventesToday)
	db.Raw(ventesQuery, ventesWeekArgs...).Row().Scan(&ventesWeekAgo)

	// Alerte si baisse de plus de 50% par rapport à la semaine dernière
	if ventesWeekAgo > 0 && ventesToday < (ventesWeekAgo*0.5) {
		alertesVentes = 1
	}

	// Calcul des alertes de performance (marge faible, rotation lente, etc.)
	var alertesPerformance int64

	// Vérifier la marge globale
	var totalRevenue, totalCost float64
	margeQuery := `
		SELECT 
			COALESCE(SUM(cl.quantity * p.prix_vente), 0) as total_revenue,
			COALESCE(SUM(cl.quantity * p.prix_achat), 0) as total_cost
		FROM commandes c
		LEFT JOIN commande_lines cl ON cl.commande_uuid = c.uuid
		LEFT JOIN products p ON cl.product_uuid = p.uuid
		WHERE c.entreprise_uuid = ? AND c.created_at >= ? AND c.deleted_at IS NULL
	`

	margeArgs := []interface{}{entrepriseUUID, weekAgo}
	if posUUID != "" && posUUID != "null" {
		margeQuery += " AND c.pos_uuid = ?"
		margeArgs = append(margeArgs, posUUID)
	}

	db.Raw(margeQuery, margeArgs...).Row().Scan(&totalRevenue, &totalCost)

	// Alerte si marge < 10%
	if totalRevenue > 0 {
		marge := ((totalRevenue - totalCost) / totalRevenue) * 100
		if marge < 10 {
			alertesPerformance++
		}
	}

	// Vérifier le taux de rotation du stock
	var stockMoyen float64
	stockMoyenQuery := `
		SELECT COALESCE(SUM(stocks.quantity - stock_endommages.quantity - commande_lines.quantity), 0) 
		FROM products 
		LEFT JOIN stocks ON products.uuid = stocks.product_uuid
		LEFT JOIN stock_endommages ON products.uuid = stock_endommages.product_uuid
		LEFT JOIN commande_lines ON products.uuid = commande_lines.product_uuid
		WHERE products.entreprise_uuid = ? AND products.deleted_at IS NULL
	`
	stockMoyenArgs := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		stockMoyenQuery += " AND products.pos_uuid = ?"
		stockMoyenArgs = append(stockMoyenArgs, posUUID)
	}

	db.Raw(stockMoyenQuery, stockMoyenArgs...).Row().Scan(&stockMoyen)

	// Alerte si rotation très lente (< 1 fois par mois)
	if stockMoyen > 0 {
		tauxRotation := totalCost / stockMoyen
		if tauxRotation < 1 {
			alertesPerformance++
		}
	}

	// Construction des alertes critiques
	var alertesCritiques []models.Alert

	// Alertes critiques pour stock très faible (< 5)
	criticalStockQuery := `
		SELECT uuid, name, stock
		FROM products 
		WHERE entreprise_uuid = ? AND stock < 5 AND stock > 0 AND deleted_at IS NULL
	`
	criticalStockArgs := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		criticalStockQuery += " AND pos_uuid = ?"
		criticalStockArgs = append(criticalStockArgs, posUUID)
	}
	criticalStockQuery += " ORDER BY stock ASC LIMIT 5"

	rows, err := db.Raw(criticalStockQuery, criticalStockArgs...).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var productUUID, productName string
			var stock float64
			if err := rows.Scan(&productUUID, &productName, &stock); err == nil {
				alert := models.Alert{
					ID:          productUUID,
					Type:        "stock",
					Level:       "critical",
					Title:       "Stock Critique",
					Message:     "Stock très faible",
					ProductName: productName,
					Value:       stock,
					Threshold:   5,
					CreatedAt:   time.Now().Format("2006-01-02T15:04:05Z"),
				}
				alertesCritiques = append(alertesCritiques, alert)
			}
		}
	}
	// Calcul de la tendance des alertes (comparaison avec la semaine précédente)
	alertesTotalCurrent := alertesStock + alertesExpiration + alertesVentes + alertesPerformance

	// Calculer les alertes de la semaine précédente
	var alertesTotalPreviousWeek int64
	previousWeekStart := now.AddDate(0, 0, -14).Format("2006-01-02")
	previousWeekEnd := now.AddDate(0, 0, -7).Format("2006-01-02")

	// Alertes stock semaine précédente
	var alertesStockPrevious int64
	stockQueryPrevious := `
		SELECT COUNT(uuid) 
		FROM products 
		WHERE entreprise_uuid = ? AND stock < 10 AND deleted_at IS NULL 
		AND created_at BETWEEN ? AND ?
	`
	stockArgsPrevious := []interface{}{entrepriseUUID, previousWeekStart, previousWeekEnd}
	if posUUID != "" && posUUID != "null" {
		stockQueryPrevious += " AND pos_uuid = ?"
		stockArgsPrevious = append(stockArgsPrevious, posUUID)
	}
	db.Raw(stockQueryPrevious, stockArgsPrevious...).Row().Scan(&alertesStockPrevious)

	alertesTotalPreviousWeek = alertesStockPrevious // Simplification pour la démonstration

	// Calcul de la tendance basée sur les données réelles
	tendanceAlertes := 0.0
	if alertesTotalPreviousWeek > 0 {
		tendanceAlertes = ((float64(alertesTotalCurrent) - float64(alertesTotalPreviousWeek)) / float64(alertesTotalPreviousWeek)) * 100
	} else if alertesTotalCurrent > 0 {
		tendanceAlertes = 100.0 // Nouvelle apparition d'alertes
	}

	// Construction de la réponse
	alertsKpis := models.AlertsKpis{
		AlertesStock:       alertesStock,
		AlertesExpiration:  alertesExpiration,
		AlertesVentes:      alertesVentes,
		AlertesPerformance: alertesPerformance,
		AlertesCritiques:   alertesCritiques,
		TendanceAlertes:    tendanceAlertes,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alertes KPI récupérées avec succès",
		"data":    alertsKpis,
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
			COALESCE(SUM(stocks.quantity - stock_endommages.quantity - commande_lines.quantity), 0) as stock_disponible,
			COALESCE(SUM(stock_endommages.quantity), 0) as stock_endommage,
			COALESCE(SUM(restitutions.quantity), 0) as stock_restitution
		FROM products 
		LEFT JOIN stocks ON products.uuid = stocks.product_uuid
		LEFT JOIN stock_endommages ON products.uuid = stock_endommages.product_uuid
		LEFT JOIN restitutions ON products.uuid = restitutions.product_uuid
		LEFT JOIN commande_lines ON products.uuid = commande_lines.product_uuid
		WHERE products.entreprise_uuid = ? AND products.deleted_at IS NULL
	`

	args := []interface{}{entrepriseUUID}
	if posUUID != "" && posUUID != "null" {
		query += " AND products.pos_uuid = ?"
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
