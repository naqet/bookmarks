package utils

import "net/http"

func InternalServerError(w http.ResponseWriter, messages ...string) {
	handleHttpError(w, http.StatusInternalServerError, messages)
}

func BadRequest(w http.ResponseWriter, messages ...string) {
	handleHttpError(w, http.StatusBadRequest, messages)
}

func Unauthorized(w http.ResponseWriter, messages ...string) {
	handleHttpError(w, http.StatusUnauthorized, messages)
}

func handleHttpError(w http.ResponseWriter, status int, messages []string) {
	msg := http.StatusText(status)
	if len(messages) >= 1 {
		msg = messages[0]
	}
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
