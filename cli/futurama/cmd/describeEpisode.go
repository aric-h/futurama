package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var DescribeEpisodeName string
var SeasonIndex int
var EpisodeIndex int

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

	resp = getHttpResponse("https://en.wikipedia.org/wiki/" + urlEpisodeName)
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	plot := []string{}
wikiLoop:
	for { // loop until
		switch tokenizer.Next() {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break wikiLoop //end of the file, break out of the loop
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		case html.StartTagToken:
			token := tokenizer.Token()
			if "span" == token.Data {
				for _, attr := range token.Attr {
					if attr.Val == "Plot" { // found plot section
						for {
							switch tokenizer.Next() {
							case html.ErrorToken:
								err := tokenizer.Err()
								if err == io.EOF {
									break wikiLoop //end of the file, break out of the loop
								}
							case html.StartTagToken:
								token := tokenizer.Token()
								if "p" == token.Data || "h3" == token.Data { // start of plot paragraph
									para := ""
								paragraphLoop:
									for {
										switch tokenizer.Next() {
										case html.ErrorToken:
											err := tokenizer.Err()
											if err == io.EOF {
												break wikiLoop // end of the file, break out of the loop
											}
										case html.TextToken:
											token := tokenizer.Token()
											para = para + string(token.Data) // append plot text
										case html.EndTagToken:
											token := tokenizer.Token()
											if "p" == token.Data || "h3" == token.Data {
												plot = append(plot, para)
												break paragraphLoop
											}
										}
									}
								} else if "h2" == token.Data { // end of plot section
									break wikiLoop
								}
							}
						}
					}
				}
			}
		}
	}

	// remove edit links
	editEx := regexp.MustCompile(`\[edit\]`)
	for i, para := range plot {
		plot[i] = editEx.ReplaceAllString(para, "")
	}

	printDescription(plot, urlEpisodeName)

}

func printDescription(plot []string, urlEpisodeName string) {
	fmt.Println("\nINFO")
	fmt.Println("----")
	fmt.Print("Season: ")
	fmt.Println(SeasonIndex)
	fmt.Print("Episode: ")
	fmt.Println(EpisodeIndex)
	fmt.Print("Title: ")
	fmt.Println(DescribeEpisodeName)

	fmt.Println("\nPLOT")
	fmt.Println("----")
	for _, line := range plot {
		fmt.Println(line)
	}

	fmt.Println("LINKS")
	fmt.Println("----")
	fmt.Println("https://en.wikipedia.org/wiki/" + urlEpisodeName)
	fmt.Println("https://theinfosphere.org/" + urlEpisodeName)
	fmt.Println("https://futurama.fandom.com/wiki/" + urlEpisodeName)
}
