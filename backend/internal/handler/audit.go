package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type AuditHandler struct {
	db     *pgxpool.Pool
	audits *repository.AuditRepository
	assets *repository.AssetRepository
}

func NewAuditHandler(db *pgxpool.Pool, audits *repository.AuditRepository, assets *repository.AssetRepository) *AuditHandler {
	return &AuditHandler{db, audits, assets}
}

type createAuditRequest struct {
	Scope             string   `json:"scope" validate:"required"`
	ScopeDepartmentID *int64   `json:"scope_department_id"`
	ScopeLocation     string   `json:"scope_location"`
	FromDate          time.Time `json:"from_date" validate:"required"`
	ToDate            time.Time `json:"to_date" validate:"required,gtfield=FromDate"`
	Auditors          []int64  `json:"auditors" validate:"required,min=1"`
}

func (h *AuditHandler) Create(c echo.Context) error {
	var req createAuditRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	a := &model.Audit{
		Auditors:          req.Auditors,
		Scope:             req.Scope,
		ScopeDepartmentID: req.ScopeDepartmentID,
		ScopeLocation:     req.ScopeLocation,
		FromDate:          req.FromDate,
		ToDate:            req.ToDate,
	}

	if err := h.audits.Create(c.Request().Context(), a); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, a)
}

func (h *AuditHandler) ScopedAssets(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid audit id")
	}

	ctx := c.Request().Context()
	audit, err := h.audits.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuditNotFound) {
			return util.Fail(c, http.StatusNotFound, "audit not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	assets, err := h.audits.ScopedAssetIDs(ctx, audit)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, assets)
}

type submitReportRequest struct {
	AssetID  int64  `json:"asset_id" validate:"required"`
	Verified bool   `json:"verified"`
	Status   string `json:"status" validate:"required,oneof=VERIFIED MISSING DAMAGED"`
	Remarks  string `json:"remarks"`
}

func (h *AuditHandler) SubmitReport(c echo.Context) error {
	auditID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid audit id")
	}

	var req submitReportRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	verifierID, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	audit, err := h.audits.GetByID(c.Request().Context(), auditID)
	if err != nil {
		if errors.Is(err, repository.ErrAuditNotFound) {
			return util.Fail(c, http.StatusNotFound, "audit not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	if audit.Status != model.AuditOpen {
		return util.Fail(c, http.StatusConflict, "audit cycle is not open")
	}

	rep := &model.AuditReport{
		AuditID:            auditID,
		AssetID:             req.AssetID,
		VerificationStatus:  req.Status,
		Remarks:             req.Remarks,
		VerifiedBy:           &verifierID,
	}

	if err := h.audits.CreateReport(c.Request().Context(), rep); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, rep)
}

func (h *AuditHandler) Close(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid audit id")
	}

	ctx := c.Request().Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	missingIDs, err := h.audits.ListFlaggedAssetIDs(ctx, tx, id, model.VerificationMissing)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	damagedIDs, err := h.audits.ListFlaggedAssetIDs(ctx, tx, id, model.VerificationDamaged)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	for _, assetID := range missingIDs {
		if err := h.assets.UpdateStatus(ctx, tx, assetID, model.AssetLost); err != nil {
			return util.FailErr(c, http.StatusInternalServerError, err)
		}
	}

	audit, err := h.audits.Close(ctx, tx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuditNotOpen) {
			return util.Fail(c, http.StatusConflict, "audit is not open")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, echo.Map{
		"audit_id":             audit.ID,
		"status":               audit.Status,
		"discrepancies_found":  len(missingIDs) + len(damagedIDs),
	})
}