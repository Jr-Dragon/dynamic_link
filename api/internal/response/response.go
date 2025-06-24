package response

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jr-dragon/dynamic_link/internal/library/logs"
	"log/slog"
	"net/http"
)

type Response struct {
	// Code is the HTTP status code
	Code int `json:"-"`
	// Message is the message body
	Message string `json:"message"`
	// Data is the response data, can be any type
	Data any `json:"data,omitempty"`
}

var dftResponse = &Response{
	Code:    http.StatusOK,
	Message: "ok",
}

func Err(err error) *Response {
	if err == nil {
		return dftResponse
	}

	slog.Error("", logs.Err(err))

	if errors.Is(err, fiber.ErrUnprocessableEntity) {
		return &Response{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if e := new(validator.ValidationErrors); errors.As(err, e) {
		return &Response{Code: http.StatusBadRequest, Message: e.Error()}
	}

	return &Response{Code: http.StatusInternalServerError, Message: "unexpected error"}
}

func Data(data any) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	}
}
