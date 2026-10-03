package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/OrationKnight21/server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) numberofRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	counter := cfg.fileserverHits.Load()
	htmlTemplate := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	fmt.Fprintf(w, htmlTemplate, counter)
}
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print("error occured opening the DB")
	}
	dbQueries := database.New(db)
	apiCfg := &apiConfig{db: dbQueries, platform: platform, jwtSecret: jwtSecret, polkaKey: polkaKey}
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.numberofRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetToZero)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerNewChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerEmail)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerSingleChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeToken)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerChirpRed)
	svr := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	svr.ListenAndServe()
}
