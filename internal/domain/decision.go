package domain

import "time"

func NewDecision(p Policy, key string, allowed bool, remaining int64, reset, timeNow time.Time) Decision {
	reason := "allowed"
	if !allowed {
		reason = "limit_exceeded"
	}
	return Decision{Allowed: allowed, PolicyID: p.ID, Algorithm: string(p.Algorithm), Key: key, Limit: p.Limit, Remaining: remaining, ResetAt: reset, Reason: reason, CreatedAt: timeNow}
}
func (d Decision) StatusCode() int {
	if d.Allowed {
		return 200
	}
	return 429
}
func (d Decision) RetryAfter(now time.Time) int64 {
	if d.Allowed || d.ResetAt.Before(now) {
		return 0
	}
	return int64(d.ResetAt.Sub(now).Seconds()) + 1
}
func (d Decision) Headers(now time.Time) map[string]string {
	return map[string]string{"X-RateLimit-Limit": itoa(d.Limit), "X-RateLimit-Remaining": itoa(d.Remaining), "X-RateLimit-Reset": d.ResetAt.UTC().Format(time.RFC3339), "Retry-After": itoa(d.RetryAfter(now))}
}
func itoa(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := ""
	if v < 0 {
		neg = "-"
		v = -v
	}
	b := make([]byte, 0, 20)
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	return neg + string(b)
}
func clamp(v, min, max int64) int64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
