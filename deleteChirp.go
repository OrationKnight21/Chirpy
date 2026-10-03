package main

import (
	"net/http"

	"github.com/OrationKnight21/server/internal/auth"
	"github.com/OrationKnight21/server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	ChirpId, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid path", err)
		return
	}
	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not found", err)
		return
	}
	userId, err := auth.ValidateJWT(userToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized user", err)
		return
	}
	chirpRet, err := cfg.db.GetSingleChirp(r.Context(), ChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if userId == chirpRet.UserID {
		err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
			ID:     chirpRet.ID,
			UserID: chirpRet.UserID,
		})
		if err != nil {
			respondWithError(w, http.StatusForbidden, "forbidden request made", err)
			return
		}
	} else {
		respondWithError(w, http.StatusForbidden, "forbidden request made", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
