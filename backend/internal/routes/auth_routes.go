package routes

import "github.com/labstack/echo/v4"

func registerPublicAuthRoutes(router *echo.Echo, handlers PublicRouteHandlers) {
	router.POST("/api/auth/login", handlers.Login)
}

func registerProtectedAuthRoutes(router *echo.Group, handlers ProtectedRouteHandlers) {
	router.GET("/api/me", handlers.Me)
	router.POST("/api/auth/logout", handlers.Logout)
}
