package repositories

import (
	"context"
	"time"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"
)

func (r *Repository) DashboardStats(ctx context.Context, user models.User) (models.DashboardStats, error) {
	var counts []struct {
		Status string
		Count  int
	}
	query := r.db.WithContext(ctx).
		Model(&models.Submission{}).
		Select("status, count(*) AS count")
	if shouldScopeToOwner(user) {
		query = query.Where("owner_id = ?", user.ID)
	}
	err := query.Group("status").Scan(&counts).Error
	if err != nil {
		return models.DashboardStats{}, err
	}
	stats := models.DashboardStats{ByStatus: map[string]int{}, GeneratedAt: time.Now().UTC(), RedisCacheState: "miss"}
	for _, row := range counts {
		stats.ByStatus[row.Status] = row.Count
		stats.Total += row.Count
	}
	stats.AwaitingReview = stats.ByStatus[string(workflow.Submitted)]
	stats.NeedsRequester = stats.ByStatus[string(workflow.ChangesRequired)] + stats.ByStatus[string(workflow.Draft)]
	stats.Completed = stats.ByStatus[string(workflow.Approved)] + stats.ByStatus[string(workflow.Rejected)] + stats.ByStatus[string(workflow.Withdrawn)]
	return stats, nil
}
