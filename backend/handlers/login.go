package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"backend/db"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("❌ Invalid request body:", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	var user User
	var hashedPassword string
	var isVerified bool
	err := db.DB.QueryRow(
		"SELECT email, password, is_verified FROM users WHERE phone=$1",
		req.Phone,
	).Scan(&user.Email, &hashedPassword, &isVerified)

	if err != nil {
		log.Println("❌ User not found:", err)
		http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
		return
	}

	if !isVerified {
		log.Println("❌ User not verified")
		http.Error(w, `{"error": "Please verify your account first"}`, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		log.Println("❌ Invalid password")
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	response := LoginResponse{
		Message: "Login successful",
		Email:   user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}