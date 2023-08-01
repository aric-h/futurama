/*
Copyright © 2023 Aric Hansen <aric.p.hansen@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get quote or list of episodes",
	Long: `Get a Futurama quote from:
  - a random episode in a random season
  - a random episode in a user-defined season
  - a user-defined episode
  - a random episode in a random season from a user-defined character
  
Get a list of episodes from:
  - a user-defined season
  - the entire series
  
Get a list of supported character names to use with the 'get quote' command  `,
	Example: `  futurama get quote (no flags = randomized season and episode)
  futurama get quote --episode "Space Pilot 3000"
  futurama get episodes --season 2
  futurama get episodes --all
  futurama get characters
  `,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"quote", "episode"},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
