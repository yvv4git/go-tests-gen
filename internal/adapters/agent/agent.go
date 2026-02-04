package agent

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
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
	mem := memory.NewConversationBuffer()

	agent := agents.NewConversationalAgent(
		a.llm, []tools.Tool{
			a.tools.fileCat,
			a.tools.fileUpdate,
		},
		agents.WithMemory(mem),
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
You are a Go testing expert. Your task is to write a unit test for a specific function in a Go file.

INPUT DATA:
- File path: %s
- Function code to test:
%s

TOOLS AVAILABLE:
1. file_cat - Read file contents
   Format: {"path": "<absolute_path>"}
2. file_update - Create or update file
   Format: {"path": "<absolute_path>", "content": "<file_content>"}

INSTRUCTIONS:
1. First, use file_cat to read the source file and understand the package structure
2. Write unit tests using the standard Go testing package
3. Create a complete _test.go file for the specified function
4. Use table-driven tests where appropriate
5. Cover edge cases and error conditions
6. Use the same package name as the source file
7. Output ONLY the complete test file content, no explanations, no markdown backticks

IMPORTANT: When calling tools, always use the exact JSON format shown above.

Write the complete unit test file now:`, params.FilePath, params.FnCode)
}
