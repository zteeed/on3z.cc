package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"llil.gq/go/database"
	"net/http"
)

type Route struct {
	Name    string
	Pattern string
	Handler http.Handler
}

type Routes []Route

var (
	appVersion = "0.0.1"
	version    = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version",
		Help: "Version information about this binary",
		ConstLabels: map[string]string{
			"version": appVersion,
		},
	})

	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of all HTTP requests",
	}, []string{"code", "handler", "method"})
)

func NewRouter(db database.Database, baseUrl string) *http.ServeMux {
	router := http.NewServeMux()

	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpRequestDuration)
	r.MustRegister(version)

	var routes = Routes{
		Route{
			"DataShortenHandler",
			"/data/shorten",
			promhttp.InstrumentHandlerDuration(
				httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "/data/shorten"}),
				promhttp.InstrumentHandlerCounter(httpRequestsTotal, &DataShortenHandler{db, baseUrl}),
			),
		},
		Route{
			"RootHandler",
			"/",
			promhttp.InstrumentHandlerDuration(
				httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "/"}),
				promhttp.InstrumentHandlerCounter(httpRequestsTotal, &RootHandler{db}),
			),
		},
		Route{
			"Prometheus",
			"/metrics",
			promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
		},
	}

	for _, route := range routes {
		handler := Logger(route.Handler, route.Name)
		router.Handle(route.Pattern, handler)
	}

	return router
}
