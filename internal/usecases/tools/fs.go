package tools

import (
	"context"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type FSTools struct {
	executor ports.FSProcessor
}

func NewFSTools(executor ports.FSProcessor) *FSTools {
	return &FSTools{
		executor: executor,
	}
}

func (f *FSTools) FileCat(ctx context.Context, params *ports.ParamsFileCat) (*ports.ResultFileCat, error) {
	return f.executor.FileCat(ctx, params)
}

func (f *FSTools) FileUpdate(ctx context.Context, params *ports.ParamsFileUpdate) (*ports.ResultFileUpdate, error) {
	return f.executor.FileUpdate(ctx, params)
}
