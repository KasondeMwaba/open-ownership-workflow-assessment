package repositories

import (
	"testing"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"
)

func TestShouldScopeToOwner(t *testing.T) {
	tests := []struct {
		name string
		user models.User
		want bool
	}{
		{name: "requester creator", user: models.User{Role: workflow.Requester, Permissions: []string{"submissions:create"}}, want: true},
		{name: "custom creator", user: models.User{Role: workflow.Role("analyst"), Permissions: []string{"submissions:create"}}, want: true},
		{name: "reviewer", user: models.User{Role: workflow.Reviewer, Permissions: []string{"submissions:review"}}, want: false},
		{name: "admin", user: models.User{Role: workflow.Admin}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldScopeToOwner(tt.user); got != tt.want {
				t.Fatalf("shouldScopeToOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}
