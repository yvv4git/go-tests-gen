package agent

import (
	"context"
	"encoding/json"

	"github.com/yvv4git/go-tests-gen/internal/usecases/tools"
)

type Scanner struct {
	entity tools.Scanner
}

func NewScanner(entity tools.Scanner) *Scanner {
	return &Scanner{
		entity: entity,
	}
}

func (s Scanner) Name() string {
	return "scanner"
}

func (s Scanner) Description() string {
	return `Scans recursive the go code directory for functions, 
	which is not covered by the tests. Input: JSON with path field`
}

func (s Scanner) Call(ctx context.Context, input string) (string, error) {
	var params struct {
		Path string `json:"path"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", err
	}

	uncoveredFunctions, err := s.entity.ScanDir(ctx, params.Path)
	if err != nil {
		return "", err
	}

	encodedResult, err := json.Marshal(uncoveredFunctions)
	if err != nil {
		return "", err
	}

	return string(encodedResult), nil
}
