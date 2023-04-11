package server

import (
	"fmt"

	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	*http.ServeMux
}

// newCorsHandler creates a new CORS handler
func newCorsHandler() (*cors.Cors, error) {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	}), nil
}

func listenServe(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

func (s *Server) ConnectServer(path string, h http.Handler, port string) error {
	c, err := newCorsHandler()
	if err != nil {
		return fmt.Errorf("failed to set CORS: %v", err)
	}
	s.ServeMux = http.NewServeMux()
	s.Handle(path, c.Handler(h))

	// Start the server
	srv := &http.Server{
		Addr:    port,
		Handler: h2c.NewHandler(s, &http2.Server{}),
	}

	go listenServe(srv)

	return nil
}
