package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/falasefemi2/hms/internal/dto"
	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
)

type DepartmentService struct {
	repo *repository.DepartmentRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		repo: repo,
	}
}

func (ds *DepartmentService) CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	if err := ds.validateCreateRequest(req); err != nil {
		return nil, err
	}

	dept := CreateRequestToModel(req)
	created, err := ds.repo.CreateDepartment(ctx, dept)
	if err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return ModelToDepartmentResponse(created), nil
}

func (ds *DepartmentService) GetDepartmentByID(ctx context.Context, deptID string) (*dto.DepartmentResponse, error) {
	dept, err := ds.repo.GetByID(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	if !dept.IsActive {
		return nil, errors.New("department is inactive")
	}

	return ModelToDepartmentResponse(dept), nil
}

func (ds *DepartmentService) GetAllDepartments(ctx context.Context, req *dto.PaginationRequest) (*dto.DepartmentListResponse, error) {
	if err := ds.validatePagination(req); err != nil {
		return nil, err
	}

	params := PaginationToRepositoryParams(req)
	resutl, err := ds.repo.GetAll(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments: %w", err)
	}
	return ModelsToListResponse(resutl.Data, resutl.TotalCount, req.Page, req.PageSize), nil
}

func (ds *DepartmentService) UpdateDepartment(ctx context.Context, deptID string, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	if err := ds.validateUpdateRequest(req); err != nil {
		return nil, err
	}
	exisitng, err := ds.repo.GetByID(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	if !exisitng.IsActive {
		return nil, errors.New("cannot update an inactive department")
	}

	repoReq := UpdateRequestToRepositoryRequest(req)
	updated, err := ds.repo.UpdateDepartment(ctx, deptID, repoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return ModelToDepartmentResponse(updated), nil
}

func (ds *DepartmentService) DeleteDepartment(ctx context.Context, deptID string) error {
	existing, err := ds.repo.GetByID(ctx, deptID)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}
	if !existing.IsActive {
		return errors.New("department is already deleted")
	}

	err = ds.repo.DeleteDepartment(ctx, deptID)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return nil
}

func ModelToDepartmentResponse(dept *models.Department) *dto.DepartmentResponse {
	if dept == nil {
		return nil
	}

	return &dto.DepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Description: dept.Description,
		IsActive:    dept.IsActive,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
	}
}

func ModelsToListResponse(departments []*models.Department, totalCount int, page int, pageSize int) *dto.DepartmentListResponse {
	responses := make([]dto.DepartmentResponse, len(departments))

	for i, dept := range departments {
		responses[i] = *ModelToDepartmentResponse(dept)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.DepartmentListResponse{
		Data:       responses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func CreateRequestToModel(req *dto.CreateDepartmentRequest) *models.Department {
	if req == nil {
		return nil
	}

	return &models.Department{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true, // New departments are active by default
	}
}

func UpdateRequestToRepositoryRequest(req *dto.UpdateDepartmentRequest) *repository.UpdateDepartmentRequest {
	if req == nil {
		return nil
	}

	return &repository.UpdateDepartmentRequest{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
}

func PaginationToRepositoryParams(req *dto.PaginationRequest) repository.PaginationParams {
	offset := (req.Page - 1) * req.PageSize

	return repository.PaginationParams{
		Limit:  req.PageSize,
		Offset: offset,
	}
}

func (s *DepartmentService) validateCreateRequest(req *dto.CreateDepartmentRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	// Validate name
	if err := s.validateName(req.Name); err != nil {
		return err
	}

	// Validate description
	if len(req.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}

	return nil
}

func (s *DepartmentService) validateUpdateRequest(req *dto.UpdateDepartmentRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	// At least one field must be provided
	if req.Name == nil && req.Description == nil && req.IsActive == nil {
		return errors.New("at least one field must be provided for update")
	}

	// Validate name if provided
	if req.Name != nil {
		if err := s.validateName(*req.Name); err != nil {
			return err
		}
	}

	// Validate description if provided
	if req.Description != nil {
		if len(*req.Description) > 500 {
			return errors.New("description cannot exceed 500 characters")
		}
	}

	return nil
}

func (s *DepartmentService) validateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("name cannot be empty")
	}

	if len(name) > 255 {
		return errors.New("name cannot exceed 255 characters")
	}

	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	return nil
}

func (s *DepartmentService) validatePagination(req *dto.PaginationRequest) error {
	if req == nil {
		return errors.New("pagination request cannot be nil")
	}

	if req.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if req.PageSize < 1 {
		return errors.New("page_size must be greater than 0")
	}

	if req.PageSize > 100 {
		return errors.New("page_size cannot exceed 100")
	}

	return nil
}
