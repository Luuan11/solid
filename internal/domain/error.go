package domain

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBookNotFound      = NewDomainError("BOOK_NOT_FOUND", "book not found", http.StatusNotFound)
	ErrBookAlreadyExists = NewDomainError("BOOK_ALREADY_EXISTS", "book already exists", http.StatusConflict)
	ErrInvalidInput      = NewDomainError("INVALID_INPUT", "invalid input", http.StatusBadRequest)
)

type DomainError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func NewDomainError(code, message string, statusCode int) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) WithError(err error) *DomainError {
	return &DomainError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Err:        err,
	}
}

func (e *DomainError) WithMessage(msg string) *DomainError {
	return &DomainError{
		Code:       e.Code,
		Message:    msg,
		StatusCode: e.StatusCode,
		Err:        e.Err,
	}
}

func GetStatusCode(err error) int {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.StatusCode
	}
	return http.StatusInternalServerError
}
