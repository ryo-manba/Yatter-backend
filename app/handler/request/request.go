package request

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// Read path parameter `id`
func IDOf(r *http.Request) (int64, error) {
	ids := chi.URLParam(r, "id")

	if ids == "" {
		return -1, errors.Errorf("id was not presence")
	}

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return -1, errors.Errorf("id was not number")
	}

	return id, nil
}

// Read path parameter `username`
func UsernameOf(r *http.Request) (string, error) {
	username := chi.URLParam(r, "username")

	if username == "" {
		return "", errors.Errorf("username was not presence")
	}
	return username, nil
}

// Read query parameter `key` and return it as int64
func QueryInt64(r *http.Request, key string, defaultValue int64) (int64, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, errors.Errorf("query parameter '%s' was not a number", key)
	}
	if parsedValue < 0 {
		return -1, errors.Errorf("query parameter '%s' must not be a negative number", key)
	}

	return parsedValue, nil
}

// Read query parameter `key` and return it as bool
func QueryBool(r *http.Request, key string, defaultValue bool) (bool, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.Errorf("query parameter '%s' was not a boolean", key)
	}

	return parsedValue, nil
}
