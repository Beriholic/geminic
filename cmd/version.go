package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const verison = "0.3.2"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of the geminic",
	Long:  `print the version of the geminic`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("geminic version %s\n", verison)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
