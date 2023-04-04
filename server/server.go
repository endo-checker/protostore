package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func setCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowCredentials: true,
	})
}

type Server struct {
	*http.ServeMux
}

func (s *Server) ConnectServer(path string, h http.Handler, port string) {
	c := setCORS()
	s.ServeMux = http.NewServeMux()
	s.Handle(path, c.Handler(h))

	// Start the server
	srv := &http.Server{
		Addr:    port,
		Handler: h2c.NewHandler(s, &http2.Server{}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to gracefully shutdown server: %v", err)
	}
}
