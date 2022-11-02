package cmd

import (
	pickle "github.com/rdcx/pickle/pkg"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate all functions from config",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pickle.Gen(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
