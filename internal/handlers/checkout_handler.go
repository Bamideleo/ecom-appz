package handlers

import (
	"ecom-appz/internal/helper"
	"ecom-appz/internal/services"
	"ecom-appz/internal/utils"
	"net/http"
)

type CheckoutHandler struct {
	Service *services.CheckoutService
}

func NewCheckoutHandler(service *services.CheckoutService) *CheckoutHandler{
	return &CheckoutHandler{Service: service}
}


func (h *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request){
	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	order, err := h.Service.Checkout(userID)
	if err != nil{
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, order, http.StatusCreated)
}