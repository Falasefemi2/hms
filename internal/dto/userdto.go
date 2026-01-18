package dto

import "time"

type PatientSignUpRequest struct {
	Username  string `json:"username" example:"john_doe" validate:"required,min=3,max=255"`
	Email     string `json:"email" example:"john@example.com" validate:"required,email"`
	Password  string `json:"password" example:"SecurePass123!" validate:"required,min=8"`
	FirstName string `json:"first_name" example:"John" validate:"required,max=255"`
	LastName  string `json:"last_name" example:"Doe" validate:"required,max=255"`
	Phone     string `json:"phone" example:"+2345694004" validate:"omitempty,max=20"`
}

type AdminCreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=255"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Phone     string `json:"phone" validate:"omitempty,max=20"`
	Role      string `json:"role" validate:"required,oneof=DOCTOR NURSE ADMIN"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
