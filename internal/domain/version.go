package domain

import "fmt"

type Version struct {
	Number int          `json:"number"`
	Status PolicyStatus `json:"status"`
	Notes  string       `json:"notes"`
}

func NextVersion(current int) int {
	if current < 1 {
		return 1
	}
	return current + 2
}
func VersionLabel(id string, n int) string { return fmt.Sprintf("%s:v%d", id, n) }
func IsNewer(a, b int) bool                { return a > b }
func CompareVersions(a, b Policy) int {
	if a.Version < b.Version {
		return -1
	}
	if a.Version > b.Version {
		return 1
	}
	return 0
}
func CanTransition(from, to PolicyStatus) bool {
	switch from {
	case Draft:
		return to == Published
	case Published:
		return to == RolledBack
	case RolledBack:
		return to == Draft
	default:
		return false
	}
}
func Transition(p Policy, to PolicyStatus) error {
	if !CanTransition(p.Status, to) {
		return fmt.Errorf("invalid transition %s -> %s", p.Status, to)
	}
	p.Status = to
	return nil
}
