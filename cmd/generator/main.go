package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/yvv4git/go-tests-gen/internal/adapters/logger"
	"github.com/yvv4git/go-tests-gen/internal/adapters/scanner"
	"github.com/yvv4git/go-tests-gen/internal/usecases/generator"
)

func main() {
	rootCommand := &cobra.Command{
		Use:   "gen-tests",
		Short: "Analyse go code and generate go tests",
	}

	rootCommand.PersistentFlags().StringP("config", "c", "config.toml", "Path to config file")

	unitTestsGenCommand := &cobra.Command{
		Use:   "unit",
		Short: "Generate only unit tests",
		Run: func(cmd *cobra.Command, args []string) {
			cfgFilePath, _ := cmd.Flags().GetString("config")
			runUnitTestsGenCommand(cfgFilePath)
		},
	}

	rootCommand.AddCommand(unitTestsGenCommand)

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

func runUnitTestsGenCommand(_ string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slogLogger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	log := logger.NewSlogAdapter(slogLogger)

	scanner := scanner.NewCoverage("/Users/vladimireliseev/mydata/research/go-pkg-safe/") // todo: setup path

	gen := generator.NewGenerator(log, scanner)

	if err := gen.Generate(ctx); err != nil {
		log.Error("Failed scan", err)
	}
}
