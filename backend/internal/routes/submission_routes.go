package routes

import "github.com/labstack/echo/v4"

func registerSubmissionRoutes(router *echo.Group, handlers ProtectedRouteHandlers) {
	router.GET("/api/submissions", handlers.ListSubmissions)
	router.POST("/api/submissions", handlers.CreateSubmission)
	router.GET("/api/submissions/:id", handlers.GetSubmission)
	router.PUT("/api/submissions/:id", handlers.UpdateSubmission)
	router.POST("/api/submissions/:id/transition", handlers.TransitionSubmission)
	router.GET("/api/submissions/:id/audit", handlers.SubmissionAudit)
}
