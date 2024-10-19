package handler

import (
	"agora/auth/authctx"
	"net/http"
)

func Login(ctx *authctx.AuthContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
		return
	}
}
