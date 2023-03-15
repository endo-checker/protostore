package server

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/cors"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func HttpGrpcMux(httpHandler http.Handler, grpcServer *grpc.Server) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if allowedOrigin(r.Header.Get("Origin")) {
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
			}
			if r.Method == "OPTIONS" {
				return
			}

			httpHandler.ServeHTTP(w, r)
		}
	})
}

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
