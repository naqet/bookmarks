package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"naqet/bookmarks/utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func NewMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(utils.AUTHORIZATION)

			if err != nil || cookie == nil {
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusSeeOther)
				// utils.Unauthorized(w)
				return
			}

			tokenString := cookie.Value

			secret, ok := os.LookupEnv("JWT_SECRET")
			if !ok {
				slog.Error("JWT secret is not set")
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusSeeOther)
				// utils.InternalServerError(w)
				return
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Invalid JWT token")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusSeeOther)
				// utils.Unauthorized(w)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok {
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusSeeOther)
				// utils.Unauthorized(w)
				return
			}

			id, err := claims.GetSubject()

			if err != nil {
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusSeeOther)
				// utils.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), utils.USER_ID_CTX_KEY, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
