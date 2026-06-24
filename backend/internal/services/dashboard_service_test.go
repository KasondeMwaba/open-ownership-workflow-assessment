package services

import (
	"testing"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"
)

func TestDashboardCacheKeyForUser(t *testing.T) {
	requester := models.User{BaseModel: models.BaseModel{ID: "requester-1"}, Role: workflow.Requester, Permissions: []string{"submissions:create", "dashboard:view"}}
	reviewer := models.User{BaseModel: models.BaseModel{ID: "reviewer-1"}, Role: workflow.Reviewer, Permissions: []string{"submissions:review", "dashboard:view"}}
	admin := models.User{BaseModel: models.BaseModel{ID: "admin-1"}, Role: workflow.Admin}

	if got := dashboardCacheKeyForUser(requester); got != "dashboard:stats:user:requester-1:v2" {
		t.Fatalf("requester should use scoped cache key, got %q", got)
	}
	if got := dashboardCacheKeyForUser(reviewer); got != dashboardCacheKey {
		t.Fatalf("reviewer should use global cache key, got %q", got)
	}
	if got := dashboardCacheKeyForUser(admin); got != dashboardCacheKey {
		t.Fatalf("admin should use global cache key, got %q", got)
	}
}
