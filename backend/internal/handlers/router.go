package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"openownership-workflow/backend/internal/config"
	"openownership-workflow/backend/internal/routes"
	"openownership-workflow/backend/internal/services"
)

type API struct {
	cfg         config.Config
	auth        *services.AuthService
	users       *services.UserService
	access      *services.AccessService
	audit       *services.AuditService
	submissions *services.SubmissionService
	dashboard   *services.DashboardService
	logger      *slog.Logger
}

type Dependencies struct {
	Config      config.Config
	Auth        *services.AuthService
	Users       *services.UserService
	Access      *services.AccessService
	Audit       *services.AuditService
	Submissions *services.SubmissionService
	Dashboard   *services.DashboardService
	Logger      *slog.Logger
}

func NewRouter(deps Dependencies) http.Handler {
	api := API{
		cfg:         deps.Config,
		auth:        deps.Auth,
		users:       deps.Users,
		access:      deps.Access,
		audit:       deps.Audit,
		submissions: deps.Submissions,
		dashboard:   deps.Dashboard,
		logger:      deps.Logger,
	}

	router := echo.New()
	router.HideBanner = true
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{deps.Config.CORSOrigin},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Use(api.requestLogger)

	routes.RegisterPublic(router, routes.PublicRouteHandlers{
		Health: api.health,
		Login:  api.login,
	})

	protected := router.Group("")
	protected.Use(api.requireAuth)
	protected.Use(api.auditActivity)
	routes.RegisterProtected(protected, routes.ProtectedRouteHandlers{
		Me:                   api.me,
		Logout:               api.logout,
		ListVisibleAudit:     api.listVisibleAudit,
		DashboardStats:       api.dashboardStats,
		ListUsers:            api.listUsers,
		CreateUser:           api.createUser,
		UpdateUser:           api.updateUser,
		SetUserStatus:        api.setUserStatus,
		ListPermissions:      api.listPermissions,
		CreatePermission:     api.createPermission,
		ListRoles:            api.listRoles,
		CreateRole:           api.createRole,
		UpdateRole:           api.updateRole,
		ListSystemAudit:      api.listSystemAudit,
		ListSubmissions:      api.listSubmissions,
		CreateSubmission:     api.createSubmission,
		GetSubmission:        api.getSubmission,
		UpdateSubmission:     api.updateSubmission,
		TransitionSubmission: api.transitionSubmission,
		SubmissionAudit:      api.auditEvents,
	})

	return router
}

func (api API) health(c echo.Context) error {
	return writeJSON(c, http.StatusOK, map[string]string{"status": "ok"})
}
