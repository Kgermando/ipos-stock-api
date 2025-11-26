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
	"github.com/kgermando/ipos-stock-api/controllers/livraisons"
	"github.com/kgermando/ipos-stock-api/controllers/livreurs"
	"github.com/kgermando/ipos-stock-api/controllers/plats"
	"github.com/kgermando/ipos-stock-api/controllers/pos"
	"github.com/kgermando/ipos-stock-api/controllers/products"
	"github.com/kgermando/ipos-stock-api/controllers/reservations"
	"github.com/kgermando/ipos-stock-api/controllers/stocks"
	tablebox "github.com/kgermando/ipos-stock-api/controllers/tableBox"
	"github.com/kgermando/ipos-stock-api/controllers/users"

	"github.com/kgermando/ipos-stock-api/controllers/zones"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {

	api := app.Group("/api", logger.New())

	// ============================================================
	// AUTHENTIFICATION ROUTES (Public)
	// ============================================================
	a := api.Group("/auth")

	// Public authentication routes
	a.Post("/register", auth.Register)
	a.Post("/login", auth.Login)
	a.Post("/forgot-password", auth.Forgot)
	a.Get("/verify-reset-token/:token", auth.VerifyResetToken)
	a.Post("/reset/:token", auth.ResetPassword)

	// Enterprise management in auth context
	a.Post("/entreprise/create", entreprises.CreateEntreprise)
	a.Put("/entreprise/update/:uuid", entreprises.UpdateEntreprise)

	// Protected authentication routes
	// app.Use(middlewares.IsAuthenticated)
	a.Get("/user", auth.AuthUser)
	a.Put("/profil/info", auth.UpdateInfo)
	a.Put("/change-password", auth.ChangePassword)
	a.Post("/logout", auth.Logout)

	// ============================================================
	// DASHBOARD ROUTES
	// ============================================================
	dash := api.Group("/dashboard")
	main := dash.Group("/main")
	main.Get("/stats", dashboard.GetDashboardStats)
	main.Get("/sales-chart", dashboard.GetSalesChartData)
	main.Get("/plat-chart", dashboard.GetPlatChartData)
	main.Get("/product-chart", dashboard.GetProductChartData)
	main.Get("/stock-alerts", dashboard.GetStockAlerts)
	main.Get("/stock-rotation", dashboard.GetStockRotationData)
	main.Get("/plat-statistics", dashboard.GetPlatStatistics)
	main.Get("/livraison-statistics", dashboard.GetLivraisonStatistics)
	main.Get("/livraison-zones", dashboard.GetLivraisonZonesData)
	main.Get("/livreur-performance", dashboard.GetLivreurPerformance)
	main.Get("/caisse-statistics", dashboard.GetCaisseStatistics)
	main.Get("/flux-tresorerie", dashboard.GetFluxTresorerieData)
	main.Get("/repartition-transactions", dashboard.GetRepartitionTransactionsData)
	main.Get("/top-transactions", dashboard.GetTopTransactions)
	main.Get("/analyse-categories", dashboard.GetAnalyseCategories)
	main.Get("/previsions-tresorerie", dashboard.GetPrevisionsTresorerie)
	main.Get("/top-caisses", dashboard.GetTopCaisses)

	// ============================================================
	// ENTREPRISE ROUTES
	// ============================================================
	e := api.Group("/entreprises")
	e.Get("/all/paginate", entreprises.GetPaginatedEntreprise)
	e.Get("/all", entreprises.GetAllEntreprises)
	e.Post("/create", entreprises.CreateEntreprise)
	e.Get("/get/:uuid", entreprises.GetEntreprise)
	e.Put("/update/:uuid", entreprises.UpdateEntreprise)
	e.Delete("/delete/:uuid", entreprises.DeleteEntreprise)

	// ============================================================
	// ABONNEMENTS ROUTES
	// ============================================================
	ab := api.Group("/abonnements")
	ab.Get("/all/paginate/:entreprise_uuid", abonnements.GetPaginatedAbonnementsEntreprise)
	ab.Get("/all/paginate", abonnements.GetPaginatedAbonnements)
	ab.Get("/all", abonnements.GetAllAbonnements)
	ab.Get("/current", abonnements.GetAbonnementActuel)
	ab.Get("/expiring", abonnements.GetAbonnementsExpirant)
	ab.Get("/statistics", abonnements.GetStatistiquesAbonnements)
	ab.Get("/verify/:uuid", abonnements.VerifierValiditeAbonnement)
	ab.Get("/get/:uuid", abonnements.GetAbonnement)
	ab.Post("/create", abonnements.CreateAbonnement)
	ab.Put("/update-statut/:uuid", abonnements.UpdateStatutAbonnement)
	ab.Put("/update/:uuid", abonnements.UpdateAbonnement)
	ab.Delete("/delete/:uuid", abonnements.DeleteAbonnement)

	// ============================================================
	// USERS ROUTES
	// ============================================================
	u := api.Group("/users")
	u.Get("/all/paginate/nosearch", users.GetPaginatedNoSerach)
	u.Get("/all/paginate", users.GetPaginatedUsersSupport)
	u.Get("/:entreprise_uuid/:pos_uuid/all/paginate", users.GetPaginatedUserByPosUUID)
	u.Get("/:entreprise_uuid/all/paginate", users.GetPaginatedUsers)
	u.Get("/all/:entreprise_uuid", users.GetAllUsersById)
	u.Get("/all", users.GetAllUsers)
	u.Post("/create", users.CreateUser)
	u.Get("/get/:uuid", users.GetUser)
	u.Put("/update/:uuid", users.UpdateUser)
	u.Delete("/delete/:uuid", users.DeleteUser)

	// ============================================================
	// POS ROUTES
	// ============================================================
	p := api.Group("/pos")
	p.Get("/all/paginate", pos.GetPaginatedPos)
	p.Get("/:entreprise_uuid/all/paginate", pos.GetPaginatedPosByUUID)
	p.Get("/:entreprise_uuid/all", pos.GetAllPosByUUId)
	p.Post("/create", pos.CreatePos)
	p.Get("/get/:uuid", pos.GetPos)
	p.Put("/update/:uuid", pos.UpdatePos)
	p.Delete("/delete/:uuid", pos.DeletePos)

	// ============================================================
	// CAISSES ROUTES
	// ============================================================
	cais := api.Group("/caisses")
	cais.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", caisses.GetDataSynchronisation)
	cais.Get("/:entreprise_uuid/:pos_uuid/all", caisses.GetAllCaisseByPos)
	cais.Get("/:entreprise_uuid/all/total", caisses.GetTotalAllCaisses)
	cais.Get("/:entreprise_uuid/all", caisses.GetAllCaisses)
	cais.Post("/create", caisses.CreateCaisse)
	cais.Get("/get/:uuid", caisses.GetCaisse)
	cais.Put("/update/:uuid", caisses.UpdateCaisse)
	cais.Delete("/delete/:uuid", caisses.DeleteCaisse)

	// ============================================================
	// CAISSE ITEMS ROUTES
	// ============================================================
	caisseItem := api.Group("/caisse-items")
	caisseItem.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", caisses.GetDataSynchronisationCaisseItem)
	caisseItem.Get("/:entreprise_uuid/:caisse_uuid/all/paginate", caisses.GetPaginatedCaisseItems)
	caisseItem.Get("/:entreprise_uuid/:caisse_uuid/all", caisses.GetAllCaisseItems)
	caisseItem.Post("/create", caisses.CreateCaisseItem)
	caisseItem.Get("/get/:uuid", caisses.GetCaisseItem)
	caisseItem.Put("/update/:uuid", caisses.UpdateCaisseItem)
	caisseItem.Delete("/delete/:uuid", caisses.DeleteCaisseItem)

	// ============================================================
	// PRODUCTS ROUTES
	// ============================================================
	pr := api.Group("/products")
	pr.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", products.GetDataSynchronisation)
	pr.Post("/:entreprise_uuid/:pos_uuid/upload-excel", products.UploadProductsFromExcel)
	pr.Get("/:entreprise_uuid/:pos_uuid/all/search", products.GetAllProductBySearch)
	pr.Get("/:entreprise_uuid/:pos_uuid/all/paginate", products.GetPaginatedProductByPosUUID)
	pr.Get("/:entreprise_uuid/:pos_uuid/all", products.GetAllProducts)
	pr.Get("/:entreprise_uuid/all/paginate", products.GetPaginatedProductEntreprise)
	pr.Get("/excel-format-info", products.GetExcelFormatInfo)
	pr.Get("/excel-template", products.GenerateProductExcelTemplate)
	pr.Post("/create", products.CreateProduct)
	pr.Get("/get/:uuid", products.GetProduct)
	pr.Put("/update/stock-endommage/:uuid", products.UpdateProductStockEndommage)
	pr.Put("/update/restitution/:uuid", products.UpdateProductRestitution)
	pr.Put("/update/stock/:uuid", products.UpdateProductStockDispo)
	pr.Put("/update/:uuid", products.UpdateProduct)
	pr.Delete("/delete/:uuid", products.DeleteProduct)

	// ============================================================
	// PLATS ROUTES
	// ============================================================
	pl := api.Group("/plats")
	pl.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", plats.GetDataSynchronisation)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/available", plats.GetAvailablePlats)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/search", plats.GetAllPlatBySearch)
	pl.Get("/:entreprise_uuid/:pos_uuid/all/paginate", plats.GetPaginatedPlatByPosUUID)
	pl.Get("/:entreprise_uuid/:pos_uuid/all", plats.GetAllPlats)
	pl.Get("/:entreprise_uuid/all/paginate", plats.GetPaginatedPlatEntreprise)
	pl.Post("/create", plats.CreatePlat)
	pl.Get("/get/:uuid", plats.GetPlat)
	pl.Put("/update/availability/:uuid", plats.UpdatePlatAvailability)
	pl.Put("/update/:uuid", plats.UpdatePlat)
	pl.Delete("/delete/:uuid", plats.DeletePlat)

	// ============================================================
	// TABLEBOX ROUTES
	// ============================================================
	tb := api.Group("/tablebox")
	tb.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", tablebox.GetDataSynchronisation)
	tb.Get("/:entreprise_uuid/:pos_uuid/all/search", tablebox.GetAllTableBoxBySearch)
	tb.Get("/:entreprise_uuid/:pos_uuid/all/paginate", tablebox.GetPaginatedTableBoxByPosUUID)
	tb.Get("/:entreprise_uuid/:pos_uuid/all", tablebox.GetAllTableBoxs)
	tb.Get("/:entreprise_uuid/:pos_uuid/category/:category", tablebox.GetTableBoxsByCategory)
	tb.Get("/:entreprise_uuid/:pos_uuid/statut/:statut", tablebox.GetTableBoxsByStatut)
	tb.Get("/:entreprise_uuid/all/paginate", tablebox.GetPaginatedTableBoxEntreprise)
	tb.Post("/create", tablebox.CreateTableBox)
	tb.Get("/get/:uuid", tablebox.GetTableBox)
	tb.Put("/update/:uuid", tablebox.UpdateTableBox)
	tb.Delete("/delete/:uuid", tablebox.DeleteTableBox)

	// ============================================================
	// RESERVATIONS ROUTES
	// ============================================================
	r := api.Group("/reservations")
	r.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", reservations.GetDataSynchronisation)
	r.Get("/:entreprise_uuid/:pos_uuid/all/search", reservations.GetAllReservationBySearch)
	r.Get("/:entreprise_uuid/:pos_uuid/all/paginate", reservations.GetPaginatedReservationByPosUUID)
	r.Get("/:entreprise_uuid/:pos_uuid/all", reservations.GetAllReservations)
	r.Get("/:entreprise_uuid/:pos_uuid/status/:status", reservations.GetReservationsByStatus)
	r.Get("/:entreprise_uuid/:pos_uuid/date/:date", reservations.GetReservationsByDate)
	r.Get("/:entreprise_uuid/:pos_uuid/table/:table", reservations.GetReservationsByTable)
	r.Get("/:entreprise_uuid/all/paginate", reservations.GetPaginatedReservationEntreprise)
	r.Post("/create", reservations.CreateReservation)
	r.Get("/get/:uuid", reservations.GetReservation)
	r.Put("/update/:uuid", reservations.UpdateReservation)
	r.Delete("/delete/:uuid", reservations.DeleteReservation)

	// ============================================================
	// STOCKS ROUTES
	// ============================================================
	s := api.Group("/stocks")
	s.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", stocks.GetDataSynchronisationStock)
	s.Get("/all/paginate/:product_uuid", stocks.GetPaginatedStock)
	s.Get("/all/total/:product_uuid", stocks.GetTotalStock)
	s.Get("/all/get/:product_uuid", stocks.GetStockMargeBeneficiaire)
	s.Get("/all/:product_uuid", stocks.GetAllStocks)
	s.Post("/create", stocks.CreateStock)
	s.Get("/get/:uuid", stocks.GetStock)
	s.Put("/update/:uuid", stocks.UpdateStock)
	s.Delete("/delete/:uuid", stocks.DeleteStock)

	// ============================================================
	// STOCK ENDOMMAGES ROUTES
	// ============================================================
	se := api.Group("/stock-endommages")
	se.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", stocks.GetDataSynchronisationStockEndommage)
	se.Get("/all/paginate/:product_uuid", stocks.GetPaginatedStockEndommage)
	se.Get("/all/total/:product_uuid", stocks.GetTotalStockEndommage)
	se.Get("/all/:product_uuid", stocks.GetAllStockEndommages)
	se.Post("/create", stocks.CreateStockEndommage)
	se.Get("/get/:uuid", stocks.GetStockEndommage)
	se.Put("/update/:uuid", stocks.UpdateStockEndommage)
	se.Delete("/delete/:uuid", stocks.DeleteStockEndommage)

	// ============================================================
	// RESTITUTIONS ROUTES
	// ============================================================
	re := api.Group("/restitutions")
	re.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", stocks.GetDataSynchronisationRestitution)
	re.Get("/all/paginate/:product_uuid", stocks.GetPaginatedRestitution)
	re.Get("/all/total/:product_uuid", stocks.GetTotalRestitution)
	re.Get("/all/:product_uuid", stocks.GetAllRestitutions)
	re.Post("/create", stocks.CreateRestitution)
	re.Get("/get/:uuid", stocks.GetRestitution)
	re.Put("/update/:uuid", stocks.UpdateRestitution)
	re.Delete("/delete/:uuid", stocks.DeleteRestitution)

	// ============================================================
	// CLIENTS ROUTES
	// ============================================================
	cl := api.Group("/clients")
	cl.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", clients.GetDataSynchronisation)
	cl.Get("/:entreprise_uuid/:pos_uuid/all/paginate", clients.GetPaginatedClient)
	cl.Get("/:entreprise_uuid/all", clients.GetAllClients)
	cl.Post("/uploads", clients.UploadCsvDataClient)
	cl.Post("/create", clients.CreateClient)
	cl.Get("/get/:uuid", clients.GetClient)
	cl.Put("/update/:uuid", clients.UpdateClient)
	cl.Delete("/delete/:uuid", clients.DeleteClient)

	// ============================================================
	// FOURNISSEURS ROUTES
	// ============================================================
	fs := api.Group("/fournisseurs")
	fs.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", fournisseurs.GetDataSynchronisation)
	fs.Get("/:entreprise_uuid/:pos_uuid/all/paginate", fournisseurs.GetPaginatedFournisseur)
	fs.Get("/:entreprise_uuid/all", fournisseurs.GetAllFournisseurs)
	fs.Post("/create", fournisseurs.CreateFournisseur)
	fs.Get("/get/:uuid", fournisseurs.GetFournisseur)
	fs.Put("/update/:uuid", fournisseurs.UpdateFournisseur)
	fs.Delete("/delete/:uuid", fournisseurs.DeleteFournisseur)

	// ============================================================
	// ZONES ROUTES
	// ============================================================
	z := api.Group("/zones")
	z.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", zones.GetDataSynchronisation)
	z.Get("/:entreprise_uuid/:pos_uuid/all/paginate", zones.GetPaginatedZone)
	z.Get("/:entreprise_uuid/all", zones.GetAllZones)
	z.Post("/uploads", zones.UploadCsvDataZone)
	z.Post("/create", zones.CreateZone)
	z.Get("/get/:uuid", zones.GetZone)
	z.Put("/update/:uuid", zones.UpdateZone)
	z.Delete("/delete/:uuid", zones.DeleteZone)

	// ============================================================
	// LIVREURS ROUTES
	// ============================================================
	lv := api.Group("/livreurs")
	lv.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", livreurs.GetDataSynchronisation)
	lv.Get("/:entreprise_uuid/:pos_uuid/all/paginate", livreurs.GetPaginatedLivreur)
	lv.Get("/:entreprise_uuid/:pos_uuid/type/:type", livreurs.GetLivreursByType)
	lv.Get("/:entreprise_uuid/all", livreurs.GetAllLivreurs)
	lv.Post("/uploads", livreurs.UploadCsvDataLivreur)
	lv.Post("/create", livreurs.CreateLivreur)
	lv.Get("/get/:uuid", livreurs.GetLivreur)
	lv.Put("/update/:uuid", livreurs.UpdateLivreur)
	lv.Delete("/delete/:uuid", livreurs.DeleteLivreur)

	// ============================================================
	// LIVRAISONS ROUTES
	// ============================================================
	liv := api.Group("/livraisons")
	liv.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", livraisons.GetDataSynchronisation)
	liv.Get("/:entreprise_uuid/:pos_uuid/all/paginate", livraisons.GetPaginatedLivraison)
	liv.Get("/:entreprise_uuid/all", livraisons.GetAllLivraisons)
	liv.Post("/uploads", livraisons.UploadCsvDataLivraison)
	liv.Post("/create", livraisons.CreateLivraison)
	liv.Get("/get/:uuid", livraisons.GetLivraison)
	liv.Put("/update/:uuid", livraisons.UpdateLivraison)
	liv.Delete("/delete/:uuid", livraisons.DeleteLivraison)

	// ============================================================
	// COMMANDES ROUTES
	// ============================================================
	cmd := api.Group("/commandes")
	cmd.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", commandes.GetDataSynchronisation)
	cmd.Get("/:entreprise_uuid/:pos_uuid/all/paginate", commandes.GetPaginatedCommandePOS)
	cmd.Get("/:entreprise_uuid/:pos_uuid/all", commandes.GetAllCommandes)
	cmd.Get("/:entreprise_uuid/all/paginate", commandes.GetPaginatedCommandeEntreprise)
	cmd.Post("/create", commandes.CreateCommande)
	cmd.Get("/get/:uuid", commandes.GetCommande)
	cmd.Put("/update/:uuid", commandes.UpdateCommande)
	cmd.Delete("/delete/:uuid", commandes.DeleteCommande)

	// ============================================================
	// COMMANDE LINES ROUTES
	// ============================================================
	cmdl := api.Group("/commande-lines")
	cmdl.Get("/:entreprise_uuid/:pos_uuid/all/synchronisation", commandes.GetDataSynchronisationCommandeLine)
	cmdl.Get("/all/paginate/:commande_uuid", commandes.GetPaginatedCommandeLineByID)
	cmdl.Get("/all/total/:product_uuid", commandes.GetTotalCommandeLine)
	cmdl.Get("/all/:commande_uuid", commandes.GetAllCommandeLineByUUId)
	cmdl.Get("/all", commandes.GetAllCommandeLines)
	cmdl.Post("/create", commandes.CreateCommandeLine)
	cmdl.Get("/get/:uuid", commandes.GetCommandeLine)
	cmdl.Put("/update/:uuid", commandes.UpdateCommandeLine)
	cmdl.Delete("/delete/:uuid", commandes.DeleteCommandeLine)

}
