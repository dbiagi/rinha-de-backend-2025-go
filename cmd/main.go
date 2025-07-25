package main

import (
	"github.com/spf13/cobra"

	"rinha2025/cmd/command"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "rinha",
		Short: "Rinha app CLI.",
	}

	rootCmd.AddCommand(command.NewServeCommand())
	rootCmd.AddCommand(command.NewWorkerCommand())
	rootCmd.Execute()
}
