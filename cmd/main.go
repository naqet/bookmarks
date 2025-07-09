package main

import (
	"context"
	"fmt"
	"log/slog"
	"naqet/bookmarks/infra/database"
	"naqet/bookmarks/views/pages"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := pages.Index().Render(context.Background(), w); err != nil {
			fmt.Println(err)
			w.Write([]byte("Error"))
		}
	})

	http.ListenAndServe(":3000", r)
}
