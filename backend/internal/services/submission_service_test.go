package services

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestValidateSubmissionPayload(t *testing.T) {
	validData := json.RawMessage(`{"company":"Acme Ltd","jurisdiction":"Zambia","registrationNumber":"PACRA-1","beneficialOwners":[{"name":"Jane Doe","ownershipPercent":42.5,"controlType":"shares"}]}`)

	tests := []struct {
		name    string
		payload SubmissionPayload
		wantErr bool
	}{
		{name: "valid", payload: SubmissionPayload{Title: "Title", Summary: "Summary", Data: validData}},
		{name: "missing title", payload: SubmissionPayload{Summary: "Summary", Data: validData}, wantErr: true},
		{name: "invalid json", payload: SubmissionPayload{Title: "Title", Summary: "Summary", Data: json.RawMessage(`{`)}, wantErr: true},
		{name: "missing company", payload: SubmissionPayload{Title: "Title", Summary: "Summary", Data: json.RawMessage(`{"jurisdiction":"Zambia","registrationNumber":"PACRA-1","beneficialOwners":[{"name":"Jane Doe","ownershipPercent":42.5,"controlType":"shares"}]}`)}, wantErr: true},
		{name: "invalid ownership percent", payload: SubmissionPayload{Title: "Title", Summary: "Summary", Data: json.RawMessage(`{"company":"Acme Ltd","jurisdiction":"Zambia","registrationNumber":"PACRA-1","beneficialOwners":[{"name":"Jane Doe","ownershipPercent":120,"controlType":"shares"}]}`)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateSubmissionPayload(tt.payload)
			if tt.wantErr && !errors.Is(err, ErrInvalidSubmission) {
				t.Fatalf("expected ErrInvalidSubmission, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
