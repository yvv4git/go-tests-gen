package ports

import (
	"context"
)

type UncoveredFunc struct {
	Package string
	File    string
	Name    string
	Line    int
	FnCode  string
}

type (
	Scanner interface {
		Scan(ctx context.Context, dir string) error
		GetUncoveredFunctions() []UncoveredFunc
	}
)
