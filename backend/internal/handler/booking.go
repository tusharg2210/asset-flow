package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type BookingHandler struct {
	db       *pgxpool.Pool
	bookings *repository.BookingRepository
}

func NewBookingHandler(db *pgxpool.Pool, bookings *repository.BookingRepository) *BookingHandler {
	return &BookingHandler{db, bookings}
}

func (h *BookingHandler) ListForAsset(c echo.Context) error {
	assetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid asset id")
	}

	slots, err := h.bookings.ListSlotsByAsset(c.Request().Context(), assetID)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	return util.Success(c, http.StatusOK, slots)
}

type createBookingRequest struct {
	AssetID   int64     `json:"asset_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}

func (h *BookingHandler) Create(c echo.Context) error {
	var req createBookingRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	bookedBy, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	ctx := c.Request().Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, req.AssetID); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	overlaps, err := h.bookings.HasOverlap(ctx, tx, req.AssetID, req.StartTime, req.EndTime)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	if overlaps {
		return util.Fail(c, http.StatusConflict, "requested time slot overlaps an existing booking")
	}

	booking := &model.ResourceBooking{AssetID: req.AssetID, BookedBy: bookedBy}
	slot := &model.BookingSlot{StartTime: req.StartTime, EndTime: req.EndTime}

	if err := h.bookings.Create(ctx, tx, booking, slot); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusCreated, echo.Map{
		"booking_id":     booking.ID,
		"asset_id":       booking.AssetID,
		"booking_status": booking.BookingStatus,
		"slots": []echo.Map{
			{"start_time": slot.StartTime, "end_time": slot.EndTime},
		},
	})
}

type updateBookingStatusRequest struct {
	BookingStatus string `json:"booking_status" validate:"required,oneof=UPCOMING ONGOING COMPLETED CANCELLED"`
}

func (h *BookingHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid booking id")
	}

	var req updateBookingStatusRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	b, err := h.bookings.UpdateStatus(c.Request().Context(), id, req.BookingStatus)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, echo.Map{
		"booking_id":     b.ID,
		"booking_status": b.BookingStatus,
	})
}