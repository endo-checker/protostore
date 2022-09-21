package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	daprpb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	gw "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const defPort = "8080"

// Handler is implemented by a service's grpc 'handler' to register with the
// server and grpc-gateway mux
type Server interface {
	// Permissions is used by the authz interceptor to control access to a method.
	// permissions.
	Permissions(method string) []string
	// RegisterGRPC registers a service's handler (server) with the gRPC server.
	RegisterGRPC(srv *grpc.Server, h interface{})
	// RegisterHTTP registers http servers for a service using a grpc-gateway
	// generated function.
	RegisterHTTP(ctx context.Context, mux *gw.ServeMux, endpoint string, opts []grpc.DialOption) error
}

// Server type is used by Server receiver methods
type server struct {
	port    string
	servers []Server
}

// ServerOption is used for startup options passed on initialisation.
type ServerOption func(*server)

// NewServer creates a new Server instance
func NewServer(opts ...ServerOption) *server {
	// check to see if Dapr APP_PORT has been set
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defPort
	}

	// apply options
	s := server{port: port}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

func (s *server) Serve() error {
	addr := ":" + s.port

	grpcSrv := grpc.NewServer()
	defer grpcSrv.Stop()
	reflection.Register(grpcSrv)

	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// map HTTP headers to gRPC context
	hm := gw.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		switch key {
		case "X-Token-C-Tenant", "X-Token-C-User", "Permissions":
			return key, true
		default:
			return gw.DefaultHeaderMatcher(key)
		}
	})

	// JSON marshaling/unmarshaling
	mo := gw.WithMarshalerOption("*", &gw.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: false,
		},
		// UnmarshalOptions: protojson.UnmarshalOptions{
		// 	DiscardUnknown: false,
		// },
	})

	// multiplex gRPC and HTTP on same port
	httpMux := gw.NewServeMux(hm, mo)
	mux := httpGrpcMux(httpMux, grpcSrv)
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// register grpc servers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, srv := range s.servers {
		srv.RegisterGRPC(grpcSrv, srv)
		srv.RegisterHTTP(ctx, httpMux, addr, dopts)

		// register Dapr server
		if cb, ok := srv.(daprpb.AppCallbackServer); ok {
			daprpb.RegisterAppCallbackServer(grpcSrv, cb)
		}
	}

	log.Println("service starting  on", addr)

	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func WithAppPort(port string) ServerOption {
	return func(s *server) {
		s.port = port
	}
}

func WithServer(srv Server) ServerOption {
	return func(s *server) {
		s.servers = append(s.servers, srv)
	}
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func httpGrpcMux(httpHandler http.Handler, grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if allowedOrigin(r.Header.Get("Origin")) {
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
			}
			if r.Method == "OPTIONS" {
				return
			}
			httpHandler.ServeHTTP(w, r)
		}
	})
}
