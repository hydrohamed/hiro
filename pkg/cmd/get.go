package cmd

import (
	"github.com/samsamihd/hiro/pkg/hiro"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"get"},
	Short:   "Download a URL",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hiro.Get(args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
