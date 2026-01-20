package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type PatientHandlers struct {
	patientService *service.PatientService
}

func NewPatientHandlers(patientService *service.PatientService) *PatientHandlers {
	return &PatientHandlers{
		patientService: patientService,
	}
}

// PatientProfile creates a patient profile.
// @Summary Create patient profile
// @Description Create patient profile. Requires valid JWT token with PATIENT role
// @Tags Patient Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.PatientSignUp true "Patient Profile details"
// @Success 201 {object} dto.PatientResponse "Patient Profile created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error - invalid input or invalid role"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - patient role required"
// @Failure 409 {object} dto.ErrorResponse "Conflict - patient already exists for this user"
// @Router /patients/patientprofile [post]
func (p *PatientHandlers) PatientProfile(w http.ResponseWriter, r *http.Request) {
	var req dto.PatientSignUp
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Clean the date string by removing ordinal suffixes
	dateStr := strings.Replace(req.DateOfBirth, "th", "", 1)
	dateStr = strings.Replace(dateStr, "st", "", 1)
	dateStr = strings.Replace(dateStr, "nd", "", 1)
	dateStr = strings.Replace(dateStr, "rd", "", 1)

	// Parse the date string
	dob, err := time.Parse("2 January 2006", dateStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid date of birth format")
		return
	}

	patient := &models.Patient{
		PatientID:             uuid.New(),
		UserID:                userID,
		DateOfBirth:           dob,
		Gender:                req.Gender,
		BloodGroup:            req.BloodGroup,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
		MedicalHistory:        req.MedicalHistory,
	}

	patientProfile, err := p.patientService.PatientProfile(r.Context(), patient)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "patient profile already exists") {
			utils.WriteError(w, http.StatusConflict, errorMsg)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, errorMsg)
		return
	}

	response := &dto.PatientResponse{
		PatientID:             patientProfile.PatientID,
		UserID:                patientProfile.UserID,
		DateOfBirth:           patientProfile.DateOfBirth,
		Gender:                patientProfile.Gender,
		BloodGroup:            patientProfile.BloodGroup,
		EmergencyContactName:  patientProfile.EmergencyContactName,
		EmergencyContactPhone: patientProfile.EmergencyContactPhone,
		MedicalHistory:        patientProfile.MedicalHistory,
		CreatedAt:             patientProfile.CreatedAt,
		UpdatedAt:             patientProfile.UpdatedAt,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}
