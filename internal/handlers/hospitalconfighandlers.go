package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type HospitalConfigHandler struct {
	hospitalConfigService *service.HospitalConfigService
}

func NewHospitalConfigHandler(hospitalConfigService *service.HospitalConfigService) *HospitalConfigHandler {
	return &HospitalConfigHandler{
		hospitalConfigService: hospitalConfigService,
	}
}

// CreateHospitalConfig godoc
// @Summary Create a new hospital configuration
// @Description Create a new hospital configuration. Requires valid JWT token with ADMIN role
// @Tags Hospital Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateHospitalConfigRequest true "Hospital configuration details"
// @Success 201 {object} dto.HospitalConfigResponse "Hospital configuration created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error - invalid input"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Router /admin/hospital-configs [post]
func (h *HospitalConfigHandler) CreateHospitalConfig(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateHospitalConfigRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	config := &models.HospitalConfig{
		ConfigID:                      uuid.New(),
		WorkingHoursStart:             strings.TrimSpace(req.WorkingHoursStart),
		WorkingHoursEnd:               strings.TrimSpace(req.WorkingHoursEnd),
		AppointmentDurationMinutes:    req.AppointmentDurationMinutes,
		MaxSameDayCancellationHours:   req.MaxSameDayCancellationHours,
		EnablePatientSelfRegistration: true, // default
	}

	if req.EnablePatientSelfRegistration != nil {
		config.EnablePatientSelfRegistration = *req.EnablePatientSelfRegistration
	}

	createdConfig, err := h.hospitalConfigService.CreateHospitalConfig(r.Context(), config)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.HospitalConfigResponse{
		ConfigID:                      createdConfig.ConfigID.String(),
		WorkingHoursStart:             createdConfig.WorkingHoursStart,
		WorkingHoursEnd:               createdConfig.WorkingHoursEnd,
		AppointmentDurationMinutes:    createdConfig.AppointmentDurationMinutes,
		MaxSameDayCancellationHours:   createdConfig.MaxSameDayCancellationHours,
		EnablePatientSelfRegistration: createdConfig.EnablePatientSelfRegistration,
		CreatedAt:                     createdConfig.CreatedAt,
		UpdatedAt:                     createdConfig.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

// GetHospitalConfig godoc
// @Summary Get a hospital configuration by ID
// @Description Get a specific hospital configuration. Requires valid JWT token with ADMIN role
// @Tags Hospital Configuration
// @Produce json
// @Security BearerAuth
// @Param id path string true "Hospital configuration ID"
// @Success 200 {object} dto.HospitalConfigResponse "Hospital configuration retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID format"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Failure 404 {object} dto.ErrorResponse "Hospital configuration not found"
// @Router /admin/hospital-configs/{id} [get]
func (h *HospitalConfigHandler) GetHospitalConfig(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	configID, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid config id")
		return
	}

	config, err := h.hospitalConfigService.GetHospitalConfigByID(r.Context(), configID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.WriteError(w, http.StatusNotFound, "hospital configuration not found")
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.HospitalConfigResponse{
		ConfigID:                      config.ConfigID.String(),
		WorkingHoursStart:             config.WorkingHoursStart,
		WorkingHoursEnd:               config.WorkingHoursEnd,
		AppointmentDurationMinutes:    config.AppointmentDurationMinutes,
		MaxSameDayCancellationHours:   config.MaxSameDayCancellationHours,
		EnablePatientSelfRegistration: config.EnablePatientSelfRegistration,
		CreatedAt:                     config.CreatedAt,
		UpdatedAt:                     config.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetAllHospitalConfigs godoc
// @Summary Get all hospital configurations
// @Description Get all hospital configurations. Requires valid JWT token with ADMIN role
// @Tags Hospital Configuration
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.HospitalConfigResponse "Hospital configurations retrieved successfully"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Router /admin/hospital-configs [get]
func (h *HospitalConfigHandler) GetAllHospitalConfigs(w http.ResponseWriter, r *http.Request) {
	configs, err := h.hospitalConfigService.GetAllHospitalConfigs(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var responses []dto.HospitalConfigResponse
	for _, config := range configs {
		responses = append(responses, dto.HospitalConfigResponse{
			ConfigID:                      config.ConfigID.String(),
			WorkingHoursStart:             config.WorkingHoursStart,
			WorkingHoursEnd:               config.WorkingHoursEnd,
			AppointmentDurationMinutes:    config.AppointmentDurationMinutes,
			MaxSameDayCancellationHours:   config.MaxSameDayCancellationHours,
			EnablePatientSelfRegistration: config.EnablePatientSelfRegistration,
			CreatedAt:                     config.CreatedAt,
			UpdatedAt:                     config.UpdatedAt,
		})
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

// UpdateHospitalConfig godoc
// @Summary Update a hospital configuration
// @Description Update an existing hospital configuration. Requires valid JWT token with ADMIN role
// @Tags Hospital Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Hospital configuration ID"
// @Param request body dto.UpdateHospitalConfigRequest true "Updated hospital configuration details"
// @Success 200 {object} dto.HospitalConfigResponse "Hospital configuration updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error - invalid input or ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Failure 404 {object} dto.ErrorResponse "Hospital configuration not found"
// @Router /admin/hospital-configs/{id} [put]
func (h *HospitalConfigHandler) UpdateHospitalConfig(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	configID, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid config id")
		return
	}

	var req dto.UpdateHospitalConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	config := &models.HospitalConfig{
		ConfigID:                      configID,
		WorkingHoursStart:             strings.TrimSpace(req.WorkingHoursStart),
		WorkingHoursEnd:               strings.TrimSpace(req.WorkingHoursEnd),
		AppointmentDurationMinutes:    req.AppointmentDurationMinutes,
		MaxSameDayCancellationHours:   req.MaxSameDayCancellationHours,
		EnablePatientSelfRegistration: true, // default
	}

	if req.EnablePatientSelfRegistration != nil {
		config.EnablePatientSelfRegistration = *req.EnablePatientSelfRegistration
	}

	updatedConfig, err := h.hospitalConfigService.UpdateHospitalConfig(r.Context(), config)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &dto.HospitalConfigResponse{
		ConfigID:                      updatedConfig.ConfigID.String(),
		WorkingHoursStart:             updatedConfig.WorkingHoursStart,
		WorkingHoursEnd:               updatedConfig.WorkingHoursEnd,
		AppointmentDurationMinutes:    updatedConfig.AppointmentDurationMinutes,
		MaxSameDayCancellationHours:   updatedConfig.MaxSameDayCancellationHours,
		EnablePatientSelfRegistration: updatedConfig.EnablePatientSelfRegistration,
		CreatedAt:                     updatedConfig.CreatedAt,
		UpdatedAt:                     updatedConfig.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// DeleteHospitalConfig godoc
// @Summary Delete a hospital configuration
// @Description Delete an existing hospital configuration. Requires valid JWT token with ADMIN role
// @Tags Hospital Configuration
// @Produce json
// @Security BearerAuth
// @Param id path string true "Hospital configuration ID"
// @Success 204 "Hospital configuration deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID format"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Failure 404 {object} dto.ErrorResponse "Hospital configuration not found"
// @Router /admin/hospital-configs/{id} [delete]
func (h *HospitalConfigHandler) DeleteHospitalConfig(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	configID, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid config id")
		return
	}

	err = h.hospitalConfigService.DeleteHospitalConfig(r.Context(), configID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
