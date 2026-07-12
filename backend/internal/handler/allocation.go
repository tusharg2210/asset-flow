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

type AllocationHandler struct {
	db          *pgxpool.Pool
	assets      *repository.AssetRepository
	allocations *repository.AllocationRepository
	transfers   *repository.TransferRepository
}

func NewAllocationHandler(
	db *pgxpool.Pool,
	assets *repository.AssetRepository,
	allocations *repository.AllocationRepository,
	transfers *repository.TransferRepository,
) *AllocationHandler {
	return &AllocationHandler{db, assets, allocations, transfers}
}

type allocateRequest struct {
	AssetID            int64      `json:"asset_id" validate:"required"`
	ToUserID           int64      `json:"to_user_id" validate:"required"`
	ExpectedReturnDate *time.Time `json:"expected_return_date"`
	Reason             string     `json:"reason"`
}

func (h *AllocationHandler) Allocate(c echo.Context) error {
	var req allocateRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	asset, err := h.assets.GetByIDForUpdate(ctx, tx, req.AssetID)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			return util.Fail(c, http.StatusNotFound, "asset not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if asset.Status != model.AssetAvailable {
		active, aErr := h.allocations.GetActiveByAssetID(ctx, tx, asset.ID)
		if aErr == nil {
			return c.JSON(http.StatusConflict, echo.Map{
				"success":          false,
				"error":            "asset is already allocated",
				"currently_held_by": active.ToUserID,
				"suggest_transfer":  true,
			})
		}
		return util.Fail(c, http.StatusConflict, "asset is not available for allocation")
	}

	alloc := &model.Allocation{
		AssetID:            req.AssetID,
		ToUserID:            req.ToUserID,
		ExpectedReturnDate:  req.ExpectedReturnDate,
		Reason:              req.Reason,
	}
	if err := h.allocations.Create(ctx, tx, alloc); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if err := h.assets.UpdateStatus(ctx, tx, asset.ID, model.AssetAllocated); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.SuccessMsg(c, http.StatusCreated, "Asset allocated successfully.", echo.Map{
		"allocation_id": alloc.ID,
		"asset_id":      alloc.AssetID,
		"status":        "Success",
	})
}

type returnRequest struct {
	ConditionNotes string `json:"condition_notes"`
}

func (h *AllocationHandler) Return(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid allocation id")
	}

	var req returnRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}

	ctx := c.Request().Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	alloc, err := h.allocations.MarkReturned(ctx, tx, id, req.ConditionNotes)
	if err != nil {
		if errors.Is(err, repository.ErrAllocationNotFound) {
			return util.Fail(c, http.StatusNotFound, "allocation not found")
		}
		if errors.Is(err, repository.ErrAllocationReturned) {
			return util.Fail(c, http.StatusConflict, "allocation already returned")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if err := h.assets.UpdateStatus(ctx, tx, alloc.AssetID, model.AssetAvailable); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if req.ConditionNotes != "" {
		if _, err := tx.Exec(ctx,
			`UPDATE assets SET condition = $2, updated_at = now() WHERE id = $1`,
			alloc.AssetID, req.ConditionNotes,
		); err != nil {
			return util.FailErr(c, http.StatusInternalServerError, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, echo.Map{
		"allocation_id":     alloc.ID,
		"actual_return_date": alloc.ActualReturnDate,
		"asset_status":      model.AssetAvailable,
	})
}

type createTransferRequest struct {
	AssetID    int64  `json:"asset_id" validate:"required"`
	FromUserID int64  `json:"from_user_id" validate:"required"`
	ToUserID   int64  `json:"to_user_id" validate:"required"`
	Reason     string `json:"reason"`
}

func (h *AllocationHandler) CreateTransfer(c echo.Context) error {
	var req createTransferRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	requestedBy, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	t := &model.TransferRequest{
		AssetID:     req.AssetID,
		FromUserID:  req.FromUserID,
		ToUserID:    req.ToUserID,
		RequestedBy: requestedBy,
		Reason:      req.Reason,
	}

	if err := h.transfers.Create(c.Request().Context(), t); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, echo.Map{
		"transfer_id": t.ID,
		"status":      t.Status,
	})
}

type updateTransferStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=APPROVED REJECTED"`
}

func (h *AllocationHandler) UpdateTransferStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid transfer id")
	}

	var req updateTransferStatusRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	approverID, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	ctx := c.Request().Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	t, err := h.transfers.GetByIDForUpdate(ctx, tx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTransferNotFound) {
			return util.Fail(c, http.StatusNotFound, "transfer not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	updated, err := h.transfers.UpdateStatus(ctx, tx, id, req.Status, approverID)
	if err != nil {
		if errors.Is(err, repository.ErrTransferNotRequested) {
			return util.Fail(c, http.StatusConflict, "transfer is not pending")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	var newAllocationID *int64
	if req.Status == model.TransferApproved {
		oldAlloc, err := h.allocations.GetActiveByAssetID(ctx, tx, t.AssetID)
		if err == nil {
			if _, rErr := h.allocations.MarkReturned(ctx, tx, oldAlloc.ID, "Transferred"); rErr != nil {
				return util.FailErr(c, http.StatusInternalServerError, rErr)
			}
		}

		newAlloc := &model.Allocation{
			AssetID:    t.AssetID,
			FromUserID: &t.FromUserID,
			ToUserID:   t.ToUserID,
			Reason:     "Transfer #" + strconv.FormatInt(t.ID, 10),
		}
		if err := h.allocations.Create(ctx, tx, newAlloc); err != nil {
			return util.FailErr(c, http.StatusInternalServerError, err)
		}
		newAllocationID = &newAlloc.ID

		if err := h.assets.UpdateStatus(ctx, tx, t.AssetID, model.AssetAllocated); err != nil {
			return util.FailErr(c, http.StatusInternalServerError, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	resp := echo.Map{"transfer_id": updated.ID, "status": updated.Status}
	if newAllocationID != nil {
		resp["new_allocation_id"] = *newAllocationID
	}
	return util.Success(c, http.StatusOK, resp)
}