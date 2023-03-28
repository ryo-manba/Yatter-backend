package httperror

import (
	"log"
	"net/http"
)

// Response with given status code
func Error(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// Response with Bad Request (400)
func BadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// Response with Not Found (404)
func NotFound(w http.ResponseWriter, err error) {
	log.Printf("[NotFound] %+v", err)

	Error(w, http.StatusNotFound)
}

// Response with Internal Server Error (500)
func InternalServerError(w http.ResponseWriter, err error) {
	log.Printf("[InternalServerError] %+v", err)

	Error(w, http.StatusInternalServerError)
}
