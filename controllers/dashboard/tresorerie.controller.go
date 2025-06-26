package dashboard

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"gorm.io/gorm"
)

// ===============================
// FONCTIONS UTILITAIRES
// ===============================

// Mapper la période frontend vers les dates de la base de données
func mapPeriodToDateRange(period string) (startDate, endDate time.Time) {
	now := time.Now()

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, now.Location())
	case "this_week":
		// Début de semaine (Lundi)
		weekday := now.Weekday()
		if weekday == 0 { // Dimanche
			weekday = 7
		}
		startDate = now.AddDate(0, 0, -int(weekday-1))
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		endDate = now
	case "this_month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	case "last_3_months":
		startDate = now.AddDate(0, -3, 0)
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		endDate = now
	case "last_6_months":
		startDate = now.AddDate(0, -6, 0)
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		endDate = now
	case "last_year":
		startDate = now.AddDate(-1, 0, 0)
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		endDate = now
	default: // ce_mois par défaut
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	}

	return
}

// Calculer la variation par rapport à la période précédente
func calculateVariation(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0.0
	}
	return ((current - previous) / previous) * 100
}

// Générer des alertes basées sur les données
func generateAlerts(kpiData models.KpiData, topCaisses []models.TopCaisse) []models.TresorerieAlert {
	var alerts []models.TresorerieAlert

	// Alerte solde critique (moins de 50000 CDF ou 500 USD)
	if kpiData.SoldeTotal < 50000 {
		alerts = append(alerts, models.TresorerieAlert{
			ID:      "solde-critique",
			Type:    "danger",
			Icon:    "alert-triangle",
			Title:   "Solde Critique",
			Message: "Le solde total des caisses est en dessous du seuil de sécurité.",
		})
	}

	// Alerte performance dégradée
	caissesFaiblePerformance := 0
	for _, caisse := range topCaisses {
		if caisse.Performance < 60 {
			caissesFaiblePerformance++
		}
	}

	if caissesFaiblePerformance > 0 {
		alerts = append(alerts, models.TresorerieAlert{
			ID:      "performance-faible",
			Type:    "warning",
			Icon:    "trending-down",
			Title:   "Performance Dégradée",
			Message: fmt.Sprintf("%d caisse(s) ont une performance inférieure à 60%%.", caissesFaiblePerformance),
		})
	}

	// Alerte déséquilibre entrées/sorties
	if kpiData.TotalSorties > kpiData.TotalEntrees*0.9 {
		alerts = append(alerts, models.TresorerieAlert{
			ID:      "desequilibre-flux",
			Type:    "warning",
			Icon:    "balance",
			Title:   "Déséquilibre des Flux",
			Message: "Les sorties représentent plus de 90% des entrées.",
		})
	}

	return alerts
}

// Générer des actions recommandées
func generateActionsRecommandees(alerts []models.TresorerieAlert, metriques models.TresorerieMetriques) []models.ActionRecommandee {
	var actions []models.ActionRecommandee

	// Actions basées sur les alertes
	for _, alert := range alerts {
		switch alert.ID {
		case "solde-critique":
			actions = append(actions, models.ActionRecommandee{
				ID:          "action-1",
				Description: "Effectuer un apport de fonds dans les caisses principales",
			})
		case "performance-faible":
			actions = append(actions, models.ActionRecommandee{
				ID:          "action-2",
				Description: "Analyser et optimiser les caisses à faible performance",
			})
		case "desequilibre-flux":
			actions = append(actions, models.ActionRecommandee{
				ID:          "action-3",
				Description: "Revoir la politique de gestion des flux financiers",
			})
		}
	}

	// Actions générales
	actions = append(actions, []models.ActionRecommandee{
		{
			ID:          "action-general-1",
			Description: "Mettre en place un système d'alertes automatiques",
		},
		{
			ID:          "action-general-2",
			Description: "Effectuer une réconciliation hebdomadaire des comptes",
		},
		{
			ID:          "action-general-3",
			Description: "Optimiser la répartition des fonds entre les caisses",
		},
	}...)

	return actions
}

// ===============================
// ENDPOINTS PRINCIPAUX
// ===============================

// GetDashboardTresorerie - Endpoint principal pour récupérer toutes les données du dashboard
func GetDashboardTresorerie(c *fiber.Ctx) error {
	db := database.DB

	// Récupération des paramètres
	userUUID := c.Query("user_uuid")
	posUUID := c.Query("pos_uuid", "")
	period := c.Query("period", "this_month")

	if userUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID utilisateur requis",
		})
	}

	// Récupération de l'entreprise de l'utilisateur
	var entrepriseUUID string
	err := db.Raw("SELECT entreprise_uuid FROM users WHERE uuid = ? AND deleted_at IS NULL", userUUID).
		Row().Scan(&entrepriseUUID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Utilisateur non trouvé",
		})
	}

	// Calcul des dates selon la période
	startDate, endDate := mapPeriodToDateRange(period)

	// Chargement des données
	kpiData, err := loadKpiData(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des KPI: " + err.Error(),
		})
	}

	topCaisses, err := loadTopCaisses(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des top caisses: " + err.Error(),
		})
	}

	metriques, err := loadMetriques(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des métriques: " + err.Error(),
		})
	}

	statsFooter, err := loadStatsFooter(db, entrepriseUUID, posUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des stats footer: " + err.Error(),
		})
	}

	// Génération des alertes et actions
	alerts := generateAlerts(kpiData, topCaisses)
	actionsRecommandees := generateActionsRecommandees(alerts, metriques)

	// Construction de la réponse
	dashboardData := models.DashboardTresorerieData{
		KpiData:             kpiData,
		TopCaisses:          topCaisses,
		Alerts:              alerts,
		Metriques:           metriques,
		ActionsRecommandees: actionsRecommandees,
		StatsFooter:         statsFooter,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    dashboardData,
	})
}

