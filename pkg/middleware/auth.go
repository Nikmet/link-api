package middleware

import (
	"context"
	"go-advanced/configs"
	"go-advanced/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func IsAuthed(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, data := jwt.NewJWT(conf.Auth.Secret).Parse(token)

		if !valid || data == nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Check if email is present in the token data
		if data.Email == "" {
			http.Error(w, "Invalid token: email claim missing", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
