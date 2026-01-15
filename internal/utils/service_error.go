package utils

import (
	"errors"
	"net/http"
)

func HandleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, ErrUnauthorized.Error())

	case errors.Is(err, ErrForbidden):
		WriteError(w, http.StatusForbidden, ErrForbidden.Error())

	case errors.Is(err, ErrInvalidToken),
		errors.Is(err, ErrExpiredToken):
		WriteError(w, http.StatusUnauthorized, err.Error())

	case errors.Is(err, ErrMissingUserID),
		errors.Is(err, ErrMissingUserRole),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrWeakPassword),
		errors.Is(err, ErrInvalidEmail):
		WriteError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrAlreadyExists),
		errors.Is(err, ErrConflict):
		WriteError(w, http.StatusConflict, err.Error())

	case errors.Is(err, ErrNotFound):
		WriteError(w, http.StatusNotFound, err.Error())

	default:
		WriteError(w, http.StatusInternalServerError, ErrInternal.Error())
	}
}
