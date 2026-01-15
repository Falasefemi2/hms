package utils

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token expired")
	ErrMissingUserID   = errors.New("user id not found in context")
	ErrMissingUserRole = errors.New("user role not found in context")
	ErrInvalidInput    = errors.New("invalid input")
	ErrWeakPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrNotFound        = errors.New("resource not found")
	ErrAlreadyExists   = errors.New("resource already exists")
	ErrConflict        = errors.New("resource conflict")
	ErrInternal        = errors.New("internal server error")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
)
