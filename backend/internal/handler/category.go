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

type CategoryHandler struct {
	categories *repository.AssetCategoryRepository
}

func NewCategoryHandler(categories *repository.AssetCategoryRepository) *CategoryHandler {
	return &CategoryHandler{categories}
}

func (h *CategoryHandler) List(c echo.Context) error {
	categories, err := h.categories.List(c.Request().Context())
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	return util.Success(c, http.StatusOK, categories)
}

type createCategoryRequest struct {
	Name               string `json:"name" validate:"required,min=2,max=100"`
	CustomFieldsSchema []byte `json:"custom_fields_schema"`
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var req createCategoryRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	cat := &model.AssetCategory{Name: req.Name, CustomFieldsSchema: req.CustomFieldsSchema}
	if err := h.categories.Create(c.Request().Context(), cat); err != nil {
		if errors.Is(err, repository.ErrCategoryExists) {
			return util.Fail(c, http.StatusConflict, "category name already exists")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	return util.Success(c, http.StatusCreated, cat)
}

type updateCategoryRequest struct {
	Name               *string `json:"name"`
	CustomFieldsSchema []byte  `json:"custom_fields_schema"`
}

func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid category id")
	}

	var req updateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}

	cat, err := h.categories.Update(c.Request().Context(), id, req.Name, req.CustomFieldsSchema)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return util.Fail(c, http.StatusNotFound, "category not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	return util.Success(c, http.StatusOK, cat)
}