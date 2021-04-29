package web_server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"life-unlimited/podcastination-server/stores"
	"log"
	"net/http"
	"time"
)

type Config struct {
	StaticDir string
	Addr      string
}

type WebServer struct {
	config  Config
	stores  *stores.Stores
	running bool
	stop    chan struct{}
}

func NewServer(config Config, stores *stores.Stores) *WebServer {
	return &WebServer{
		config:  config,
		stores:  stores,
		running: false,
		stop:    make(chan struct{}),
	}
}

func (s *WebServer) Start() error {
	if s.running {
		return fmt.Errorf("web server already running")
	}

	go s.run()

	s.running = true
	return nil
}

func (s *WebServer) run() {
	r := mux.NewRouter()
	// Enable CORS.
	r.Use(middleware)

	// Static file handling.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.config.StaticDir))))
	// Not found handler with cors.
	r.NotFoundHandler = middleware(http.NotFoundHandler())

	s.populateRESTRoutes(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.config.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start web_server.
	go func() { log.Fatal(srv.ListenAndServe()) }()

	// Wait for stop command.
	_ = <-s.stop
	log.Println("shutting down web web_server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("could not shutdown web server: %v", err)
	}
}

// middleware activates cross site stuff and avoids caching.
//
// Cors stuff taken from https://asanchez.dev/blog/cors-golang-options/.
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers.
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		// Avoid caching.
		w.Header().Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next.
		next.ServeHTTP(w, r)
		return
	})
}

func (s *WebServer) Stop() error {
	if !s.running {
		return fmt.Errorf("web server already stopped")
	}
	s.running = false
	s.stop <- struct{}{}
	return nil
}
