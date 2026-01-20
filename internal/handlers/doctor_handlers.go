package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type DoctorHandler struct {
	doctorService *service.DoctorService
}

func NewDoctorHandler(doctorService *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
	}
}

// CreateDoctor godoc
// @Summary Create a new doctor
// @Description Create a new doctor. Requires valid JWT token with ADMIN role
// @Tags Doctor Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DoctorSignUpRequest true "Doctor creation details"
// @Success 201 {object} dto.DoctorResponse "Doctor created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error - invalid input or invalid role"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Failure 409 {object} dto.ErrorResponse "Conflict - doctor already exists for this user"
// @Router /admin/doctors [post]
func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	var req dto.DoctorSignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	doctor := &models.Doctor{
		DocotrID:        uuid.New(), // Generate UUID here
		UserID:          userID,
		Specialization:  strings.TrimSpace(req.Specialization),
		LicenseNumber:   strings.TrimSpace(req.LicenseNumber),
		DepartmentID:    departmentID,
		ConsultationFee: req.ConsultationFee,
	}

	createdDoctor, err := h.doctorService.CreateDoctor(r.Context(), doctor)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "doctor already exists") {
			utils.WriteError(w, http.StatusConflict, errorMsg)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, errorMsg)
		return
	}

	response := &dto.DoctorResponse{
		DoctorID:        createdDoctor.DocotrID.String(),
		UserID:          createdDoctor.UserID.String(),
		Specialization:  createdDoctor.Specialization,
		LicenseNumber:   createdDoctor.LicenseNumber,
		DepartmentID:    createdDoctor.DepartmentID.String(),
		ConsultationFee: createdDoctor.ConsultationFee,
		IsAvailable:     createdDoctor.IsAvailable,
		CreatedAt:       createdDoctor.CreatedAt,
		UpdatedAt:       createdDoctor.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}
