package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"openownership-workflow/backend/internal/config"
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

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{deps.Config.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Use(api.requestLogger)

	router.Get("/healthz", api.health)
	router.Post("/api/auth/login", api.login)

	router.Group(func(protected chi.Router) {
		protected.Use(api.requireAuth)
		protected.Use(api.auditActivity)
		protected.Get("/api/me", api.me)
		protected.Post("/api/auth/logout", api.logout)
		protected.Get("/api/audit", api.listVisibleAudit)
		protected.Get("/api/dashboard", api.dashboardStats)
		protected.Get("/api/admin/users", api.listUsers)
		protected.Post("/api/admin/users", api.createUser)
		protected.Put("/api/admin/users/{id}", api.updateUser)
		protected.Post("/api/admin/users/{id}/status", api.setUserStatus)
		protected.Get("/api/admin/permissions", api.listPermissions)
		protected.Post("/api/admin/permissions", api.createPermission)
		protected.Get("/api/admin/roles", api.listRoles)
		protected.Post("/api/admin/roles", api.createRole)
		protected.Put("/api/admin/roles/{id}", api.updateRole)
		protected.Get("/api/admin/audit", api.listSystemAudit)
		protected.Get("/api/submissions", api.listSubmissions)
		protected.Post("/api/submissions", api.createSubmission)
		protected.Get("/api/submissions/{id}", api.getSubmission)
		protected.Put("/api/submissions/{id}", api.updateSubmission)
		protected.Post("/api/submissions/{id}/transition", api.transitionSubmission)
		protected.Get("/api/submissions/{id}/audit", api.auditEvents)
	})

	return router
}

func (api API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
