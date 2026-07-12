package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type DepartmentHandler struct {
	departments *repository.DepartmentRepository
}

func NewDepartmentHandler(departments *repository.DepartmentRepository) *DepartmentHandler {
	return &DepartmentHandler{departments}
}

func (h *DepartmentHandler) List(c echo.Context) error {
	status := c.QueryParam("status")
	departments, err := h.departments.List(c.Request().Context(), status)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	return util.Success(c, http.StatusOK, departments)
}

type createDepartmentRequest struct {
	Name               string `json:"name" validate:"required,min=2,max=100"`
	ParentDepartmentID *int64 `json:"parent_department_id"`
	HeadID             *int64 `json:"head_id"`
}

func (h *DepartmentHandler) Create(c echo.Context) error {
	var req createDepartmentRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	d := &model.Department{
		Name:               req.Name,
		ParentDepartmentID: req.ParentDepartmentID,
		HeadID:             req.HeadID,
		Status:             model.StatusActive,
	}

	if err := h.departments.Create(c.Request().Context(), d); err != nil {
		if errors.Is(err, repository.ErrDepartmentExists) {
			return util.Fail(c, http.StatusConflict, "department name already exists")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, d)
}

type updateDepartmentRequest struct {
	Name               *string `json:"name"`
	HeadID             *int64  `json:"head_id"`
	ParentDepartmentID *int64  `json:"parent_department_id"`
	Status             *string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}

func (h *DepartmentHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid department id")
	}

	var req updateDepartmentRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	d, err := h.departments.Update(c.Request().Context(), id, req.Name, req.HeadID, req.ParentDepartmentID, req.Status)
	if err != nil {
		if errors.Is(err, repository.ErrDepartmentNotFound) {
			return util.Fail(c, http.StatusNotFound, "department not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, d)
}