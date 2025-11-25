package models

import "time"

// DashboardStats représente les statistiques principales du dashboard
type DashboardStats struct {
	TotalArticles                  int64   `json:"totalArticles"`
	ArticlesRuptureStock           int64   `json:"articlesRuptureStock"`
	ArticlesRuptureStockPercentage int     `json:"articlesRuptureStockPercentage"`
	TotalVentes                    int64   `json:"totalVentes"`
	TotalVentesPercentage          int     `json:"totalVentesPercentage"`
	TotalMontantVendu              float64 `json:"totalMontantVendu"`
	TotalMontantVenduPercentage    int     `json:"totalMontantVenduPercentage"`
}

// SalesChartData représente les données pour le graphique de ventes
type SalesChartData struct {
	Dates          []string  `json:"dates"`
	TotalCommandes []int64   `json:"totalCommandes"`
	MontantVendu   []float64 `json:"montantVendu"`
	GainObtenu     []float64 `json:"gainObtenu"`
}

// PlatChartData représente les données pour le graphique donut des plats
type PlatChartData struct {
	Labels      []string  `json:"labels"`
	Series      []float64 `json:"series"`
	Percentages []float64 `json:"percentages"`
}

// ProductChartData représente les données pour le graphique donut des produits
type ProductChartData struct {
	Labels      []string  `json:"labels"`
	Series      []float64 `json:"series"`
	Percentages []float64 `json:"percentages"`
}

// StockAlert représente une alerte de stock
type StockAlert struct {
	UUID       string  `json:"uuid"`
	Name       string  `json:"name"`
	Reference  string  `json:"reference"`
	UniteVente string  `json:"unite_vente"`
	Stock      float64 `json:"stock"`
	AlertType  string  `json:"alertType"` // "rupture" ou "avertissement"
	Image      string  `json:"image"`
	PrixVente  float64 `json:"prix_vente"`
}

// StockRotationData représente les données de rotation de stock
type StockRotationData struct {
	ProductNames  []string  `json:"productNames"`
	RotationRates []float64 `json:"rotationRates"`
	Categories    []string  `json:"categories"`
	Colors        []string  `json:"colors"`
}

// PlatStatistics représente les statistiques des plats
type PlatStatistics struct {
	TotalPlats       int64   `json:"totalPlats"`
	TotalClients     int64   `json:"totalClients"`
	QuantitesVendues int64   `json:"quantitesVendues"`
	ChiffresAffaires float64 `json:"chiffresAffaires"`
}

// LivraisonStats représente les statistiques des livraisons
type LivraisonStats struct {
	TotalLivraisons      int64   `json:"totalLivraisons"`
	EnCours              int64   `json:"enCours"`
	Effectuees           int64   `json:"effectuees"`
	Annulees             int64   `json:"annulees"`
	EnCoursPercentage    float64 `json:"enCoursPercentage"`
	EffectueesPercentage float64 `json:"effectueesPercentage"`
	AnnuleesPercentage   float64 `json:"annuleesPercentage"`
	TotalRevenu          float64 `json:"totalRevenu"`
	RevenuMoyen          float64 `json:"revenuMoyen"`
	TauxReussite         float64 `json:"tauxReussite"`
}

// LivraisonZoneData représente les données des zones de livraison
type LivraisonZoneData struct {
	ZoneName         string  `json:"zoneName"`
	NombreLivraisons int64   `json:"nombreLivraisons"`
	Revenu           float64 `json:"revenu"`
}

// LivreurPerformance représente les performances des livreurs
type LivreurPerformance struct {
	UUID            string  `json:"uuid"`
	Name            string  `json:"name"`
	TotalLivraisons int64   `json:"totalLivraisons"`
	Effectuees      int64   `json:"effectuees"`
	EnCours         int64   `json:"enCours"`
	Annulees        int64   `json:"annulees"`
	TauxReussite    float64 `json:"tauxReussite"`
}

 
// TopTransactions représente les tops transactions groupées
type TopTransactions struct {
	TopEntrees []TopTransaction `json:"topEntrees"`
	TopSorties []TopTransaction `json:"topSorties"`
}



// CaisseStatistics représente les statistiques principales de la caisse
type CaisseStatistics struct {
	// Statistiques principales
	SoldeCaisse         float64 `json:"soldeCaisse"`
	TotalEntrees        float64 `json:"totalEntrees"`
	TotalSorties        float64 `json:"totalSorties"`
	NombreTransactions  int     `json:"nombreTransactions"`
	MontantDebut        float64 `json:"montantDebut"` // Montant d'ouverture / dépôt initial

	// Moyennes et ratios
	MoyenneEntree      float64 `json:"moyenneEntree"`
	MoyenneSortie      float64 `json:"moyenneSortie"`
	RatioEntreeSortie  float64 `json:"ratioEntreeSortie"`
	TauxLiquidite      float64 `json:"tauxLiquidite"`

	// Évolution
	EvolutionJournaliere float64 `json:"evolutionJournaliere"`
	EvolutionPercentage  float64 `json:"evolutionPercentage"`
	Tendance             string  `json:"tendance"` // 'hausse' | 'baisse' | 'stable'

	// Analyse temporelle
	JourLePlusActif         string  `json:"jourLePlusActif"`
	HeureLaPlusActive       string  `json:"heureLaPlusActive"`
	NombreTransactionsParJour float64 `json:"nombreTransactionsParJour"`
}

// FluxTresorerieData représente les données pour le graphique de flux de trésorerie
type FluxTresorerieData struct {
	Dates   []string  `json:"dates"`
	Entrees []float64 `json:"entrees"`
	Sorties []float64 `json:"sorties"`
	Soldes  []float64 `json:"soldes"`
}

// RepartitionTransactionsData représente les données pour le graphique de répartition
type RepartitionTransactionsData struct {
	Labels      []string  `json:"labels"`
	Values      []float64 `json:"values"`
	Percentages []float64 `json:"percentages"`
	Colors      []string  `json:"colors"`
}

// TopTransaction représente une transaction importante (entrée ou sortie)
type TopTransaction struct {
	Libelle   string    `json:"libelle"`
	Montant   float64   `json:"montant"`
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
	Reference string    `json:"reference"`
}

// TopTransactionsResponse représente la réponse pour les top transactions
type TopTransactionsResponse struct {
	TopEntrees []TopTransaction `json:"topEntrees"`
	TopSorties []TopTransaction `json:"topSorties"`
}

// CategorieAnalysis représente l'analyse d'une catégorie de transactions
type CategorieAnalysis struct {
	Categorie           string  `json:"categorie"`
	TotalMontant        float64 `json:"totalMontant"`
	NombreTransactions  int     `json:"nombreTransactions"`
	Pourcentage         float64 `json:"pourcentage"`
	Moyenne             float64 `json:"moyenne"`
	Tendance            string  `json:"tendance"` // 'croissante' | 'décroissante' | 'stable'
}

// PrevisionTresorerie représente une prévision de trésorerie
type PrevisionTresorerie struct {
	Date            string  `json:"date"`
	PrevisionEntree float64 `json:"previsionEntree"`
	PrevisionSortie float64 `json:"previsionSortie"`
	PrevisionSolde  float64 `json:"previsionSolde"`
	Confiance       float64 `json:"confiance"` // Niveau de confiance de la prévision (0-100)
}

// EvolutionData structure interne pour calculer l'évolution
type EvolutionData struct {
	Montant    float64
	Percentage float64
	Tendance   string
}
