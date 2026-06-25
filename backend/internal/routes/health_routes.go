package routes

import "github.com/labstack/echo/v4"

func registerHealthRoutes(router *echo.Echo, handlers PublicRouteHandlers) {
	router.GET("/healthz", handlers.Health)
}
