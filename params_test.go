package rerouter

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchParamsWithOptional(t *testing.T) {
	pattern := regexp.MustCompile("^/(index.html)?$")

	params, ok := matchParams(pattern, "/")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{}, params)

	params, ok = matchParams(pattern, "/index.html")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{}, params)

	params, ok = matchParams(pattern, "/about.html")
	assert.False(t, ok)
	assert.Equal(t, map[string]any(nil), params)
}

func TestMatchParamsWithPartial(t *testing.T) {
	pattern := regexp.MustCompile("/about")

	params, ok := matchParams(pattern, "/about")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{}, params)

	params, ok = matchParams(pattern, "/info/about")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{}, params)

	params, ok = matchParams(pattern, "/about/info")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{}, params)
}

func TestMatchParamsWithWildCardCaptureInMiddle(t *testing.T) {
	pattern := regexp.MustCompile("^/things/(?P<name>.*)/info$")

	params, ok := matchParams(pattern, "/things/one/info")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{"name": "one"}, params)

	params, ok = matchParams(pattern, "/things/one/info/info")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{"name": "one/info"}, params)
}

func TestMatchParamsWithOptionalWildCardCaptureAtEnd(t *testing.T) {
	pattern := regexp.MustCompile("^/places(/(?P<name>.*))?$")

	params, ok := matchParams(pattern, "/places")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{"name": ""}, params)

	params, ok = matchParams(pattern, "/places/earth")
	assert.True(t, ok)
	assert.Equal(t, map[string]any{"name": "earth"}, params)
}

func TestSetThenGetParams(t *testing.T) {
	params := map[string]any{"one": 1, "two": 2}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	req = setParams(req, params)
	assert.Equal(t, params, Params(req))
}
