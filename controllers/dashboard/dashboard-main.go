package dashboard

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
)

// ============================= ENDPOINTS =============================

// GetDashboardStats retourne les statistiques principales du dashboard
func GetDashboardStats(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	stats := calculateDashboardStats(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(stats)
}

// GetSalesChartData retourne les données pour le graphique de ventes
func GetSalesChartData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de début invalide",
			})
		}
		endDate, err = time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de fin invalide",
			})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les dates de début et de fin sont requises",
		})
	}

	chartData := getSalesChartData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(chartData)
}

// GetPlatChartData retourne les données pour le graphique donut des plats
func GetPlatChartData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de début invalide",
			})
		}
		endDate, err = time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de fin invalide",
			})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les dates de début et de fin sont requises",
		})
	}

	chartData := getPlatChartData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(chartData)
}

// GetProductChartData retourne les données pour le graphique donut des produits
func GetProductChartData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de début invalide",
			})
		}
		endDate, err = time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de fin invalide",
			})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les dates de début et de fin sont requises",
		})
	}

	chartData := getProductChartData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(chartData)
}

// GetStockAlerts retourne les alertes de stock
func GetStockAlerts(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	alerts := getStockAlerts(entrepriseUUID, posUUID)
	return c.JSON(alerts)
}

// GetStockRotationData retourne les données de rotation de stock
func GetStockRotationData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	rotationData := getStockRotationData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(rotationData)
}

// GetPlatStatistics retourne les statistiques des plats
func GetPlatStatistics(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	stats := getPlatStatistics(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(stats)
}

// GetLivraisonStatistics retourne les statistiques des livraisons
func GetLivraisonStatistics(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	stats := getLivraisonStatistics(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(stats)
}

// GetLivraisonZonesData retourne les données des zones de livraison
func GetLivraisonZonesData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	zonesData := getLivraisonZonesData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(zonesData)
}

// GetLivreurPerformance retourne les performances des livreurs
func GetLivreurPerformance(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	performance := getLivreurPerformance(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(performance)
}

// GetCaisseStatistics retourne les statistiques de la caisse
func GetCaisseStatistics(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	stats := getCaisseStatistics(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(stats)
}

// GetFluxTresorerieData retourne les données de flux de trésorerie
func GetFluxTresorerieData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de début invalide",
			})
		}
		endDate, err = time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Format de date de fin invalide",
			})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les dates de début et de fin sont requises",
		})
	}

	fluxData := getFluxTresorerieData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(fluxData)
}

// GetRepartitionTransactionsData retourne la répartition des transactions
func GetRepartitionTransactionsData(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	repartitionData := getRepartitionTransactionsData(entrepriseUUID, posUUID, startDate, endDate)
	return c.JSON(repartitionData)
}

// GetTopTransactions retourne les meilleures transactions
func GetTopTransactions(c *fiber.Ctx) error {
	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.Query("limit", "5")

	if entrepriseUUID == "" || posUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Les paramètres entreprise_uuid et pos_uuid sont requis",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	var startDate, endDate *time.Time
	if startDateStr != "" && endDateStr != "" {
		start, err1 := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr)
		end, err2 := time.Parse("2006-01-02T15:04:05Z07:00", endDateStr)
		if err1 == nil && err2 == nil {
			startDate = &start
			endDate = &end
		}
	}

	topTransactions := getTopTransactions(entrepriseUUID, posUUID, startDate, endDate, limit)
	return c.JSON(topTransactions)
}

// ============================= FONCTIONS DE CALCUL =============================

