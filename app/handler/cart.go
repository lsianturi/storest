package handler

import (
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lsianturi/storest/app/model"
)

func GetAllCarts(db *gorm.DB, w http.ResponseWriter, r *http.Request)  {
	carts := []model.Cart{}
	db.Preload("Items").Preload("Items.Product").Preload("DeliveryAddr").Find(&carts)
	respondJSON(w, http.StatusOK, carts)
}

func GetCart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	Cart := getCartOr404(db, id, w, r)
	if Cart == nil {
		return
	}
	respondJSON(w, http.StatusOK, Cart)
}

func getCartOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Cart {
	cart := model.Cart{}
	if err := db.Preload("Items").Preload("Items.Product").Preload("DeliveryAddr").Find(&cart, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &cart
}

func SaveCart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cart := model.Cart{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cart); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	
	tx := db.Begin()
	if err := tx.Save(&cart).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, item := range cart.Items {
		tx.Table("products").Where("id = ?", item.ProductID).UpdateColumn("inventory", gorm.Expr("inventory - ?", item.Quantity))
	}

	tx.Commit()
	respondJSON(w, http.StatusCreated, cart)
}

func UpdateCart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["cartid"])
	cart := getCart(db, id, w, r)
	if cart == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cart); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&cart).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, cart)
}

func DeleteCart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cartid, _ := strconv.Atoi(vars["cartid"])
	cart := model.Cart{}
	if err := db.Preload("Items").Find(&cart, cartid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := db.Delete(&cart).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func DeleteCartItem(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cartid, _ := strconv.Atoi(vars["cartid"])
	itemid, _ := strconv.Atoi(vars["itemid"])
	cart := model.Cart{}
	if err := db.Preload("Items", "id = ?", itemid).Find(&cart, cartid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	tx := db.Begin()
	for _, item := range cart.Items {
		if err := tx.Delete(&item).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tx.Table("products").Where("id = ?", item.ProductID).UpdateColumn("inventory", gorm.Expr("inventory + ?", item.Quantity))
	}
	tx.Commit()

	respondJSON(w, http.StatusNoContent, nil)
}

func getCart(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Cart {
	cart := model.Cart{}
	if err := db.Find(&cart, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &cart
}