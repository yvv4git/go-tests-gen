package main

import (
	"context"
	stdLog "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/yvv4git/go-tests-gen/internal/adapters/agent"
	"github.com/yvv4git/go-tests-gen/internal/adapters/fs"
	"github.com/yvv4git/go-tests-gen/internal/adapters/logger"
	"github.com/yvv4git/go-tests-gen/internal/adapters/scanner"
	"github.com/yvv4git/go-tests-gen/internal/config"
	"github.com/yvv4git/go-tests-gen/internal/usecases/generator"
	"github.com/yvv4git/go-tests-gen/internal/usecases/tools"
)

func main() {
	rootCommand := &cobra.Command{
		Use:   "gen-tests",
		Short: "Analyse go code and generate go tests",
	}

	rootCommand.PersistentFlags().StringP("config", "c", "config.toml", "Path to config file")
	rootCommand.PersistentFlags().StringP("path", "p", ".", "Path to directory to scan")

	unitTestsGenCommand := &cobra.Command{
		Use:   "unit",
		Short: "Generate only unit tests",
		Run: func(cmd *cobra.Command, args []string) {
			cfgFilePath, _ := cmd.Flags().GetString("config")
			scanPath, _ := cmd.Flags().GetString("path")
			runUnitTestsGenCommand(cfgFilePath, scanPath)
		},
	}

	rootCommand.AddCommand(unitTestsGenCommand)

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

func runUnitTestsGenCommand(cfgFilePath, path string) {
	var cfg config.Config

	if err := config.Load(cfgFilePath, &cfg); err != nil {
		stdLog.Fatalf("Failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slogLogger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	log := logger.NewSlogAdapter(slogLogger)

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: nil, // no proxy
		},
		Timeout: 30 * time.Second, // todo: make it configurable
	}

	llm, err := openai.New(
		openai.WithBaseURL(cfg.LLM.OpenAI.URL),
		openai.WithToken(cfg.LLM.OpenAI.Token),
		openai.WithModel(cfg.LLM.OpenAI.Model),
		openai.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatalf("Failed setup llm: %v", err)
	}

	scannerAdapter := scanner.NewScanner()
	fsAdapter := fs.NewFS()
	fsInbound := tools.NewFSTools(fsAdapter)
	fileCatTool := agent.NewFileCat(fsInbound)
	fileUpdateTool := agent.NewFileUpdate(fsInbound)

	agent, err := agent.NewAgentBuilder().
		SetLLM(llm).
		SetToolFileCat(fileCatTool).
		SetToolFileUpdate(fileUpdateTool).
		SetOptions(agent.AgentOptions{
			Temperature: cfg.LLM.Temperature,
			MaxTokens:   cfg.LLM.MaxTokens,
		}).Build()
	if err != nil {
		log.Fatalf("Failes setup agent: %v", err)
	}

	gen := generator.NewGenerator(log, scannerAdapter, agent)

	if err := gen.Generate(ctx, path); err != nil {
		log.Error("Failed run generate", "error", err)
	}
}
