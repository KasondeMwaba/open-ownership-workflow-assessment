package workflow

import (
	"errors"
	"fmt"
)

var ErrTransitionNotAllowed = errors.New("transition is not allowed")

var transitions = map[Status]map[Status][]Role{
	Draft: {
		Submitted: {Requester, Admin},
		Withdrawn: {Requester, Admin},
	},
	ChangesRequired: {
		Submitted: {Requester, Admin},
		Withdrawn: {Requester, Admin},
	},
	Submitted: {
		ChangesRequired: {Reviewer, Admin},
		Approved:        {Reviewer, Admin},
		Rejected:        {Reviewer, Admin},
	},
}

func CanTransition(from, to Status, role Role) bool {
	next, ok := transitions[from]
	if !ok {
		return false
	}
	roles, ok := next[to]
	if !ok {
		return false
	}
	for _, allowed := range roles {
		if allowed == role {
			return true
		}
	}
	return false
}

func ValidateTransition(from, to Status, role Role) error {
	if CanTransition(from, to, role) {
		return nil
	}
	return fmt.Errorf("%w: %s cannot move %s to %s", ErrTransitionNotAllowed, role, from, to)
}
