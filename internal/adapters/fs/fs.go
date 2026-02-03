package fs

import (
	"context"
	"os"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type FS struct{}

func NewFS() *FS {
	return &FS{}
}

func (f *FS) FileCat(ctx context.Context, params *ports.ParamsFileCat) (*ports.ResultFileCat, error) {
	content, err := os.ReadFile(params.FilePath)
	if err != nil {
		return nil, err
	}

	return &ports.ResultFileCat{
		Content: string(content),
	}, nil
}

func (f *FS) FileUpdate(ctx context.Context, params *ports.ParamsFileUpdate) (*ports.ResultFileUpdate, error) {
	err := os.WriteFile(params.FilePath, []byte(params.Content), os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &ports.ResultFileUpdate{}, nil
}
