package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/JLee871/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits  atomic.Int32
	databaseQueries *database.Queries
	platform        string
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func main() {
	serverMutex := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMutex,
	}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{}
	apiCfg.databaseQueries = dbQueries
	apiCfg.platform = os.Getenv("PLATFORM")

	sH := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	serverMutex.Handle("/app/", apiCfg.middlewareMetricsInc(sH))

	rH := http.HandlerFunc(readyHandler)
	serverMutex.Handle("GET /api/healthz", rH)

	cH := http.HandlerFunc(apiCfg.countHandler)
	serverMutex.Handle("GET /admin/metrics", cH)

	resetH := http.HandlerFunc(apiCfg.resetHandler)
	serverMutex.Handle("POST /admin/reset", resetH)

	chirpH := http.HandlerFunc(apiCfg.chirpHandler)
	serverMutex.Handle("POST /api/chirps", chirpH)

	userH := http.HandlerFunc(apiCfg.userHandler)
	serverMutex.Handle("POST /api/users", userH)

	log.Fatal(server.ListenAndServe())
}
