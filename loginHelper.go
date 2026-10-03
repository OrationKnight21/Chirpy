package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/OrationKnight21/server/internal/auth"
	"github.com/OrationKnight21/server/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	tokenString := auth.MakeRefreshToken()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	duration := time.Hour * 24 * 60
	rToken, err := cfg.db.CreateToken(r.Context(), database.CreateTokenParams{Token: tokenString,
		CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, ExpiresAt: time.Now().Add(duration), RevokedAt: sql.NullTime{}})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token", err)
		return
	}
	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access JWT", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: rToken.Token})
}
