package app

import (
	"llil.gq/go/database"
	"net/http"
	"strings"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

func NewRouter(db database.Database, baseUrl string) *http.ServeMux {
	router := http.NewServeMux()
	var routes = Routes{
		Route{
			"DataShortenHandler",
			strings.ToUpper("Post"),
			"/data/shorten",
			&DataShortenHandler{db, baseUrl},
		},
		Route{
			"RootHandler",
			strings.ToUpper("Get"),
			"/",
			&RootHandler{db},
		},
	}
	for _, route := range routes {
		handler := Logger(route.Handler, route.Name)
		router.Handle(route.Pattern, handler)
	}

	return router
}
