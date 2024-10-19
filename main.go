package main

import (
	"agora/auth/authctx"
	"agora/auth/handler"
	"log"
	"net/http"
)

func main() {
	ctx, err := authctx.New()
	if err != nil {
		log.Fatalf("Error occured when launching a server: %v", err)
	}

	http.HandleFunc("/api/v1/register", handler.Register(ctx))
	http.HandleFunc("/api/v1/login", handler.Login(ctx))

	log.Printf("Server is running at localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error occurred when launching a server: %v", err)
	}
}
