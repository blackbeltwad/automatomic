package middleware

import (
	"context"
	"net/http"
	"strings"

	"automatomic/internal/auth"
	"automatomic/internal/model"
)

type contextKey string

const ClaimsContextKey contextKey = "userClaims"

// JWTMiddleware checks for access_token in cookies first, then Authorization header as fallback
func JWTMiddleware(jwtMgr *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string

			// 1. Try extracting token from HTTP-Only cookie
			if cookie, err := r.Cookie("access_token"); err == nil {
				tokenStr = cookie.Value
			}

			// 2. Fallback: Try Authorization header (Bearer <token>)
			if tokenStr == "" {
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if tokenStr == "" {
				http.Error(w, `{"error":"missing authentication token"}`, http.StatusUnauthorized)
				return
			}

			// 3. Verify token signature and claims
			claims, err := jwtMgr.Verify(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"unauthorized: invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// 4. Inject claims into request context for downstream handlers & RequireScope
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireScope enforces scope permissions using the claims injected by JWTMiddleware
func RequireScope(requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*model.Claims)
			if !ok || claims == nil {
				http.Error(w, `{"error":"unauthorized context"}`, http.StatusUnauthorized)
				return
			}

			hasScope := false
			for _, scope := range claims.Scopes {
				if scope == requiredScope || scope == model.ScopeAdmin {
					hasScope = true
					break
				}
			}

			if !hasScope {
				http.Error(w, `{"error":"forbidden: insufficient permissions"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}