// ===============================
// FONCTIONS DE CHARGEMENT DES DONNÉES
// ===============================

// Charger les données KPI
func loadKpiData(db *gorm.DB, entrepriseUUID, posUUID string, startDate, endDate time.Time) (models.KpiData, error) {
	var kpiData models.KpiData

	// Construction de la requête selon le POS
	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS total_entrees,
			COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS total_sorties,
			COALESCE(SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END), 0) AS total_fond_caisse,
			COUNT(CASE WHEN type_transaction = 'Entree' THEN 1 END) AS nombre_entrees,
			COUNT(CASE WHEN type_transaction = 'Sortie' THEN 1 END) AS nombre_sorties
		FROM caisse_items
		WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND pos_uuid = ?"
		args = []interface{}{entrepriseUUID, startDate, endDate, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{entrepriseUUID, startDate, endDate}
	}

	row := db.Raw(query, args...).Row()
	err := row.Scan(
		&kpiData.TotalEntrees,
		&kpiData.TotalSorties,
		&kpiData.TotalFondDeCaisse,
		&kpiData.NombreEntrees,
		&kpiData.NombreSorties,
	)

	if err != nil {
		return kpiData, err
	}

	// Calcul du solde total
	kpiData.SoldeTotal = kpiData.TotalEntrees - kpiData.TotalSorties + kpiData.TotalFondDeCaisse

	// Calcul de la variation (période précédente)
	prevStartDate := startDate.AddDate(0, -1, 0)
	prevEndDate := endDate.AddDate(0, -1, 0)

	var prevSoldeTotal float64
	prevQuery := strings.Replace(query, "created_at BETWEEN ? AND ?", "created_at BETWEEN ? AND ?", 1)
	prevArgs := make([]interface{}, len(args))
	copy(prevArgs, args)
	prevArgs[1] = prevStartDate
	prevArgs[2] = prevEndDate

	prevRow := db.Raw(prevQuery, prevArgs...).Row()
	var prevEntrees, prevSorties, prevFond float64
	var dummy1, dummy2 int

	err = prevRow.Scan(&prevEntrees, &prevSorties, &prevFond, &dummy1, &dummy2)
	if err == nil {
		prevSoldeTotal = prevEntrees - prevSorties + prevFond
		kpiData.VariationSolde = calculateVariation(kpiData.SoldeTotal, prevSoldeTotal)
	}

	return kpiData, nil
}

// Charger les top caisses
func loadTopCaisses(db *gorm.DB, entrepriseUUID, posUUID string, startDate, endDate time.Time) ([]models.TopCaisse, error) {
	var topCaisses []models.TopCaisse

	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			c.uuid, 
			c.name,
			p.name AS pos_name,
			COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) AS total_entrees,
			COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0) AS total_sorties,
			COUNT(ci.uuid) AS nombre_transactions,
			(COALESCE(c.montant_debut, 0) + 
			 COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) - 
			 COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0)) AS solde
		FROM caisses c
		LEFT JOIN pos p ON p.uuid = c.pos_uuid
		LEFT JOIN caisse_items ci ON ci.caisse_uuid = c.uuid 
			AND ci.created_at BETWEEN ? AND ? 
			AND ci.deleted_at IS NULL
		WHERE c.entreprise_uuid = ? AND c.deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND c.pos_uuid = ?"
		args = []interface{}{startDate, endDate, entrepriseUUID, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{startDate, endDate, entrepriseUUID}
	}

	query += `
		GROUP BY c.uuid, c.name, p.name, c.montant_debut
		ORDER BY solde DESC
		LIMIT 10`

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return topCaisses, err
	}
	defer rows.Close()

	for rows.Next() {
		var caisse models.TopCaisse
		err := rows.Scan(
			&caisse.UUID,
			&caisse.Name,
			&caisse.PosName,
			&caisse.TotalEntrees,
			&caisse.TotalSorties,
			&caisse.NombreTransactions,
			&caisse.Solde,
		)
		if err != nil {
			continue
		}

		// Calcul de la performance (pourcentage de rentabilité)
		totalFlux := caisse.TotalEntrees + caisse.TotalSorties
		if totalFlux > 0 {
			caisse.Performance = ((caisse.TotalEntrees - caisse.TotalSorties) / totalFlux) * 100
		}

		topCaisses = append(topCaisses, caisse)
	}

	return topCaisses, nil
}

