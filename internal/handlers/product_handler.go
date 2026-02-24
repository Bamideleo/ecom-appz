package handlers

import (
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	Repo repositories.ProductRepository
}


func NewProductHandler(repo repositories.ProductRepository) *ProductHandler{
	return  &ProductHandler{Repo: repo}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request){
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err !=nil{
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Create(&product); err !=nil{
		utils.JSONError(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, product, http.StatusCreated)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request){
	products, err := h.Repo.FindAll()

	if err != nil{
		utils.JSONError(w, "Could not fetch products", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, products, http.StatusOK)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request){
	id, _ :=strconv.Atoi(chi.URLParam(r, "id"))

	product, err := h.Repo.FindAll(id)
	if err != nil{
		utils.JSONError(w, "Product not found", http.StatusNotFound)
	}

	utils.JSONResponse(w, product, http.StatusOK)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request){
	id, _ :=strconv.Atoi(chi. URLParam(r, "id"))

	var product models.Product
	if err :=json.NewDecoder(r.Body).Decode(&product); err !=nil{
		utils.JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	product.ID = id
	if err := h.Repo.Update(&product); err != nil{
		utils.JSONError(w, "Could not update product", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{
		"message": "Product update",
	
	}, http.StatusOK)
}


func (h *ProductHandler)Delete(w http.ResponseWriter, r *http.Request){
	id, _ :=strconv.Atoi(chi. URLParam(r, "id"))
	if err := h.Repo.Delete(id); err !=nil{
		utils.JSONError(w, "Could not delete product", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{
		"message":"Product deleted",
	}, http.StatusOK)
}