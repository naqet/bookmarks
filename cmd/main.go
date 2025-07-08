package main

import (
	"context"
	"fmt"
	"naqet/bookmarks/views/pages"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
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
