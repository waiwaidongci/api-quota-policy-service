package domain

import "time"

func FixedWindowAllow(s *WindowState, limit int64, window time.Duration, now time.Time) (bool, int64, time.Time) {
	if s.StartedAt.IsZero() {
		s.StartedAt = now
	}
	if now.Sub(s.StartedAt) >= window {
		s.StartedAt = now
		s.Count = 0
	}
	reset := s.StartedAt.Add(window)
	if s.Count >= limit {
		return false, 0, reset
	}
	s.Count++
	return true, limit - s.Count, reset
}
func SlidingWindowAllow(s *WindowState, limit int64, window time.Duration, now time.Time) (bool, int64, time.Time) {
	cut := now.Add(-window)
	first := 0
	for first < len(s.Hits) && s.Hits[first].Before(cut) {
		first++
	}
	s.Hits = s.Hits[first:]
	reset := now.Add(window)
	if len(s.Hits) > 0 {
		reset = s.Hits[0].Add(window)
	}
	if int64(len(s.Hits)) >= limit {
		return false, 0, reset
	}
	s.Hits = append(s.Hits, now)
	return true, limit - int64(len(s.Hits)), reset
}
func TokenBucketAllow(s *WindowState, burst int64, rate float64, now time.Time) (bool, int64, time.Time) {
	RefillTokens(s, burst, rate, now)
	if s.Tokens < 1 {
		return false, int64(s.Tokens), now.Add(time.Second)
	}
	s.Tokens--
	return true, int64(s.Tokens), now.Add(time.Second)
}
func Expiry(now time.Time, window time.Duration) time.Time { return now.Add(window) }
func IsExpired(now, expiry time.Time) bool                 { return !expiry.IsZero() && !now.Before(expiry) }
func WindowSeconds(d time.Duration) int64                  { return int64(d / time.Second) }
