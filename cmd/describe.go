package cmd

import (
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a Futurama episode (powered by Infosphere)",
	Long:  "Describe the plot of a given Futurama episode",
	Example: `  futurama describe episode --name "Space Pilot 3000"
  `,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
