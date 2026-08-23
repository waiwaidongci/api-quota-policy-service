package domain

import "time"

func (p Policy) Validate() error {
	if p.Name == "" || p.Limit <= 0 || p.WindowSeconds <= 0 {
		return ErrInvalid
	}
	if p.Algorithm != FixedWindow && p.Algorithm != SlidingWindow && p.Algorithm != TokenBucket {
		return ErrUnsupportedAlgorithm
	}
	if p.Algorithm == TokenBucket && p.Burst <= 0 {
		return ErrInvalid
	}
	return nil
}

func (p Policy) Window() time.Duration { return time.Duration(p.WindowSeconds) * time.Second }
