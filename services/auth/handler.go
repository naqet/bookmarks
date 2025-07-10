package auth

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	db   *sql.DB
	vali *validator.Validate
}

func Init(r chi.Router, db *sql.DB, vali *validator.Validate) *authHandler {
	h := &authHandler{db, vali}
	r.Get("/login", h.loginPage)
	r.Get("/signup", h.signUpPage)

	r.Post("/api/auth/login", h.login)
	r.Post("/api/auth/signup", h.signUp)
	r.Post("/api/auth/logout", h.logout)

	return h
}
