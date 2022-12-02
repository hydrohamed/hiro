package cmd

import (
	"github.com/samsamihd/hiro/pkg/hiro"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import",
	Aliases: []string{"import"},
	Short:   "Import a URL to queue",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hiro.Import(args[0], queue)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
