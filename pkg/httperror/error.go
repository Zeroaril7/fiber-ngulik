package httperror

import "net/http"

type CommonErrorData struct {
	Code         int    `json:"code"`
	ResponseCode int    `json:"responseCode,omitempty"`
	Message      string `json:"message"`
}

type ErrorString struct {
	code    int
	message string
}

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.message
}

func (e ErrorString) Message() string {
	return e.message
}

func NewErrorString(code int, defaultMessage, customMessage string) ErrorString {
	if customMessage == "" {
		customMessage = defaultMessage
	}
	return ErrorString{
		code:    code,
		message: customMessage,
	}
}

func NewBadRequest(msg string) ErrorString {
	return NewErrorString(http.StatusBadRequest, "Bad Request", msg)
}

func NewUnauthorized(msg string) ErrorString {
	return NewErrorString(http.StatusUnauthorized, "Unauthorized", msg)
}

func NewForbidden(msg string) ErrorString {
	return NewErrorString(http.StatusForbidden, "Forbidden", msg)
}

func NewNotFound(msg string) ErrorString {
	return NewErrorString(http.StatusNotFound, "Not Found", msg)
}

func NewConflict(msg string) ErrorString {
	return NewErrorString(http.StatusConflict, "Conflict", msg)
}

func NewInternalServerError(msg string) ErrorString {
	return NewErrorString(http.StatusInternalServerError, "Internal Server Error", msg)
}

func BadRequest(msg string) error {
	return NewBadRequest(msg)
}

func Unauthorized(msg string) error {
	return NewUnauthorized(msg)
}

func Forbidden(msg string) error {
	return NewForbidden(msg)
}

func NotFound(msg string) error {
	return NewNotFound(msg)
}

func Conflict(msg string) error {
	return NewConflict(msg)
}

func InternalServerError(msg string) error {
	return NewInternalServerError(msg)
}
