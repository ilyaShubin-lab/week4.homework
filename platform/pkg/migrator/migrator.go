package migrator

import "context"

type Migrator interface {
	Up(ctx context.Context) error
}
