package main

import (
	"encoding/json"
	"net/http"

	"github.com/OrationKnight21/server/internal/auth"
	"github.com/OrationKnight21/server/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized user", err)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate jwt", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode params", err)
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusNotImplemented, "could not hash", err)
		return
	}
	UpdatedUser, err := cfg.db.ChangeUserCreds(r.Context(), database.ChangeUserCredsParams{Email: params.Email, HashedPassword: hashedPass, ID: userId})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not update creds", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        UpdatedUser.ID,
			CreatedAt: UpdatedUser.CreatedAt,
			UpdatedAt: UpdatedUser.UpdatedAt,
			Email:     UpdatedUser.Email,
		},
	})
}
