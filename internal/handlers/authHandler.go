package handlers

import (
	"context"
	"ecom-appz/internal/auth"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"encoding/json"
	"net/http"
	"time"
)

type AuthHandler struct {
	UserRepo *repositories.UserRepository
	RefreshRepo repositories.RefreshRepository
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

	utils.JSONResponse(w, map[string]string{
		"message": "User registered successfully",
	}, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var req AuthRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	
	user, err := h.UserRepo.GetByEmail(context.Background(), req.Email)

	if err != nil || !user.CheckPassword(req.Password){
		RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	accessToken, _ := auth.GenerateToken(user.ID, user.Role)

	refreshToken, exp, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.JSONError(w, "Could not generate refresh token", http.StatusInternalServerError)
		return
	}

	err = h.RefreshRepo.Store(&models.RefreshToken{
		UserId:    user.ID,
		Token:     refreshToken,
		ExpiresAt: exp,
		CreatedAt: time.Now(),
	})
	if err != nil {
		utils.JSONError(w, "Could not store refresh token", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, http.StatusOK)
}


func(h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request){
	var body struct{
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err !=nil{
		utils.JSONError(w,"Invalid request", http.StatusBadRequest)
		return
	}
	storedToken, err := h.RefreshRepo.Find(body.RefreshToken)

	if err != nil || storedToken.ExpiresAt.Before(time.Now()){
		utils.JSONError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.GetByID(r.Context(),storedToken.UserId)

	if err != nil{
		utils.JSONError(w,"User not found", http.StatusUnauthorized)
		return
	}
	newAccess, _:= auth.GenerateToken(user.ID, user.Role)

	utils.JSONResponse(w, map[string]string{
		"access_token":newAccess,
	}, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request){
	var body struct{
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err !=nil{
		utils.JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	h.RefreshRepo.Delete(body.RefreshToken)

	utils.JSONResponse(w, map[string]string{
		"message":"Logged out successfully",
	}, http.StatusOK)
	
}