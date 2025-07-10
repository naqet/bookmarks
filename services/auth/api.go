package auth

import (
	"database/sql"
	"errors"
	"log/slog"
	"naqet/bookmarks/infra/database"
	"naqet/bookmarks/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
	}

	data := request{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}

	if err := h.vali.Struct(&data); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	user := database.User{}
	if err := h.db.QueryRow("select id, password from users where username = $1", data.Username).Scan(&user.ID, &user.Password); err != nil {
		slog.Error(err.Error())
		utils.BadRequest(w, "Invalid username or password")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

	if err != nil {
		utils.BadRequest(w, "Invalid username or password")
		return
	}

	expiration := time.Now().Add(time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"time": expiration.Unix(),
	})

	secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		slog.Error("JWT_SECRET env not set")
		utils.InternalServerError(w)
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("signing token failed", slog.Any("error", err))
		utils.InternalServerError(w)
		return
	}

	cookie := http.Cookie{
		Name:     utils.AUTHORIZATION,
		Value:    tokenString,
		Path:     "/",
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	// if utils.IsHtmxRequest(r) {
	// 	utils.AddHtmxRedirect(w, "/dashboard")
	// }

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username        string `validate:"required"`
		Password        string `validate:"required"`
		PasswordConfirm string `validate:"required,eqfield=Password"`
	}

	data := request{
		Username:        r.PostFormValue("username"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password-confirm"),
	}

	if err := h.vali.Struct(&data); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	var exists bool
	if err := h.db.QueryRow("select exists(select 1 from users where username = $1)", data.Username).Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.InternalServerError(w, err.Error())
			return
		}
	}

	if exists {
		utils.BadRequest(w, "Username taken")
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("bcrypting password failed", slog.Any("error", err))
		utils.InternalServerError(w)
		return
	}

	if _, err = h.db.Exec("insert into users (username, password) values ($1, $2)", data.Username, pass); err != nil {
		slog.Error("inserting user into db failed", slog.Any("error", err))
		utils.InternalServerError(w)
		return
	}

	// if utils.IsHtmxRequest(r) {
	// 	utils.AddHtmxRedirect(w, "/auth/login")
	// }

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
}

func (h *authHandler) logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     utils.AUTHORIZATION,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	// if utils.IsHtmxRequest(r) {
	// 	utils.AddHtmxRedirect(w, "/")
	// }
	//
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
