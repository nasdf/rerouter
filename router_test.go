package rerouter

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterMatch(t *testing.T) {
	pattern := regexp.MustCompile("^/test$")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, map[string]any{}, Params(r))
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rw := httptest.NewRecorder()

	router := New()
	router.Handle(pattern, handler)
	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusTeapot, rw.Result().StatusCode)
}

func TestRouterMatchWithParams(t *testing.T) {
	pattern := regexp.MustCompile("^/test/(?P<name>.*)$")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, map[string]any{"name": "stuff"}, Params(r))
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/stuff", nil)
	rw := httptest.NewRecorder()

	router := New()
	router.Handle(pattern, handler)
	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusTeapot, rw.Result().StatusCode)
}

func TestRouterNoMatch(t *testing.T) {
	pattern := regexp.MustCompile("^/test$")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	rw := httptest.NewRecorder()

	router := New()
	router.Handle(pattern, handler)
	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Result().StatusCode)
}
