package auth

import (
	"naqet/bookmarks/views/pages"
	"net/http"
)

func (h *authHandler) loginPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.Login().Render(r.Context(), w); err != nil {
		w.Write([]byte("Error"))
	}
}

func (h *authHandler) signUpPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.SignUp().Render(r.Context(), w); err != nil {
		w.Write([]byte("Error"))
	}
}