// Charger les métriques financières
func loadMetriques(db *gorm.DB, entrepriseUUID, posUUID string, startDate, endDate time.Time) (models.TresorerieMetriques, error) {
	var metriques models.TresorerieMetriques

	// Calcul des métriques complexes basées sur les données réelles
	var totalEntrees, totalSorties float64
	var nombreTransactions int

	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS total_entrees,
			COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS total_sorties,
			COUNT(*) AS nombre_transactions
		FROM caisse_items
		WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND pos_uuid = ?"
		args = []interface{}{entrepriseUUID, startDate, endDate, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{entrepriseUUID, startDate, endDate}
	}

	row := db.Raw(query, args...).Row()
	err := row.Scan(&totalEntrees, &totalSorties, &nombreTransactions)
	if err != nil {
		return metriques, err
	}

	// Calcul des métriques
	daysDiff := endDate.Sub(startDate).Hours() / 24
	if daysDiff == 0 {
		daysDiff = 1
	}

	// Vélocité : transactions par jour
	metriques.Velocite = float64(nombreTransactions) / daysDiff

	// Ratio de liquidité : capacité à couvrir les sorties
	if totalSorties > 0 {
		metriques.RatioLiquidite = (totalEntrees / totalSorties) * 100
	} else {
		metriques.RatioLiquidite = 100
	}

	// Efficacité : montant moyen par transaction
	if nombreTransactions > 0 {
		metriques.Efficacite = (totalEntrees + totalSorties) / float64(nombreTransactions)
	}

	// Calcul de la variation de vélocité et du taux de croissance (période précédente)
	prevStartDate := startDate.AddDate(0, 0, -int(daysDiff))
	prevEndDate := startDate.AddDate(0, 0, -1)

	var prevTransactions int
	var prevEntrees, prevSorties float64

	prevArgs := make([]interface{}, len(args))
	copy(prevArgs, args)
	prevArgs[1] = prevStartDate
	prevArgs[2] = prevEndDate

	prevRow := db.Raw(query, prevArgs...).Row()
	err = prevRow.Scan(&prevEntrees, &prevSorties, &prevTransactions)
	if err == nil {
		prevVelocite := float64(prevTransactions) / daysDiff
		metriques.VelociteVariation = calculateVariation(metriques.Velocite, prevVelocite)

		prevSolde := prevEntrees - prevSorties
		currentSolde := totalEntrees - totalSorties
		metriques.TauxCroissance = calculateVariation(currentSolde, prevSolde)
	}

	return metriques, nil
}

// Charger les statistiques du footer
func loadStatsFooter(db *gorm.DB, entrepriseUUID, posUUID string) (models.TresorerieStatsFooter, error) {
	var statsFooter models.TresorerieStatsFooter

	// Données pour aujourd'hui
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999, today.Location())

	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			COUNT(*) AS transactions_aujourdhui,
			COALESCE(AVG(montant), 0) AS montant_moyen_transaction
		FROM caisse_items
		WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND pos_uuid = ?"
		args = []interface{}{entrepriseUUID, startOfDay, endOfDay, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{entrepriseUUID, startOfDay, endOfDay}
	}

	row := db.Raw(query, args...).Row()
	err := row.Scan(
		&statsFooter.TransactionsAujourdhui,
		&statsFooter.MontantMoyenTransaction,
	)
	if err != nil {
		return statsFooter, err
	}

	// Temps moyen de traitement (simulé entre 2 et 5 secondes)
	statsFooter.TempsMoyenTraitement = 2.5 + (float64(statsFooter.TransactionsAujourdhui%10) * 0.25)

	// Dernière mise à jour
	statsFooter.DerniereMiseAJour = time.Now()

	return statsFooter, nil
}

// ===============================
// ENDPOINTS POUR LES GRAPHIQUES
// ===============================

// GetEvolutionChartData - Données pour le graphique d'évolution
func GetEvolutionChartData(c *fiber.Ctx) error {
	db := database.DB

	// Récupération des paramètres
	userUUID := c.Query("user_uuid")
	posUUID := c.Query("pos_uuid", "")
	period := c.Query("period", "this_month")

	if userUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID utilisateur requis",
		})
	}

	// Récupération de l'entreprise de l'utilisateur
	var entrepriseUUID string
	err := db.Raw("SELECT entreprise_uuid FROM users WHERE uuid = ? AND deleted_at IS NULL", userUUID).
		Row().Scan(&entrepriseUUID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Utilisateur non trouvé",
		})
	}

	// Calcul des dates selon la période
	startDate, endDate := mapPeriodToDateRange(period)

	// Génération des données d'évolution
	evolutionData, err := generateEvolutionData(db, entrepriseUUID, posUUID, startDate, endDate, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors de la génération des données d'évolution: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    evolutionData,
	})
}

// GetRepartitionChartData - Données pour le graphique de répartition (Pie Chart)
func GetRepartitionChartData(c *fiber.Ctx) error {
	db := database.DB

	// Récupération des paramètres
	userUUID := c.Query("user_uuid")
	posUUID := c.Query("pos_uuid", "")
	period := c.Query("period", "this_month")

	if userUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID utilisateur requis",
		})
	}

	// Récupération de l'entreprise de l'utilisateur
	var entrepriseUUID string
	err := db.Raw("SELECT entreprise_uuid FROM users WHERE uuid = ? AND deleted_at IS NULL", userUUID).
		Row().Scan(&entrepriseUUID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Utilisateur non trouvé",
		})
	}

	// Calcul des dates selon la période
	startDate, endDate := mapPeriodToDateRange(period)

	// Génération des données de répartition
	repartitionData, err := generateRepartitionData(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors de la génération des données de répartition: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    repartitionData,
	})
}

// ===============================
// FONCTIONS DE GÉNÉRATION DES GRAPHIQUES
// ===============================

