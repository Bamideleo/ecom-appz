package handlers

import (
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	Repo repositories.CategoryRepository
}

func NewCategoryHandler(repo repositories.CategoryRepository) *CategoryHandler{
	return &CategoryHandler{Repo: repo}
}


func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request){
	category, err := h.Repo.FindAll()

	if err != nil{
		utils.JSONError(w, "Could not fetch products", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, category, http.StatusOK)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request){
	var category models.Category
	if err:= json.NewDecoder(r.Body).Decode(&category); err != nil{
		utils.JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Create(&category); err !=nil{
		utils.JSONError(w, "Could not create category", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, category, http.StatusCreated)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	if err != nil {
		utils.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category.ID = id

	if err := h.Repo.Update(&category); err != nil {
		utils.JSONError(w, "Could not update category", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{
		"message": "Category updated successfully",
	}, http.StatusOK)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	if err != nil {
		utils.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		utils.JSONError(w, "Could not delete category", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{
		"message": "Category deleted successfully",
	}, http.StatusOK)
}


func (h *CategoryHandler) DetachProduct(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	if err != nil {
		utils.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	

	// productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DetachProduct(categoryID, productID); err != nil {
		utils.JSONError(w, "Could not detach product", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{
		"message": "Product detached from category successfully",
	}, http.StatusOK)
}

func (h *CategoryHandler) AttachProduct(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	categoryID, _ := strconv.Atoi(idStr)
	productID, _ := strconv.Atoi(idStr)

	if err := h.Repo.AttachProduct(categoryID, productID); err !=nil{
		utils.JSONError(w, "Could not attach product", http.StatusInternalServerError)
		return
	}
	utils. JSONResponse(w, map[string]string{
		"message": "Product attached to category",
	}, http.StatusOK)

}
