package cmd

import (
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a Futurama episode (powered by Wikipedia)",
	Long:  "Describe the plot of a user-defined Futurama episode",
	Example: `  futurama describe episode --name "Space Pilot 3000"
  `,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