// Générer les données d'évolution
func generateEvolutionData(db *gorm.DB, entrepriseUUID, posUUID string, startDate, endDate time.Time, period string) (models.EvolutionChartData, error) {
	var evolutionData models.EvolutionChartData

	// Déterminer l'intervalle selon la période
	var groupBy string

	switch period {
	case "today":
		groupBy = "HOUR(created_at)"
	case "this_week":
		groupBy = "DATE(created_at)"
	case "this_month":
		groupBy = "DATE(created_at)"
	default:
		groupBy = "DATE(created_at)"
	}

	// Construction de la requête selon le POS
	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			DATE(created_at) as date,
			SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) as entrees,
			SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) as sorties,
			SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END) as fond_caisses
		FROM caisse_items
		WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND pos_uuid = ?"
		args = []interface{}{entrepriseUUID, startDate, endDate, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{entrepriseUUID, startDate, endDate}
	}

	query += ` GROUP BY ` + groupBy + ` ORDER BY date`
	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return evolutionData, err
	}
	defer rows.Close()

	var dailyData []models.DayData
	var cumulativeSolde float64

	for rows.Next() {
		var data models.DayData
		err := rows.Scan(
			&data.Date,
			&data.Entrees,
			&data.Sorties,
			&data.FondCaisses,
		)
		if err != nil {
			continue
		}
		dailyData = append(dailyData, data)
	}

	// Construction des séries de données
	for _, data := range dailyData {
		// Format de la date selon la période
		var dateStr string
		switch period {
		case "today":
			dateStr = data.Date.Format("15:04")
		case "this_week":
			dateStr = data.Date.Format("Mon")
		default:
			dateStr = data.Date.Format("02/01")
		}

		evolutionData.Dates = append(evolutionData.Dates, dateStr)
		evolutionData.EntreesData = append(evolutionData.EntreesData, data.Entrees)
		evolutionData.SortiesData = append(evolutionData.SortiesData, data.Sorties)
		evolutionData.FondCaissesData = append(evolutionData.FondCaissesData, data.FondCaisses)

		// Calcul du solde cumulatif
		cumulativeSolde += data.Entrees - data.Sorties + data.FondCaisses
		evolutionData.SoldeData = append(evolutionData.SoldeData, cumulativeSolde)
	}

	return evolutionData, nil
}

// Générer les données de répartition
func generateRepartitionData(db *gorm.DB, entrepriseUUID, posUUID string, startDate, endDate time.Time) (models.RepartitionChartData, error) {
	var repartitionData models.RepartitionChartData

	// Construction de la requête selon le POS
	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS total_entrees,
			COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS total_sorties,
			COALESCE(SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END), 0) AS total_fond_caisse
		FROM caisse_items
		WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND pos_uuid = ?"
		args = []interface{}{entrepriseUUID, startDate, endDate, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{entrepriseUUID, startDate, endDate}
	}

	var totalEntrees, totalSorties, totalFondCaisse float64
	row := db.Raw(query, args...).Row()
	err := row.Scan(&totalEntrees, &totalSorties, &totalFondCaisse)
	if err != nil {
		return repartitionData, err
	}

	// Construction des données de répartition
	repartitionData.Series = []float64{totalEntrees, totalSorties, totalFondCaisse}

	// Calcul des pourcentages pour les labels
	total := totalEntrees + totalSorties + totalFondCaisse
	if total > 0 {
		entreesPercent := (totalEntrees / total) * 100
		sortiesPercent := (totalSorties / total) * 100
		fondPercent := (totalFondCaisse / total) * 100

		repartitionData.Labels = []string{
			fmt.Sprintf("Entrées (%.1f%%)", entreesPercent),
			fmt.Sprintf("Sorties (%.1f%%)", sortiesPercent),
			fmt.Sprintf("Fond de Caisse (%.1f%%)", fondPercent),
		}
	} else {
		repartitionData.Labels = []string{"Entrées (0%)", "Sorties (0%)", "Fond de Caisse (0%)"}
	}

	return repartitionData, nil
}

// ===============================
// ENDPOINTS COMPLEMENTAIRES
// ===============================

// GetFinancialSummary - Endpoint simplifié pour un résumé financier rapide
func GetFinancialSummary(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid", "")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID entreprise requis",
		})
	}

	// Période par défaut : ce mois
	startDate, endDate := mapPeriodToDateRange("this_month")

	// Chargement des KPI
	kpiData, err := loadKpiData(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des données: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    kpiData,
	})
}

// GetTopCaisses - Endpoint pour récupérer uniquement le top des caisses
func GetTopCaisses(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid", "")
	period := c.Query("period", "this_month")
	limitStr := c.Query("limit", "10")

	if entrepriseUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID entreprise requis",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// Calcul des dates selon la période
	startDate, endDate := mapPeriodToDateRange(period)

	// Chargement des top caisses avec limite personnalisée
	topCaisses, err := loadTopCaissesWithLimit(db, entrepriseUUID, posUUID, startDate, endDate, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des top caisses: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    topCaisses,
	})
}

