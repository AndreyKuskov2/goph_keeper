// Package http_handlers provides HTTP request handlers for the GophKeeper server.
// It includes handlers for user management, data storage, and utility endpoints.
package http_handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type responseBad struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// responseError sends an error response with the specified status code and message.
func responseError(w http.ResponseWriter, r *http.Request, code int, message string) {
	resp := responseBad{
		Code:  code,
		Error: message,
	}
	render.Status(r, code)
	render.JSON(w, r, resp)
}

// responseOK sends a successful response with the specified status code, data, and message.
func responseOK(w http.ResponseWriter, r *http.Request, code int, data any, message string) {
	resp := response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	render.Status(r, code)
	render.JSON(w, r, resp)
}
