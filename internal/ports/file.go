package ports

import "context"

// FileCat
type (
	ParamsFileCat struct {
		FilePath string
	}

	ResultFileCat struct {
		Content string
	}
)

// FileUpdate
type (
	ParamsFileUpdate struct {
		FilePath string
		Content  string
	}

	ResultFileUpdate struct{}
)

type FSProcessor interface {
	FileCat(ctx context.Context, params *ParamsFileCat) (*ResultFileCat, error)
	FileUpdate(ctx context.Context, params *ParamsFileUpdate) (*ResultFileUpdate, error)
}
