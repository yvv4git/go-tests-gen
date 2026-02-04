package generator

import (
	"context"
	"fmt"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

//
// Usecase - Generator
//

type Generator struct {
	log     ports.Logger
	agent   ports.Agent
	scanner ports.Scanner
}

func NewGenerator(log ports.Logger, scanner ports.Scanner, agent ports.Agent) *Generator {
	return &Generator{
		log:     log,
		scanner: scanner,
		agent:   agent,
	}
}

func (g *Generator) Generate(ctx context.Context, path string) error {
	// todo: implement
	// 1. Find in projects functions without tests code coverage.
	// 2. Find files, read files where func placed.
	// 3. Find tests file for functions.
	// 4. Send files to LLM and get tests code.
	// 5. Get code from LLM, save to file.
	if err := g.scanner.Scan(ctx, path); err != nil {
		return fmt.Errorf("scan dir[%s]: %w", path, err)
	}

	uncoveredFns := g.scanner.GetUncoveredFunctions()

	// for _, fn := range uncoveredFns {

	// }
	fn := uncoveredFns[len(uncoveredFns)-1] // get last fn

	err := g.agent.GenerateUnitTests(ctx, &ports.ParamsGenerateUnitTests{
		FilePath: fn.File,
	})
	if err != nil {
		return fmt.Errorf("generate unit test with LLM agent: %w", err)
	}

	return nil
}
