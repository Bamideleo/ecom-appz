package handlers

import (
	"ecom-appz/internal/helper"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CartHandler struct {
	Repo repositories.CartRepository
}

func NewCartHandler(repo repositories.CartRepository) *CartHandler{
	return &CartHandler{Repo: repo}
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request){
	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	var req struct{
		ProductID int `json:"product_id"`
		Quantity int `json:"quantity"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	cartID, err := h.Repo.GetOrCreateCart(userID)
	if err != nil{
		utils.JSONError(w, "Failed toget cart", 500)
		return
	}

	err = h.Repo.AddItem(cartID, req.ProductID, req.Quantity)
	if err != nil{
		utils.JSONError(w, "Failed to add item", 500)
		return
	}
	utils.JSONResponse(w, "Item added to cart", 200)
}

func (h *CartHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request){
	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	var req struct{
		ProductID int `json:"product_id"`
		Quantity int `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	cartID, _ := h.Repo.GetOrCreateCart(userID)
	err := h.Repo.UpdateQuantity(cartID, req.ProductID, req.Quantity)
	if err != nil{
		utils.JSONError(w, err.Error(), 400)
		return
	}
	utils.JSONResponse(w, "Cart updated", 200)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request){
	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

   productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	
	cartID, _:= h.Repo.GetOrCreateCart(userID)
	err = h.Repo.RemoveItem(cartID, productID)

	if err != nil{
		utils.JSONError(w, "Failed to remove item", 500)
		return
	}
	utils.JSONResponse(w, "Item removed", 200)
}