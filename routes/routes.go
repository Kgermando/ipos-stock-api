package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/controllers/abonnements"
	"github.com/kgermando/ipos-stock-api/controllers/auth"
	"github.com/kgermando/ipos-stock-api/controllers/caisses"
	"github.com/kgermando/ipos-stock-api/controllers/clients"
	"github.com/kgermando/ipos-stock-api/controllers/commandes"
	"github.com/kgermando/ipos-stock-api/controllers/dashboard"
	"github.com/kgermando/ipos-stock-api/controllers/entreprises"
	"github.com/kgermando/ipos-stock-api/controllers/fournisseurs"
	"github.com/kgermando/ipos-stock-api/controllers/plats"
	"github.com/kgermando/ipos-stock-api/controllers/pos"
	"github.com/kgermando/ipos-stock-api/controllers/products"
	"github.com/kgermando/ipos-stock-api/controllers/stocks"
	tablebox "github.com/kgermando/ipos-stock-api/controllers/tableBox"
	"github.com/kgermando/ipos-stock-api/controllers/users"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {

	api := app.Group("/api", logger.New())

	// Authentification controller
	a := api.Group("/auth")
	a.Post("/register", auth.Register)
	a.Post("/login", auth.Login)
	a.Post("/forgot-password", auth.Forgot)
	a.Get("/verify-reset-token/:token", auth.VerifyResetToken)
	a.Post("/reset/:token", auth.ResetPassword)

	a.Post("/entreprise/create", entreprises.CreateEntreprise)
	a.Put("/entreprise/update/:uuid", entreprises.UpdateEntreprise)

	// app.Use(middlewares.IsAuthenticated)

	a.Get("/user", auth.AuthUser)
	a.Put("/profil/info", auth.UpdateInfo)
	a.Put("/change-password", auth.ChangePassword)
	a.Post("/logout", auth.Logout)

	// Entreprise controller
	e := api.Group("/entreprises")
	e.Get("/all", entreprises.GetAllEntreprises)
	e.Get("/all/paginate", entreprises.GetPaginatedEntreprise)
	e.Get("/get/:uuid", entreprises.GetEntreprise)
	e.Post("/create", entreprises.CreateEntreprise)
	e.Put("/update/:uuid", entreprises.UpdateEntreprise)
	e.Delete("/delete/:uuid", entreprises.DeleteEntreprise)

	// Abonnements controller
	ab := api.Group("/abonnements")
	ab.Get("/all", abonnements.GetAllAbonnements)
	ab.Get("/all/paginate", abonnements.GetPaginatedAbonnements)
	ab.Get("/all/paginate/:entreprise_uuid", abonnements.GetPaginatedAbonnementsEntreprise)
	ab.Get("/get/:uuid", abonnements.GetAbonnement)
	ab.Post("/create", abonnements.CreateAbonnement)
	ab.Put("/update/:uuid", abonnements.UpdateAbonnement)
	ab.Delete("/delete/:uuid", abonnements.DeleteAbonnement)
	ab.Put("/update-statut/:uuid", abonnements.UpdateStatutAbonnement)
	ab.Get("/current", abonnements.GetAbonnementActuel)
	ab.Get("/verify/:uuid", abonnements.VerifierValiditeAbonnement)
	ab.Get("/expiring", abonnements.GetAbonnementsExpirant)
	ab.Get("/statistics", abonnements.GetStatistiquesAbonnements)

	// Users controller
	u := api.Group("/users")
	u.Get("/all", users.GetAllUsers)
	u.Get("/all/paginate", users.GetPaginatedUsersSupport)
	u.Get("/all/:entreprise_uuid", users.GetAllUsersById)
	u.Get("/all/paginate/nosearch", users.GetPaginatedNoSerach)
	u.Get("/:entreprise_uuid/all/paginate", users.GetPaginatedUsers)
	u.Get("/:entreprise_uuid/:pos_uuid/all/paginate", users.GetPaginatedUserByPosUUID)
	u.Get("/get/:uuid", users.GetUser)
	u.Post("/create", users.CreateUser)
	u.Put("/update/:uuid", users.UpdateUser)
	u.Delete("/delete/:uuid", users.DeleteUser)

	// POS controller
	p := api.Group("/pos")
	// p.Get("/all", pos.GetAllPoss)
	p.Get("/all/paginate", pos.GetPaginatedPos)
	p.Get("/all/:entreprise_uuid", pos.GetAllPosByUUId)
	p.Get("/:entreprise_uuid/all/paginate", pos.GetPaginatedPosByUUID)
	p.Get("/get/:uuid", pos.GetPos)
	p.Post("/create", pos.CreatePos)
	p.Put("/update/:uuid", pos.UpdatePos)
	p.Delete("/delete/:uuid", pos.DeletePos)

	// Caisses controller
	cais := api.Group("/caisses")
	cais.Get("/:entreprise_uuid/all/total", caisses.GetTotalAllCaisses)
	cais.Get("/:entreprise_uuid/all", caisses.GetAllCaisses)
	cais.Get("/:entreprise_uuid/:pos_uuid/all", caisses.GetAllCaisseByPos)
	cais.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", caisses.GetDataSynchronisation)
	cais.Get("/get/:uuid", caisses.GetCaisse)
	cais.Post("/create", caisses.CreateCaisse)
	cais.Put("/update/:uuid", caisses.UpdateCaisse)
	cais.Delete("/delete/:uuid", caisses.DeleteCaisse)

	// Caisse item Controller
	caisseItem := api.Group("/caisse-items")
	caisseItem.Get("/:entreprise_uuid/:caisse_uuid/all/paginate", caisses.GetPaginatedCaisseItems)
	caisseItem.Get("/:entreprise_uuid/:caisse_uuid/all", caisses.GetAllCaisseItems)
	caisseItem.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", caisses.GetDataSynchronisationCaisseItem)
	caisseItem.Get("/get/:uuid", caisses.GetCaisseItem)
	caisseItem.Post("/create", caisses.CreateCaisseItem)
	caisseItem.Put("/update/:uuid", caisses.UpdateCaisseItem)
	caisseItem.Delete("/delete/:uuid", caisses.DeleteCaisseItem)

	// Product controller
	pr := api.Group("/products")
	pr.Get("/:entreprise_uuid/all/paginate", products.GetPaginatedProductEntreprise)
	pr.Get("/:entreprise_uuid/:pos_uuid/all", products.GetAllProducts)
	pr.Get("/:entreprise_uuid/:pos_uuid/all/paginate", products.GetPaginatedProductByPosUUID)
	pr.Get("/:entreprise_uuid/:pos_uuid/all/search", products.GetAllProductBySearch)
	pr.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", products.GetDataSynchronisation)
	pr.Get("/get/:uuid", products.GetProduct)
	pr.Get("/excel-template", products.GenerateProductExcelTemplate)
	pr.Get("/excel-format-info", products.GetExcelFormatInfo)
	pr.Post("/create", products.CreateProduct)
	pr.Post("/:entreprise_uuid/:pos_uuid/upload-excel", products.UploadProductsFromExcel)
	pr.Put("/update/:uuid", products.UpdateProduct)
	pr.Put("/update/stock/:uuid", products.UpdateProductStockDispo)
	pr.Put("/update/stock-endommage/:uuid", products.UpdateProductStockEndommage)
	pr.Put("/update/restitution/:uuid", products.UpdateProductRestitution)
	pr.Delete("/delete/:uuid", products.DeleteProduct)

	// Plat controller
	pl := api.Group("/plats")
	pl.Get("/:entreprise_uuid/all/paginate", plats.GetPaginatedPlatEntreprise)
	pl.Get("/:entreprise_uuid/:pos_uuid/all", plats.GetAllPlats)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/available", plats.GetAvailablePlats)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/paginate", plats.GetPaginatedPlatByPosUUID)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/search", plats.GetAllPlatBySearch)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", plats.GetDataSynchronisation)
	pl.Get("/get/:uuid", plats.GetPlat)
	pl.Post("/create", plats.CreatePlat)
	pl.Put("/update/:uuid", plats.UpdatePlat)
	pl.Put("/update/availability/:uuid", plats.UpdatePlatAvailability)
	pl.Delete("/delete/:uuid", plats.DeletePlat)

	// TableBox controller
	tb := api.Group("/tablebox")
	tb.Get("/:entreprise_uuid/all/paginate", tablebox.GetPaginatedTableBoxEntreprise)
	tb.Get("/:entreprise_uuid/:pos_uuid/all", tablebox.GetAllTableBoxs)
	tb.Get("/:entreprise_uuid/:pos_uuid/all/paginate", tablebox.GetPaginatedTableBoxByPosUUID)
	tb.Get("/:entreprise_uuid/:pos_uuid/all/search", tablebox.GetAllTableBoxBySearch)
	tb.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", tablebox.GetDataSynchronisation)
	tb.Get("/:entreprise_uuid/:pos_uuid/category/:category", tablebox.GetTableBoxsByCategory)
	tb.Get("/get/:uuid", tablebox.GetTableBox)
	tb.Post("/create", tablebox.CreateTableBox)
	tb.Put("/update/:uuid", tablebox.UpdateTableBox)
	tb.Delete("/delete/:uuid", tablebox.DeleteTableBox)

	// Stock controller
	s := api.Group("/stocks")
	s.Get("/all/paginate/:product_uuid", stocks.GetPaginatedStock)
	s.Get("/all/total/:product_uuid", stocks.GetTotalStock)
	s.Get("/all/get/:product_uuid", stocks.GetStockMargeBeneficiaire)
	s.Get("/all/:product_uuid", stocks.GetAllStocks)
	s.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", stocks.GetDataSynchronisationStock)
	s.Get("/get/:uuid", stocks.GetStock)
	s.Post("/create", stocks.CreateStock)
	s.Put("/update/:uuid", stocks.UpdateStock)
	s.Delete("/delete/:uuid", stocks.DeleteStock)

	// StockEndommage controller
	se := api.Group("/stock-endommages")
	se.Get("/all/paginate/:product_uuid", stocks.GetPaginatedStockEndommage)
	se.Get("/all/total/:product_uuid", stocks.GetTotalStockEndommage)
	se.Get("/all/:product_uuid", stocks.GetAllStockEndommages)
	se.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", stocks.GetDataSynchronisationStockEndommage)
	se.Get("/get/:uuid", stocks.GetStockEndommage)
	se.Post("/create", stocks.CreateStockEndommage)
	se.Put("/update/:uuid", stocks.UpdateStockEndommage)
	se.Delete("/delete/:uuid", stocks.DeleteStockEndommage)

	// Restitution controller
	re := api.Group("/restitutions")
	re.Get("/all/paginate/:product_uuid", stocks.GetPaginatedRestitution)
	re.Get("/all/total/:product_uuid", stocks.GetTotalRestitution)
	re.Get("/all/:product_uuid", stocks.GetAllRestitutions)
	re.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", stocks.GetDataSynchronisationRestitution)
	re.Get("/get/:uuid", stocks.GetRestitution)
	re.Post("/create", stocks.CreateRestitution)
	re.Put("/update/:uuid", stocks.UpdateRestitution)
	re.Delete("/delete/:uuid", stocks.DeleteRestitution)

	// Client controller
	cl := api.Group("/clients")
	cl.Get("/:entreprise_uuid/all", clients.GetAllClients)
	cl.Get("/:entreprise_uuid/:pos_uuid/all/paginate", clients.GetPaginatedClient)
	cl.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", clients.GetDataSynchronisation)
	cl.Get("/get/:uuid", clients.GetClient)
	cl.Post("/create", clients.CreateClient)
	cl.Post("/uploads", clients.UploadCsvDataClient)
	cl.Put("/update/:uuid", clients.UpdateClient)
	cl.Delete("/delete/:uuid", clients.DeleteClient)

	// Fournisseur controller
	fs := api.Group("/fournisseurs")
	fs.Get("/:entreprise_uuid/all", fournisseurs.GetAllFournisseurs)
	fs.Get("/:entreprise_uuid/:pos_uuid/all/paginate", fournisseurs.GetPaginatedFournisseur)
	fs.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", fournisseurs.GetDataSynchronisation)
	fs.Get("/get/:uuid", fournisseurs.GetFournisseur)
	fs.Post("/create", fournisseurs.CreateFournisseur)
	fs.Put("/update/:uuid", fournisseurs.UpdateFournisseur)
	fs.Delete("/delete/:uuid", fournisseurs.DeleteFournisseur)

	// Commande controller
	cmd := api.Group("/commandes")
	cmd.Get("/:entreprise_uuid/all/paginate", commandes.GetPaginatedCommandeEntreprise)
	cmd.Get("/:entreprise_uuid/:pos_uuid/all", commandes.GetAllCommandes)
	cmd.Get("/:entreprise_uuid/:pos_uuid/all/paginate", commandes.GetPaginatedCommandePOS)
	cmd.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", commandes.GetDataSynchronisation)
	cmd.Get("/get/:uuid", commandes.GetCommande)
	cmd.Post("/create", commandes.CreateCommande)
	cmd.Put("/update/:uuid", commandes.UpdateCommande)
	cmd.Delete("/delete/:uuid", commandes.DeleteCommande)

	// Commande line controller
	cmdl := api.Group("/commande-lines")
	cmdl.Get("/all", commandes.GetAllCommandeLines)
	cmdl.Get("/all/paginate/:commande_uuid", commandes.GetPaginatedCommandeLineByID)
	cmdl.Get("/all/total/:product_uuid", commandes.GetTotalCommandeLine)
	cmdl.Get("/all/:commande_uuid", commandes.GetAllCommandeLineByUUId)
	cmdl.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", commandes.GetDataSynchronisationCommandeLine)
	cmdl.Get("/get/:uuid", commandes.GetCommandeLine)
	cmdl.Post("/create", commandes.CreateCommandeLine)
	cmdl.Put("/update/:uuid", commandes.UpdateCommandeLine)
	cmdl.Delete("/delete/:uuid", commandes.DeleteCommandeLine)

	// Dashboard
	// Tresorerie - Nouveau système complet de dashboard
	dash := api.Group("/dashboard")

	tresorerie := dash.Group("/tresoreries")
	// Endpoint principal pour récupérer toutes les données du dashboard
	tresorerie.Get("/", dashboard.GetDashboardTresorerie)

	// Endpoints pour les graphiques
	tresorerie.Get("/evolution-chart", dashboard.GetEvolutionChartData)
	tresorerie.Get("/repartition-chart", dashboard.GetRepartitionChartData)
	// Endpoints complémentaires
	tresorerie.Get("/financial-summary", dashboard.GetFinancialSummary)
	tresorerie.Get("/top-caisses", dashboard.GetTopCaisses)
	tresorerie.Get("/alerts-recommendations", dashboard.GetAlertsAndRecommendations)
	tresorerie.Get("/flux-chart", dashboard.SetupFluxChart)
	tresorerie.Get("/performance-caisse", dashboard.PerformanceCaisse)

	// KPI Dashboard - Nouveaux endpoints améliorés
	kpi := dash.Group("/kpi")

	// Endpoints principaux pour le frontend Angular
	kpi.Get("/global-kpis", dashboard.GlobalKpis)
	kpi.Get("/evolution-vente", dashboard.GetEvolutionVente)
	kpi.Get("/performance-vente", dashboard.GetPerformanceVente)
	kpi.Get("/best-selling-product", dashboard.GetBestSellingProduct)
	kpi.Get("/stock-kpis", dashboard.GetStockKpis)
	kpi.Get("/stock-faible", dashboard.GetStockFaible)
	kpi.Get("/stock-chart", dashboard.SetupStockChart)
	kpi.Get("/alerts-kpis", dashboard.GetAlertsKpis)

	// Endpoints legacy pour compatibilité
	kpi.Get("/global", dashboard.GlobalKpiSummary)
	kpi.Get("/best-selling-prroduct", dashboard.BestSellingProduct)

}
