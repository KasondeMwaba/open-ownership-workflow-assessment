package routes

import "github.com/labstack/echo/v4"

func registerAuditRoutes(router *echo.Group, handlers ProtectedRouteHandlers) {
	router.GET("/api/audit", handlers.ListVisibleAudit)
	router.GET("/api/admin/audit", handlers.ListSystemAudit)
}
