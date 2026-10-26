package mockengine

import (
	"net/http"

	"mock-platform/internal/config"
)

type Engine struct {
	routes []config.Route
}

func New(cfg *config.Config) *Engine {
	return &Engine{
		routes: cfg.Routes,
	}
}

func (e *Engine) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	for _, route := range e.routes {
		if route.Method == r.Method &&
			matchPath(route.Path, r.URL.Path) {
			e.writeResponse(w, route.Response)
			return
		}
	}

	http.NotFound(w, r)
}

func (e *Engine) writeResponse(
	w http.ResponseWriter,
	response config.Response,
) {
	for name, value := range response.Headers {
		w.Header().Set(name, value)
	}

	w.WriteHeader(response.Status)

	if len(response.Body) == 0 {
		return
	}

	_, _ = w.Write(response.Body)
}
