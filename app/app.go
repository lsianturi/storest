package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	"github.com/lsianturi/storest/app/handler"
	"github.com/lsianturi/storest/app/model"
	"github.com/lsianturi/storest/config"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Assets http.FileSystem
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@tcp(10.15.2.115:3306)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	
	a.Assets = http.Dir(config.AssetDir)
	a.Router = mux.NewRouter().StrictSlash(true)
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	a.Router.PathPrefix("/api/v1/images/").Handler(http.StripPrefix("/api/v1/images/", http.FileServer(a.Assets)))

	sub := a.Router.PathPrefix("/api/v1").Subrouter()
	// Routing for handling the projects
	sub.Methods("POST").Path("/login").HandlerFunc(a.Login)
	sub.Methods("GET").Path("/products").HandlerFunc(a.GetAllProducts)
	sub.Methods("GET").Path("/products/{id:[0-9]+}").HandlerFunc(a.GetProduct)
	
	sub.Methods("GET").Path("/customers").HandlerFunc(a.GetAllCustomer)
	sub.Methods("POST").Path("/customers").HandlerFunc(a.CreateCustomer)
	sub.Methods("PUT").Path("/customers/{custid:[0-9]+}").HandlerFunc(a.UpdateCustomer)
	sub.Methods("GET").Path("/customers/{custid:[0-9]+}").HandlerFunc(a.GetCustomer)
	sub.Methods("GET").Path("/customers/{custid:[0-9]+}/{cartid:[0-9]+}").HandlerFunc(a.GetCustomerCart)
	sub.Methods("GET").Path("/customers/{custid:[0-9]+}/{cartid:[0-9]+}/{itemid:[0-9]+}").HandlerFunc(a.GetCustomerCartItem)
	
	sub.Methods("POST").Path("/carts").HandlerFunc(a.SaveCart)
	sub.Methods("PUT").Path("/carts/{cartid:[0-9]+}").HandlerFunc(a.UpdateCart)
	sub.Methods("DELETE").Path("/carts/{cartid:[0-9]+}").HandlerFunc(a.DeleteCart)
	sub.Methods("DELETE").Path("/carts/{cartid:[0-9]+}/{itemid:[0-9]+}").HandlerFunc(a.DeleteCartItem)
}

/*
** Projects Handlers
 */
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	handler.Login(a.DB, w, r)
}

func (a *App) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	handler.GetAllProducts(a.DB, w, r)
}
func (a *App) GetProduct(w http.ResponseWriter, r *http.Request) {
	handler.GetProduct(a.DB, w, r)
}

/*
** Customer Handlers
*/

func (a *App) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	handler.CreateCustomer(a.DB, w, r)
}

func (a *App) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCustomer(a.DB, w, r)
}

func (a *App) DeleteCart(w http.ResponseWriter, r *http.Request) {
	handler.DeleteCart(a.DB, w, r)
}

func (a *App) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	handler.DeleteCartItem(a.DB, w, r)
}

func (a *App) GetAllCustomer(w http.ResponseWriter, r *http.Request)  {
	handler.GetAllCustomer(a.DB, w, r)
}
func (a *App) GetCustomer(w http.ResponseWriter, r *http.Request) {
	handler.GetCustomer(a.DB, w, r)
}

func (a *App) GetCustomerCart(w http.ResponseWriter, r *http.Request) {
	handler.GetCustomerCart(a.DB, w, r)
}

func (a *App) GetCustomerCartItem(w http.ResponseWriter, r *http.Request) {
	handler.GetCustomerCartItem(a.DB, w, r)
}
/*
** Cart Handlers
*/

func (a *App) GetAllCarts(w http.ResponseWriter, r *http.Request)  {
	handler.GetAllCarts(a.DB, w, r)
}
func (a *App) GetCart(w http.ResponseWriter, r *http.Request) {
	handler.GetCart(a.DB, w, r)
}

func (a *App) SaveCart(w http.ResponseWriter, r *http.Request) {
	handler.SaveCart(a.DB, w, r)
}

func (a *App) UpdateCart(w http.ResponseWriter, r *http.Request)  {
	handler.UpdateCart(a.DB, w, r)
}
// Run the app on it's router
func (a *App) Run(host string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Request-Method", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Fatal(http.ListenAndServe(host, handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
}
