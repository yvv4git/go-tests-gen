package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) runGoUnitTest(ctx context.Context, sourceFilePath, testFuncName string) error {
	// dir is the package directory (e.g., /path/to/project/pkg)
	dir := filepath.Dir(sourceFilePath)
	// cmd.Dir is the project root (parent of package directory)
	projectRoot := filepath.Dir(dir)

	// pkgPath is relative path from project root to package (e.g., ./pkg)
	pkgPath := "." + string(filepath.Separator) + filepath.Base(dir)

	// Create command: go test -run TestName ./pkg
	cmd := exec.CommandContext(ctx, "go", "test", "-run", testFuncName, pkgPath)
	cmd.Dir = projectRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Include stdout and stderr in error message
		stderrStr := stderr.String()
		if stderrStr == "" {
			stderrStr = "(no stderr output)"
		}
		stdoutStr := stdout.String()
		if stdoutStr == "" {
			stdoutStr = "(no stdout output)"
		}
		return fmt.Errorf("go test failed: stdout=%s, stderr=%s, err=%w", stdoutStr, stderrStr, err)
	}

	return nil
}

func (e *Executor) RunGoUnitTest(ctx context.Context, params *ports.ParamsRunGoUnitTest) (*ports.ResultRunGoUnitTest, error) {
	if err := e.runGoUnitTest(ctx, params.SourceFilePath, params.TestFnName); err != nil {
		return nil, err
	}

	return &ports.ResultRunGoUnitTest{}, nil
}
