package serrors

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrBadInput      = errors.New("bad input")
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal")
	ErrForbidden     = errors.New("forbidden")
	ErrAuthorization = errors.New("unauthorized")
)

func GetHttpCodeFrom(err error) int {
	switch {
	case errors.Is(err, ErrBadInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrAuthorization):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func GetErrorByTypeAndLog(err error) error {
	log.Println(err)
	switch {
	case errors.Is(err, ErrBadInput):
		return nil
	case errors.Is(err, ErrNotFound):
		return nil
	case errors.Is(err, ErrInternal):
		return err
	case errors.Is(err, ErrForbidden):
		return nil
	case errors.Is(err, ErrAuthorization):
		return nil
	default:
		return err
	}
}
