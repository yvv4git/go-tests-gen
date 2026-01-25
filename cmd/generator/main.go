package main

import (
	"os"

	"github.com/spf13/cobra"
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

func runUnitTestsGenCommand(cfgFilePath string) {

}
