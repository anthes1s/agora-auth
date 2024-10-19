package handler

import (
	"agora/auth/authctx"
	"agora/auth/hash"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

// TODO: Refactor this, please
// TODO: Also, define and return propper error JSONs if request fails
func Register(ctx *authctx.AuthContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		req := RegisterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		matched := ctx.ValidateUsername.MatchString(req.Username)
		if !matched {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		conn, err := ctx.Pool.Acquire(context.Background())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer conn.Release()

		var username string
		query := `SELECT username FROM users WHERE username = $1`
		err = conn.QueryRow(context.Background(), query, req.Username).Scan(&username)
		switch err {
		case pgx.ErrNoRows:
			hashedPassword, err := hash.Generate(req.Password)
			if err != nil {
				log.Println("Error during hashing", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			_, err = conn.Query(context.Background(), "INSERT INTO users (username, hash) VALUES ($1, $2)", req.Username, hashedPassword)
			if err != nil {
				log.Println("Error during creating user", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		case nil:
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res := RegisterResponse{
			Success: true,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&res)
		return
	}
}
