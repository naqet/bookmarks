package dashboard

import (
	"log/slog"
	"naqet/bookmarks/infra/database"
	"naqet/bookmarks/utils"
	"naqet/bookmarks/views/pages"
	"net/http"
)

func (h *dashboardHandler) homePage(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)
	if !ok {
		utils.Unauthorized(w)
		return
	}
	marks := []database.Bookmark{}
	res, err := h.db.Query("select id, title, url, tags, description, read, created_at from bookmarks where owner_id = $1", userId)

	if err != nil {
		slog.Error("couldn't prepare query for selecting bookmarks", slog.Any("error", err))
		utils.InternalServerError(w)
		return
	}

	for res.Next() {
		mark := database.Bookmark{}
		err := res.Scan(&mark.ID, &mark.Title, &mark.Url, &mark.Tags, &mark.Description, &mark.Read, &mark.CreatedAt)

		if err != nil {
			slog.Error("couldn't scan bookmark", slog.Any("error", err))
			utils.InternalServerError(w)
			return
		}

		marks = append(marks, mark)
	}

	if err := pages.Index(marks).Render(r.Context(), w); err != nil {
		w.Write([]byte("Error"))
	}
}

func (h *dashboardHandler) settingsPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.Settings().Render(r.Context(), w); err != nil {
		w.Write([]byte("Error"))
	}
}
