/*
Copyright © 2023 Aric Hansen <aric.p.hansen@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// go run -ldflags "-X futurama/cmd.Version=x.x.x" main.go version
var Version string = "v0.1.X" // set to git tag during automated build

// getCmd represents the get command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Display the version for Futurama CLI.",
	Long:    "Display the version for Futurama CLI.",
	Example: `  futurama version`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	asciiArt := [6]string{
		"  __       _                                     _____ _      _____ ",
		" / _|     | |                                   / ____| |    |_   _|",
		"| |_ _   _| |_ _   _ _ __ __ _ _ __ ___   __ _ | |    | |      | |  ",
		"|  _| | | | __| | | | '__/ _` | '_ ` _ \\ / _` || |    | |      | |  ",
		"| | | |_| | |_| |_| | | | (_| | | | | | | (_| || |____| |____ _| |_ ",
		"|_|  \\__,_|\\__|\\__,_|_|  \\__,_|_| |_| |_|\\__,_| \\_____|______|_____|",
	}

	if term.IsTerminal(0) {
		width, height, err := term.GetSize(0)
		if err != nil {
			return
		}
		if height > 6 && width > 78 { // make sure there's enough room for ascii art
			for _, i := range asciiArt {
				fmt.Println(i)
			}
		}
		fmt.Println("futurama cli " + Version)
	}
}
