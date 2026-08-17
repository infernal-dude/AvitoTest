package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const signKey = "avitoTheBest"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.Header.Get("Authorization")
		if rawToken == "" || !strings.HasPrefix(rawToken, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(rawToken, "Bearer ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(signKey), nil
		})

		// if err != nil || !token.Valid {
		// 	http.Error(w, "expired or invalid token", http.StatusBadRequest)
		// 	return
		// }
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "expired or invalid token пидр", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_type", claims.UserType)
		next(w, r.WithContext(ctx))
	})
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserType string    `json:"user_type"`
	jwt.RegisteredClaims
}
