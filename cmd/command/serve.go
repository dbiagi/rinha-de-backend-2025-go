package command

import (
	"rinha2025/internal/config"
	internalhttp "rinha2025/internal/http"

	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	return &cobra.Command{
		Use: "serve",
		Run: runServe,
	}
}

func runServe(cmd *cobra.Command, args []string) {
	env, _ := cmd.Flags().GetString("env")
	if env == "" {
		env = config.DevelopmentEnv
	}
	c := config.LoadConfig(env)
	server := internalhttp.NewServer(c)
	server.Start()
}
