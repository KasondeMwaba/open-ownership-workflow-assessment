package workflow

import "testing"

func TestValidateTransitionRoleRules(t *testing.T) {
	tests := []struct {
		name    string
		from    Status
		to      Status
		role    Role
		wantErr bool
	}{
		{"requester submits draft", Draft, Submitted, Requester, false},
		{"requester cannot approve", Submitted, Approved, Requester, true},
		{"reviewer requests changes", Submitted, ChangesRequired, Reviewer, false},
		{"terminal is locked", Approved, Submitted, Admin, true},
		{"admin can reject", Submitted, Rejected, Admin, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(tt.from, tt.to, tt.role)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
