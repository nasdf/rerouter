package rerouter

import (
	"context"
	"net/http"
	"regexp"
)

type contextKey string

var paramsContextKey = contextKey("params")

// Params returns the params for the given request.
func Params(req *http.Request) map[string]any {
	val, _ := req.Context().Value(paramsContextKey).(map[string]any)
	return val
}

// setParams sets the params for the given request.
func setParams(req *http.Request, params map[string]any) *http.Request {
	ctx := context.WithValue(req.Context(), paramsContextKey, params)
	return req.WithContext(ctx)
}

// matchParams returns a bool indicating a match and a key value map of params.
func matchParams(pattern *regexp.Regexp, path string) (map[string]any, bool) {
	matches := pattern.FindStringSubmatch(path)
	if len(matches) == 0 {
		return nil, false
	}
	params := make(map[string]any)
	for _, name := range pattern.SubexpNames() {
		i := pattern.SubexpIndex(name)
		if i >= 0 && i < len(matches) {
			params[name] = matches[pattern.SubexpIndex(name)]
		}
	}
	return params, true
}
