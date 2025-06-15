package models

// TopProduct représente un produit avec ses statistiques de vente/stock
type TopProduct struct {
	UUID      string  `json:"uuid"`
	Name      string  `json:"name"`
	Quantite  float64 `json:"quantite"`
	Valeur    float64 `json:"valeur"`
	Stock     float64 `json:"stock"`
	Variation float64 `json:"variation"`
}

// StockKpiData représente les KPI liés au stock
type StockKpiData struct {
	TotalStock       float64 `json:"totalStock"`
	StockValeur      float64 `json:"stockValeur"`
	StockEndommage   float64 `json:"stockEndommage"`
	StockAlertes     int64   `json:"stockAlertes"`
	StockRestitution float64 `json:"stockRestitution"`
	TauxRotation     float64 `json:"tauxRotation"`
	StockVariation   float64 `json:"stockVariation"`
}

// Sale représente les données de vente par période
type Sale struct {
	Date     string  `json:"date"`
	Commande float64 `json:"commande"`
	Vente    float64 `json:"vente"`
	Tva      float64 `json:"tva"`
}

// GlobalKpis représente les KPI globaux du dashboard
type GlobalKpis struct {
	TotalRevenue               float64 `json:"totalRevenue"`
	TotalRevenueVariation      float64 `json:"totalRevenueVariation"`
	TotalCommandes             int64   `json:"totalCommandes"`
	TotalCommandesVariation    float64 `json:"totalCommandesVariation"`
	TotalProduits              int64   `json:"totalProduits"`
	TotalProduitsVariation     float64 `json:"totalProduitsVariation"`
	AverageOrderValue          float64 `json:"averageOrderValue"`
	AverageOrderValueVariation float64 `json:"averageOrderValueVariation"`
}

// SalesKpis représente les KPI de vente
type SalesKpis struct {
	VentesAujourdhui    float64      `json:"ventesAujourdhui"`
	VentesHier          float64      `json:"ventesHier"`
	VentesSemaine       float64      `json:"ventesSemaine"`
	VentesMois          float64      `json:"ventesMois"`
	VariationJour       float64      `json:"variationJour"`
	VariationSemaine    float64      `json:"variationSemaine"`
	VariationMois       float64      `json:"variationMois"`
	BestSellingProducts []TopProduct `json:"bestSellingProducts"`
}

// StockKpis représente les KPI de stock
type StockKpis struct {
	TotalStock          float64      `json:"totalStock"`
	StockValeur         float64      `json:"stockValeur"`
	StockEndommage      float64      `json:"stockEndommage"`
	StockAlertes        int64        `json:"stockAlertes"`
	StockRestitution    float64      `json:"stockRestitution"`
	TauxRotation        float64      `json:"tauxRotation"`
	StockVariation      float64      `json:"stockVariation"`
	ProduitsStockFaible []TopProduct `json:"produitsStockFaible"`
}

// AlertsKpis représente les KPI d'alertes
type AlertsKpis struct {
	AlertesStock       int64   `json:"alertesStock"`
	AlertesExpiration  int64   `json:"alertesExpiration"`
	AlertesVentes      int64   `json:"alertesVentes"`
	AlertesPerformance int64   `json:"alertesPerformance"`
	AlertesCritiques   []Alert `json:"alertesCritiques"`
	TendanceAlertes    float64 `json:"tendanceAlertes"`
}

// Alert représente une alerte du système
type Alert struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Level       string  `json:"level"`
	Title       string  `json:"title"`
	Message     string  `json:"message"`
	ProductName string  `json:"productName,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Threshold   float64 `json:"threshold,omitempty"`
	CreatedAt   string  `json:"createdAt"`
}

// PerformanceKpis représente les KPI de performance
type PerformanceKpis struct {
	MargeGlobale       float64 `json:"margeGlobale"`
	MargeBrute         float64 `json:"margeBrute"`
	MargeNette         float64 `json:"margeNette"`
	CroissanceVentes   float64 `json:"croissanceVentes"`
	PerformancePOS     float64 `json:"performancePOS"`
	EfficaciteStock    float64 `json:"efficaciteStock"`
	SatisfactionClient float64 `json:"satisfactionClient"`
}

// TrendData représente les données de tendance
type TrendData struct {
	Label     string  `json:"label"`
	Value     float64 `json:"value"`
	Variation float64 `json:"variation"`
}

// PerformanceVente représente les données de performance de vente
type PerformanceVente struct {
	TotalRevenue float64 `json:"totalRevenue"`
	TotalCost    float64 `json:"totalCost"`
	MargeBrute   float64 `json:"margeBrute"`
	MargeGlobale float64 `json:"margeGlobale"`
	Performance  float64 `json:"performance"`
}

// EvolutionVente représente l'évolution des ventes
type EvolutionVente struct {
	VentesAujourdhui float64 `json:"ventesAujourdhui"`
	VentesHier       float64 `json:"ventesHier"`
	VentesSemaine    float64 `json:"ventesSemaine"`
	VentesMois       float64 `json:"ventesMois"`
	VariationJour    float64 `json:"variationJour"`
	VariationSemaine float64 `json:"variationSemaine"`
	VariationMois    float64 `json:"variationMois"`
}

// ChartData représente les données pour les graphiques
type ChartData struct {
	Series []float64 `json:"series"`
	Labels []string  `json:"labels"`
	Colors []string  `json:"colors"`
}
