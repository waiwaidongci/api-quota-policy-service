package repository

import "context"

type Tx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
type Migrator interface {
	Up(context.Context) error
	Down(context.Context) error
}
type Health interface {
	Health(context.Context) error
	Ready(context.Context) error
}
