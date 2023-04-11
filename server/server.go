package server

import (
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	*http.ServeMux
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Connect-Protocol-Version,Accept,Authorization,Content-Type,X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func listenServe(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

func (s *Server) ConnectServer(path string, h http.Handler, port string) error {
	r := mux.NewRouter()
	r.Use(cors)
	r.Handle(path, h).Methods(http.MethodPost, http.MethodOptions)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	if err := listenServe(srv); err != nil {
		log.Fatal(err)
	}
	listenServe(srv)

	return nil
}
