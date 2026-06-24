package models

import "time"

type DashboardStats struct {
	Total           int            `json:"total"`
	ByStatus        map[string]int `json:"byStatus"`
	AwaitingReview  int            `json:"awaitingReview"`
	NeedsRequester  int            `json:"needsRequester"`
	Completed       int            `json:"completed"`
	GeneratedAt     time.Time      `json:"generatedAt"`
	RedisCacheState string         `json:"redisCacheState"`
}
