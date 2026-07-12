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

type AssetHandler struct {
	assets       *repository.AssetRepository
	allocations  *repository.AllocationRepository
	maintenance  *repository.MaintenanceRepository
	categories   *repository.AssetCategoryRepository
}

func NewAssetHandler(
	assets *repository.AssetRepository,
	allocations *repository.AllocationRepository,
	maintenance *repository.MaintenanceRepository,
	categories *repository.AssetCategoryRepository,
) *AssetHandler {
	return &AssetHandler{assets, allocations, maintenance, categories}
}

func (h *AssetHandler) List(c echo.Context) error {
	page, limit := util.PageLimit(c)

	var categoryID *int64
	if v := c.QueryParam("category_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return util.Fail(c, http.StatusBadRequest, "invalid category_id")
		}
		categoryID = &id
	}

	f := repository.AssetFilter{
		Tag:          c.QueryParam("tag"),
		SerialNumber: c.QueryParam("serial_number"),
		QRCode:       c.QueryParam("qr_code"),
		CategoryID:   categoryID,
		Status:       c.QueryParam("status"),
		Location:     c.QueryParam("location"),
		Page:         page,
		Limit:        limit,
	}

	assets, total, err := h.assets.List(c.Request().Context(), f)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, model.PaginatedResponse[model.Asset]{
		Data: assets, Page: page, Limit: limit, Total: total, TotalPages: util.TotalPages(total, limit),
	})
}

type createAssetRequest struct {
	Name              string   `json:"name" validate:"required,min=2,max=200"`
	CategoryID        *int64   `json:"category_id"`
	SerialNumber      *string  `json:"serial_number"`
	Location          string   `json:"location" validate:"required"`
	Condition         string   `json:"condition" validate:"required"`
	PhotosDocs        []string `json:"photos_docs"`
	IsSharable        bool     `json:"is_sharable"`
	IsBookable        bool     `json:"is_bookable"`
	AcquisitionDate   string   `json:"acquisition_date"` 
	AcquisitionCost   float64  `json:"acquisition_cost" validate:"gte=0"`
}

func (h *AssetHandler) Create(c echo.Context) error {
	var req createAssetRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	tag, err := h.generateNextTag(c)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	a := &model.Asset{
		Tag:             tag,
		Name:            req.Name,
		CategoryID:      req.CategoryID,
		SerialNumber:    req.SerialNumber,
		Status:          model.AssetAvailable,
		Location:        req.Location,
		Condition:       req.Condition,
		PhotosDocs:      req.PhotosDocs,
		IsSharable:      req.IsSharable,
		IsBookable:      req.IsBookable,
		AcquisitionCost: req.AcquisitionCost,
	}

	if err := h.assets.Create(c.Request().Context(), a); err != nil {
		if errors.Is(err, repository.ErrAssetTagExists) {
			return util.Fail(c, http.StatusConflict, "asset tag collision, please retry")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, a)
}

func (h *AssetHandler) generateNextTag(c echo.Context) (string, error) {
	_, total, err := h.assets.List(c.Request().Context(), repository.AssetFilter{Limit: 1})
	if err != nil {
		return "", err
	}
	return "AF-" + padTag(total+1), nil
}

func padTag(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 4 {
		s = "0" + s
	}
	return s
}

func (h *AssetHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid asset id")
	}

	a, err := h.assets.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return util.Fail(c, http.StatusNotFound, "asset not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, a)
}

func (h *AssetHandler) History(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid asset id")
	}
	ctx := c.Request().Context()

	allocations, err := h.allocations.ListByAsset(ctx, id)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	maintenanceRecords, err := h.maintenance.ListByAsset(ctx, id)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, echo.Map{
		"allocations":         allocations,
		"maintenance_records": maintenanceRecords,
	})
}