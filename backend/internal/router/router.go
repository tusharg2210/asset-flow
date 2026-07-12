package router

import (
	"github.com/labstack/echo/v4"

	"asset-flow/internal/config"
	"asset-flow/internal/handler"
	"asset-flow/internal/middleware"
	"asset-flow/internal/model"
)

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

func Register(e *echo.Echo, cfg *config.Config, h *handler.Handlers) {
	secret := cfg.JWT.JWTSecret

	api := e.Group("/api")

	authPublic := api.Group("/auth")
	authPublic.POST("/signup", h.Auth.Signup)
	authPublic.POST("/login", h.Auth.Login)
	authPublic.POST("/refresh", h.Auth.Refresh)

	authed := api.Group("", middleware.RequireAuth(secret))

	authed.GET("/auth/me", h.Auth.Me)
	authed.POST("/auth/logout", h.Auth.Logout)

	authed.GET("/dashboard/metrics", h.Dashboard.Metrics)
	authed.GET("/dashboard/alerts", h.Dashboard.Alerts)

	authed.GET("/assets", h.Asset.List)
	authed.GET("/assets/:id", h.Asset.GetByID)
	authed.GET("/assets/:id/history", h.Asset.History)
	authed.GET("/assets/:id/bookings", h.Booking.ListForAsset)

	authed.POST("/bookings", h.Booking.Create)
	authed.PUT("/bookings/:id/status", h.Booking.UpdateStatus)

	authed.GET("/maintenance", h.Maintenance.List)
	authed.POST("/maintenance", h.Maintenance.Create)

	authed.POST("/transfers", h.Allocation.CreateTransfer)

	authed.GET("/audits/:id/assets", h.Audit.ScopedAssets)
	authed.POST("/audits/:id/reports", h.Audit.SubmitReport)

	assetMgr := api.Group("", middleware.RequireAuth(secret),
		middleware.RequireRoles(model.RoleAssetManager, model.RoleAdmin))

	assetMgr.POST("/assets", h.Asset.Create)
	assetMgr.POST("/allocations", h.Allocation.Allocate)
	assetMgr.PUT("/allocations/:id/return", h.Allocation.Return)
	assetMgr.POST("/maintenance/:id/workflow", h.Maintenance.AdvanceWorkflow)

	approvers := api.Group("", middleware.RequireAuth(secret),
		middleware.RequireRoles(model.RoleAssetManager, model.RoleDepartmentHead, model.RoleAdmin))

	approvers.PUT("/transfers/:id/status", h.Allocation.UpdateTransferStatus)

	managers := api.Group("", middleware.RequireAuth(secret),
		middleware.RequireRoles(model.RoleAssetManager, model.RoleDepartmentHead, model.RoleAdmin))

	managers.GET("/reports", h.Report.Get)

	admin := api.Group("", middleware.RequireAuth(secret), middleware.RequireRoles(model.RoleAdmin))

	admin.GET("/departments", h.Department.List)
	admin.POST("/departments", h.Department.Create)
	admin.PUT("/departments/:id", h.Department.Update)

	admin.GET("/categories", h.Category.List) 
	admin.POST("/categories", h.Category.Create)
	admin.PUT("/categories/:id", h.Category.Update)

	admin.GET("/users", h.User.List)
	admin.PUT("/users/:id/role", h.User.UpdateRole)

	admin.POST("/audits", h.Audit.Create)
	admin.PUT("/audits/:id/close", h.Audit.Close)

	admin.GET("/logs", h.Log.List)
}