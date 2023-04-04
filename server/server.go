package server

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	*http.ServeMux
}

func setCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowCredentials: true,
	})
}

func listenServe(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

func (s *Server) ConnectServer(ctx context.Context, path string, h http.Handler, port string) error {
	c := setCORS()
	s.ServeMux = http.NewServeMux()
	s.Handle(path, c.Handler(h))

	// Start the server
	srv := &http.Server{
		Addr:    port,
		Handler: h2c.NewHandler(s, &http2.Server{}),
	}

	go listenServe(srv)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	select {
	case <-ctx.Done():
		fmt.Println("shutting down server...")
		// Give the server 5 seconds to gracefully shutdown
		ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			return fmt.Errorf("failed to gracefully shutdown server: %v", err)
		}

	case sig := <-quit:
		fmt.Printf("received signal %s", sig)
		// Give the server 5 seconds to gracefully shutdown
		ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			return fmt.Errorf("failed to gracefully shutdown server: %v", err)
		}
	}

	return nil
}
