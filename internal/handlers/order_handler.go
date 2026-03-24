package handlers

import (
	"ecom-appz/internal/helper"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/services"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type OrderHandler struct {
	Repo repositories.OrderRepository
	Service services.OrderService
}

func NewOrdertHandler(repo repositories.OrderRepository) *OrderHandler{
	return &OrderHandler{Repo: repo}
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	var req struct{
		Status string `json:"status"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	status := models.OrderStatus(req.Status)
	err = h.Service.UpdateOrderStatus(orderID, status)
	if err != nil{
		utils.JSONError(w, err.Error(), 400)
		return
	}
	utils.JSONResponse(w, "Order status update", 200)
}

func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {

	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	orders, err := h.Service.GetUserOrders(userID)
	if err != nil {
		utils.JSONError(w, "Could not fetch orders", 500)
		return
	}

	utils.JSONResponse(w, orders, 200)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {

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
	order, err := h.Service.GetOrder(id)
	if err != nil {
		utils.JSONError(w, "Order not found", 404)
		return
	}

	utils.JSONResponse(w, order, 200)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	orders, err := h.Service.GetAllOrders()
	if err != nil {
		utils.JSONError(w, "Failed to fetch orders", 500)
		return
	}

	utils.JSONResponse(w, orders, 200)
}

