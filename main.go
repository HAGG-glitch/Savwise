package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"savwise-ai/internal/config"
	appdb "savwise-ai/internal/db"
	"savwise-ai/internal/handlers"
	"savwise-ai/internal/middleware"
	"savwise-ai/internal/services"
)

func main() {
	cfg := config.Load()

	database, err := appdb.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()

	if err := appdb.RunMigrations(database, "migrations"); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	h := &handlers.Handler{
		DB:   database,
		Cfg:  cfg,
		Groq: services.NewGroqService(cfg.GroqAPIKey, cfg.GroqModel),
	}

	mux := http.NewServeMux()
	registerPages(mux)
	h.RegisterRoutes(mux)

	addr := ":" + cfg.AppPort
	log.Printf("SavWise AI server running at http://localhost%s", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      middleware.Chain(limitBody(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
		next.ServeHTTP(w, r)
	})
}

func registerPages(mux *http.ServeMux) {
	webDir := "web"
	fileServer := http.FileServer(http.Dir(webDir))
	mux.Handle("/assets/", http.StripPrefix("/", fileServer))
	mux.Handle("/css/", http.StripPrefix("/", fileServer))
	mux.Handle("/js/", http.StripPrefix("/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/")
		if !strings.HasSuffix(path, ".html") {
			path += ".html"
		}
		full := filepath.Join(webDir, path)
		if _, err := os.Stat(full); err != nil {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, full)
	})
}
