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
	const filepathRoot = "."
	const port = "8080"

	serverMutex := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
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

	sH := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	serverMutex.Handle("/app/", apiCfg.middlewareMetricsInc(sH))

	serverMutex.HandleFunc("GET /api/healthz", readyHandler)

	serverMutex.HandleFunc("GET /admin/metrics", apiCfg.countHandler)

	serverMutex.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	serverMutex.HandleFunc("POST /api/chirps", apiCfg.postchirpHandler)
	serverMutex.HandleFunc("GET /api/chirps", apiCfg.getallchirpsHandler)
	serverMutex.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getsinglechirpHandler)

	serverMutex.HandleFunc("POST /api/users", apiCfg.postuserHandler)

	serverMutex.HandleFunc("POST /api/login", apiCfg.loginHandler)

	log.Fatal(server.ListenAndServe())
}
