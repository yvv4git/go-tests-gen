package tools

import (
	"context"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type ExecTools struct {
	executor ports.Executor
}

func NewExecTools(executor ports.Executor) *ExecTools {
	return &ExecTools{
		executor: executor,
	}
}

func (e *ExecTools) RunGoUnitTest(ctx context.Context, params *ports.ParamsRunGoUnitTest) error {
	_, err := e.executor.RunGoUnitTest(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
