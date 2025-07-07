package main

import (
	"github.com/spf13/cobra"

	"rinha2025/cmd/command"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "shopping-bag",
		Short: "Shopping bag app CLI.",
	}

	rootCmd.AddCommand(command.NewServeCommand())
	rootCmd.Execute()
}
