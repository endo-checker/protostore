package auth

import (
	"context"
	_ "embed"
	"net/http"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Authz(p string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allow := authorise(r.Context(), p); !allow {
				http.Error(w, "insufficient permission", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authorise(ctx context.Context, permission string) bool {
	_, claims, _ := FromContext(ctx)
	ps, ok := claims["permissions"].([]interface{})
	if !ok {
		return false
	}

	allow := false
	for _, p := range ps {
		if p.(string) == permission {
			allow = true
			break
		}
	}

	return allow
}

//go:embed authz.rego
var policy []byte

// UnaryAuthZInterceptor uses OPA to check user's permissions against the permissions
// required by a gRPC method.
func AuthzInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// return handler(ctx, req)

		// get permissions claim
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return resp, status.Errorf(codes.Aborted, "%s", "no incoming context")
		}
		var userPerms []string
		if perms := md.Get("permissions"); len(perms) > 0 {
			userPerms = strings.Split(perms[0], ",")
		}

		// get permissions from handler/server
		var perms []string
		p, ok := info.Server.(interface{ Permissions(string) []string })
		if ok {
			perms = p.Permissions(info.FullMethod)
			// return resp, status.Errorf(codes.Aborted, "%s", "could not get permissions")
		}

		inputs := map[string]interface{}{
			"method":      info.FullMethod,
			"methodPerms": perms, //p.Permissions(info.FullMethod),
			"userPerms":   userPerms,
		}

		rego := rego.New(
			rego.Query("data.authz.allow"),
			rego.Module("policy.rego", string(policy)),
			rego.Input(inputs),
		)

		// evaluate policy
		res, err := rego.Eval(ctx)
		if err != nil {
			return resp, status.Errorf(codes.Aborted, err.Error())
		}

		if res.Allowed() {
			return handler(ctx, req)
		}

		// default is deny
		return resp, status.Errorf(codes.PermissionDenied, "unauthorized")
	}
}
