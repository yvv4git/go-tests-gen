package tools

import (
	"context"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type Scanner struct {
	executor ports.Scanner
}

func NewScanner(executor ports.Scanner) *Scanner {
	return &Scanner{
		executor: executor,
	}
}

func (s *Scanner) ScanDir(ctx context.Context, dir string) ([]ports.UncoveredFunc, error) {
	if err := s.executor.Scan(ctx); err != nil {
		return nil, err
	}

	uncoveredFunctions := s.executor.GetUncoveredFunctions()

	return uncoveredFunctions, nil
}
