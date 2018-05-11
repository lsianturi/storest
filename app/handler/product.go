package handler

import (
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lsianturi/storest/app/model"
)

func GetAllProducts(db *gorm.DB, w http.ResponseWriter, r *http.Request)  {
	products := []model.Product{}
	db.Find(&products)
	respondJSON(w, http.StatusOK, products)
}

func GetProduct(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	product := getProductOr404(db, id, w, r)
	if product == nil {
		return
	}
	respondJSON(w, http.StatusOK, product)
}


func getProductOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Product {
	product := model.Product{}
	if err := db.First(&product, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &product
}
