package agent

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/yvv4git/go-tests-gen/internal/ports"

	"github.com/tmc/langchaingo/tools"
)

type AgentTools struct {
	fileCat    *FileCatTool
	fileUpdate *FileUpdateTool
	runGoTest  *RunGoTestTool
}

type AgentOptions struct {
	Temperature float64
	MaxTokens   int
}

type Agent struct {
	llm   *openai.LLM
	tools AgentTools
	opts  AgentOptions
}

func (a *Agent) GenerateUnitTests(ctx context.Context, params *ports.ParamsGenerateUnitTests) error {
	agent := agents.NewOneShotAgent(
		a.llm,
		[]tools.Tool{
			a.tools.fileCat,
			a.tools.fileUpdate,
			a.tools.runGoTest,
		},
	)

	executor := agents.NewExecutor(agent)

	response, err := executor.Call(
		ctx,
		map[string]any{"input": a.setupPromptGenUnitTest(params)},
		chains.WithTemperature(a.opts.Temperature),
		chains.WithMaxTokens(a.opts.MaxTokens),
	)
	if err != nil {
		return fmt.Errorf("execute command: %w", err)
	}

	fmt.Println("Resp: ", response["output"])

	// todo: implement
	return nil
}

func (a *Agent) setupPromptGenUnitTest(params *ports.ParamsGenerateUnitTests) string {
	return fmt.Sprintf(`
You are a Go testing expert. Your goal is to write unit tests for uncovered functions in a Go file.

Source file path: %s

Your task:
1. Read the source file using the available tools
2. Analyze the code and identify functions that need tests
3. Write complete unit tests using Go testing package
4. Save tests to a _test.go file in the same directory
5. Verify tests work correctly

Begin now.`, params.FilePath)
}
