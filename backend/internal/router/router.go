package router

import (
	"github.com/labstack/echo/v4"

	"asset-flow/internal/config"
	"asset-flow/internal/handler"
	"asset-flow/internal/middleware"
	"asset-flow/internal/model"
)

// Handlers bundles every handler the router needs. Built once in cmd/main.go
// and passed in — keeps router.go free of repository/db wiring.
type Handlers struct {
	Auth        *handler.AuthHandler
	Dashboard   *handler.DashboardHandler
	Department  *handler.DepartmentHandler
	Category    *handler.CategoryHandler
	User        *handler.UserHandler
	Asset       *handler.AssetHandler
	Allocation  *handler.AllocationHandler
	Booking     *handler.BookingHandler
	Maintenance *handler.MaintenanceHandler
	Audit       *handler.AuditHandler
	Report      *handler.ReportHandler
	Log         *handler.LogHandler
}

// Register mounts every route under /api. CORS, request logging, recover,
// etc. are assumed to already be attached to `e` by internal/server before
// this is called.
func Register(e *echo.Echo, cfg *config.Config, h *handler.Handlers) {
	secret := cfg.JWT.JWTSecret

	api := e.Group("/api")

	// ---------------------------------------------------------------
	// Public — no token required
	// ---------------------------------------------------------------
	authPublic := api.Group("/auth")
	authPublic.POST("/signup", h.Auth.Signup)
	authPublic.POST("/login", h.Auth.Login)
	authPublic.POST("/refresh", h.Auth.Refresh)

	// ---------------------------------------------------------------
	// Authenticated — any logged-in role
	// ---------------------------------------------------------------
	authed := api.Group("", middleware.RequireAuth(secret))

	authed.GET("/auth/me", h.Auth.Me)
	authed.POST("/auth/logout", h.Auth.Logout)

	authed.GET("/dashboard/metrics", h.Dashboard.Metrics)
	authed.GET("/dashboard/alerts", h.Dashboard.Alerts)

	// Asset directory: browsing is open to everyone once authenticated;
	// mutation is role-gated further down.
	authed.GET("/assets", h.Asset.List)
	authed.GET("/assets/:id", h.Asset.GetByID)
	authed.GET("/assets/:id/history", h.Asset.History)
	authed.GET("/assets/:id/bookings", h.Booking.ListForAsset)

	// Bookings: Employee/Dept Head/Asset Manager can all book shared resources.
	authed.POST("/bookings", h.Booking.Create)
	authed.PUT("/bookings/:id/status", h.Booking.UpdateStatus)

	// Maintenance: any holder can raise a request; approval is gated below.
	authed.GET("/maintenance", h.Maintenance.List)
	authed.POST("/maintenance", h.Maintenance.Create)

	// Transfers: the current holder or a requester initiates; approval is
	// gated below. "Initiates return/transfer requests" is explicitly an
	// Employee capability per the spec.
	authed.POST("/transfers", h.Allocation.CreateTransfer)

	// Audits: any authenticated user can view scoped assets and submit a
	// report — narrowing this to "assigned auditors only" needs a DB check
	// against audits.auditors, which the handler doesn't currently enforce.
	// Flagging as a gap rather than silently under- or over-restricting.
	authed.GET("/audits/:id/assets", h.Audit.ScopedAssets)
	authed.POST("/audits/:id/reports", h.Audit.SubmitReport)

	// ---------------------------------------------------------------
	// Asset Manager (+ Admin as a superset)
	// ---------------------------------------------------------------
	assetMgr := api.Group("", middleware.RequireAuth(secret),
		middleware.RequireRoles(model.RoleAssetManager, model.RoleAdmin))

	assetMgr.POST("/assets", h.Asset.Create)
	assetMgr.POST("/allocations", h.Allocation.Allocate)
	assetMgr.PUT("/allocations/:id/return", h.Allocation.Return) // "Approves asset returns and condition check-in notes"
	assetMgr.POST("/maintenance/:id/workflow", h.Maintenance.AdvanceWorkflow)

	// ---------------------------------------------------------------
	// Transfer approval: Asset Manager or Department Head (+ Admin)
	// ---------------------------------------------------------------
	approvers := api.Group("", middleware.RequireAuth(secret),
		middleware.RequireRoles(model.RoleAssetManager, model.RoleDepartmentHead, model.RoleAdmin))

	approvers.PUT("/transfers/:id/status", h.Allocation.UpdateTransferStatus)

	// ---------------------------------------------------------------
	// Managers (Asset Manager, Department Head) + Admin — analytics
	// ---------------------------------------------------------------
	managers := api.Group("", middleware.RequireAuth(secret),
		middleware.RequireRoles(model.RoleAssetManager, model.RoleDepartmentHead, model.RoleAdmin))

	managers.GET("/reports", h.Report.Get)

	// ---------------------------------------------------------------
	// Admin only — org setup, role assignment, audit lifecycle, logs
	// ---------------------------------------------------------------
	admin := api.Group("", middleware.RequireAuth(secret), middleware.RequireRoles(model.RoleAdmin))

	admin.GET("/departments", h.Department.List)
	admin.POST("/departments", h.Department.Create)
	admin.PUT("/departments/:id", h.Department.Update)

	admin.GET("/categories", h.Category.List) // not in the API doc; supports Screen 3 Tab B
	admin.POST("/categories", h.Category.Create)
	admin.PUT("/categories/:id", h.Category.Update)

	admin.GET("/users", h.User.List)
	admin.PUT("/users/:id/role", h.User.UpdateRole)

	admin.POST("/audits", h.Audit.Create)
	admin.PUT("/audits/:id/close", h.Audit.Close)

	admin.GET("/logs", h.Log.List)
}