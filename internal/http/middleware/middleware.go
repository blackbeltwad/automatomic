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

func JWTMiddleware(jwtMgr *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"invalid bearer token format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtMgr.Verify(parts[1])
			if err != nil {
				http.Error(w, `{"error":"unauthorized token validation failed"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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