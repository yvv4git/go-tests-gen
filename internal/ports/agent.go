package ports

import "context"

type Agent interface {
	GenerateUnitTests(ctx context.Context, path string) error
}
