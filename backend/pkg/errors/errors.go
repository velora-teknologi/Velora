package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ValidationError   ErrorType = "VALIDATION_ERROR"
	NotFoundError     ErrorType = "NOT_FOUND"
	UnauthorizedError ErrorType = "UNAUTHORIZED"
	ForbiddenError    ErrorType = "FORBIDDEN"
	InternalError     ErrorType = "INTERNAL_ERROR"
	ConflictError     ErrorType = "CONFLICT"
)

type CustomError struct {
	Type       ErrorType
	Message    string
	StatusCode int
	Details    map[string]interface{}
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewValidationError(message string) *CustomError {
	return &CustomError{
		Type:       ValidationError,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    make(map[string]interface{}),
	}
}

func NewNotFoundError(resource string) *CustomError {
	return &CustomError{
		Type:       NotFoundError,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
		Details:    make(map[string]interface{}),
	}
}

func NewUnauthorizedError(message string) *CustomError {
	return &CustomError{
		Type:       UnauthorizedError,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Details:    make(map[string]interface{}),
	}
}

func NewForbiddenError(message string) *CustomError {
	return &CustomError{
		Type:       ForbiddenError,
		Message:    message,
		StatusCode: http.StatusForbidden,
		Details:    make(map[string]interface{}),
	}
}

func NewInternalError(message string) *CustomError {
	return &CustomError{
		Type:       InternalError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Details:    make(map[string]interface{}),
	}
}

func NewConflictError(message string) *CustomError {
	return &CustomError{
		Type:       ConflictError,
		Message:    message,
		StatusCode: http.StatusConflict,
		Details:    make(map[string]interface{}),
	}
}
