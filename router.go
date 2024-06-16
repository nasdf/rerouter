package rerouter

import (
	"net/http"
	"regexp"
)

var _ (http.Handler) = (*router)(nil)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type router struct {
	routes []route
}

// New returns a new router.
func New() *router {
	return &router{}
}

// Handle registers a new handler that will match requests with the given pattern.
func (r *router) Handle(pattern *regexp.Regexp, handler http.Handler) {
	r.routes = append(r.routes, route{pattern, handler})
}

func (r *router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		params, ok := matchParams(route.pattern, req.URL.Path)
		if !ok {
			continue // route did not match
		}
		route.handler.ServeHTTP(rw, setParams(req, params))
		return // route handled
	}
	// if no matches found return a 404
	rw.WriteHeader(http.StatusNotFound)
}
