package cmd

import (
	"fmt"
	"os"

	"github.com/Beriholic/geminic/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set the config file",
	Long:  `Set the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Create(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var childCmd = &cobra.Command{
	Use:   "local",
	Short: "Create local config file",
	Long:  "Create local config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.CreateLocal(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	configCmd.AddCommand(childCmd)
	rootCmd.AddCommand(configCmd)
}
