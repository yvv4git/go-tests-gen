package generator

import (
	"context"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

//
// Usecase - Generator
//

type Generator struct {
	log   ports.Logger
	agent ports.Agent
}

func NewGenerator(log ports.Logger, agent ports.Agent) *Generator {
	return &Generator{
		log:   log,
		agent: agent,
	}
}

func (g *Generator) Generate(ctx context.Context, path string) error {
	// todo: implement
	// 1. Find in projects functions without tests code coverage.
	// 2. Find files, read files where func placed.
	// 3. Find tests file for functions.
	// 4. Send files to LLM and get tests code.
	// 5. Get code from LLM, save to file.
	return g.agent.GenerateUnitTests(ctx, path)
}
