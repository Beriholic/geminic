package cmd

import (
	"fmt"
	"os"

	"github.com/Beriholic/geminic/internal"
	"github.com/spf13/cobra"
)

var userCommit string = ""

func init() {
	rootCmd.Flags().StringVarP(&userCommit, "commit", "c", "", "commit message")
}

var rootCmd = &cobra.Command{
	Use:   "geminic",
	Short: "Using Gemini to Write Git Commits ",
	Long:  `Using Gemini to Write Git Commits `,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		err := internal.GeneratorCommit(ctx, userCommit)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
