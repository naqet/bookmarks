package dashboard;

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type dashboardHandler struct {
	db   *sql.DB
	vali *validator.Validate
}

func Init(r chi.Router, db *sql.DB, vali *validator.Validate) *dashboardHandler {
	h := &dashboardHandler{db, vali}

	r.Get("/", h.homePage)
	r.Get("/settings", h.settingsPage)

	return h
}
