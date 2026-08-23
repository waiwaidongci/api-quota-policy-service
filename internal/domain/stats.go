package domain

import "time"

type DecisionStats struct {
	Total    uint64
	Allowed  uint64
	Rejected uint64
	LastAt   time.Time
}

func (s *DecisionStats) Record(d Decision) {
	s.Total++
	if d.Allowed {
		s.Allowed++
	} else {
		s.Rejected++
	}
	s.LastAt = d.CreatedAt
}
func (s DecisionStats) Rate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Allowed) / float64(s.Total)
}
func (s DecisionStats) RejectionRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Rejected) / float64(s.Total)
}
func (s DecisionStats) Healthy(threshold float64) bool { return s.RejectionRate() <= threshold }
func (s DecisionStats) Merge(other DecisionStats) DecisionStats {
	r := DecisionStats{Total: s.Total + other.Total, Allowed: s.Allowed + other.Allowed, Rejected: s.Rejected + other.Rejected, LastAt: s.LastAt}
	if other.LastAt.After(r.LastAt) {
		r.LastAt = other.LastAt
	}
	return r
}
func (s DecisionStats) Reset() { s.Total = 0; s.Allowed = 0; s.Rejected = 0; s.LastAt = time.Time{} }
func Percentile(values []time.Duration, p float64) time.Duration {
	if len(values) == 0 {
		return 0
	}
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	idx := int(float64(len(values)-1) * p)
	return values[idx]
}
func Average(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}
	var n time.Duration
	for _, v := range values {
		n += v
	}
	return n / time.Duration(len(values))
}
