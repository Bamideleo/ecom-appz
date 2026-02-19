package handlers

import (
	"ecom-appz/internal/helper"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
)


// type UpdateProfileRequest struct {
//     FullName *string `json:"full_name"`
//     Email    *string `json:"email"`
//     Phone    *string `json:"phone"`
// }


type ProfileHandler struct {
	UserRepo repositories.UserRepository
}



func NewProfileHandler(userRepo repositories.UserRepository) *ProfileHandler{
	return &ProfileHandler{UserRepo: userRepo}
}




func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	user, err := h.UserRepo.GetByID(r.Context(), userID)
	if err != nil {
		utils.JSONError(w, "User not found", http.StatusNotFound)
		return
	}

	utils.JSONResponse(w, user, http.StatusOK)
}


func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := helper.GetUserClaims(r.Context())
	if !ok {
		utils.JSONError(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.UserRepo.UpdateProfile(userID, &req); err != nil {
		utils.JSONError(w, "Could not update profile", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{"message": "Profile updated successfully"}, http.StatusOK)
}
