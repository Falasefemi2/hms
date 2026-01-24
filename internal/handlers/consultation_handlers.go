package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type ConsultationHandler struct {
	consultationService *service.ConsultationService
}

func NewConsultationHandler(consultationService *service.ConsultationService) *ConsultationHandler {
	return &ConsultationHandler{
		consultationService: consultationService,
	}
}

// CreateConsultation godoc
// @Summary Create a new consultation
// @Description Create a new consultation for a completed appointment
// @Tags Consultation Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateConsultationRequest true "Consultation creation details"
// @Success 201 {object} dto.ConsultationResponse "Consultation created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Router /consultations [post]
func (h *ConsultationHandler) CreateConsultation(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateConsultationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appointmentID, err := uuid.Parse(req.AppointmentID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid appointment id")
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

	consultation := &models.Consultation{
		ConsultationID: uuid.New(),
		AppointmentID:  appointmentID,
		PatientID:      patientID,
		DoctorID:       doctorID,
		Diagnosis:      req.Diagnosis,
		Notes:          req.Notes,
		IsEditable:     true,
	}

	createdConsultation, err := h.consultationService.CreateConsultation(r.Context(), consultation)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.ConsultationResponse{
		ConsultationID: createdConsultation.ConsultationID.String(),
		AppointmentID:  createdConsultation.AppointmentID.String(),
		PatientID:      createdConsultation.PatientID.String(),
		DoctorID:       createdConsultation.DoctorID.String(),
		Diagnosis:      createdConsultation.Diagnosis,
		Notes:          createdConsultation.Notes,
		CreatedAt:      createdConsultation.CreatedAt,
		IsEditable:     createdConsultation.IsEditable,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

// GetConsultation godoc
// @Summary Get consultation by ID
// @Description Get consultation details by ID
// @Tags Consultation Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Consultation ID"
// @Success 200 {object} dto.ConsultationResponse "Consultation details"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID"
// @Failure 404 {object} dto.ErrorResponse "Consultation not found"
// @Router /consultations/{id} [get]
func (h *ConsultationHandler) GetConsultation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	consultationID, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid consultation id")
		return
	}

	consultation, err := h.consultationService.GetConsultationByID(r.Context(), consultationID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "consultation not found")
		return
	}

	response := &dto.ConsultationResponse{
		ConsultationID: consultation.ConsultationID.String(),
		AppointmentID:  consultation.AppointmentID.String(),
		PatientID:      consultation.PatientID.String(),
		DoctorID:       consultation.DoctorID.String(),
		Diagnosis:      consultation.Diagnosis,
		Notes:          consultation.Notes,
		CreatedAt:      consultation.CreatedAt,
		IsEditable:     consultation.IsEditable,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// UpdateConsultation godoc
// @Summary Update a consultation
// @Description Update consultation details
// @Tags Consultation Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Consultation ID"
// @Param request body dto.UpdateConsultationRequest true "Consultation update details"
// @Success 200 {object} dto.ConsultationResponse "Consultation updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 404 {object} dto.ErrorResponse "Consultation not found"
// @Router /consultations/{id} [put]
func (h *ConsultationHandler) UpdateConsultation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	consultationID, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid consultation id")
		return
	}

	var req dto.UpdateConsultationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	consultation := &models.Consultation{
		ConsultationID: consultationID,
		Diagnosis:      req.Diagnosis,
		Notes:          req.Notes,
	}

	updatedConsultation, err := h.consultationService.UpdateConsultation(r.Context(), consultation)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.ConsultationResponse{
		ConsultationID: updatedConsultation.ConsultationID.String(),
		AppointmentID:  updatedConsultation.AppointmentID.String(),
		PatientID:      updatedConsultation.PatientID.String(),
		DoctorID:       updatedConsultation.DoctorID.String(),
		Diagnosis:      updatedConsultation.Diagnosis,
		Notes:          updatedConsultation.Notes,
		CreatedAt:      updatedConsultation.CreatedAt,
		IsEditable:     updatedConsultation.IsEditable,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
