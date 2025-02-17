package httputils

import (
	"net/http"
)

type HttpError struct {
	Message    string
	StatusCode int
}

func (err *HttpError) Error() string {
	return http.StatusText(err.StatusCode) + ": " + err.Message
}

// HandleHTTPErrors centralizes error handling for http-related operations
func HandleHTTPErrors(w http.ResponseWriter, err error) {
	if err != nil {
		if e, ok := err.(*HttpError); ok {
			WriteJSONError(w, e.Error(), e.StatusCode)
			return
		} else {
			WriteJSONError(w, e.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewBadRequestError(message string) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewConflictError(message string) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewForbiddenError(message string) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(message string) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewUnprocessableEntityError(message string) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
	}
}
