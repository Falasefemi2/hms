package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) SignUpPatient(w http.ResponseWriter, r *http.Request) {
	var req dto.PatientSignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid bad request")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)

	createdUser, err := u.userService.CreatePatientUser(
		r.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.Phone,
	)

	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "email already registered") || strings.Contains(errorMsg, "username already taken") {
			utils.WriteError(w, http.StatusConflict, errorMsg)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, errorMsg)
		return
	}

	response := u.userToResponse(createdUser)

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.AdminCreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Role = strings.TrimSpace(req.Role)

	createdUser, err := u.userService.CreateAdminUser(
		r.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Role,
	)

	if err != nil {
		// Check for specific validation errors
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "email already registered") ||
			strings.Contains(errorMsg, "username already taken") {
			utils.WriteError(w, http.StatusConflict, errorMsg)
			return
		}
		if strings.Contains(errorMsg, "invalid role") ||
			strings.Contains(errorMsg, "must self-register") {
			utils.WriteError(w, http.StatusBadRequest, errorMsg)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, errorMsg)
		return
	}
	response := u.userToResponse(createdUser)

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "userID required")
		return
	}

	user, err := u.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	response := u.userToResponse(user)
	utils.WriteJSON(w, http.StatusOK, response)
}

func (u *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		// TODO LATER
	}
	users, err := u.userService.ListUsers(r.Context(), limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to retrieve users")
		return
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *u.userToResponse(user)
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

func (u *UserHandler) userToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}
