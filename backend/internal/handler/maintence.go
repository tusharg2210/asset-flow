package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type MaintenanceHandler struct {
	db          *pgxpool.Pool
	maintenance *repository.MaintenanceRepository
	assets      *repository.AssetRepository
}

func NewMaintenanceHandler(db *pgxpool.Pool, maintenance *repository.MaintenanceRepository, assets *repository.AssetRepository) *MaintenanceHandler {
	return &MaintenanceHandler{db, maintenance, assets}
}

func (h *MaintenanceHandler) List(c echo.Context) error {
	page, limit := util.PageLimit(c)

	var assetID *int64
	if v := c.QueryParam("asset_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return util.Fail(c, http.StatusBadRequest, "invalid asset_id")
		}
		assetID = &id
	}

	f := repository.MaintenanceFilter{
		Status:  c.QueryParam("status"),
		AssetID: assetID,
		Page:    page,
		Limit:   limit,
	}

	records, total, err := h.maintenance.List(c.Request().Context(), f)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, model.PaginatedResponse[model.Maintenance]{
		Data: records, Page: page, Limit: limit, Total: total, TotalPages: util.TotalPages(total, limit),
	})
}

type createMaintenanceRequest struct {
	AssetID     int64    `json:"asset_id" validate:"required"`
	Priority    string   `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH CRITICAL"`
	Description string   `json:"description" validate:"required"`
	Images      []string `json:"images"`
}

func (h *MaintenanceHandler) Create(c echo.Context) error {
	var req createMaintenanceRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	raisedBy, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	m := &model.Maintenance{
		AssetID:     req.AssetID,
		RaisedBy:    &raisedBy,
		Priority:    req.Priority,
		Description: req.Description,
		Images:      req.Images,
	}

	if err := h.maintenance.Create(c.Request().Context(), m); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, m)
}

type maintenanceWorkflowRequest struct {
	Status       string `json:"status" validate:"required,oneof=APPROVED REJECTED TECHNICIAN_ASSIGNED IN_PROGRESS RESOLVED"`
	Description  string `json:"description"`
	TechnicianID *int64 `json:"technician_id"`
}

func (h *MaintenanceHandler) AdvanceWorkflow(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid maintenance id")
	}

	var req maintenanceWorkflowRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	updatedBy, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	ctx := c.Request().Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	m, err := h.maintenance.UpdateStatus(ctx, tx, id, req.Status, req.TechnicianID)
	if err != nil {
		if errors.Is(err, repository.ErrMaintenanceNotFound) {
			return util.Fail(c, http.StatusNotFound, "maintenance request not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if err := h.maintenance.CreateWorkflowStep(ctx, tx, &model.MaintenanceWorkflow{
		MaintenanceID: id,
		Status:        req.Status,
		Description:   req.Description,
		UpdatedBy:     &updatedBy,
	}); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	var assetUpdatedTo string
	switch req.Status {
	case model.MaintenanceApproved:
		assetUpdatedTo = model.AssetUnderMaintenance
	case model.MaintenanceResolved:
		assetUpdatedTo = model.AssetAvailable
	}
	if assetUpdatedTo != "" {
		if err := h.assets.UpdateStatus(ctx, tx, m.AssetID, assetUpdatedTo); err != nil {
			return util.FailErr(c, http.StatusInternalServerError, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	resp := echo.Map{"maintenance_id": m.ID, "new_status": m.Status}
	if assetUpdatedTo != "" {
		resp["asset_updated_to"] = assetUpdatedTo
	}
	return util.Success(c, http.StatusOK, resp)
}