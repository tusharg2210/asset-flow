package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type NotificationHandler struct {
	repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) List(c echo.Context) error {
	userID, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	page, limit := util.PageLimit(c)
	unreadOnly := c.QueryParam("unread") == "true"

	notifications, total, err := h.repo.ListByUser(c.Request().Context(), userID, unreadOnly, page, limit)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, map[string]any{
		"data":  notifications,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
	userID, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid notification id")
	}

	err = h.repo.MarkRead(c.Request().Context(), id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			return util.Fail(c, http.StatusNotFound, "notification not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.SuccessMsg(c, http.StatusOK, "notification marked as read", nil)
}
