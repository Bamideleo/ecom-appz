package handlers

import (
	"context"
	"ecom-appz/internal/auth"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	UserRepo *repositories.UserRepository
}

type AuthRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request){
	var req AuthRequest
	_ =json.NewDecoder(r.Body).Decode(&req)

	user := &models.User{
		Email: req.Email,
		Role: "user",
		IsActive: true,
	}

	if err := user.HashPassword(req.Password); err !=nil{
		RespondError(w, http.StatusInternalServerError, "password error")
		return
	}

	if err := h.UserRepo.Create(context.Background(), user); err !=nil{
		RespondError(w, http.StatusBadRequest, "user already exist")
		return	
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var req AuthRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	
	user, err := h.UserRepo.GetByEmail(context.Background(), req.Email)

	if err != nil || !user.CheckPassword(req.Password){
		RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, _ := auth.GenerateToken(user.ID, user.Role)

	json.NewEncoder(w).Encode(map[string]string{
		"token":token,
	})
}