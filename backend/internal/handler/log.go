package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type LogHandler struct {
	logs *repository.ActivityLogRepository
}

func NewLogHandler(logs *repository.ActivityLogRepository) *LogHandler {
	return &LogHandler{logs}
}

func (h *LogHandler) List(c echo.Context) error {
	page, limit := util.PageLimit(c)

	var userID *int64
	if v := c.QueryParam("user_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return util.Fail(c, http.StatusBadRequest, "invalid user_id")
		}
		userID = &id
	}

	logs, total, err := h.logs.List(c.Request().Context(), repository.ActivityLogFilter{
		UserID:     userID,
		EntityType: c.QueryParam("entity_type"),
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	return util.Success(c, http.StatusOK, model.PaginatedResponse[model.ActivityLog]{
		Data: logs, Page: page, Limit: limit, Total: total, TotalPages: util.TotalPages(total, limit),
	})
}