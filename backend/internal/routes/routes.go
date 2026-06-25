package routes

import (
	"github.com/labstack/echo/v4"
)

// Routes only map URLs to handlers. Middleware and service wiring stay in the API composition layer.
type PublicRouteHandlers struct {
	Health echo.HandlerFunc
	Login  echo.HandlerFunc
}

type ProtectedRouteHandlers struct {
	Me                   echo.HandlerFunc
	Logout               echo.HandlerFunc
	ListVisibleAudit     echo.HandlerFunc
	DashboardStats       echo.HandlerFunc
	ListUsers            echo.HandlerFunc
	CreateUser           echo.HandlerFunc
	UpdateUser           echo.HandlerFunc
	SetUserStatus        echo.HandlerFunc
	ListPermissions      echo.HandlerFunc
	CreatePermission     echo.HandlerFunc
	ListRoles            echo.HandlerFunc
	CreateRole           echo.HandlerFunc
	UpdateRole           echo.HandlerFunc
	ListSystemAudit      echo.HandlerFunc
	ListSubmissions      echo.HandlerFunc
	CreateSubmission     echo.HandlerFunc
	GetSubmission        echo.HandlerFunc
	UpdateSubmission     echo.HandlerFunc
	TransitionSubmission echo.HandlerFunc
	SubmissionAudit      echo.HandlerFunc
}

func RegisterPublic(router *echo.Echo, handlers PublicRouteHandlers) {
	registerHealthRoutes(router, handlers)
	registerPublicAuthRoutes(router, handlers)
}

func RegisterProtected(router *echo.Group, handlers ProtectedRouteHandlers) {
	registerProtectedAuthRoutes(router, handlers)
	registerAuditRoutes(router, handlers)
	registerDashboardRoutes(router, handlers)
	registerAdminRoutes(router, handlers)
	registerSubmissionRoutes(router, handlers)
}
