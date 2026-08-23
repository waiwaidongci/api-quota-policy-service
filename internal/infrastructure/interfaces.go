package infrastructure

import "context"

type Database interface {
	Ping(context.Context) error
	Close() error
}
type Redis interface {
	Ping(context.Context) error
	Close() error
}
type EventBus interface {
	Publish(context.Context, []byte) error
	Close() error
}
type Metrics interface {
	Inc(string)
	Observe(string, float64)
}
