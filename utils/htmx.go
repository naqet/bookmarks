package utils

import "net/http"

func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func SetHtmxRedirect(w http.ResponseWriter, destination string) {
	w.Header().Add("HX-Redirect", destination)
}

func SetHtmxEventAfterSwap(w http.ResponseWriter, event string) {
	w.Header().Add("HX-Trigger-After-Swap", event)
}

func SetHtmxRefresh(w http.ResponseWriter) {
	w.Header().Add("HX-Refresh", "true")
}
