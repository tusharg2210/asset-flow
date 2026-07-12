package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/util"
)


type ReportHandler struct{}

func NewReportHandler() *ReportHandler { return &ReportHandler{} }

func (h *ReportHandler) Get(c echo.Context) error {
	return util.Fail(c, http.StatusNotImplemented, "reports aggregation not yet implemented")
}