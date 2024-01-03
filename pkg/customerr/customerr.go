package customerr

import (
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message: %v - detail: %v", e.Message, e.Details)
}

func New(message string, details any) error {
	return &Error{
		Message: message,
		Details: details,
	}
}

func WithStatus(httpStatusCode int, message string, details any) error {
	return &Error{
		StatusCode: httpStatusCode,
		Message:    message,
		Details:    details,
	}
}

func StatusCode(err error) int {
	e, ok := err.(*Error)
	if !ok {
		return http.StatusInternalServerError
	}

	return e.StatusCode
}
