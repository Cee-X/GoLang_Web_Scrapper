package main

import (
	"fmt"
	"net/http"

	"github.com/Cee-X/rssagg/internal/auth"
	"github.com/Cee-X/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func(apCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Count't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}