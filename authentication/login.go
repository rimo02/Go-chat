package authentication

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Username == "rimo" && req.Password == "123" {
		json.NewEncoder(w).Encode(LoginResponse{
			Success: true,
			Message: "You have sucessfully logged in",
		})
	} else {
		json.NewEncoder(w).Encode(LoginResponse{Success: false, Message: "Invalid username"})
	}
}
