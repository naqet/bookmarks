package marks

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type marksHandler struct {
	db   *sql.DB
	vali *validator.Validate
}

func Init(r chi.Router, db *sql.DB, vali *validator.Validate) *marksHandler {
	h := &marksHandler{db, vali}

	r.Get("/api/marks", h.getMarks)
	r.Post("/api/marks", h.createMark)

	r.Get("/api/marks/get-info", h.getInfo)

	r.Get("/api/marks/{id}", h.getMark)
	r.Put("/api/marks/{id}", h.editMark)
	r.Delete("/api/marks/{id}", h.deleteMark)

	r.Get("/api/marks/edit-modal/{id}", h.getEditMarkModal)
	return h
}
