package cmd

import (
	"github.com/samsamihd/hiro/pkg/hiro"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"start"},
	Short:   "Add a URL to queue",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		hiro.Start(queue, false)
	},
}

func init() {
	startCmd.Flags().StringVarP(&queue, "queue", "q", "", "Start specific queue")
	rootCmd.AddCommand(startCmd)
}
