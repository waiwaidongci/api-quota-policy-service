package domain

import "sort"

type ByPriority []Policy

func (p ByPriority) Len() int { return len(p) }
func (p ByPriority) Less(i, j int) bool {
	if p[i].Priority == p[j].Priority {
		return p[i].UpdatedAt.Before(p[j].UpdatedAt)
	}
	return p[i].Priority > p[j].Priority
}
func (p ByPriority) Swap(i, j int)     { p[i], p[j] = p[j], p[i] }
func SortPolicies(p []Policy) []Policy { sort.Sort(ByPriority(p)); return p }
func HighestPriority(p []Policy) (Policy, bool) {
	if len(p) == 0 {
		return Policy{}, false
	}
	SortPolicies(p)
	return p[0], true
}
func ActivePolicies(p []Policy) []Policy {
	o := make([]Policy, 0, len(p))
	for _, v := range p {
		if v.Status == Published {
			o = append(o, v)
		}
	}
	return SortPolicies(o)
}
func SameScope(a, b Policy) bool { return a.Rule == b.Rule && a.Algorithm == b.Algorithm }
func IsPublished(p Policy) bool  { return p.Status == Published }
func CanPublish(p Policy) bool   { return p.Status == Draft || p.Status == RolledBack }
func CanRollback(p Policy) bool  { return p.Status == Published }
