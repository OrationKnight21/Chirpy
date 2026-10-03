package main

import (
	"encoding/json"
	"net/http"

	"github.com/OrationKnight21/server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized user", err)
		return
	}
	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode", err)
		return
	}
	if params.Event == "user.upgraded" {
		_, err := cfg.db.ChangeChirpRedStatus(r.Context(), params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "could not change status", err)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
