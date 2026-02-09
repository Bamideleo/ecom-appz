package handlers

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Message string `json:"message"`
}

func RespondError(w http.ResponseWriter, status int, msg string){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{Message: msg})
}