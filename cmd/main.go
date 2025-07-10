package main

import (
	"log/slog"
	"naqet/bookmarks/infra/database"
	"naqet/bookmarks/services/auth"
	"naqet/bookmarks/services/dashboard"
	"naqet/bookmarks/services/marks"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("env variables couldn't be loaded", slog.Any("error", err))
		os.Exit(1)
	}

	db := database.Init()
	if err := database.Migrate(db); err != nil {
		slog.Error("migration failed", slog.Any("error", err))
		os.Exit(1)
	}

	vali := validator.New(validator.WithRequiredStructEnabled())

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	auth.Init(r, db, vali)

	// Protected
	r.Group(func(g chi.Router) {
		g.Use(auth.NewMiddleware(db))
		marks.Init(g, db, vali)
		dashboard.Init(g, db, vali)
	})

	http.ListenAndServe(":3000", r)
}
