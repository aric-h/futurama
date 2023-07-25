package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var DescribeEpisodeName string
var SeasonIndex int

var describeEpisodeCmd = &cobra.Command{
	Use:   "episode",
	Short: "Describe a Futurama episode (powered by Wikipedia)",
	Long:  "Describe the plot of a given Futurama episode",
	Example: `  futurama describe episode --name "Space Pilot 3000"
  `,
	// Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		err, SeasonIndex = validateEpisodeName(DescribeEpisodeName)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			cmd.Help()
		} else {
			describeEpisode()
		}
	},
}

func init() {
	describeCmd.AddCommand(describeEpisodeCmd)
	describeEpisodeCmd.Flags().StringVarP(&DescribeEpisodeName, "name", "n", "", "Episode name (use list-episodes command for assistance)")
}

func describeEpisode() {
	var resp *http.Response

	urlEpisodeName := strings.Replace(strings.Replace(DescribeEpisodeName, "'", "%27", -1), " ", "_", -1)
	if DescribeEpisodeName == "A Farewell to Arms" {
		urlEpisodeName = urlEpisodeName + "_(Futurama)"
	}
	wikiUrl := "https://en.wikipedia.org/wiki/" + urlEpisodeName
	infosphereUrl := "https://theinfosphere.org/" + urlEpisodeName
	resp = getHttpResponse(wikiUrl)
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	description := []string{}
	fmt.Println("starting wiki loop")
infosphereLoop:
	for { // loop until
		switch tokenizer.Next() {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break infosphereLoop //end of the file, break out of the loop
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		case html.EndTagToken:
			token := tokenizer.Token()
			if "table" == token.Data {
				fmt.Println("found table tag")
				switch tokenizer.Next() {
				case html.ErrorToken:
					err := tokenizer.Err()
					if err == io.EOF {
						break infosphereLoop //end of the file, break out of the loop
					}
					log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
				case html.StartTagToken:
					token := tokenizer.Token()
					fmt.Println(token.Data)
					if "p" == token.Data { // start of description section
						fmt.Println("found description section")
						descLine := ""
						for {
							switch tokenizer.Next() {
							case html.ErrorToken:
								err := tokenizer.Err()
								if err == io.EOF {
									break infosphereLoop //end of the file, break out of the loop
								}
								log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
							case html.TextToken: // assemble text of episode description
								token := tokenizer.Token()
								descLine = descLine + token.Data
							case html.StartTagToken:
								token := tokenizer.Token()
								if "p" == token.Data { // new paragraph; start new descLine
									description = append(description, descLine)
									fmt.Println(descLine)
									descLine = ""
								} else if "div" == token.Data { // end of description section
									description = append(description, descLine)
									break infosphereLoop
								}
							}
						}
					}
				}
			}
		}
	}

	printDescription(description, infosphereUrl)

}

func printDescription(description []string, url string) {
	fmt.Print("Season: ")
	fmt.Println(SeasonIndex)
	fmt.Print("Episode: ")
	fmt.Println(DescribeEpisodeName)
	fmt.Println("----")
	for _, line := range description {
		fmt.Println(line)
	}
	fmt.Println("----")
	fmt.Println(url)
}
