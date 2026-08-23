package adapter

import "time"

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}
type SystemClock struct{}

func (SystemClock) Now() time.Time                  { return time.Now() }
func (SystemClock) Since(t time.Time) time.Duration { return time.Since(t) }

type FixedClock struct{ T time.Time }

func (c *FixedClock) Now() time.Time                  { return c.T }
func (c *FixedClock) Since(t time.Time) time.Duration { return c.T.Sub(t) }
func (c *FixedClock) Advance(d time.Duration)         { c.T = c.T.Add(d) }