// calculateDashboardStats calcule les statistiques principales du dashboard
func calculateDashboardStats(entrepriseUUID, posUUID string, startDate, endDate *time.Time) models.DashboardStats {
	db := database.DB

	// 1. Total articles
	var totalArticles int64
	db.Model(&models.Product{}).Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID).Count(&totalArticles)

	// 2. Articles en rupture de stock
	var articlesRuptureStock int64
	db.Model(&models.Product{}).Where("entreprise_uuid = ? AND pos_uuid = ? AND stock <= 0", entrepriseUUID, posUUID).Count(&articlesRuptureStock)

	// 3. Total ventes (nombre de produits vendus uniquement)
	query := db.Table("commande_lines cl").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "product")

	if startDate != nil && endDate != nil {
		query = query.Where("c.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	var totalVentes int64
	query.Select("COALESCE(SUM(cl.quantity), 0)").Scan(&totalVentes)

	// 4. Total montant vendu (uniquement pour les produits)
	var totalMontantVendu float64
	subquery := db.Table("commande_lines cl").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Joins("JOIN products p ON cl.product_uuid = p.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "product")

	if startDate != nil && endDate != nil {
		subquery = subquery.Where("c.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	subquery.Select("COALESCE(SUM(cl.quantity * p.prix_vente), 0)").Scan(&totalMontantVendu)

	// Calcul des pourcentages
	var articlesRuptureStockPercentage int
	if totalArticles > 0 {
		articlesRuptureStockPercentage = int(math.Round(float64(articlesRuptureStock) / float64(totalArticles) * 100))
	}

	var totalVentesPercentage int
	if totalArticles > 0 {
		totalVentesPercentage = int(math.Round(float64(totalVentes) / float64(totalArticles) * 100))
	}

	var totalMontantVenduPercentage int
	if totalArticles > 0 {
		// Assumons 1000 comme prix moyen pour le calcul du pourcentage
		totalMontantVenduPercentage = int(math.Round(totalMontantVendu / (float64(totalArticles) * 1000) * 100))
	}

	return models.DashboardStats{
		TotalArticles:                  totalArticles,
		ArticlesRuptureStock:           articlesRuptureStock,
		ArticlesRuptureStockPercentage: articlesRuptureStockPercentage,
		TotalVentes:                    totalVentes,
		TotalVentesPercentage:          totalVentesPercentage,
		TotalMontantVendu:              totalMontantVendu,
		TotalMontantVenduPercentage:    totalMontantVenduPercentage,
	}
}

// getSalesChartData récupère les données pour le graphique de ventes
func getSalesChartData(entrepriseUUID, posUUID string, startDate, endDate time.Time) models.SalesChartData {
	db := database.DB

	// Vérifier si c'est le même jour
	isOneDay := startDate.Format("2006-01-02") == endDate.Format("2006-01-02")

	timeData := make(map[string]struct {
		commandes int64
		montant   float64
		gain      float64
	})

	var timeKeys []string

	if isOneDay {
		// Grouper par heure pour une seule journée
		for hour := 0; hour < 24; hour++ {
			hourKey := fmt.Sprintf("%02d", hour)
			timeKeys = append(timeKeys, hourKey)
			timeData[hourKey] = struct {
				commandes int64
				montant   float64
				gain      float64
			}{0, 0, 0}
		}
	} else {
		// Grouper par jour pour plusieurs jours
		currentDate := startDate
		for currentDate.Before(endDate) || currentDate.Equal(endDate) {
			dateKey := currentDate.Format("2006-01-02")
			timeKeys = append(timeKeys, dateKey)
			timeData[dateKey] = struct {
				commandes int64
				montant   float64
				gain      float64
			}{0, 0, 0}
			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	// Récupérer les commandes payées avec des produits
	var results []struct {
		CreatedAt time.Time
		Quantity  int64
		PrixVente float64
		PrixAchat float64
	}

	query := db.Table("commande_lines cl").
		Select("c.created_at, cl.quantity, p.prix_vente, p.prix_achat").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Joins("JOIN products p ON cl.product_uuid = p.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "product").
		Where("c.created_at BETWEEN ? AND ?", startDate, endDate)

	query.Scan(&results)

	// Traiter les résultats
	for _, result := range results {
		var timeKey string
		if isOneDay {
			timeKey = fmt.Sprintf("%02d", result.CreatedAt.Hour())
		} else {
			timeKey = result.CreatedAt.Format("2006-01-02")
		}

		if data, exists := timeData[timeKey]; exists {
			data.commandes += result.Quantity
			chiffresAffaires := result.PrixVente * float64(result.Quantity)
			cout := result.PrixAchat * float64(result.Quantity)
			data.montant += chiffresAffaires
			data.gain += chiffresAffaires - cout
			timeData[timeKey] = data
		}
	}

	// Préparer les données pour le graphique
	var chartLabels []string
	var totalCommandes []int64
	var montantVendu []float64
	var gainObtenu []float64

	for _, key := range timeKeys {
		if isOneDay {
			chartLabels = append(chartLabels, key+"h")
		} else {
			date, _ := time.Parse("2006-01-02", key)
			chartLabels = append(chartLabels, date.Format("02/01"))
		}

		data := timeData[key]
		totalCommandes = append(totalCommandes, data.commandes)
		montantVendu = append(montantVendu, math.Round(data.montant*100)/100)
		gainObtenu = append(gainObtenu, math.Round(data.gain*100)/100)
	}

	return models.SalesChartData{
		Dates:          chartLabels,
		TotalCommandes: totalCommandes,
		MontantVendu:   montantVendu,
		GainObtenu:     gainObtenu,
	}
}

// getPlatChartData récupère les données pour le graphique donut des plats
func getPlatChartData(entrepriseUUID, posUUID string, startDate, endDate time.Time) models.PlatChartData {
	db := database.DB

	var results []struct {
		Name     string
		Montant  float64
		Quantity int64
	}

	query := db.Table("commande_lines cl").
		Select("pl.name, SUM(cl.quantity * pl.prix) as montant, SUM(cl.quantity) as quantity").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Joins("JOIN plats pl ON cl.plat_uuid = pl.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "plat").
		Where("c.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("pl.uuid, pl.name").
		Order("montant DESC").
		Limit(10)

	query.Scan(&results)

	if len(results) == 0 {
		return models.PlatChartData{
			Labels:      []string{},
			Series:      []float64{},
			Percentages: []float64{},
		}
	}

	// Calculer le total pour les pourcentages
	var totalMontant float64
	for _, result := range results {
		totalMontant += result.Montant
	}

	var labels []string
	var series []float64
	var percentages []float64

	for _, result := range results {
		labels = append(labels, result.Name)
		series = append(series, math.Round(result.Montant*100)/100)
		if totalMontant > 0 {
			percentages = append(percentages, math.Round((result.Montant/totalMontant)*10000)/100)
		} else {
			percentages = append(percentages, 0)
		}
	}

	return models.PlatChartData{
		Labels:      labels,
		Series:      series,
		Percentages: percentages,
	}
}

// getProductChartData récupère les données pour le graphique donut des produits
func getProductChartData(entrepriseUUID, posUUID string, startDate, endDate time.Time) models.ProductChartData {
	db := database.DB

	var results []struct {
		Name     string
		Montant  float64
		Quantity int64
	}

	query := db.Table("commande_lines cl").
		Select("p.name, SUM(cl.quantity * p.prix_vente) as montant, SUM(cl.quantity) as quantity").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Joins("JOIN products p ON cl.product_uuid = p.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "product").
		Where("c.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("p.uuid, p.name").
		Order("montant DESC").
		Limit(10)

	query.Scan(&results)

	if len(results) == 0 {
		return models.ProductChartData{
			Labels:      []string{},
			Series:      []float64{},
			Percentages: []float64{},
		}
	}

	// Calculer le total pour les pourcentages
	var totalMontant float64
	for _, result := range results {
		totalMontant += result.Montant
	}

	var labels []string
	var series []float64
	var percentages []float64

	for _, result := range results {
		labels = append(labels, result.Name)
		series = append(series, math.Round(result.Montant*100)/100)
		if totalMontant > 0 {
			percentages = append(percentages, math.Round((result.Montant/totalMontant)*10000)/100)
		} else {
			percentages = append(percentages, 0)
		}
	}

	return models.ProductChartData{
		Labels:      labels,
		Series:      series,
		Percentages: percentages,
	}
}

// getStockAlerts récupère les produits en alerte de stock
func getStockAlerts(entrepriseUUID, posUUID string) []models.StockAlert {
	db := database.DB

	var products []models.Product
	db.Where("entreprise_uuid = ? AND pos_uuid = ? AND stock <= ?", entrepriseUUID, posUUID, 5).
		Order("stock ASC").
		Find(&products)

	var alerts []models.StockAlert
	for _, product := range products {
		alertType := "avertissement"
		if product.Stock <= 0 {
			alertType = "rupture"
		}

		alerts = append(alerts, models.StockAlert{
			UUID:       product.UUID,
			Name:       product.Name,
			Reference:  product.Reference,
			UniteVente: product.UniteVente,
			Stock:      product.Stock,
			AlertType:  alertType,
			Image:      product.Image,
			PrixVente:  product.PrixVente,
		})
	}

	return alerts
}

// getStockRotationData calcule le taux de rotation de stock pour les produits
func getStockRotationData(entrepriseUUID, posUUID string, startDate, endDate *time.Time) models.StockRotationData {
	db := database.DB

	// Définir la période d'analyse (par défaut les 12 derniers mois)
	var periodEnd time.Time
	var periodStart time.Time

	if endDate != nil {
		periodEnd = *endDate
	} else {
		periodEnd = time.Now()
	}

	if startDate != nil {
		periodStart = *startDate
	} else {
		periodStart = periodEnd.AddDate(-1, 0, 0) // 12 mois en arrière
	}

	// Récupérer les produits avec stock
	var products []models.Product
	db.Where("entreprise_uuid = ? AND pos_uuid = ? AND stock > 0", entrepriseUUID, posUUID).Find(&products)

	if len(products) == 0 {
		return models.StockRotationData{
			ProductNames:  []string{"Aucune donnée"},
			RotationRates: []float64{0},
			Categories:    []string{"Pas de rotation"},
			Colors:        []string{"#6c757d"},
		}
	}

	// Récupérer les ventes par produit sur la période
	salesByProduct := make(map[string]float64)

	var salesResults []struct {
		ProductUUID string
		TotalSales  float64
	}

	db.Table("commande_lines cl").
		Select("cl.product_uuid, SUM(cl.quantity) as total_sales").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "product").
		Where("c.created_at BETWEEN ? AND ?", periodStart, periodEnd).
		Group("cl.product_uuid").
		Scan(&salesResults)

	for _, result := range salesResults {
		salesByProduct[result.ProductUUID] = result.TotalSales
	}

	// Calculer le taux de rotation pour chaque produit
	type productRotation struct {
		name         string
		rotationRate float64
		category     string
		color        string
	}

	var productRotations []productRotation
	colorPalette := []string{
		"#28a745", // Vert - Rotation excellente (≥6)
		"#6f42c1", // Violet - Rotation très bonne (4-6)
		"#007bff", // Bleu - Rotation bonne (2-4)
		"#fd7e14", // Orange - Rotation moyenne (1-2)
		"#dc3545", // Rouge - Rotation faible (<1)
	}

	for _, product := range products {
		quantiteVendue := salesByProduct[product.UUID]
		stockMoyen := product.Stock // Simplification : on utilise le stock actuel

		if stockMoyen > 0 && quantiteVendue > 0 {
			tauxRotation := quantiteVendue / stockMoyen

			// Catégorisation du taux de rotation
			var category string
			var colorIndex int

			if tauxRotation >= 6 {
				category = "Excellente (≥6)"
				colorIndex = 0
			} else if tauxRotation >= 4 {
				category = "Très bonne (4-6)"
				colorIndex = 1
			} else if tauxRotation >= 2 {
				category = "Bonne (2-4)"
				colorIndex = 2
			} else if tauxRotation >= 1 {
				category = "Moyenne (1-2)"
				colorIndex = 3
			} else {
				category = "Faible (<1)"
				colorIndex = 4
			}

			productRotations = append(productRotations, productRotation{
				name:         product.Name,
				rotationRate: math.Round(tauxRotation*100) / 100,
				category:     category,
				color:        colorPalette[colorIndex],
			})
		}
	}

	// Trier par taux de rotation décroissant et prendre les 10 premiers
	for i := 0; i < len(productRotations)-1; i++ {
		for j := i + 1; j < len(productRotations); j++ {
			if productRotations[j].rotationRate > productRotations[i].rotationRate {
				productRotations[i], productRotations[j] = productRotations[j], productRotations[i]
			}
		}
	}

	if len(productRotations) > 10 {
		productRotations = productRotations[:10]
	}

	if len(productRotations) == 0 {
		return models.StockRotationData{
			ProductNames:  []string{"Aucune donnée"},
			RotationRates: []float64{0},
			Categories:    []string{"Pas de rotation"},
			Colors:        []string{"#6c757d"},
		}
	}

	var productNames []string
	var rotationRates []float64
	var categories []string
	var colors []string

	for _, pr := range productRotations {
		productNames = append(productNames, pr.name)
		rotationRates = append(rotationRates, pr.rotationRate)
		categories = append(categories, pr.category)
		colors = append(colors, pr.color)
	}

	return models.StockRotationData{
		ProductNames:  productNames,
		RotationRates: rotationRates,
		Categories:    categories,
		Colors:        colors,
	}
}

// getPlatStatistics calcule les statistiques des plats
func getPlatStatistics(entrepriseUUID, posUUID string, startDate, endDate *time.Time) models.PlatStatistics {
	db := database.DB

	// Total plats
	var totalPlats int64
	db.Model(&models.Plat{}).Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID).Count(&totalPlats)

	// Total clients uniques ayant commandé des plats
	var totalClients int64
	subquery := db.Table("commande_lines cl").
		Select("c.client_uuid").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "plat")

	if startDate != nil && endDate != nil {
		subquery = subquery.Where("c.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	subquery.Group("c.client_uuid").Count(&totalClients)

	// Quantités vendues (plats)
	var quantitesVendues int64
	platQuery := db.Table("commande_lines cl").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "plat")

	if startDate != nil && endDate != nil {
		platQuery = platQuery.Where("c.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	platQuery.Select("COALESCE(SUM(cl.quantity), 0)").Scan(&quantitesVendues)

	// Chiffre d'affaires (plats)
	var chiffresAffaires float64
	caQuery := db.Table("commande_lines cl").
		Joins("JOIN commandes c ON cl.commande_uuid = c.uuid").
		Joins("JOIN plats pl ON cl.plat_uuid = pl.uuid").
		Where("c.entreprise_uuid = ? AND c.pos_uuid = ? AND c.status = ? AND cl.item_type = ?", entrepriseUUID, posUUID, "paid", "plat")

	if startDate != nil && endDate != nil {
		caQuery = caQuery.Where("c.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	caQuery.Select("COALESCE(SUM(cl.quantity * pl.prix), 0)").Scan(&chiffresAffaires)

	return models.PlatStatistics{
		TotalPlats:       totalPlats,
		TotalClients:     totalClients,
		QuantitesVendues: quantitesVendues,
		ChiffresAffaires: math.Round(chiffresAffaires*100) / 100,
	}
}

// getLivraisonStatistics calcule les statistiques des livraisons
func getLivraisonStatistics(entrepriseUUID, posUUID string, startDate, endDate *time.Time) models.LivraisonStats {
	db := database.DB

	// Construire la requête de base
	baseQuery := db.Model(&models.Livraison{}).Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID)

	if startDate != nil && endDate != nil {
		baseQuery = baseQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	// Total livraisons
	var totalLivraisons int64
	baseQuery.Count(&totalLivraisons)

	// Livraisons par statut
	var enCours, effectuees, annulees int64
	baseQuery.Where("statut = ?", "En cours").Count(&enCours)

	baseQueryCopy := db.Model(&models.Livraison{}).Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID)
	if startDate != nil && endDate != nil {
		baseQueryCopy = baseQueryCopy.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	baseQueryCopy.Where("statut = ?", "Effectuée").Count(&effectuees)

	baseQueryCopy2 := db.Model(&models.Livraison{}).Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID)
	if startDate != nil && endDate != nil {
		baseQueryCopy2 = baseQueryCopy2.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	baseQueryCopy2.Where("statut = ?", "Annulée").Count(&annulees)

	// Calculer les pourcentages
	var enCoursPercentage, effectueesPercentage, annuleesPercentage float64
	if totalLivraisons > 0 {
		enCoursPercentage = math.Round(float64(enCours)/float64(totalLivraisons)*10000) / 100
		effectueesPercentage = math.Round(float64(effectuees)/float64(totalLivraisons)*10000) / 100
		annuleesPercentage = math.Round(float64(annulees)/float64(totalLivraisons)*10000) / 100
	}

	// Calculer le revenu total (frais de livraison)
	var totalRevenu float64
	var livraisons []models.Livraison

	livraisonQuery := db.Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID)
	if startDate != nil && endDate != nil {
		livraisonQuery = livraisonQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	livraisonQuery.Find(&livraisons)

	// Pour cet exemple, on utilise une valeur fixe de 5.0 pour les frais de livraison
	// Dans une implémentation réelle, il faudrait ajouter un champ price ou cost au modèle Zone ou Livraison
	for range livraisons {
		totalRevenu += 5.0 // Frais de livraison fixe
	}

	// Calculer le revenu moyen par livraison
	var revenuMoyen float64
	if totalLivraisons > 0 {
		revenuMoyen = math.Round(totalRevenu/float64(totalLivraisons)*100) / 100
	}

	// Calculer le taux de réussite
	var tauxReussite float64
	if totalLivraisons > 0 {
		tauxReussite = math.Round(float64(effectuees)/float64(totalLivraisons)*10000) / 100
	}

	return models.LivraisonStats{
		TotalLivraisons:      totalLivraisons,
		EnCours:              enCours,
		Effectuees:           effectuees,
		Annulees:             annulees,
		EnCoursPercentage:    enCoursPercentage,
		EffectueesPercentage: effectueesPercentage,
		AnnuleesPercentage:   annuleesPercentage,
		TotalRevenu:          math.Round(totalRevenu*100) / 100,
		RevenuMoyen:          revenuMoyen,
		TauxReussite:         tauxReussite,
	}
}

// getLivraisonZonesData récupère les données des zones de livraison (Top 5)
func getLivraisonZonesData(entrepriseUUID, posUUID string, startDate, endDate *time.Time) []models.LivraisonZoneData {
	db := database.DB

	var results []struct {
		ZoneName         string
		NombreLivraisons int64
		Revenu           float64
	}

	query := db.Table("livraisons l").
		Select("z.name as zone_name, COUNT(*) as nombre_livraisons, COUNT(*) * 5.0 as revenu").
		Joins("JOIN zones z ON l.zone_uuid = z.uuid").
		Where("l.entreprise_uuid = ? AND l.pos_uuid = ?", entrepriseUUID, posUUID)

	if startDate != nil && endDate != nil {
		query = query.Where("l.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Group("z.uuid, z.name").
		Order("nombre_livraisons DESC").
		Limit(5).
		Scan(&results)

	var zonesData []models.LivraisonZoneData
	for _, result := range results {
		zonesData = append(zonesData, models.LivraisonZoneData{
			ZoneName:         result.ZoneName,
			NombreLivraisons: result.NombreLivraisons,
			Revenu:           math.Round(result.Revenu*100) / 100,
		})
	}

	return zonesData
}

// getLivreurPerformance récupère les performances des livreurs (Top 5)
func getLivreurPerformance(entrepriseUUID, posUUID string, startDate, endDate *time.Time) []models.LivreurPerformance {
	db := database.DB

	var results []struct {
		UUID            string
		Name            string
		TotalLivraisons int64
		Effectuees      int64
		EnCours         int64
		Annulees        int64
	}

	query := db.Table("livraisons l").
		Select(`lr.uuid, lr.name, 
				COUNT(*) as total_livraisons,
				SUM(CASE WHEN l.statut = 'Effectuée' THEN 1 ELSE 0 END) as effectuees,
				SUM(CASE WHEN l.statut = 'En cours' THEN 1 ELSE 0 END) as en_cours,
				SUM(CASE WHEN l.statut = 'Annulée' THEN 1 ELSE 0 END) as annulees`).
		Joins("JOIN livreurs lr ON l.livreur_uuid = lr.uuid").
		Where("l.entreprise_uuid = ? AND l.pos_uuid = ?", entrepriseUUID, posUUID)

	if startDate != nil && endDate != nil {
		query = query.Where("l.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Group("lr.uuid, lr.name").
		Order("effectuees DESC").
		Limit(5).
		Scan(&results)

	var livreurPerf []models.LivreurPerformance
	for _, result := range results {
		var tauxReussite float64
		if result.TotalLivraisons > 0 {
			tauxReussite = math.Round(float64(result.Effectuees)/float64(result.TotalLivraisons)*10000) / 100
		}

		livreurPerf = append(livreurPerf, models.LivreurPerformance{
			UUID:            result.UUID,
			Name:            result.Name,
			TotalLivraisons: result.TotalLivraisons,
			Effectuees:      result.Effectuees,
			EnCours:         result.EnCours,
			Annulees:        result.Annulees,
			TauxReussite:    tauxReussite,
		})
	}

	return livreurPerf
}

// getCaisseStatistics calcule les statistiques de la caisse
func getCaisseStatistics(entrepriseUUID, posUUID string, startDate, endDate *time.Time) models.CaisseStatistics {
	db := database.DB

	// Récupérer toutes les caisses du POS
	var caisses []models.Caisse
	db.Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID).Find(&caisses)

	if len(caisses) == 0 {
		return getEmptyCaisseStatistics()
	}

	// Extraire les UUIDs des caisses
	var caisseUUIDs []string
	for _, caisse := range caisses {
		caisseUUIDs = append(caisseUUIDs, caisse.UUID)
	}

	// Construire la requête de base pour les items de caisse
	baseQuery := db.Model(&models.CaisseItem{}).Where("caisse_uuid IN ?", caisseUUIDs)

	if startDate != nil && endDate != nil {
		baseQuery = baseQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	// Calculer les totaux
	var results struct {
		TotalEntrees       float64
		TotalSorties       float64
		NombreTransactions int64
	}

	baseQuery.Select(`
		SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) as total_entrees,
		SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) as total_sorties,
		COUNT(*) as nombre_transactions
	`).Scan(&results)

	soldeCaisse := results.TotalEntrees - results.TotalSorties

	// Calculer les moyennes
	var moyenneEntree, moyenneSortie float64

	var nombreEntrees, nombreSorties int64
	baseQuery.Where("type_transaction = ?", "Entree").Count(&nombreEntrees)

	baseQueryCopy := db.Model(&models.CaisseItem{}).Where("caisse_uuid IN ?", caisseUUIDs)
	if startDate != nil && endDate != nil {
		baseQueryCopy = baseQueryCopy.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	baseQueryCopy.Where("type_transaction = ?", "Sortie").Count(&nombreSorties)

	if nombreEntrees > 0 {
		moyenneEntree = results.TotalEntrees / float64(nombreEntrees)
	}
	if nombreSorties > 0 {
		moyenneSortie = results.TotalSorties / float64(nombreSorties)
	}

	// Calculer les ratios
	var ratioEntreeSortie float64
	if results.TotalSorties > 0 {
		ratioEntreeSortie = results.TotalEntrees / results.TotalSorties
	}

	var tauxLiquidite float64
	if results.TotalEntrees > 0 {
		tauxLiquidite = (soldeCaisse / results.TotalEntrees) * 100
	}

	// Calculer l'évolution (simulée pour cet exemple)
	evolution := calculateCaisseEvolution(entrepriseUUID, posUUID, startDate, endDate)

	// Analyser le jour le plus actif
	jourLePlusActif := getJourLePlusActif(caisseUUIDs, startDate, endDate)

	// Analyser l'heure la plus active
	heureLaPlusActive := getHeureLaPlusActive(caisseUUIDs, startDate, endDate)

	// Calculer le nombre moyen de transactions par jour
	var nombreTransactionsParJour float64
	if startDate != nil && endDate != nil {
		nombreJours := endDate.Sub(*startDate).Hours() / 24
		if nombreJours > 0 {
			nombreTransactionsParJour = float64(results.NombreTransactions) / nombreJours
		}
	}

	return models.CaisseStatistics{
		SoldeCaisse:               soldeCaisse,
		TotalEntrees:              results.TotalEntrees,
		TotalSorties:              results.TotalSorties,
		NombreTransactions:        results.NombreTransactions,
		MoyenneEntree:             moyenneEntree,
		MoyenneSortie:             moyenneSortie,
		RatioEntreeSortie:         ratioEntreeSortie,
		TauxLiquidite:             tauxLiquidite,
		EvolutionJournaliere:      evolution.montant,
		EvolutionPercentage:       evolution.percentage,
		Tendance:                  evolution.tendance,
		JourLePlusActif:           jourLePlusActif,
		HeureLaPlusActive:         heureLaPlusActive,
		NombreTransactionsParJour: math.Round(nombreTransactionsParJour*10) / 10,
	}
}

// getFluxTresorerieData récupère les données de flux de trésorerie
func getFluxTresorerieData(entrepriseUUID, posUUID string, startDate, endDate time.Time) models.FluxTresorerieData {
	db := database.DB

	// Récupérer les caisses
	var caisses []models.Caisse
	db.Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID).Find(&caisses)

	if len(caisses) == 0 {
		return models.FluxTresorerieData{
			Dates:   []string{},
			Entrees: []float64{},
			Sorties: []float64{},
			Soldes:  []float64{},
		}
	}

	var caisseUUIDs []string
	for _, caisse := range caisses {
		caisseUUIDs = append(caisseUUIDs, caisse.UUID)
	}

	// Déterminer si on affiche par heure ou par jour
	isOneDay := startDate.Format("2006-01-02") == endDate.Format("2006-01-02")

	if isOneDay {
		return getFluxParHeure(caisseUUIDs, startDate)
	} else {
		return getFluxParJour(caisseUUIDs, startDate, endDate)
	}
}

// getRepartitionTransactionsData récupère la répartition des transactions
func getRepartitionTransactionsData(entrepriseUUID, posUUID string, startDate, endDate *time.Time) models.RepartitionTransactionsData {
	db := database.DB

	// Récupérer les caisses
	var caisses []models.Caisse
	db.Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID).Find(&caisses)

	if len(caisses) == 0 {
		return models.RepartitionTransactionsData{
			Labels:      []string{},
			Values:      []float64{},
			Percentages: []float64{},
			Colors:      []string{},
		}
	}

	var caisseUUIDs []string
	for _, caisse := range caisses {
		caisseUUIDs = append(caisseUUIDs, caisse.UUID)
	}

	// Récupérer les transactions groupées par libellé
	var results []struct {
		Libelle string
		Montant float64
	}

	query := db.Table("caisse_items ci").
		Select("ci.libelle, SUM(ci.montant) as montant").
		Where("ci.caisse_uuid IN ?", caisseUUIDs)

	if startDate != nil && endDate != nil {
		query = query.Where("ci.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Group("ci.libelle").
		Order("montant DESC").
		Limit(8).
		Scan(&results)

	if len(results) == 0 {
		return models.RepartitionTransactionsData{
			Labels:      []string{},
			Values:      []float64{},
			Percentages: []float64{},
			Colors:      []string{},
		}
	}

	// Calculer le total
	var total float64
	for _, result := range results {
		total += result.Montant
	}

	var labels []string
	var values []float64
	var percentages []float64
	var colors []string

	colorPalette := []string{
		"#28a745", "#007bff", "#ffc107", "#dc3545", "#6f42c1",
		"#20c997", "#fd7e14", "#6c757d",
	}

	for i, result := range results {
		labels = append(labels, result.Libelle)
		values = append(values, result.Montant)
		if total > 0 {
			percentages = append(percentages, (result.Montant/total)*100)
		} else {
			percentages = append(percentages, 0)
		}
		colors = append(colors, colorPalette[i%len(colorPalette)])
	}

	return models.RepartitionTransactionsData{
		Labels:      labels,
		Values:      values,
		Percentages: percentages,
		Colors:      colors,
	}
}

// getTopTransactions récupère les meilleures transactions (entrées et sorties)
func getTopTransactions(entrepriseUUID, posUUID string, startDate, endDate *time.Time, limit int) models.TopTransactions {
	db := database.DB

	// Récupérer les caisses
	var caisses []models.Caisse
	db.Where("entreprise_uuid = ? AND pos_uuid = ?", entrepriseUUID, posUUID).Find(&caisses)

	if len(caisses) == 0 {
		return models.TopTransactions{
			TopEntrees: []models.TopTransaction{},
			TopSorties: []models.TopTransaction{},
		}
	}

	var caisseUUIDs []string
	for _, caisse := range caisses {
		caisseUUIDs = append(caisseUUIDs, caisse.UUID)
	}

	// Top entrées
	var topEntrees []models.CaisseItem
	entreeQuery := db.Where("caisse_uuid IN ? AND type_transaction = ?", caisseUUIDs, "Entree")
	if startDate != nil && endDate != nil {
		entreeQuery = entreeQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	entreeQuery.Order("montant DESC").Limit(limit).Find(&topEntrees)

	// Top sorties
	var topSorties []models.CaisseItem
	sortieQuery := db.Where("caisse_uuid IN ? AND type_transaction = ?", caisseUUIDs, "Sortie")
	if startDate != nil && endDate != nil {
		sortieQuery = sortieQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	sortieQuery.Order("montant DESC").Limit(limit).Find(&topSorties)

	// Convertir en format de réponse
	var topEntreesResponse []models.TopTransaction
	for _, entree := range topEntrees {
		topEntreesResponse = append(topEntreesResponse, models.TopTransaction{
			Libelle:   entree.Libelle,
			Montant:   entree.Montant,
			Type:      entree.TypeTransaction,
			Date:      entree.CreatedAt,
			Reference: entree.Reference,
		})
	}

	var topSortiesResponse []models.TopTransaction
	for _, sortie := range topSorties {
		topSortiesResponse = append(topSortiesResponse, models.TopTransaction{
			Libelle:   sortie.Libelle,
			Montant:   sortie.Montant,
			Type:      sortie.TypeTransaction,
			Date:      sortie.CreatedAt,
			Reference: sortie.Reference,
		})
	}

	return models.TopTransactions{
		TopEntrees: topEntreesResponse,
		TopSorties: topSortiesResponse,
	}
}

// ============================= FONCTIONS UTILITAIRES =============================

// getEmptyCaisseStatistics retourne des statistiques de caisse vides
func getEmptyCaisseStatistics() models.CaisseStatistics {
	return models.CaisseStatistics{
		SoldeCaisse:               0,
		TotalEntrees:              0,
		TotalSorties:              0,
		NombreTransactions:        0,
		MoyenneEntree:             0,
		MoyenneSortie:             0,
		RatioEntreeSortie:         0,
		TauxLiquidite:             0,
		EvolutionJournaliere:      0,
		EvolutionPercentage:       0,
		Tendance:                  "stable",
		JourLePlusActif:           "N/A",
		HeureLaPlusActive:         "N/A",
		NombreTransactionsParJour: 0,
	}
}

// calculateCaisseEvolution calcule l'évolution de la caisse
func calculateCaisseEvolution(entrepriseUUID, posUUID string, startDate, endDate *time.Time) struct {
	montant    float64
	percentage float64
	tendance   string
} {
	// Pour simplifier, on retourne des valeurs par défaut
	// Dans une implémentation complète, on comparerait avec la période précédente
	return struct {
		montant    float64
		percentage float64
		tendance   string
	}{
		montant:    0,
		percentage: 0,
		tendance:   "stable",
	}
}

// getJourLePlusActif trouve le jour le plus actif
func getJourLePlusActif(caisseUUIDs []string, startDate, endDate *time.Time) string {
	db := database.DB

	if len(caisseUUIDs) == 0 {
		return "N/A"
	}

	var results []struct {
		DayOfWeek int64
		Count     int64
	}

	query := db.Table("caisse_items ci").
		Select("EXTRACT(DOW FROM ci.created_at) as day_of_week, COUNT(*) as count").
		Where("ci.caisse_uuid IN ?", caisseUUIDs)

	if startDate != nil && endDate != nil {
		query = query.Where("ci.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Group("EXTRACT(DOW FROM ci.created_at)").
		Order("count DESC").
		Limit(1).
		Scan(&results)

	if len(results) == 0 {
		return "N/A"
	}

	jours := []string{"Dimanche", "Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi"}
	dayIndex := results[0].DayOfWeek
	if dayIndex >= 0 && dayIndex < 7 {
		return jours[dayIndex]
	}

	return "N/A"
}

// getHeureLaPlusActive trouve l'heure la plus active
func getHeureLaPlusActive(caisseUUIDs []string, startDate, endDate *time.Time) string {
	db := database.DB

	if len(caisseUUIDs) == 0 {
		return "N/A"
	}

	var results []struct {
		Hour  int64
		Count int64
	}

	query := db.Table("caisse_items ci").
		Select("EXTRACT(HOUR FROM ci.created_at) as hour, COUNT(*) as count").
		Where("ci.caisse_uuid IN ?", caisseUUIDs)

	if startDate != nil && endDate != nil {
		query = query.Where("ci.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Group("EXTRACT(HOUR FROM ci.created_at)").
		Order("count DESC").
		Limit(1).
		Scan(&results)

	if len(results) == 0 {
		return "N/A"
	}

	return fmt.Sprintf("%dh00", results[0].Hour)
}

// getFluxParHeure récupère le flux de trésorerie par heure pour une journée
func getFluxParHeure(caisseUUIDs []string, date time.Time) models.FluxTresorerieData {
	db := database.DB

	var dates []string
	var entrees []float64
	var sorties []float64
	var soldes []float64

	var soldeCumule float64

	for hour := 0; hour < 24; hour++ {
		var result struct {
			TotalEntree float64
			TotalSortie float64
		}

		startHour := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		endHour := startHour.Add(time.Hour - time.Nanosecond)

		db.Table("caisse_items ci").
			Select(`
				SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END) as total_entree,
				SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END) as total_sortie
			`).
			Where("ci.caisse_uuid IN ? AND ci.created_at BETWEEN ? AND ?", caisseUUIDs, startHour, endHour).
			Scan(&result)

		soldeCumule += result.TotalEntree - result.TotalSortie

		dates = append(dates, fmt.Sprintf("%dh", hour))
		entrees = append(entrees, result.TotalEntree)
		sorties = append(sorties, result.TotalSortie)
		soldes = append(soldes, soldeCumule)
	}

	return models.FluxTresorerieData{
		Dates:   dates,
		Entrees: entrees,
		Sorties: sorties,
		Soldes:  soldes,
	}
}

// getFluxParJour récupère le flux de trésorerie par jour
func getFluxParJour(caisseUUIDs []string, startDate, endDate time.Time) models.FluxTresorerieData {
	db := database.DB

	var dates []string
	var entrees []float64
	var sorties []float64
	var soldes []float64

	var soldeCumule float64
	currentDate := startDate

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		var result struct {
			TotalEntree float64
			TotalSortie float64
		}

		startDay := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
		endDay := startDay.Add(24*time.Hour - time.Nanosecond)

		db.Table("caisse_items ci").
			Select(`
				SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END) as total_entree,
				SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END) as total_sortie
			`).
			Where("ci.caisse_uuid IN ? AND ci.created_at BETWEEN ? AND ?", caisseUUIDs, startDay, endDay).
			Scan(&result)

		soldeCumule += result.TotalEntree - result.TotalSortie

		dates = append(dates, currentDate.Format("02/01"))
		entrees = append(entrees, result.TotalEntree)
		sorties = append(sorties, result.TotalSortie)
		soldes = append(soldes, soldeCumule)

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return models.FluxTresorerieData{
		Dates:   dates,
		Entrees: entrees,
		Sorties: sorties,
		Soldes:  soldes,
	}
}
