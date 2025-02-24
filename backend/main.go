package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"backend/db"
	"backend/handlers"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := mux.NewRouter()
	router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST") 
	router.HandleFunc("/verify-otp", handlers.VerifyOTPHandler).Methods("POST") 
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("🚀 Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}