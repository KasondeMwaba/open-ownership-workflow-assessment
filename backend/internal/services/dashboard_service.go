package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/repositories"
	"openownership-workflow/backend/internal/workflow"
)

const dashboardCacheKey = "dashboard:stats:global:v2"

type DashboardService struct {
	repo  *repositories.Repository
	cache *redis.Client
}

func NewDashboardService(repo *repositories.Repository, cache *redis.Client) *DashboardService {
	return &DashboardService{repo: repo, cache: cache}
}

func (s *DashboardService) Stats(ctx context.Context, user models.User) (models.DashboardStats, error) {
	cacheKey := dashboardCacheKeyForUser(user)
	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
			var stats models.DashboardStats
			if json.Unmarshal([]byte(cached), &stats) == nil {
				stats.RedisCacheState = "hit"
				return stats, nil
			}
		}
	}

	stats, err := s.repo.DashboardStats(ctx, user)
	if err != nil {
		return models.DashboardStats{}, err
	}
	if s.cache != nil {
		if payload, err := json.Marshal(stats); err == nil {
			_ = s.cache.Set(ctx, cacheKey, payload, 30*time.Second).Err()
		}
	}
	return stats, nil
}

func (s *DashboardService) Invalidate(ctx context.Context) {
	if s.cache != nil {
		keys, err := s.cache.Keys(ctx, "dashboard:stats:*").Result()
		if err == nil && len(keys) > 0 {
			_ = s.cache.Del(ctx, keys...).Err()
		}
	}
}

func dashboardCacheKeyForUser(user models.User) string {
	if shouldUseOwnerScopedStats(user) {
		return fmt.Sprintf("dashboard:stats:user:%s:v2", user.ID)
	}
	return dashboardCacheKey
}

func shouldUseOwnerScopedStats(user models.User) bool {
	return user.HasPermission("submissions:create") && !user.HasPermission("submissions:review") && user.Role != workflow.Admin
}
