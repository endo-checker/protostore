package main

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

// connect handler setup
func ConnectServer(path string, h http.Handler, port string) {
	c := setCORS()

	mux := http.NewServeMux()

	mux.Handle(path, h)
	handler := c.Handler(mux)

	http.ListenAndServe(
		port,
		h2c.NewHandler(handler, &http2.Server{}),
	)
}