package server 

import (
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func setCORS() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowCredentials: true,
	})
	return c
}

// connect handler setup, need to initiate http.ServeMux as well as variables
func (s Server) ConnectServer(path string, h http.Handler, port string) {
	c := setCORS()

	mux := http.NewServeMux()

	mux.Handle(path, h)
	handler := c.Handler(mux)

	http.ListenAndServe(
		port,
		h2c.NewHandler(handler, &http2.Server{}),
	)
}

type Server struct {
	*http.ServeMux
}
