package workflow

type Status string

const (
	Draft           Status = "draft"
	Submitted       Status = "submitted"
	ChangesRequired Status = "changes_required"
	Approved        Status = "approved"
	Rejected        Status = "rejected"
	Withdrawn       Status = "withdrawn"
)

type Role string

const (
	Requester Role = "requester"
	Reviewer  Role = "reviewer"
	Admin     Role = "admin"
)

func Terminal(status Status) bool {
	return status == Approved || status == Rejected || status == Withdrawn
}
