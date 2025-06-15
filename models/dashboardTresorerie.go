package models

import "time"

// ===============================
// STRUCTURES DE DONNÉES DASHBOARD TRÉSORERIE
// ===============================

// Structures pour les données KPI
type KpiData struct {
	SoldeTotal        float64 `json:"solde_total"`
	TotalEntrees      float64 `json:"total_entrees"`
	TotalSorties      float64 `json:"total_sorties"`
	TotalFondDeCaisse float64 `json:"total_fond_de_caisse"`
	NombreEntrees     int     `json:"nombre_entrees"`
	NombreSorties     int     `json:"nombre_sorties"`
	VariationSolde    float64 `json:"variation_solde"`
}

// Structures pour les Top Caisses
type TopCaisse struct {
	UUID               string  `json:"uuid"`
	Name               string  `json:"name"`
	Solde              float64 `json:"solde"`
	TotalEntrees       float64 `json:"total_entrees"`
	TotalSorties       float64 `json:"total_sorties"`
	NombreTransactions int     `json:"nombre_transactions"`
	Performance        float64 `json:"performance"`
}

// Structures pour les Alerts
type TresorerieAlert struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Structures pour les Métriques
type TresorerieMetriques struct {
	Velocite          float64 `json:"velocite"`
	VelociteVariation float64 `json:"velocite_variation"`
	RatioLiquidite    float64 `json:"ratio_liquidite"`
	TauxCroissance    float64 `json:"taux_croissance"`
	Efficacite        float64 `json:"efficacite"`
}

// Structures pour les Actions Recommandées
type ActionRecommandee struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Structures pour les Stats Footer
type TresorerieStatsFooter struct {
	TransactionsAujourdhui  int       `json:"transactions_aujourdhui"`
	MontantMoyenTransaction float64   `json:"montant_moyen_transaction"`
	TempsMoyenTraitement    float64   `json:"temps_moyen_traitement"`
	DerniereMiseAJour       time.Time `json:"derniere_mise_a_jour"`
}

// Structure principale du dashboard
type DashboardTresorerieData struct {
	KpiData             KpiData               `json:"kpi_data"`
	TopCaisses          []TopCaisse           `json:"top_caisses"`
	Alerts              []TresorerieAlert     `json:"alerts"`
	Metriques           TresorerieMetriques   `json:"metriques"`
	ActionsRecommandees []ActionRecommandee   `json:"actions_recommandees"`
	StatsFooter         TresorerieStatsFooter `json:"stats_footer"`
}

// Structures pour les graphiques
type EvolutionChartData struct {
	Dates           []string  `json:"dates"`
	SoldeData       []float64 `json:"solde_data"`
	EntreesData     []float64 `json:"entrees_data"`
	SortiesData     []float64 `json:"sorties_data"`
	FondCaissesData []float64 `json:"fond_caisses_data"`
}

type RepartitionChartData struct {
	Series []float64 `json:"series"`
	Labels []string  `json:"labels"`
}

// Structure pour les données d'un jour (utilisée dans les calculs internes)
type DayData struct {
	Date        time.Time `json:"date"`
	Entrees     float64   `json:"entrees"`
	Sorties     float64   `json:"sorties"`
	FondCaisses float64   `json:"fond_caisses"`
}

// Structure pour la trésorerie (réutilisable)
type Tresorerie struct {
	TotalEntrees float64 `json:"total_entrees"`
	TotalSorties float64 `json:"total_sorties"`
	SoldeNet     float64 `json:"solde_net"`
	Transactions int     `json:"transactions"`
}

// Structure pour les données de trésorerie par date
type TresorerieDateData struct {
	CreatedAt time.Time `json:"created_at"`
	Entrees   float64   `json:"entrees"`
	Sorties   float64   `json:"sorties"`
}

// Structure pour les données de trésorerie par mois
type TresorerieMoisData struct {
	Mois    string  `json:"mois"`
	Entrees float64 `json:"entrees"`
	Sorties float64 `json:"sorties"`
}

// Structure pour les données de trésorerie avec détails complets pour l'évolution
type TresorerieEvolutionData struct {
	Date        time.Time `json:"date"`
	Entrees     float64   `json:"entrees"`
	Sorties     float64   `json:"sorties"`
	FondCaisses float64   `json:"fond_caisses"`
}
