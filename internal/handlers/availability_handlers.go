package handlers

import (
	"encoding/json"

	"github.com/google/uuid"

	"net/http"
	"strings"
	"time"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type AvailabilityHandlers struct {
	availabilityService *service.AvailabilityService
}

func NewAvailabilityHandlers(availabilityService *service.AvailabilityService) *AvailabilityHandlers {
	return &AvailabilityHandlers{
		availabilityService: availabilityService,
	}
}

var validDays = map[string]bool{
	"Monday":    true,
	"Tuesday":   true,
	"Wednesday": true,
	"Thursday":  true,
	"Friday":    true,
	"Saturday":  true,
	"Sunday":    true,
}

func (a *AvailabilityHandlers) validateTimeFormat(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	return err
}

func (a *AvailabilityHandlers) validateDayOfWeek(day string) bool {
	return validDays[day]
}

// CreateAvailability creates doctor availability
// @Summary Create doctor availability
// @Description Create doctor availability slot. Requires valid JWT token with ADMIN role
// @Tags Doctor Availability
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AvailabilityRequest true "Availability details"
// @Success 201 {object} dto.AvailabilityResponse "Availability created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Router /admin/doctors/availability [post]
func (a *AvailabilityHandlers) CreateAvailability(w http.ResponseWriter, r *http.Request) {
	var req dto.AvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !a.validateDayOfWeek(req.DayOfWeek) {
		utils.WriteError(w, http.StatusBadRequest, "invalid day of week. use: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday")
		return
	}

	if err := a.validateTimeFormat(req.StartTime); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid start time format. use HH:MM (24-hour format)")
		return
	}

	if err := a.validateTimeFormat(req.EndTime); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid end time format. use HH:MM (24-hour format)")
		return
	}

	startTime, _ := time.Parse("15:04", req.StartTime)
	endTime, _ := time.Parse("15:04", req.EndTime)
	if startTime.After(endTime) || startTime.Equal(endTime) {
		utils.WriteError(w, http.StatusBadRequest, "start time must be before end time")
		return
	}

	if req.MaxAppointments <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "max appointments must be greater than 0")
		return
	}

	availability := &models.Availability{
		AvailabilityID: uuid.New(),
		DoctorID:       req.DoctorID,
		DayOfWeek:      req.DayOfWeek,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		MaxAppointment: req.MaxAppointments,
	}

	createdAvailability, err := a.availabilityService.CreateDoctorAvailability(r.Context(), availability)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "not found") {
			utils.WriteError(w, http.StatusNotFound, errorMsg)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, errorMsg)
		return
	}

	response := &dto.AvailabilityResponse{
		AvailabilityID:  createdAvailability.AvailabilityID,
		DoctorID:        createdAvailability.DoctorID,
		DayOfWeek:       createdAvailability.DayOfWeek,
		StartTime:       createdAvailability.StartTime,
		EndTime:         createdAvailability.EndTime,
		MaxAppointments: createdAvailability.MaxAppointment,
		CreatedAt:       createdAvailability.CreatedAt,
		UpdatedAt:       createdAvailability.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}
