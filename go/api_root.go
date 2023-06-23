package app

import (
	"on3z.cc/go/database"
	"net/http"
	"strings"
)

type RootHandler struct {
	database database.Database
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	result, err := database.SelectShortURL(h.database, shortURL)
	if err != nil {
		w.Header().Set("Location", "/404")
		w.WriteHeader(302)
		return
	}
	w.Header().Set("Location", result.LongURL)
	w.WriteHeader(301)
}
