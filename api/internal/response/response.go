package response

import (
	"net/http"
)

type Response struct {
	// Code is the HTTP status code
	Code int `json:"-"`
	// Message is the message body
	Message string `json:"message"`
}

var dftResponse = &Response{
	Code:    http.StatusOK,
	Message: "ok",
}

func New(err error) *Response {
	if err == nil {
		return dftResponse
	}

	return &Response{Code: http.StatusInternalServerError, Message: "unexpected error"}
}
