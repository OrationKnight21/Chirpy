package main

import (
	"net/http"
	"time"

	"github.com/OrationKnight21/server/internal/auth"
)

func (cfg *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type RefreshToken struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find the refresh token", err)
		return
	}
	uId, err := cfg.db.GetUserRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized user", err)
		return
	}
	userJwt, err := auth.MakeJWT(uId, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access JWT", err)
		return
	}
	respondWithJSON(w, http.StatusOK, RefreshToken{
		Token: userJwt,
	})
}
