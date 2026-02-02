package agent

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"

	"github.com/tmc/langchaingo/tools"
)

type AgentTools struct {
	scanner *Scanner
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

func (a *Agent) GenerateUnitTests(ctx context.Context, path string) error {
	mem := memory.NewConversationBuffer()

	agent := agents.NewConversationalAgent(
		a.llm, []tools.Tool{
			a.tools.scanner,
		},
		agents.WithMemory(mem),
	)

	executor := agents.NewExecutor(agent)

	response, err := executor.Call(
		ctx,
		map[string]any{"input": fmt.Sprintf("Scan the project code for functions not covered by tests in dir: %s", path)},
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
