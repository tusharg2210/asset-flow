package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type DashboardHandler struct {
	assets       *repository.AssetRepository
	allocations  *repository.AllocationRepository
	bookings     *repository.BookingRepository
	transfers    *repository.TransferRepository
	maintenance  *repository.MaintenanceRepository
}

func NewDashboardHandler(
	assets *repository.AssetRepository,
	allocations *repository.AllocationRepository,
	bookings *repository.BookingRepository,
	transfers *repository.TransferRepository,
	maintenance *repository.MaintenanceRepository,
) *DashboardHandler {
	return &DashboardHandler{assets, allocations, bookings, transfers, maintenance}
}


func (h *DashboardHandler) Metrics(c echo.Context) error {
	ctx := c.Request().Context()

	_, availableTotal, err := h.assets.List(ctx, repository.AssetFilter{Status: "AVAILABLE", Limit: 1})
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	_, allocatedTotal, err := h.assets.List(ctx, repository.AssetFilter{Status: "ALLOCATED", Limit: 1})
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	maintenanceToday, err := h.maintenance.CountDueToday(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	activeBookings, err := h.bookings.CountActive(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	pendingTransfers, err := h.transfers.CountPending(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	overdue, err := h.allocations.ListOverdue(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, echo.Map{
		"assetsAvailable":  availableTotal,
		"assetsAllocated":  allocatedTotal,
		"maintenanceToday": maintenanceToday,
		"activeBookings":   activeBookings,
		"pendingTransfers": pendingTransfers,
		"upcomingReturns":  len(overdue), // TODO: split into "upcoming" (not yet due) vs "overdue" — see Alerts below
	})
}

func (h *DashboardHandler) Alerts(c echo.Context) error {
	ctx := c.Request().Context()

	overdue, err := h.allocations.ListOverdue(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	pendingTransfers, err := h.transfers.CountPending(ctx)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, echo.Map{
		"overdueAllocations": overdue, // service layer should enrich with asset_tag / held_by / days_overdue before this reaches JSON
		"pendingTransfers":   pendingTransfers,
	})
}