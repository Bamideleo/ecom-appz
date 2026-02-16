package handlers

import (
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
)

type ProfileHandler struct {
	UserRepo repositories.UserRepository
}



func NewProfileHandler(userRepo repositories.UserRepository) *ProfileHandler{
	return &ProfileHandler{UserRepo: userRepo}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("userID").(string)

	user, err := h.UserRepo.GetByID(r.Context(), userID)

	if err !=nil{
		utils.JSONError(w, "User not found", http.StatusOK)
		return
	}

	utils.JSONResponse(w, user, http.StatusOK)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value("userId").(string)

	var updateData models.User
	if err :=json.NewDecoder(r.Body).Decode(&updateData); err !=nil{
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateData.ID = userId

	if err:= h.UserRepo.UpdateProfile(&updateData); err != nil{
	utils.JSONError(w, "Could not update profile", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{
		"message":"Profile update successfully",

	}, http.StatusOK)


}