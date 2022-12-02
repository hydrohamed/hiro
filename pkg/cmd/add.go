package cmd

import (
	"github.com/samsamihd/hiro/pkg/hiro"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"add"},
	Short:   "Add a URL to queue",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hiro.Add(args[0], queue)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
