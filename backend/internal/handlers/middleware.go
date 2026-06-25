package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"openownership-workflow/backend/internal/auth"
	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/services"
)

const userContextKey = "user"

func (api API) requestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		r := c.Request()
		api.logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
		return err
	}
}

func (api API) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return writeError(c, http.StatusUnauthorized, "missing bearer token")
		}
		claims, err := auth.ParseToken(api.cfg.JWTSecret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			return writeError(c, http.StatusUnauthorized, "invalid token")
		}
		user, err := api.auth.FindUserByID(r.Context(), claims.UserID)
		if err != nil {
			return writeError(c, http.StatusUnauthorized, "user no longer exists")
		}
		if !user.IsActive {
			return writeError(c, http.StatusForbidden, "account is disabled")
		}
		c.Set(userContextKey, user)
		return next(c)
	}
}

func (api API) auditActivity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := currentUser(c)
		r := c.Request()
		start := time.Now()
		err := next(c)
		duration := time.Since(start).Milliseconds()
		if shouldSkipActivityAudit(r.URL.Path) {
			return err
		}
		statusCode := c.Response().Status
		if statusCode == 0 {
			statusCode = http.StatusOK
			if err != nil {
				statusCode = http.StatusInternalServerError
			}
		}
		api.audit.RecordActivityEvent(r.Context(), services.ActivityAuditInput{
			ActorID:    user.ID,
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			StatusCode: statusCode,
			Success:    statusCode < http.StatusBadRequest,
			DurationMs: duration,
			IPAddress:  clientIP(r),
			UserAgent:  r.UserAgent(),
			Browser:    browserName(r.UserAgent()),
			Metadata: map[string]any{
				"contentLength": r.ContentLength,
				"referer":       r.Referer(),
			},
		})
		return err
	}
}

func shouldSkipActivityAudit(path string) bool {
	return path == "/api/me"
}

func currentUser(c echo.Context) models.User {
	user, _ := c.Get(userContextKey).(models.User)
	return user
}
