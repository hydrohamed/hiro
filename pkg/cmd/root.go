package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var queue string

var rootCmd = &cobra.Command{
	Use:     "hiro",
	Version: "0.0.1",
	Short:   "A powerful, fast, and lightweight download manager and accelerator.",
	Long: `A powerful, fast, and lightweight download manager and accelerator built with
			love by samsamihd in Go.
			Complete documentation is available at https://github.com/samsamihd/hiro`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&queue, "queue", "q", "Add download task to specific queue")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
