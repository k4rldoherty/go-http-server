package main

// go standard library
import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/k4rldoherty/go-http-server/internal/database"
	"github.com/k4rldoherty/go-http-server/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	const filePathRoot = "."
	const port = "8080"

	godotenv.Load("../../.env")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")

	cfg := server.APIConfig{
		FileServerHits: atomic.Int32{},
		DBQueries:      dbQueries,
		Platform:       platform,
	}

	fileServerHandler := http.FileServer((http.Dir(filePathRoot)))
	// a serveMux directs traffic to the relevant handler function - like a controller in .NET
	serveMux := http.NewServeMux()
	// sets /app as the root directory, even when you dont actually have an app directory on your server
	serveMux.Handle(
		"/app/",
		http.StripPrefix("/app", cfg.MiddlewareMetricsInc(fileServerHandler)),
	)

	// util endpoints to check health etc
	serveMux.HandleFunc("GET /api/healthz", server.HealthzHandler)
	serveMux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	serveMux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	// server obejct to listen on port 8080
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %v and listening on port %v", filePathRoot, port)
	// starts the server
	log.Fatal(server.ListenAndServe())
}
