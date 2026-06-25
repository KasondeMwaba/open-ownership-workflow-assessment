package dto

import (
	"encoding/json"

	"openownership-workflow/backend/internal/workflow"
)

type SubmissionRequest struct {
	Title   string          `json:"title"`
	Summary string          `json:"summary"`
	Data    json.RawMessage `json:"data"`
}

type TransitionSubmissionRequest struct {
	Status  workflow.Status `json:"status"`
	Comment string          `json:"comment"`
}

type PaginatedResponse[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
