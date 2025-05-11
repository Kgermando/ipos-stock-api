package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/ipos-stock-api/controllers/auth"
	"github.com/kgermando/ipos-stock-api/controllers/caisses"
	"github.com/kgermando/ipos-stock-api/controllers/clients"
	"github.com/kgermando/ipos-stock-api/controllers/commandes"
	"github.com/kgermando/ipos-stock-api/controllers/entreprises"
	"github.com/kgermando/ipos-stock-api/controllers/fournisseurs"
	"github.com/kgermando/ipos-stock-api/controllers/pos"
	"github.com/kgermando/ipos-stock-api/controllers/products"
	"github.com/kgermando/ipos-stock-api/controllers/stocks"
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
	a.Post("/reset/:token", auth.ResetPassword)

	// app.Use(middlewares.IsAuthenticated)

	a.Get("/user", auth.AuthUser)
	a.Put("/profil/info", auth.UpdateInfo)
	a.Put("/change-password", auth.ChangePassword)
	a.Post("/logout", auth.Logout)

	// Users controller
	u := api.Group("/users")
	u.Get("/all", users.GetAllUsers)
	u.Get("/all/paginate", users.GetPaginatedUsers)
	u.Get("/all/paginate/nosearch", users.GetPaginatedNoSerach)
	u.Get("/get/:uuid", users.GetUser)
	u.Post("/create", users.CreateUser)
	u.Put("/update/:uuid", users.UpdateUser)
	u.Delete("/delete/:uuid", users.DeleteUser)

	// Entreprise controller
	e := api.Group("/entreprises")
	e.Get("/all", entreprises.GetAllEntreprises)
	e.Get("/all/paginate", entreprises.GetPaginatedEntreprise)
	e.Get("/get/:uuid", entreprises.GetEntreprise)
	e.Post("/create", entreprises.CreateEntreprise)
	e.Put("/update/:uuid", entreprises.UpdateEntreprise)
	e.Delete("/delete/:uuid", entreprises.DeleteEntreprise)

	// POS controller
	p := api.Group("/pos")
	p.Get("/all", pos.GetAllPoss)
	p.Get("/all/paginate", pos.GetPaginatedPos)
	p.Get("/all/paginate/:entreprise_uuid", pos.GetPaginatedPosByUUID)
	p.Get("/all/:entreprise_uuid", pos.GetAllPosByUUId)
	p.Get("/get/:uuid", pos.GetPos)
	p.Post("/create", pos.CreatePos)
	p.Put("/update/:uuid", pos.UpdatePos)
	p.Delete("/delete/:uuid", pos.DeletePos)

	// Caisses controller
	cais := api.Group("/caisses")
	cais.Get("/:code_entreprise/all/total", caisses.GetTotalAllCaisses)
	cais.Get("/:code_entreprise/all", caisses.GetAllCaisses)
	cais.Get("/:code_entreprise/:pos_uuid/all", caisses.GetAllCaisseByPos)
	cais.Get("/get/:uuid", caisses.GetCaisse)
	cais.Post("/create", caisses.CreateCaisse)
	cais.Put("/update/:uuid", caisses.UpdateCaisse)
	cais.Delete("/delete/:uuid", caisses.DeleteCaisse)

	// Caisse item Controller
	caisseItem := api.Group("/caisse-items")
	caisseItem.Get("/:code_entreprise/:caisse_uuid/all/paginate", caisses.GetPaginatedCaisseItems)
	caisseItem.Get("/:code_entreprise/:caisse_uuid/all", caisses.GetAllCaisseItems)
	caisseItem.Get("/get/:uuid", caisses.GetCaisseItem)
	caisseItem.Post("/create", caisses.CreateCaisseItem)
	caisseItem.Put("/update/:uuid", caisses.UpdateCaisseItem)
	caisseItem.Delete("/delete/:uuid", caisses.DeleteCaisseItem)

	// Product controller
	pr := api.Group("/products")
	pr.Get("/:code_entreprise/all/paginate", products.GetPaginatedProductEntreprise)
	pr.Get("/:code_entreprise/:pos_uuid/all", products.GetAllProducts)
	pr.Get("/:code_entreprise/:pos_uuid/all/paginate", products.GetPaginatedProductByPosUUID)
	pr.Get("/:code_entreprise/:pos_uuid/all/search", products.GetAllProductBySearch)
	pr.Get("/get/:uuid", products.GetProduct)
	pr.Post("/create", products.CreateProduct)
	pr.Put("/update/:uuid", products.UpdateProduct)
	pr.Put("/update/stock/:uuid", products.UpdateProductStockDispo)
	pr.Put("/update/stock-endommage/:uuid", products.UpdateProductStockEndommage)
	pr.Put("/update/restitution/:uuid", products.UpdateProductRestitution)
	pr.Delete("/delete/:uuid", products.DeleteProduct)

	// Stock controller
	s := api.Group("/stocks")
	s.Get("/all/paginate/:product_uuid", stocks.GetPaginatedStock)
	s.Get("/all/total/:product_uuid", stocks.GetTotalStock)
	s.Get("/all/get/:product_uuid", stocks.GetStockMargeBeneficiaire)
	s.Get("/all/:product_uuid", stocks.GetAllStocks)
	s.Get("/get/:uuid", stocks.GetStock)
	s.Post("/create", stocks.CreateStock)
	s.Put("/update/:uuid", stocks.UpdateStock)
	s.Delete("/delete/:uuid", stocks.DeleteStock)

	// StockEndommage controller
	se := api.Group("/stock-endommages")
	se.Get("/all/paginate/:product_uuid", stocks.GetPaginatedStockEndommage)
	se.Get("/all/total/:product_uuid", stocks.GetTotalStockEndommage)
	se.Get("/all/:product_uuid", stocks.GetAllStockEndommages)
	se.Get("/get/:uuid", stocks.GetStockEndommage)
	se.Post("/create", stocks.CreateStockEndommage)
	se.Put("/update/:uuid", stocks.UpdateStockEndommage)
	se.Delete("/delete/:uuid", stocks.DeleteStockEndommage)

	// Restitution controller
	re := api.Group("/restitutions")
	re.Get("/all/paginate/:product_uuid", stocks.GetPaginatedRestitution)
	re.Get("/all/total/:product_uuid", stocks.GetTotalRestitution)
	re.Get("/all/:product_uuid", stocks.GetAllRestitutions)
	re.Get("/get/:uuid", stocks.GetRestitution)
	re.Post("/create", stocks.CreateRestitution)
	re.Put("/update/:uuid", stocks.UpdateRestitution)
	re.Delete("/delete/:uuid", stocks.DeleteRestitution)

	// Client controller
	cl := api.Group("/clients")
	cl.Get("/:code_entreprise/all", clients.GetAllClients)
	cl.Get("/:code_entreprise/all/paginate", clients.GetPaginatedClient)
	cl.Get("/get/:uuid", clients.GetClient)
	cl.Post("/create", clients.CreateClient)
	cl.Post("/uploads", clients.UploadCsvDataClient)
	cl.Put("/update/:uuid", clients.UpdateClient)
	cl.Delete("/delete/:uuid", clients.DeleteClient)

	// Fournisseur controller
	fs := api.Group("/fournisseurs")
	fs.Get("/:code_entreprise/all", fournisseurs.GetAllFournisseurs)
	fs.Get("/:code_entreprise/all/paginate", fournisseurs.GetPaginatedFournisseur)
	fs.Get("/get/:uuid", fournisseurs.GetFournisseur)
	fs.Post("/create", fournisseurs.CreateFournisseur)
	fs.Put("/update/:uuid", fournisseurs.UpdateFournisseur)
	fs.Delete("/delete/:uuid", fournisseurs.DeleteFournisseur)

	// Commande controller
	cmd := api.Group("/commandes")
	cmd.Get("/:code_entreprise/all/paginate", commandes.GetPaginatedCommandeEntreprise)
	cmd.Get("/:code_entreprise/:pos_uuid/all", commandes.GetAllCommandes)
	cmd.Get("/get/:uuid", commandes.GetCommande)
	cmd.Post("/create", commandes.CreateCommande)
	cmd.Put("/update/:uuid", commandes.UpdateCommande)
	cmd.Delete("/delete/:uuid", commandes.DeleteCommande)

	// Commande line controller
	cmdl := api.Group("/commandes-lines")
	cmdl.Get("/all", commandes.GetAllCommandeLines)
	cmdl.Get("/all/paginate/:commande_uuid", commandes.GetPaginatedCommandeLineByID)
	cmdl.Get("/all/total/:product_uuid", commandes.GetTotalCommandeLine)
	cmdl.Get("/all/:commande_uuid", commandes.GetAllCommandeLineByUUId)
	cmdl.Get("/get/:uuid", commandes.GetCommandeLine)
	cmdl.Post("/create", commandes.CreateCommandeLine)
	cmdl.Put("/update/:uuid", commandes.UpdateCommandeLine)
	cmdl.Delete("/delete/:uuid", commandes.DeleteCommandeLine)

}
