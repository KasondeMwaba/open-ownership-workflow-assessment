package routes

import "github.com/labstack/echo/v4"

func registerDashboardRoutes(router *echo.Group, handlers ProtectedRouteHandlers) {
	router.GET("/api/dashboard", handlers.DashboardStats)
}
