package marks

import (
	"database/sql"
	"errors"
	"log/slog"
	"naqet/bookmarks/infra/database"
	"naqet/bookmarks/utils"
	"naqet/bookmarks/views/components"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func (h *marksHandler) getMarks(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)
	if !ok {
		utils.Unauthorized(w)
		return
	}
	marks := []database.Bookmark{}

	queryValues := r.URL.Query()
	query := strings.ToLower(queryValues.Get("query"))

	res, err := h.db.Query("select id, title, url, tags, description, read, created_at from bookmarks where owner_id = $1 and lower(title) like $2", userId, "%"+query+"%")

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

	components.BookmarkCards(marks).Render(r.Context(), w)
}

func (h *marksHandler) getInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	url := query.Get("url")
	if err := h.vali.Var(url, "required,http_url"); err != nil {
		utils.BadRequest(w, "Invalid URL")
		return
	}
	res, err := http.Get(url)
	if err != nil {
		utils.InternalServerError(w, "Page isn't reachable")
		return
	}
	n, err := html.Parse(res.Body)
	if err != nil {
		slog.Error(err.Error())
		utils.InternalServerError(w)
		return
	}
	defer res.Body.Close()

	// TODO: add favicon handling after s3 integration
	title := findTitle(n)
	desc := findMeta(n, []string{"description", "og:description", "twitter:description"})

	if err := components.BookmarkTitleAndDesc(title, desc).Render(r.Context(), w); err != nil {
		utils.InternalServerError(w)
		return
	}
}

func (h *marksHandler) createMark(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Url         string `validate:"required,url"`
		Tags        string
		Title       string `validate:"required"`
		Description string
	}

	data := request{
		Url:         r.PostFormValue("url"),
		Tags:        r.PostFormValue("tags"),
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
	}

	if err := h.vali.Struct(&data); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)
	if !ok {
		utils.Unauthorized(w)
		return
	}

	var exists bool
	if err := h.db.QueryRow("select exists(select 1 from bookmarks where title = $1)", data.Title).Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("couldn't check if bookmark exists", slog.Any("error", err))
			utils.InternalServerError(w)
			return
		}
	}

	if exists {
		utils.BadRequest(w, "Bookmark with this title already exists")
		return
	}

	if _, err := h.db.Exec(
		"insert into bookmarks (url, tags, title, description, owner_id) values ($1, $2, $3, $4, $5)",
		data.Url,
		data.Tags,
		data.Title,
		data.Description,
		userId,
	); err != nil {
		slog.Error("couldn't insert bookmark", slog.Any("error", err))
		utils.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
}

func (h *marksHandler) getMark(w http.ResponseWriter, r *http.Request) {}

func (h *marksHandler) getEditMarkModal(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)

	if !ok {
		utils.Unauthorized(w)
		return
	}
	id := r.PathValue("id")

	mark := database.Bookmark{}

	if err := h.db.QueryRow(
		"select id, title, url, tags, description, read, created_at from bookmarks where owner_id = $1 and id = $2",
		userId,
		id,
	).Scan(&mark.ID, &mark.Title, &mark.Url, &mark.Tags, &mark.Description, &mark.Read, &mark.CreatedAt); err != nil {
		utils.BadRequest(w)
		return
	}

	components.EditBookmarkModal(mark).Render(r.Context(), w)
}

func (h *marksHandler) editMark(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)
	if !ok {
		utils.Unauthorized(w)
		return
	}

	type request struct {
		Title       string `validate:"required"`
		Url         string `validate:"required,url"`
		Tags        string
		Description string
	}

	data := request{
		Url:         r.PostFormValue("url"),
		Tags:        r.PostFormValue("tags"),
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
	}

	if err := h.vali.Struct(&data); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	id := r.PathValue("id")

	if _, err := h.db.Exec(
		"update bookmarks set title = $1, url = $2, tags = $3, description = $4 where id = $5 and owner_id = $6",
		data.Title,
		data.Url,
		data.Tags,
		data.Description,
		id,
		userId,
	); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(http.StatusText(http.StatusAccepted)))
}

func (h *marksHandler) deleteMark(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.USER_ID_CTX_KEY).(string)
	if !ok {
		utils.Unauthorized(w)
		return
	}

	id := r.PathValue("id")

	var exists bool
	if err := h.db.QueryRow("select exists(select 1 from bookmarks where id = $1 and owner_id = $2)", id, userId).Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.BadRequest(w)
			return
		}
	}

	if !exists {
		utils.BadRequest(w)
		return
	}

	if _, err := h.db.Exec("delete from bookmarks where id = $1 and owner_id = $2", id, userId); err != nil {
		utils.BadRequest(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(http.StatusText(http.StatusAccepted)))
}
