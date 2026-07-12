package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/config"
	"asset-flow/internal/repository"
)

type Handlers struct {
	Auth        *AuthHandler
	Dashboard   *DashboardHandler
	Department  *DepartmentHandler
	Category    *CategoryHandler
	User        *UserHandler
	Asset       *AssetHandler
	Allocation  *AllocationHandler
	Booking     *BookingHandler
	Maintenance *MaintenanceHandler
	Audit       *AuditHandler
	Report      *ReportHandler
	Log         *LogHandler
}

func New(repos *repository.Repositories, pool *pgxpool.Pool, cfg *config.Config) *Handlers {
	return &Handlers{
		Auth:      NewAuthHandler(repos.User, cfg.JWT),
		Dashboard: NewDashboardHandler(repos.Asset, repos.Allocation, repos.Booking, repos.Transfer, repos.Maintenance),
		Department: NewDepartmentHandler(repos.Department),
		Category:   NewCategoryHandler(repos.Category),
		User:       NewUserHandler(repos.User),
		Asset:      NewAssetHandler(repos.Asset, repos.Allocation, repos.Maintenance, repos.Category),
		Allocation: NewAllocationHandler(pool, repos.Asset, repos.Allocation, repos.Transfer),
		Booking:    NewBookingHandler(pool, repos.Booking),
		Maintenance: NewMaintenanceHandler(pool, repos.Maintenance, repos.Asset),
		Audit:      NewAuditHandler(pool, repos.Audit, repos.Asset),
		Report:     NewReportHandler(),
		Log:        NewLogHandler(repos.ActivityLog),
	}
}