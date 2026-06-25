package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"openownership-workflow/backend/internal/dto"
	"openownership-workflow/backend/internal/services"
)

func (api API) login(c echo.Context) error {
	r := c.Request()
	var payload dto.LoginRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	result, err := api.auth.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "could not issue token"
		reason := "token_error"
		if errors.Is(err, services.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
			message = "invalid email or password"
			reason = "invalid_credentials"
		} else if errors.Is(err, services.ErrUserDisabled) {
			status = http.StatusForbidden
			message = "account is disabled"
			reason = "account_disabled"
		}
		api.recordLoginAttempt(r, payload.Email, nil, false, reason)
		return writeError(c, status, message)
	}
	api.recordLoginAttempt(r, result.User.Email, &result.User.ID, true, "authenticated")
	return writeJSON(c, http.StatusOK, result)
}

func (api API) logout(c echo.Context) error {
	r := c.Request()
	user := currentUser(c)
	api.audit.RecordSessionEvent(r.Context(), services.SessionAuditInput{
		ActorID:   &user.ID,
		Email:     user.Email,
		EventType: "logout",
		Success:   true,
		IPAddress: clientIP(r),
		UserAgent: r.UserAgent(),
		Browser:   browserName(r.UserAgent()),
		Reason:    "user_signed_out",
		Metadata:  map[string]any{"method": r.Method, "path": r.URL.Path},
	})
	return writeJSON(c, http.StatusOK, dto.LogoutResponse{Status: "logged_out"})
}

func (api API) me(c echo.Context) error {
	return writeJSON(c, http.StatusOK, currentUser(c))
}

func (api API) recordLoginAttempt(r *http.Request, email string, actorID *string, success bool, reason string) {
	if actorID == nil {
		if user, err := api.auth.FindUserByEmail(r.Context(), email); err == nil {
			actorID = &user.ID
		}
	}
	api.audit.RecordSessionEvent(r.Context(), services.SessionAuditInput{
		ActorID:   actorID,
		Email:     email,
		EventType: "login",
		Success:   success,
		IPAddress: clientIP(r),
		UserAgent: r.UserAgent(),
		Browser:   browserName(r.UserAgent()),
		Reason:    reason,
		Metadata:  map[string]any{"method": r.Method, "path": r.URL.Path},
	})
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	host := r.RemoteAddr
	if index := strings.LastIndex(host, ":"); index > -1 {
		return host[:index]
	}
	return host
}

func browserName(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "edg/"):
		return "Microsoft Edge"
	case strings.Contains(ua, "chrome/") && !strings.Contains(ua, "chromium"):
		return "Chrome"
	case strings.Contains(ua, "firefox/"):
		return "Firefox"
	case strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/"):
		return "Safari"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		return "Opera"
	default:
		return "Unknown"
	}
}
