package ports

import "context"

type ParamsGenerateUnitTests struct {
	FilePath string // путь к исходному файлу с функцией
}

type Agent interface {
	GenerateUnitTests(ctx context.Context, params *ParamsGenerateUnitTests) error
}
