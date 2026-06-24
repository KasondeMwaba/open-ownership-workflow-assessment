package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"openownership-workflow/backend/internal/auth"
	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/services"
)

type contextKey string

const userContextKey contextKey = "user"

func (api API) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		api.logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	})
}

func (api API) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		claims, err := auth.ParseToken(api.cfg.JWTSecret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		user, err := api.auth.FindUserByID(r.Context(), claims.UserID)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "user no longer exists")
			return
		}
		if !user.IsActive {
			writeError(w, http.StatusForbidden, "account is disabled")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func (api API) auditActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := currentUser(r)
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(recorder, r)
		duration := time.Since(start).Milliseconds()
		if shouldSkipActivityAudit(r.URL.Path) {
			return
		}
		api.audit.RecordActivityEvent(r.Context(), services.ActivityAuditInput{
			ActorID:    user.ID,
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			StatusCode: recorder.statusCode,
			Success:    recorder.statusCode < http.StatusBadRequest,
			DurationMs: duration,
			IPAddress:  clientIP(r),
			UserAgent:  r.UserAgent(),
			Browser:    browserName(r.UserAgent()),
			Metadata: map[string]any{
				"contentLength": r.ContentLength,
				"referer":       r.Referer(),
			},
		})
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (recorder *statusRecorder) WriteHeader(statusCode int) {
	recorder.statusCode = statusCode
	recorder.ResponseWriter.WriteHeader(statusCode)
}

func shouldSkipActivityAudit(path string) bool {
	return path == "/api/me"
}

func currentUser(r *http.Request) models.User {
	user, _ := r.Context().Value(userContextKey).(models.User)
	return user
}
