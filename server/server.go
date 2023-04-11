package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Server represents a HTTP server with CORS support
type Server struct {
	*http.ServeMux
}

// newCorsHandler creates a new CORS handler
func newCorsHandler() (*cors.Cors, error) {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	}), nil
}

// listenAndServe starts the HTTP server and listens for incoming requests
func (s *Server) listenAndServe(port string) error {
	log.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, h2c.NewHandler(s, &http2.Server{}))
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server on port %s: %v", port, err)
	}
	return nil
}

// ConnectServer connects a HTTP handler to a server and starts listening for requests
func (s *Server) ConnectServer(path string, h http.Handler, port string) error {
	c, err := newCorsHandler()
	if err != nil {
		return err
	}

	s.ServeMux = http.NewServeMux()
	s.Handle(path, c.Handler(h))

	go s.listenAndServe(port)

	return nil
}
