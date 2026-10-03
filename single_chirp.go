package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerSingleChirp(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	ChirpId, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirpId", err)
		return
	}
	allChirp, err := cfg.db.GetSingleChirp(r.Context(), ChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not fetch chirps", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        allChirp.ID,
		CreatedAt: allChirp.CreatedAt,
		UpdatedAt: allChirp.UpdatedAt,
		UserId:    allChirp.UserID,
		Body:      allChirp.Body,
	})
}
