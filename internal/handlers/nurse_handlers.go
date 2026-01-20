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

type NurseHandler struct {
	nurseService *service.NurseSerivce
}

func NewNurseHandler(nurseService *service.NurseSerivce) *NurseHandler {
	return &NurseHandler{
		nurseService: nurseService,
	}
}

// CreateNurse godoc
// @Summary Create a new nurse
// @Description Create a new nurse. Requires valid JWT token with ADMIN role
// @Tags Nurse Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.NurseSignupRequest true "Nurse creation details"
// @Success 201 {object} dto.NurseResponse "Nurse created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error - invalid input or invalid role"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} dto.ErrorResponse "Forbidden - admin role required"
// @Failure 409 {object} dto.ErrorResponse "Conflict - nurse already exists for this user"
// @Router /admin/nurses [post]
func (n *NurseHandler) CreateNurse(w http.ResponseWriter, r *http.Request) {
	var req dto.NurseSignupRequest
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

	nurse := &models.Nurse{
		NurseID:       uuid.New(),
		UserID:        userID,
		DepartmentID:  departmentID,
		Shift:         strings.TrimSpace(req.Shift),
		LicenseNumber: strings.TrimSpace(req.LicenseNumber),
	}

	createdNurse, err := n.nurseService.CreateNurse(r.Context(), nurse)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "nurse already exists") {
			utils.WriteError(w, http.StatusConflict, errorMsg)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, errorMsg)
		return
	}

	response := &dto.NurseResponse{
		NurseID:       createdNurse.NurseID.String(),
		UserID:        createdNurse.NurseID.String(),
		DepartmentID:  createdNurse.DepartmentID.String(),
		Shift:         createdNurse.Shift,
		LicenseNumber: createdNurse.LicenseNumber,
		CreatedAt:     createdNurse.CreatedAt,
		UpdatedAt:     createdNurse.UpdatedAt,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}
