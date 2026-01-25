package generator

import (
	"context"

	"github.com/yvv4git/go-tests-gen/internal/generator/ports"
)

type Generator struct {
	log ports.Logger
}

func NewGenerator(log ports.Logger) *Generator {
	return &Generator{
		log: log,
	}
}

func (g *Generator) Gen(ctx context.Context) error {
	// todo: implement
	return nil
}
