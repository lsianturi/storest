package handler

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/lsianturi/storest/app/model"
)

func CreateCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	customer := model.Customer{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	hash, _ := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)

	customer.PasswordHash = string(hash)
	if err := db.Save(&customer).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, customer)
}

func UpdateCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	custid, _ := strconv.Atoi(vars["custid"])
	customer := getCustomer(db, custid, w, r)
	if customer == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&customer).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, customer)
}

func GetAllCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request)  {
	customers := []model.Customer{}
	db.Preload("Carts").Preload("Carts.Items").Preload("Carts.Items.Product").Preload("Carts.DeliveryAddr").Preload("Addresses").Find(&customers)
	respondJSON(w, http.StatusOK, customers)
}

func GetCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	custid, _ := strconv.Atoi(vars["custid"])
	Cart := getCustomerOr404(db, custid, w, r)
	if Cart == nil {
		return
	}
	respondJSON(w, http.StatusOK, Cart)
}

func GetCustomerCart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cust, _ := strconv.Atoi(vars["custid"])
	cart, _ := strconv.Atoi(vars["cartid"])
	customer := model.Customer{}
	if err := db.Preload("Addresses").Preload("Carts", "custid = ?", cart).Preload("Carts.Items").Preload("Carts.Items.Product").Preload("Carts.DeliveryAddr").Find(&customer, cust).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, customer)
}

func GetCustomerCartItem(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	custid, _ := strconv.Atoi(vars["custid"])
	cartid, _ := strconv.Atoi(vars["cartid"])
	itemid, _ := strconv.Atoi(vars["itemid"])
	customer := model.Customer{}
	if err := db.Preload("Addresses").Preload("Carts", "custid = ?", cartid).Preload("Carts.Items", "custid = ?", itemid).Preload("Carts.Items.Product").Preload("Carts.DeliveryAddr").Find(&customer, custid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, customer)
}


func getCustomerOr404(db *gorm.DB, custid int, w http.ResponseWriter, r *http.Request) *model.Customer {
	customer := model.Customer{}
	if err := db.Preload("Addresses").Preload("Carts").Preload("Carts.Items").Preload("Carts.Items.Product").Preload("Carts.DeliveryAddr").Find(&customer, custid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		fmt.Println(err)
		return nil
	}
	return &customer
}

func getCustomer(db *gorm.DB, custid int, w http.ResponseWriter, r *http.Request) *model.Customer {
	customer := model.Customer{}
	if err := db.Preload("Addresses").Find(&customer, custid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &customer
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cust := model.Customer{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cust); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	cust2 := model.Customer{}
	db.Where("email = ?", cust.Email).Find(&cust2)
	if cust2.CommonModel.ID == 0 {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cust2.PasswordHash), []byte(cust.Password)); err != nil {
		respondError(w, http.StatusNotFound, "Wrong email or password")
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	var mySigningKey = []byte("secret")
	tokenString, _ := token.SignedString(mySigningKey)
	cust2.Token = tokenString

	respondJSON(w, http.StatusOK, cust2)
}