package apperror

import "net/http"

type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	HTTPStatus int
	Code       string
	Message    string
	Details    []Detail
}

func (e *AppError) Error() string {
	return e.Message
}

func New(httpStatus int, code, message string, details ...Detail) *AppError {
	return &AppError{HTTPStatus: httpStatus, Code: code, Message: message, Details: details}
}

func Validation(message string, details ...Detail) *AppError {
	return New(http.StatusBadRequest, "VALIDATION_ERROR", message, details...)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(message string, details ...Detail) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, details...)
}

func BusinessRule(message string) *AppError {
	return New(http.StatusUnprocessableEntity, "BUSINESS_RULE_VIOLATION", message)
}

func TooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
