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

func responseError(w http.ResponseWriter, r *http.Request, code int, message string) {
	resp := responseBad{
		Code:  code,
		Error: message,
	}
	render.Status(r, code)
	render.JSON(w, r, resp)
}

func responseOK(w http.ResponseWriter, r *http.Request, code int, data any, message string) {
	resp := response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	render.Status(r, code)
	render.JSON(w, r, resp)
}
