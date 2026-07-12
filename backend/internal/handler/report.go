package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/util"
)


type ReportHandler struct{}

func NewReportHandler() *ReportHandler { return &ReportHandler{} }

func (h *ReportHandler) Get(c echo.Context) error {
	mockReports := echo.Map{
		"utilization": []echo.Map{
			{"label": "Eng", "value": 40},
			{"label": "HR", "value": 65},
			{"label": "IT", "value": 85},
			{"label": "Ops", "value": 55},
			{"label": "Sales", "value": 35},
			{"label": "Mktg", "value": 75},
		},
		"mostUsed": []echo.Map{
			{"asset": "Conference Room B2", "stat": "34 bookings this month"},
			{"asset": "Company Van", "stat": "21 trips this month"},
			{"asset": "MacBook Pro 16", "stat": "18 allocations"},
		},
		"idle": []echo.Map{
			{"asset": "Office Desk", "stat": "unused 60+ days"},
			{"asset": "Dell Monitor 27", "stat": "unused 45 days"},
		},
		"actionNeeded": []echo.Map{
			{"asset": "Company Van", "stat": "service due in 5 days"},
			{"asset": "Dell Monitor 27", "stat": "3 years old : nearing retirement"},
		},
	}
	return util.Success(c, http.StatusOK, mockReports)
}