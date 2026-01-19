package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/service"
	"github.com/falasefemi2/hms/internal/utils"
)

type DepartmentHandler struct {
	deptService *service.DepartmentService
}

func NewDeptHandler(deptService *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		deptService: deptService,
	}
}

// CreateDepartment godoc
// @Summary      Create a new department
// @Description  Creates a new department with the provided name and description
// @Tags         departments
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        request  body      dto.CreateDepartmentRequest  true  "Department details"
// @Success      201      {object}  dto.DepartmentResponse       "Department created successfully"
// @Failure      400      {object}  map[string]string            "Invalid request or validation error"
// @Failure      500      {object}  map[string]string            "Internal server error"
// @Security     Bearer
// @Router       /admin/departments [post]
func (dh *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	createdDept, err := dh.deptService.CreateDepartment(r.Context(), &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdDept)
}

// GetDepartment godoc
// @Summary      Get a department by ID
// @Description  Retrieves a single department by its ID
// @Tags         departments
// @Produce      json
// @Security BearerAuth
// @Param        id       path      string                   true  "Department ID (UUID)"
// @Success      200      {object}  dto.DepartmentResponse   "Department retrieved successfully"
// @Failure      400      {object}  map[string]string        "Invalid department ID"
// @Failure      404      {object}  map[string]string        "Department not found"
// @Failure      500      {object}  map[string]string        "Internal server error"
// @Security     Bearer
// @Router       /admin/departments/{id} [get]
func (dh *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	deptID := chi.URLParam(r, "id")

	if deptID == "" {
		utils.WriteError(w, http.StatusBadRequest, "department id required")
		return
	}

	dept, err := dh.deptService.GetDepartmentByID(r.Context(), deptID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dept)
}

// GetAllDepartments godoc
// @Summary      List all departments
// @Description  Retrieves all active departments with pagination support
// @Tags         departments
// @Produce      json
// @Security BearerAuth
// @Param        page       query     int                          false  "Page number (default 1)"                         default(1)
// @Param        page_size  query     int                          false  "Number of items per page (default 10, max 100)"  default(10)
// @Success      200        {object}  dto.DepartmentListResponse   "Departments retrieved successfully"
// @Failure      400        {object}  map[string]string            "Invalid pagination parameters"
// @Failure      500        {object}  map[string]string            "Internal server error"
// @Security     Bearer
// @Router      /admin/departments [get]
func (dh *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	pagination := &dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	// Parse page parameter
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			utils.WriteError(w, http.StatusBadRequest, "page must be a positive integer")
			return
		}
		pagination.Page = page
	}

	// Parse page_size parameter
	pageSizeStr := r.URL.Query().Get("page_size")
	if pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			utils.WriteError(w, http.StatusBadRequest, "page_size must be a positive integer")
			return
		}
		if pageSize > 100 {
			pageSize = 100 // Cap max page size
		}
		pagination.PageSize = pageSize
	}

	result, err := dh.deptService.GetAllDepartments(r.Context(), pagination)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

// UpdateDepartment godoc
// @Summary      Update a department
// @Description  Updates one or more fields of an existing department. Only provided fields will be updated.
// @Tags         departments
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        id       path      string                       true  "Department ID (UUID)"
// @Param        request  body      dto.UpdateDepartmentRequest  true  "Fields to update (all optional)"
// @Success      200      {object}  dto.DepartmentResponse       "Department updated successfully"
// @Failure      400      {object}  map[string]string            "Invalid request or validation error"
// @Failure      404      {object}  map[string]string            "Department not found"
// @Failure      500      {object}  map[string]string            "Internal server error"
// @Security     Bearer
// @Router       /admin/departments/{id} [put]
func (dh *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	deptID := chi.URLParam(r, "id")

	if deptID == "" {
		utils.WriteError(w, http.StatusBadRequest, "department id required")
		return
	}

	var req dto.UpdateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Trim whitespace if fields are provided
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		req.Description = &trimmed
	}

	updated, err := dh.deptService.UpdateDepartment(r.Context(), deptID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, updated)
}

// DeleteDepartment godoc
// @Summary      Delete a department
// @Description  Soft deletes a department by setting is_active to false
// @Tags         departments
// @Produce      json
// @Security BearerAuth
// @Param        id  path      string             true  "Department ID (UUID)"
// @Success      200 {object}  map[string]string  "Department deleted successfully"
// @Failure      400 {object}  map[string]string  "Invalid request or department already deleted"
// @Failure      404 {object}  map[string]string  "Department not found"
// @Failure      500 {object}  map[string]string  "Internal server error"
// @Security     Bearer
// @Router       /admin/departments/{id} [delete]
func (dh *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	deptID := chi.URLParam(r, "id")

	if deptID == "" {
		utils.WriteError(w, http.StatusBadRequest, "department id required")
		return
	}

	err := dh.deptService.DeleteDepartment(r.Context(), deptID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "department deleted successfully",
	})
}
