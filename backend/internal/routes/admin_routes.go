package routes

import "github.com/labstack/echo/v4"

func registerAdminRoutes(router *echo.Group, handlers ProtectedRouteHandlers) {
	router.GET("/api/admin/users", handlers.ListUsers)
	router.POST("/api/admin/users", handlers.CreateUser)
	router.PUT("/api/admin/users/:id", handlers.UpdateUser)
	router.POST("/api/admin/users/:id/status", handlers.SetUserStatus)

	router.GET("/api/admin/permissions", handlers.ListPermissions)
	router.POST("/api/admin/permissions", handlers.CreatePermission)

	router.GET("/api/admin/roles", handlers.ListRoles)
	router.POST("/api/admin/roles", handlers.CreateRole)
	router.PUT("/api/admin/roles/:id", handlers.UpdateRole)
}
