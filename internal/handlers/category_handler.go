package handlers

import (
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	Repo repositories.CategoryRepository
}

func NewCategoryHandler(repo repositories.CategoryRepository) *CategoryHandler{
	return &CategoryHandler{Repo: repo}
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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
	categoryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(chi.URLParam(r, "productID"))
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
	categoryID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))

	if err := h.Repo.AttachProduct(categoryID, productID); err !=nil{
		utils.JSONError(w, "Could not attach product", http.StatusInternalServerError)
		return
	}
	utils. JSONResponse(w, map[string]string{
		"message": "Product attached to category",
	}, http.StatusOK)

}
