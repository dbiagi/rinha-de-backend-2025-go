package command

import "github.com/spf13/cobra"

func NewHealthCheckWorker() *cobra.Command {
	return &cobra.Command{
		Use: "worker",
		Run: work,
	}
}

func work(cmd *cobra.Command, args []string) {

}
