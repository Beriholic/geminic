package cmd

import (
	"fmt"
	"os"

	"github.com/Beriholic/geminic/internal"
	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "select Gemini's model",
	Long:  `select Gemini's model`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		err := internal.UpdateGeminiModelSelect(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}
