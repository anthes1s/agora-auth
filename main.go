package main

import (
	// "agora/auth/hash"
	"agora/auth/authctx"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func main() {
	ctx, err := authctx.New()
	if err != nil {
		log.Fatalf("Error occured when launching a server: %v", err)
	}

	http.HandleFunc("/api/v1/register", Register(ctx))

	log.Printf("Server is running at localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error occurred when launching a server: %v", err)
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	JWT     string `json:"jwt"`
	Success bool   `json:"success"`
}

func Register(ctx *authctx.AuthContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Extract JSON
		req := RegisterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Verify received data
		// 1. Check if data is not empty
		if req.Username == "" || req.Password == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		// 2. Ensure that Username only consists of alphanumeric symbols
		// 2.1. Compile pattern
		pattern := ``
		re, err := regexp.Compile(pattern)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// 2.2 Verify string to the compiled pattern
		matched := re.MatchString(req.Username)
		if !matched {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}

		// Add user to the database

		res := RegisterResponse{
			JWT:     "your-secret-token",
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&res)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
	return
}
