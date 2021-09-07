package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/todo-app/internal"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/pkg/logger"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a defferred function (which will aslways be run in the envent
		// of a panic as Go unwinds the stack).
		defer func() {
			// User the built in recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Return a 500 Internal server Error response
				internal.ErrInternalServer(fmt.Errorf("%s", err), "Internal server error").Send(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// AuthenticationMiddleware purposefully returns a http.HandlerFunc rather
// than an http.handle so that it can be applied to individual routes and
// not used on every single route.
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Hello from auth middleware")

		// Get the token from the cookie
		token, err := identity.GetTokenFromCookie(r)
		if err != nil || token == "" {
			internal.ErrUnauthorized(err, "unauthorized").Send(w)
			return
		}
		// Extract the info from the token and place it in the claims var
		claims, err := identity.ExtractClaimsFromToken(token)
		if err != nil {
			internal.ErrUnauthorized(err, "unauthorized").Send(w)
			return
		}

		// place the user claims (id, email) in the context
		ctx := context.WithValue(r.Context(), identity.UserCtxKey, claims)

		// Set the context on a new request struct to pass it to next
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
