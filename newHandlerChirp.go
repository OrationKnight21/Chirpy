package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	o := r.URL.Query().Get("sort")
	chirps := []Chirp{}
	if s != "" {
		authId, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not convert", err)
			return
		}
		authorChirp, err := cfg.db.GetChirpsbyAuthor(r.Context(), authId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not receive records for the id", err)
			return
		}
		for _, val := range authorChirp {
			chirps = append(chirps, Chirp{
				ID:        val.ID,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
				Body:      val.Body,
				UserId:    val.UserID,
			})
		}

	} else {
		allChirp, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not retreive the records", err)
			return
		}

		for _, val := range allChirp {
			chirps = append(chirps, Chirp{
				ID:        val.ID,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
				Body:      val.Body,
				UserId:    val.UserID,
			})
		}
	}
	if o != "" && o == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
