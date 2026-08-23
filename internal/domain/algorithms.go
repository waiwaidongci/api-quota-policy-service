package domain

import "time"

type Counter interface {
	Allow(key string, limit int64, window time.Duration, burst int64, refill float64, now time.Time) (bool, int64, time.Time)
}

type WindowState struct {
	Count     int64
	StartedAt time.Time
	Tokens    float64
	Last      time.Time
	Hits      []time.Time
}

func refillTokens(s *WindowState, burst int64, rate float64, now time.Time) {
	if s.Last.IsZero() {
		s.Last = now
		s.Tokens = float64(burst)
		return
	}
	if rate <= 0 {
		return
	}
	s.Tokens += now.Sub(s.Last).Seconds() * rate
	if s.Tokens > float64(burst) {
		s.Tokens = float64(burst)
	}
	s.Last = now
}

// RefillTokens exposes token bucket refill behavior to infrastructure adapters.
func RefillTokens(s *WindowState, burst int64, rate float64, now time.Time) {
	refillTokens(s, burst, rate, now)
}
