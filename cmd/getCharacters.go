/*
Copyright © 2023 Aric Hansen | aric.p.hansen@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var charactersCmd = &cobra.Command{
	Use:     "characters",
	Short:   "Get list of valid characters for passing into the 'get quotes' command",
	Long:    "Get list of valid characters for passing into the 'get quotes' command",
	Example: `  futurama get characters`,
	Run: func(cmd *cobra.Command, args []string) {
		listSupportedCharacters()
	},
}

func init() {
	getCmd.AddCommand(charactersCmd)
}

func listSupportedCharacters() {
	supportedCharacters := getSupportedCharacters()

	fmt.Println()
	for _, c := range supportedCharacters {
		fmt.Println(c)
	}
}
