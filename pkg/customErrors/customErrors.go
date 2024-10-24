package customErrors

import (
	"fmt"
	"net/http"
)

type BadRequestError struct {
	Msg string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.Msg)
}

type UnauthorizedError struct {
	Msg string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Msg)
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{Msg: msg}
}

func NewUnauthorizedError(msg string) error {
	return &UnauthorizedError{Msg: msg}
}

func HandleError(err error) (int, map[string]string) {
	switch e := err.(type) {
	case *BadRequestError:
		return http.StatusBadRequest, map[string]string{"error": e.Error()}
	case *UnauthorizedError:
		return http.StatusUnauthorized, map[string]string{"error": e.Error()}
	default:
		return http.StatusInternalServerError, map[string]string{"error": "something went wrong"}
	}
}
