package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// CreateAppointment godoc
// @Summary Create a new appointment
// @Description Create a new appointment. Requires valid JWT token
// @Tags Appointment Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAppointmentRequest true "Appointment creation details"
// @Success 201 {object} dto.AppointmentResponse "Appointment created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Router /appointments [post]
func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAppointmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid patient id")
		return
	}

	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid doctor id")
		return
	}

	appointmentDate, err := time.Parse(time.RFC3339, req.AppointmentDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid appointment date format")
		return
	}

	appointment := &models.Appointment{
		AppointmentID:   uuid.New(),
		PatientID:       patientID,
		DoctorID:        doctorID,
		AppointmentDate: appointmentDate,
		DurationMinutes: req.DurationMinutes,
		Status:          "PENDING",
		Notes:           req.Notes,
	}

	createdAppointment, err := h.appointmentService.CreateAppointment(r.Context(), appointment)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.AppointmentResponse{
		AppointmentID:   createdAppointment.AppointmentID.String(),
		PatientID:       createdAppointment.PatientID.String(),
		DoctorID:        createdAppointment.DoctorID.String(),
		AppointmentDate: createdAppointment.AppointmentDate,
		DurationMinutes: createdAppointment.DurationMinutes,
		Status:          createdAppointment.Status,
		Notes:           createdAppointment.Notes,
		CreatedAt:       createdAppointment.CreatedAt,
		UpdatedAt:       createdAppointment.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

// GetAppointment godoc
// @Summary Get appointment by ID
// @Description Get appointment details by ID
// @Tags Appointment Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse "Appointment details"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID"
// @Failure 404 {object} dto.ErrorResponse "Appointment not found"
// @Router /appointments/{id} [get]
func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appointmentID, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid appointment id")
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByID(r.Context(), appointmentID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "appointment not found")
		return
	}

	response := &dto.AppointmentResponse{
		AppointmentID:   appointment.AppointmentID.String(),
		PatientID:       appointment.PatientID.String(),
		DoctorID:        appointment.DoctorID.String(),
		AppointmentDate: appointment.AppointmentDate,
		DurationMinutes: appointment.DurationMinutes,
		Status:          appointment.Status,
		Notes:           appointment.Notes,
		CreatedAt:       appointment.CreatedAt,
		UpdatedAt:       appointment.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// UpdateAppointment godoc
// @Summary Update an appointment
// @Description Update appointment details
// @Tags Appointment Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Param request body dto.UpdateAppointmentRequest true "Appointment update details"
// @Success 200 {object} dto.AppointmentResponse "Appointment updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 404 {object} dto.ErrorResponse "Appointment not found"
// @Router /appointments/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	appointmentID, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid appointment id")
		return
	}

	var req dto.UpdateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appointmentDate, err := time.Parse(time.RFC3339, req.AppointmentDate)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid appointment date format")
		return
	}

	appointment := &models.Appointment{
		AppointmentID:   appointmentID,
		AppointmentDate: appointmentDate,
		DurationMinutes: req.DurationMinutes,
		Status:          req.Status,
		Notes:           req.Notes,
	}

	updatedAppointment, err := h.appointmentService.UpdateAppointment(r.Context(), appointment)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.AppointmentResponse{
		AppointmentID:   updatedAppointment.AppointmentID.String(),
		PatientID:       updatedAppointment.PatientID.String(),
		DoctorID:        updatedAppointment.DoctorID.String(),
		AppointmentDate: updatedAppointment.AppointmentDate,
		DurationMinutes: updatedAppointment.DurationMinutes,
		Status:          updatedAppointment.Status,
		Notes:           updatedAppointment.Notes,
		CreatedAt:       updatedAppointment.CreatedAt,
		UpdatedAt:       updatedAppointment.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
