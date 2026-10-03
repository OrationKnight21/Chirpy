package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OrationKnight21/server/internal/auth"
	"github.com/OrationKnight21/server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

func (a *apiConfig) handlerEmail(w http.ResponseWriter, r *http.Request) {
	type emailReturn struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	emailVal := emailReturn{}
	err := decoder.Decode(&emailVal)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the email", err)
		return
	}
	hashed, err := auth.HashPassword(emailVal.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash the password", err)
		return
	}
	usr, err := a.db.CreateUser(r.Context(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(),
		UpdatedAt: time.Now(), Email: emailVal.Email, HashedPassword: hashed})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          usr.ID,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
		Email:       usr.Email,
		IsChirpyRed: usr.IsChirpyRed,
	})
}
