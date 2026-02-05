package executor

import (
	"context"
	"os/exec"
	"path/filepath"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) runGoUnitTest(ctx context.Context, sourceFilePath, testFuncName string) error {
	dir := filepath.Dir(sourceFilePath)

	// Create command: go test -run TestName ./path/to/package
	cmd := exec.CommandContext(ctx, "go", "test", "-run", testFuncName, "./"+dir)
	cmd.Dir = filepath.Dir(sourceFilePath) + "/.."

	return cmd.Run()
}

func (e *Executor) RunGoUnitTest(ctx context.Context, params *ports.ParamsRunGoUnitTest) (*ports.ResultRunGoUnitTest, error) {
	if err := e.runGoUnitTest(ctx, params.SourceFilePath, params.TestFnName); err != nil {
		return nil, err
	}

	return &ports.ResultRunGoUnitTest{}, nil
}
