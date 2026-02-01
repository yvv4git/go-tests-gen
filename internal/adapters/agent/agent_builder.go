package agent

import (
	"errors"

	"github.com/tmc/langchaingo/llms/openai"
)

type AgentBuilder struct {
	llm   *openai.LLM
	tools AgentTools
	opts  AgentOptions
}

func NewAgentBuilder() *AgentBuilder {
	const (
		defaultTemperature = 0.7
		defaultMaxTokens   = 1000
	)

	return &AgentBuilder{
		opts: AgentOptions{
			Temperature: defaultTemperature,
			MaxTokens:   defaultMaxTokens,
		},
	}
}

func (b *AgentBuilder) Build() (*Agent, error) {
	if b.llm == nil {
		return nil, errors.New("no set llm")
	}

	if b.tools.scanner == nil {
		return nil, errors.New("no set scanner tool")
	}

	entity := &Agent{
		llm:   b.llm,
		tools: b.tools,
		opts:  b.opts,
	}

	return entity, nil
}

func (b *AgentBuilder) SetLLM(value *openai.LLM) *AgentBuilder {
	b.llm = value
	return b
}

func (b *AgentBuilder) SetOptions(value AgentOptions) *AgentBuilder {
	b.opts = value
	return b
}

func (b *AgentBuilder) SetToolScanner(value *Scanner) *AgentBuilder {
	b.tools.scanner = value
	return b
}
