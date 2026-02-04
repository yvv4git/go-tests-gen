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
	fileCat    *FileCat
	fileUpdate *FileUpdate
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

You have access to the following tools:

1. file_cat - Read file contents
   Use this tool to read the source file first.
   Input: {"path": "full_path_to_file"}

2. file_update - Create or update file
   Use this tool to save the generated test file.
   Input: {"path": "full_path_to_test_file", "content": "complete_test_code"}

Your task:
1. First, use file_cat to read: %s
2. Analyze the code and identify functions that need tests
3. Write complete unit tests using Go testing package
4. Save tests to a _test.go file in the same directory
   - If source is /path/to/file.go, test should be /path/to/file_test.go

IMPORTANT:
- After generating test code, you MUST use file_update tool to save it
- content should be the COMPLETE test file content (no markdown backticks, no explanations)
- DO NOT output test code directly - use file_update tool

Example:
Thought: I need to read the source file first.
Action: file_cat
Action Input: {"path": "/path/to/source.go"}

Thought: Now I need to save the test file.
Action: file_update
Action Input: {"path": "/path/to/source_test.go", "content": "package pkg\n\nimport \"testing\"\n\nfunc TestFunc(t *testing.T) {}"}

Begin now.`, params.FilePath, params.FilePath)
}
