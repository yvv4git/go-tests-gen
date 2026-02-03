package ports

import "context"

type ParamsGenerateUnitTests struct {
	FilePath string
	FnCode   string
}

type Agent interface {
	GenerateUnitTests(ctx context.Context, params *ParamsGenerateUnitTests) error
}
