package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type StaticFileServerConfig struct {
	StaticDir string
	Addr      string
}

type StaticFileServer struct {
	config  StaticFileServerConfig
	running bool
	stop    chan struct{}
}

func NewStaticFileServer(config StaticFileServerConfig) *StaticFileServer {
	return &StaticFileServer{
		config:  config,
		running: false,
		stop:    make(chan struct{}),
	}
}

func (s *StaticFileServer) Start() error {
	if s.running {
		return fmt.Errorf("static file server already running")
	}

	go s.run()

	s.running = true
	return nil
}

func (s *StaticFileServer) run() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.config.StaticDir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         s.config.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *StaticFileServer) Stop() error {
	if !s.running {
		return fmt.Errorf("static file server already stopped")
	}
	s.running = false
	s.stop <- struct{}{}
	return nil
}
