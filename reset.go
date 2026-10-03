package main

import "net/http"

func (cfg *apiConfig) resetToZero(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "you are not permitted to do this", nil)
		return
	}
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not complete the action", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users deleted"))
}
