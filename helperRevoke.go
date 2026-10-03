package main

import (
	"net/http"

	"github.com/OrationKnight21/server/internal/auth"
)

func (cfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token found", err)
		return
	}
	err = cfg.db.UpdateRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke session", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