// loadTopCaissesWithLimit - Version modifiée avec limite personnalisable
func loadTopCaissesWithLimit(db *gorm.DB, entrepriseUUID, posUUID string, startDate, endDate time.Time, limit int) ([]models.TopCaisse, error) {
	var topCaisses []models.TopCaisse

	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			c.uuid, 
			c.name, 
			COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) AS total_entrees,
			COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0) AS total_sorties,
			COUNT(ci.uuid) AS nombre_transactions,
			(COALESCE(c.montant_debut, 0) + 
			 COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) - 
			 COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0)) AS solde
		FROM caisses c
		LEFT JOIN caisse_items ci ON ci.caisse_uuid = c.uuid 
			AND ci.created_at BETWEEN ? AND ? 
			AND ci.deleted_at IS NULL
		WHERE c.entreprise_uuid = ? AND c.deleted_at IS NULL`

	if posUUID != "" {
		query = baseQuery + " AND c.pos_uuid = ?"
		args = []interface{}{startDate, endDate, entrepriseUUID, posUUID}
	} else {
		query = baseQuery
		args = []interface{}{startDate, endDate, entrepriseUUID}
	}

	query += fmt.Sprintf(`
		GROUP BY c.uuid, c.name, c.montant_debut
		ORDER BY solde DESC
		LIMIT %d`, limit)

	rows, err := db.Raw(query, args...).Rows()
	if err != nil {
		return topCaisses, err
	}
	defer rows.Close()

	for rows.Next() {
		var caisse models.TopCaisse
		err := rows.Scan(
			&caisse.UUID,
			&caisse.Name,
			&caisse.TotalEntrees,
			&caisse.TotalSorties,
			&caisse.NombreTransactions,
			&caisse.Solde,
		)
		if err != nil {
			continue
		}

		// Calcul de la performance (pourcentage de rentabilité)
		totalFlux := caisse.TotalEntrees + caisse.TotalSorties
		if totalFlux > 0 {
			caisse.Performance = ((caisse.TotalEntrees - caisse.TotalSorties) / totalFlux) * 100
		}

		topCaisses = append(topCaisses, caisse)
	}

	return topCaisses, nil
}

// GetAlertsAndRecommendations - Endpoint pour récupérer alertes et recommandations
func GetAlertsAndRecommendations(c *fiber.Ctx) error {
	db := database.DB

	userUUID := c.Query("user_uuid")
	posUUID := c.Query("pos_uuid", "")
	period := c.Query("period", "this_month")

	if userUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID utilisateur requis",
		})
	}

	// Récupération de l'entreprise de l'utilisateur
	var entrepriseUUID string
	err := db.Raw("SELECT entreprise_uuid FROM users WHERE uuid = ? AND deleted_at IS NULL", userUUID).
		Row().Scan(&entrepriseUUID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Utilisateur non trouvé",
		})
	}

	// Calcul des dates selon la période
	startDate, endDate := mapPeriodToDateRange(period)

	// Chargement des données nécessaires pour générer alertes et recommandations
	kpiData, err := loadKpiData(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des KPI: " + err.Error(),
		})
	}

	topCaisses, err := loadTopCaisses(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des top caisses: " + err.Error(),
		})
	}

	metriques, err := loadMetriques(db, entrepriseUUID, posUUID, startDate, endDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Erreur lors du chargement des métriques: " + err.Error(),
		})
	}

	// Génération des alertes et recommandations
	alerts := generateAlerts(kpiData, topCaisses)
	actionsRecommandees := generateActionsRecommandees(alerts, metriques)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"alerts":               alerts,
			"actions_recommandees": actionsRecommandees},
	})
}

// GetEvolutionSolde - Endpoint pour récupérer l'évolution du solde avec simulation
func GetEvolutionSolde(c *fiber.Ctx) error {
	db := database.DB

	// Récupération des paramètres
	userUUID := c.Query("user_uuid")
	posUUID := c.Query("pos_uuid", "")

	if userUUID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "UUID utilisateur requis",
		})
	}

	// Récupération de l'entreprise de l'utilisateur
	var entrepriseUUID string
	err := db.Raw("SELECT entreprise_uuid FROM users WHERE uuid = ? AND deleted_at IS NULL", userUUID).
		Row().Scan(&entrepriseUUID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Utilisateur non trouvé",
		})
	}
	var tresorerie []models.TresorerieEvolutionData

	// Variables pour stocker les données des graphiques
	var dates []string
	var entreesData []float64
	var sortiesData []float64
	var fondDeCaisseData []float64
	var soldeData []float64

	if posUUID == "" {
		err := db.Raw(`
			SELECT 
				DATE(created_at) AS date, 
				SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) AS entrees,
				SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) AS sorties,
				SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END) AS fond_caisses
			FROM caisse_items
			WHERE entreprise_uuid = ? AND
			deleted_at IS NULL
			GROUP BY DATE(created_at)
			ORDER BY DATE(created_at)
		`, entrepriseUUID).Scan(&tresorerie).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}
	} else {
		err := db.Raw(`
        SELECT 
            DATE(created_at) AS date, 
            SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) AS entrees,
            SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) AS sorties,
			SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END) AS fond_caisses
        FROM caisse_items
		WHERE entreprise_uuid = ? AND pos_uuid = ? AND
			deleted_at IS NULL
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at)
    `, entrepriseUUID, posUUID).Scan(&tresorerie).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}
	}
	for _, t := range tresorerie {
		dates = append(dates, t.Date.Format("2006-01-02"))
		entreesData = append(entreesData, t.Entrees)
		sortiesData = append(sortiesData, t.Sorties)
		fondDeCaisseData = append(fondDeCaisseData, t.FondCaisses)
		soldeData = append(soldeData, (t.Entrees-t.Sorties)+t.FondCaisses) // Calcul du solde
	}

	results := map[string]interface{}{
		"dates":            dates,
		"soldedata":        soldeData,
		"entreesdata":      entreesData,
		"sortiesdata":      sortiesData,
		"fonddecaissedata": fondDeCaisseData,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Evolution de la trésorerie",
		"data":    results,
	})
}

func SetupRepartitionPieChart(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut") // Format attendu : YYYY-MM-DD
	dateFin := c.Query("date_fin")     // Format attendu : YYYY-MM-DD

	var totalEntrees, totalSorties, totalFondDeCaisse float64

	if posUUID == "" {
		// Récupération des sommes globales pour le PieChart avec plage de dates
		err := db.Raw(`
        SELECT 
            COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS total_entrees,
            COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS total_sorties,
			COALESCE(SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END), 0) AS total_fond_caisse
        FROM caisse_items 
        WHERE entreprise_uuid = ? 
              AND created_at BETWEEN ? AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, dateDebut, dateFin).Row().
			Scan(&totalEntrees, &totalSorties, &totalFondDeCaisse)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}
	} else {
		// Récupération des sommes globales pour le PieChart avec plage de dates
		err := db.Raw(`
        SELECT 
            COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS total_entrees,
            COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS total_sorties,
			COALESCE(SUM(CASE WHEN type_transaction = 'MontantDebut' THEN montant ELSE 0 END), 0) AS total_fond_caisse
        FROM caisse_items 
        WHERE entreprise_uuid = ? AND pos_uuid = ? 
              AND created_at BETWEEN ? AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, posUUID, dateDebut, dateFin).Row().
			Scan(&totalEntrees, &totalSorties, &totalFondDeCaisse)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}
	}

	// Calcul du total global
	totalGeneral := totalEntrees + totalSorties + totalFondDeCaisse
	if totalGeneral == 0 {
		totalGeneral = 1 // Éviter la division par zéro
	}

	// Calcul des pourcentages pour le PieChart
	repartition := map[string]float64{
		"Entrees":    (totalEntrees / totalGeneral) * 100,
		"Sorties":    (totalSorties / totalGeneral) * 100,
		"FondCaisse": (totalFondDeCaisse / totalGeneral) * 100,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Répartition de la trésorerie pour PieChart avec plage de dates",
		"data":    repartition,
	})
}

func SetupPerformanceChart(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	dateDebut := c.Query("date_debut") // Format attendu : YYYY-MM-DD
	dateFin := c.Query("date_fin")     // Format attendu : YYYY-MM-DD

	var createdAt []string
	var entreesData []float64
	var sortiesData []float64 // Définition de la structure pour stocker les résultats
	var tresorerie []models.TresorerieDateData

	if posUUID == "" {
		// Correction de la requête SQL et récupération des données
		err := db.Raw(`
        SELECT 
            DATE(created_at) AS created_at, 
            SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) AS entrees,
            SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) AS sorties 
        FROM caisse_items 
        WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND
			deleted_at IS NULL
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at)
    `, entrepriseUUID, dateDebut, dateFin).Scan(&tresorerie).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}

	} else {
		// Correction de la requête SQL et récupération des données
		err := db.Raw(`
        SELECT 
            DATE(created_at) AS created_at, 
            SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) AS entrees,
            SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) AS sorties 
        FROM caisse_items 
        WHERE entreprise_uuid = ? AND pos_uuid = ? AND created_at BETWEEN ? AND ? AND
			deleted_at IS NULL
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at)
    `, entrepriseUUID, posUUID, dateDebut, dateFin).Scan(&tresorerie).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}

	}

	for _, t := range tresorerie {
		createdAt = append(createdAt, t.CreatedAt.Format("2006-01-02"))
		entreesData = append(entreesData, t.Entrees)
		sortiesData = append(sortiesData, t.Sorties)
	}

	results := map[string]interface{}{
		"created_at":  createdAt,
		"entreesdata": entreesData,
		"sortiesdata": sortiesData,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Répartition de la trésorerie",
		"data":    results,
	})
}

func SetupFluxChart(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	// Récupérer l'année en cours dynamiquement
	currentYear := time.Now().Year()

	var mois []string
	var entreesData []float64
	var sortiesData []float64
	// Définition de la structure pour stocker les résultats
	var tresorerie []models.TresorerieMoisData

	if posUUID == "" {
		// Requête SQL optimisée pour récupérer les données par mois avec leur nom (PostgreSQL)
		err := db.Raw(`
        SELECT 
            TO_CHAR(created_at, 'Month') AS mois,
            SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) AS entrees,
            SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) AS sorties
        FROM caisse_items 
        WHERE entreprise_uuid = ? AND EXTRACT(YEAR FROM created_at) = ? AND
			deleted_at IS NULL
        GROUP BY EXTRACT(MONTH FROM created_at), TO_CHAR(created_at, 'Month')
        ORDER BY EXTRACT(MONTH FROM created_at)
    `, entrepriseUUID, currentYear).Scan(&tresorerie).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}
	} else {
		// Requête SQL optimisée pour récupérer les données par mois avec leur nom (PostgreSQL)
		err := db.Raw(`
        SELECT 
            TO_CHAR(created_at, 'Month') AS mois,
            SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END) AS entrees,
            SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END) AS sorties
        FROM caisse_items 
        WHERE entreprise_uuid = ? AND pos_uuid = ? AND
			deleted_at IS NULL
              AND EXTRACT(YEAR FROM created_at) = ?
        GROUP BY EXTRACT(MONTH FROM created_at), TO_CHAR(created_at, 'Month')
        ORDER BY EXTRACT(MONTH FROM created_at)
    `, entrepriseUUID, posUUID, currentYear).Scan(&tresorerie).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des données",
				"error":   err.Error(),
			})
		}
	}

	// Traduction des noms des mois en français
	moisMap := map[string]string{
		"January ": "Jan", "February": "Fév", "March   ": "Mar", "April   ": "Avr",
		"May     ": "Mai", "June    ": "Juin", "July    ": "Juil", "August  ": "Août",
		"September": "Sep", "October ": "Oct", "November": "Nov", "December": "Déc",
	}

	for _, t := range tresorerie {
		// Nettoyer les espaces en trop du nom du mois
		monthName := strings.TrimSpace(t.Mois)
		if frenchMonth, exists := moisMap[t.Mois]; exists {
			mois = append(mois, frenchMonth)
		} else if frenchMonth, exists := moisMap[monthName]; exists {
			mois = append(mois, frenchMonth)
		} else {
			// Fallback: utiliser les 3 premières lettres
			mois = append(mois, monthName[:3])
		}
		entreesData = append(entreesData, t.Entrees)
		sortiesData = append(sortiesData, t.Sorties)
	}

	results := map[string]interface{}{
		"mois":        mois,
		"entreesdata": entreesData,
		"sortiesdata": sortiesData,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Flux de trésorerie par mois pour l'année en cours",
		"data":    results,
	})
}

func TopCaisse(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")

	var caisses []struct {
		UUID               string  `json:"uuid"`
		Name               string  `json:"name"`
		Solde              float64 `json:"solde"`
		TotalEntrees       float64 `json:"totalEntrees"`
		TotalSorties       float64 `json:"totalSorties"`
		NombreTransactions int     `json:"nombreTransactions"`
		Performance        float64 `json:"performance"`
	}

	if posUUID == "" {
		// Requête SQL pour récupérer les données des caisses
		err := db.Raw(`
        SELECT 
            c.uuid, 
            c.name, 
            COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) AS totalEntrees,
            COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0) AS totalSorties,
            COUNT(ci.uuid) AS nombreTransactions,
            (COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) 
            - COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0)) AS solde,
            ROUND(
                (COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) 
                - COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0))
                / NULLIF(COALESCE(SUM(ci.montant), 0), 0) * 100, 2
            ) AS performance
        FROM caisses c
        LEFT JOIN caisse_items ci ON ci.caisse_uuid = c.uuid
		WHERE c.entreprise_uuid = ?  AND
			c.deleted_at IS NULL
        GROUP BY c.uuid, c.name
        ORDER BY performance DESC
    `, entrepriseUUID).Scan(&caisses).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des caisses",
				"error":   err.Error(),
			})
		}
	} else {
		// Requête SQL pour récupérer les données des caisses
		err := db.Raw(`
        SELECT 
            c.uuid, 
            c.name, 
            COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) AS totalEntrees,
            COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0) AS totalSorties,
            COUNT(ci.uuid) AS nombreTransactions,
            (COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) 
            - COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0)) AS solde,
            ROUND(
                (COALESCE(SUM(CASE WHEN ci.type_transaction = 'Entree' THEN ci.montant ELSE 0 END), 0) 
                - COALESCE(SUM(CASE WHEN ci.type_transaction = 'Sortie' THEN ci.montant ELSE 0 END), 0))
                / NULLIF(COALESCE(SUM(ci.montant), 0), 0) * 100, 2
            ) AS performance
        FROM caisses c
        LEFT JOIN caisse_items ci ON ci.caisse_uuid = c.uuid
		WHERE c.entreprise_uuid = ? AND c.pos_uuid = ?  AND
			c.deleted_at IS NULL
        GROUP BY c.uuid, c.name
        ORDER BY performance DESC
    `, entrepriseUUID, posUUID).Scan(&caisses).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des caisses",
				"error":   err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Liste des caisses avec leurs statistiques",
		"data":    caisses,
	})
}

func FinancialMetrics(c *fiber.Ctx) error {
	db := database.DB

	entrepriseUUID := c.Query("entreprise_uuid")
	posUUID := c.Query("pos_uuid")
	// Récupération des paramètres de la plage de dates
	dateDebut := c.Query("date_debut") // Format attendu : YYYY-MM-DD
	dateFin := c.Query("date_fin")     // Format attendu : YYYY-MM-DD

	// **Conversion des dates en type time.Time**
	parsedDateDebut, err := time.Parse("2006-01-02", dateDebut)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Format de date incorrect pour date_debut",
		})
	}

	parsedDateFin, err := time.Parse("2006-01-02", dateFin)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Format de date incorrect pour date_fin",
		})
	}

	var totalEntrees, totalSorties float64
	var nombreTransactions int
	var velocite, velociteVariation, ratioLiquidite, tauxCroissance, efficacite float64

	if posUUID == "" {
		// Requête SQL pour récupérer les données dans la plage de dates spécifiée
		err = db.Raw(`
        SELECT 
				COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS totalEntrees,
				COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS totalSorties,
				COUNT(*) AS nombreTransactions
			FROM caisse_items
			WHERE entreprise_uuid = ? AND created_at BETWEEN ? AND ? AND
			deleted_at IS NULL
		`, entrepriseUUID, dateDebut, dateFin, dateDebut, dateFin).Row().
			Scan(&totalEntrees, &totalSorties, &nombreTransactions)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des métriques",
				"error":   err.Error(),
			})
		}

		// **Calcul des métriques**
		totalFlux := totalEntrees + totalSorties
		solde := totalEntrees - totalSorties

		if totalFlux == 0 {
			totalFlux = 1 // Évite la division par zéro
		}

		// **Vélocité** → Nombre moyen de transactions par mois
		dureeMois := float64(parsedDateFin.Month() - parsedDateDebut.Month() + 1)
		velocite = float64(nombreTransactions) / dureeMois

		// **Vélocité Variation** → Comparaison avec la période précédente
		var velocitePeriodePrecedente float64
		_ = db.Raw(`
        SELECT COUNT(*) FROM caisse_items 
        WHERE entreprise_uuid = ? AND created_at BETWEEN DATE_SUB(?, INTERVAL ? MONTH) AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, dateDebut, int(dureeMois), dateDebut).Row().Scan(&velocitePeriodePrecedente)

		velociteVariation = ((velocite - velocitePeriodePrecedente) / (velocitePeriodePrecedente + 1)) * 100

		// **Ratio de Liquidité** → Capacité à couvrir les sorties par les entrées
		ratioLiquidite = (totalEntrees / (totalSorties + 1)) * 100

		// **Taux de Croissance** → Évolution du solde sur la période
		var soldePeriodePrecedente float64
		_ = db.Raw(`
        SELECT COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) 
               - COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0)
        FROM caisse_items
        WHERE entreprise_uuid = ? AND created_at BETWEEN DATE_SUB(?, INTERVAL ? MONTH) AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, dateDebut, int(dureeMois), dateDebut).Row().Scan(&soldePeriodePrecedente)

		tauxCroissance = ((solde - soldePeriodePrecedente) / (soldePeriodePrecedente + 1)) * 100

		// **Efficacité** → Ratio des flux financiers sur le nombre de transactions
		efficacite = (totalFlux / (float64(nombreTransactions) + 1)) * 100

	} else {
		// Requête SQL pour récupérer les données dans la plage de dates spécifiée
		err = db.Raw(`
        SELECT 
            COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) AS totalEntrees,
            COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0) AS totalSorties,
            COUNT(*) AS nombreTransactions
        FROM caisse_items
        WHERE entreprise_uuid = ? AND pos_uuid = ? AND created_at BETWEEN ? AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, posUUID, dateDebut, dateFin, dateDebut, dateFin).Row().
			Scan(&totalEntrees, &totalSorties, &nombreTransactions)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Erreur lors de la récupération des métriques",
				"error":   err.Error(),
			})
		}

		// **Calcul des métriques**
		totalFlux := totalEntrees + totalSorties
		solde := totalEntrees - totalSorties

		if totalFlux == 0 {
			totalFlux = 1 // Évite la division par zéro
		}

		// **Vélocité** → Nombre moyen de transactions par mois
		dureeMois := float64(parsedDateFin.Month() - parsedDateDebut.Month() + 1)
		velocite = float64(nombreTransactions) / dureeMois

		// **Vélocité Variation** → Comparaison avec la période précédente
		var velocitePeriodePrecedente float64
		_ = db.Raw(`
        SELECT COUNT(*) FROM caisse_items 
        WHERE entreprise_uuid = ? AND pos_uuid = ? AND 
		created_at BETWEEN DATE_SUB(?, INTERVAL ? MONTH) AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, posUUID, dateDebut, int(dureeMois), dateDebut).Row().Scan(&velocitePeriodePrecedente)

		velociteVariation = ((velocite - velocitePeriodePrecedente) / (velocitePeriodePrecedente + 1)) * 100

		// **Ratio de Liquidité** → Capacité à couvrir les sorties par les entrées
		ratioLiquidite = (totalEntrees / (totalSorties + 1)) * 100

		// **Taux de Croissance** → Évolution du solde sur la période
		var soldePeriodePrecedente float64
		_ = db.Raw(`
        SELECT COALESCE(SUM(CASE WHEN type_transaction = 'Entree' THEN montant ELSE 0 END), 0) 
               - COALESCE(SUM(CASE WHEN type_transaction = 'Sortie' THEN montant ELSE 0 END), 0)
        FROM caisse_items
        WHERE entreprise_uuid = ? AND pos_uuid = ? AND 
		created_at BETWEEN DATE_SUB(?, INTERVAL ? MONTH) AND ? AND
			deleted_at IS NULL
    `, entrepriseUUID, posUUID, dateDebut, int(dureeMois), dateDebut).Row().Scan(&soldePeriodePrecedente)

		tauxCroissance = ((solde - soldePeriodePrecedente) / (soldePeriodePrecedente + 1)) * 100

		// **Efficacité** → Ratio des flux financiers sur le nombre de transactions
		efficacite = (totalFlux / (float64(nombreTransactions) + 1)) * 100
	}

	// **Résultat final**
	metriques := map[string]float64{
		"velocite":          velocite,
		"velociteVariation": velociteVariation,
		"ratioLiquidite":    ratioLiquidite,
		"tauxCroissance":    tauxCroissance,
		"efficacite":        efficacite,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Métriques financières calculées pour la période spécifiée",
		"data":    metriques,
	})
}
