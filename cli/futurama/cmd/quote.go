/*
Copyright © 2023 Aric Hansen <aric.p.hansen@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

// vars for storing flag input
var QuoteSeason int
var QuoteEpisode string
var QuoteCharacter string
var AllQuotes bool

// quoteCmd represents the quote command
var quoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "Get random Futurama quote",
	Long: `Get a Futurama quote from:
  - a random episode in a random season
  - a random episode in a user-defined season
  - a user-defined episode
  
Or get all quotes from user-defined season or episode `,
	Example: `  futurama get-quote (no flags = randomized season and episode)
  futurama get quote --season 2
  futurama get quote --episode "Space Pilot 3000"`,
	Run: func(cmd *cobra.Command, args []string) {
		err := validateInput()
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			cmd.Help()
		} else {
			randomize()
			getQuote()
		}
	},
}

func init() {
	getCmd.AddCommand(quoteCmd)
	quoteCmd.Flags().IntVarP(&QuoteSeason, "season", "s", 0, "Season number (1-7)")
	quoteCmd.Flags().StringVarP(&QuoteEpisode, "episode", "e", "", "Episode name (use 'futurama get episodes' command for assistance)")
	quoteCmd.Flags().StringVarP(&QuoteCharacter, "character", "c", "", "Character name (e.g. 'Fry', 'Bender')")
	quoteCmd.Flags().BoolVarP(&AllQuotes, "all", "a", false, "Toggle for returning all quotes from a season or episode")
}

func validateInput() error {
	// validate season number
	invalidSeason := true
	for i := 0; i < 8; i++ {
		if i == QuoteSeason {
			invalidSeason = false
		}
	}

	if invalidSeason {
		return errors.New("Invalid season number. Please select a value from 1-7.")
	}

	// validate episode name
	if QuoteEpisode != "" {
		series := getSeries()
		invalidEpisode := true
		for i, season := range series {
			for _, ep := range season.episodes {
				if ep == QuoteEpisode {
					if QuoteSeason == 0 { // set QuoteSeason if not provided
						QuoteSeason = i + 1
						invalidEpisode = false
					} else {
						if i+1 == QuoteSeason { // if user provides season and episode, make sure they match
							invalidEpisode = false
						} else {
							return errors.New("Season/episode mismatch. Please use the `futurama get episodes` command for assistance.")
						}
					}
				}
			}
		}

		if invalidEpisode {
			return errors.New("Invalid episode name. Please use the `futurama get episodes` command for assistance.")
		}
	}

	return nil
}

func randomize() {
	series := getSeries()

	// randomize season if no input
	if QuoteSeason == 0 && QuoteEpisode == "" {
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 7
		QuoteSeason = rand.Intn(max-min+1) + min
	}

	// randomize episode if not specified
	if QuoteEpisode == "" {
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := len(series[QuoteSeason-1].episodes) - 1
		randEpisodeIndex := rand.Intn(max-min+1) + min
		QuoteEpisode = series[QuoteSeason-1].episodes[randEpisodeIndex]
	}
}

func getQuote() error {
	var season Season
	var resp *http.Response

	if QuoteSeason == 5 {
		urlEpisode := "_" + strings.Replace(strings.Replace(QuoteEpisode, "'", "%27", -1), " ", "_", -1)
		resp = getHttpResponse("https://en.wikiquote.org/wiki/Futurama:" + urlEpisode)
		// season = getSeasonFiveQuotes(resp)
	} else {
		resp = getHttpResponse("https://en.wikiquote.org/wiki/Futurama/Season_" + strconv.Itoa(QuoteSeason))
		season = getSeasonQuotes(resp)
	}

	defer resp.Body.Close()

	printQuotes(season)

	return nil
}

func getHttpResponse(url string) *http.Response {
	resp, err := http.Get(url)

	if err != nil {
		//.Fatalf() prints the error and exits the process
		// return errors.Newf("error fetching URL: %v\n", err)
		log.Fatalf("Error fetching WikiQuote URL: %v\n", err)
	}

	//check response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("WikiQuote response status code was %d\n", resp.StatusCode)
	}

	//check response content type
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("WikiQuote response content type was %s, not text/html\n", ctype)
	}

	return resp
}

func getSeasonQuotes(resp *http.Response) Season {
	var season = Season{name: "Season " + strconv.Itoa(QuoteSeason)}

	tokenizer := html.NewTokenizer(resp.Body)

	return season
}
