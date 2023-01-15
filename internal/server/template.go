package server

import (
	"encoding/json"
	"net/http"
)

type response struct {
	ErrorMessage interface{} `json:"error_message,omitempty"`
	Data         interface{} `json:"data"`
}

func RenderResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	res := response{
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	d, _ := json.Marshal(res)
	w.Write(d)
}

func RenderError(w http.ResponseWriter, statusCode int, err error) {
	res := response{
		ErrorMessage: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	d, _ := json.Marshal(res)
	w.Write(d)
}
