/*
Copyright © 2023 Aric Hansen | aric.p.hansen@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var AllEpisodes bool
var SeasonNumber int

var episodesCmd = &cobra.Command{
	Use:   "episodes",
	Short: "Get list of episodes from series or season",
	Long:  "Get list of all episodes or only episodes for a given season",
	Example: `  futurama get episodes (return all episodes if no flags provided)
  futurama get episodes --all
  futurama get episodes --season 2`,
	Run: func(cmd *cobra.Command, args []string) {
		err := listEpisodes()
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			cmd.Help()
		}
	},
}

func init() {
	getCmd.AddCommand(episodesCmd)
	episodesCmd.Flags().BoolVarP(&AllEpisodes, "all", "a", true, "Show episodes from all seasons")
	episodesCmd.Flags().IntVarP(&SeasonNumber, "season", "s", 0, "Season number (1-7)")
	episodesCmd.MarkFlagsMutuallyExclusive("season", "all")
}

func listEpisodes() error {
	if SeasonNumber != 0 { // if season provided, turn off default -a flag
		AllEpisodes = false
	}

	series := getSeries()
	if AllEpisodes {
		for _, season := range series {
			fmt.Println("#### " + season.name + " ####")
			for _, ep := range season.episodes {
				fmt.Print(ep + "\n")
			}
			fmt.Println()
		}
		return nil
	} else {
		if SeasonNumber > 0 && SeasonNumber < 8 {
			fmt.Println("#### " + series[SeasonNumber-1].name + " ####")
			for _, ep := range series[SeasonNumber-1].episodes {
				fmt.Print(ep + "\n")
			}
			return nil
		} else {
			return errors.New("Invalid season number")
		}
	}
}
