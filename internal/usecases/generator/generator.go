package generator

import (
	"context"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

//
// Usecase - Generator
//

type Generator struct {
	log     ports.Logger
	scanner ports.Scanner
}

func NewGenerator(log ports.Logger, scanner ports.Scanner) *Generator {
	return &Generator{
		log:     log,
		scanner: scanner,
	}
}

func (g *Generator) Generate(ctx context.Context) error {
	if err := g.scanner.Scan(ctx); err != nil {
		return err
	}

	fnList := g.scanner.GetUncoveredFunctions()
	_ = fnList

	// todo: implement
	// 1. Find in projects functions without tests code coverage.
	// 2. Find files, read files where func placed.
	// 3. Find tests file for functions.
	// 4. Send files to LLM and get tests code.
	// 5. Get code from LLM, save to file.
	return nil
}
