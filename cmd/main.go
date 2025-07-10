package main

import (
	"context"
	"log/slog"
	"naqet/bookmarks/infra/database"
	"naqet/bookmarks/services/auth"
	"naqet/bookmarks/services/marks"
	"naqet/bookmarks/utils"
	"naqet/bookmarks/views/pages"
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

	r.Group(func(g chi.Router) {
		g.Use(auth.NewMiddleware(db))
		g.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)
			if !ok {
				utils.Unauthorized(w)
				return
			}
			marks := []database.Bookmark{}
			res, err := db.Query("select title, url, tags, description, read, created_at from bookmarks where owner_id = $1", userId)

			if err != nil {
				slog.Error("couldn't prepare query for selecting bookmarks", slog.Any("error", err))
				utils.InternalServerError(w)
				return
			}

			for res.Next() {
				mark := database.Bookmark{}
				err := res.Scan(&mark.Title, &mark.Url, &mark.Tags, &mark.Description, &mark.Read, &mark.CreatedAt)

				if err != nil {
					slog.Error("couldn't scan bookmark", slog.Any("error", err))
					utils.InternalServerError(w)
					return
				}

				marks = append(marks, mark)
			}

			if err := pages.Index(marks).Render(context.Background(), w); err != nil {
				w.Write([]byte("Error"))
			}
		})
		marks.Init(g, db, vali)
	})

	http.ListenAndServe(":3000", r)
}
