/*
Copyright © 2023 Aric Hansen <aric.p.hansen@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mpvl/unique"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
  - a random episode in a random season from a user-defined character`,
	Example: `  futurama get-quote (no flags = randomized season and episode)
  futurama get quote --season 2
  futurama get quote --episode "Space Pilot 3000"
  futurama get quote --character "Fry"`,
	Run: func(cmd *cobra.Command, args []string) {
		err := validateInput(cmd.Flags())
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			cmd.Help()
		} else {
			// quoteLoop()
			fmt.Println(cmd.Flags().NFlag())
		}
	},
}

func init() {
	getCmd.AddCommand(quoteCmd)
	quoteCmd.Flags().IntVarP(&QuoteSeason, "season", "s", 0, "Season number (1-7)")
	quoteCmd.Flags().StringVarP(&QuoteEpisode, "episode", "e", "", "Episode name (use 'futurama get episodes' command for assistance)")
	quoteCmd.Flags().StringVarP(&QuoteCharacter, "character", "c", "", "Character name (e.g. 'Fry', 'Bender')")
	quoteCmd.Flags().BoolVarP(&AllQuotes, "all", "a", false, "Toggle for returning all quotes from a season or episode")
	quoteCmd.MarkFlagsMutuallyExclusive("season", "episode")
}

func validateFlags(flags *pflag.FlagSet) error {
	// check if multiple flags are set
	fmt.Println(&flags)

	return nil
}

func validateInput(flags *pflag.FlagSet) error {
	// validate number of flags
	err := validateFlags(flags)
	if err != nil {
		return err
	}

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
		err, QuoteSeason = validateEpisodeName(QuoteEpisode)
		if err != nil {
			return err
		}
	}

	// validate character input
	if QuoteCharacter != "" {
		supportedCharacters := getSupportedCharacters()

		invalidCharacter := true
		for _, c := range supportedCharacters {
			if strings.ToLower(QuoteCharacter) == strings.ToLower(c) {
				invalidCharacter = false
				break
			}
		}

		if invalidCharacter {
			return errors.New("Invalid character input. Please use the 'futurama get characters' command for assistance.")
		}

	}

	return nil
}

// func quoteLoop() {
// 	// capture initial season and episode input in case user-defined character can't be found in randomized episode
// 	initialEpisodeInput := QuoteEpisode
// 	initialSeasonInput := QuoteSeason

// quoteLoop:
// 	for {
// 		randomize()
// 		season := getQuotes()
// 		if QuoteCharacter != "" {
// 			err := isCharacterPresent(season)
// 		}
// 	}

// }

func isCharacterPresent(season Season) error {
	for _, ep := range season.episodes {
		for _, quote := range ep.quotes {
			for _, character := range quote.characters {
				if character == QuoteCharacter {
					return nil
				}
			}
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

func getQuotes() Season {
	var season Season
	var resp *http.Response

	if QuoteSeason == 5 {
		urlEpisode := "_" + strings.Replace(strings.Replace(QuoteEpisode, "'", "%27", -1), " ", "_", -1)
		resp = getHttpResponse("https://en.wikiquote.org/wiki/Futurama:" + urlEpisode)
		season = getSeasonFiveQuotes(resp)
	} else {
		resp = getHttpResponse("https://en.wikiquote.org/wiki/Futurama/Season_" + strconv.Itoa(QuoteSeason))
		season = getSeasonQuotes(resp)
	}

	defer resp.Body.Close()

	// printQuotes(season)
	// fmt.Println(season)

	return season
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

	// tokenize WikiQuote response
	tokenizer := html.NewTokenizer(resp.Body)

seasonLoop:
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break seasonLoop //end of the file, break out of the loop
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		case html.EndTagToken:
			if "ul" == tokenizer.Token().Data { // end of episode list; start episodes/quotes section
				for {
					switch tokenizer.Next() {
					case html.ErrorToken:
						err := tokenizer.Err()
						if err == io.EOF {
							break seasonLoop //end of the file, break out of the loop
						}
						log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
					case html.StartTagToken:
						token := tokenizer.Token()
						if "span" == token.Data { // found episode title line
							for _, attr := range token.Attr {
								if attr.Val == "External_links" { // reached end of quote page
									break seasonLoop
								}
							}
							episode := getEpisodeName(tokenizer)
							season.episodes = append(season.episodes, episode)
						}
					}
				}
			}
		}
	}

	return season
}

func getEpisodeName(tokenizer *html.Tokenizer) Episode {
	// initialize episode var
	ep := Episode{}

findEpisodeName:
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break findEpisodeName //end of the file, break out of the loop
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		case html.TextToken:
			ep.name = tokenizer.Token().Data
			break findEpisodeName
		}
	}
	ep.quotes = getEpisodeQuotes(tokenizer)
	// fmt.Println(ep)
	return ep
}

func getEpisodeQuotes(tokenizer *html.Tokenizer) []Quote {
	episodeQuotes := []Quote{}

findNextQuote:
	for {
		quote := Quote{}
	getQuoteLines:
		for {
			switch tokenizer.Next() {
			case html.ErrorToken:
				err := tokenizer.Err()
				if err == io.EOF {
					break findNextQuote //end of the file, break out of the loop
				}
				log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			case html.StartTagToken:
				switch tokenizer.Token().Data {
				case "dl", "dd": // start of quote line
					line := ""
					speaker := false
				getQuoteLine:
					for {
						switch tokenizer.Next() {
						case html.ErrorToken:
							err := tokenizer.Err()
							if err == io.EOF {
								break findNextQuote //end of the file, break out of the loop
							}
							log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
						case html.StartTagToken:
							if "b" == tokenizer.Token().Data { // bolded speaker of quote line
								speaker = true
							}
						case html.TextToken:
							token := tokenizer.Token()
							line = line + token.Data
							if speaker {
								quote.characters = append(quote.characters, token.Data)
								speaker = false
							}
						case html.EndTagToken:
							if "dd" == tokenizer.Token().Data { // end of quote line
								quote.lines = append(quote.lines, line)
								break getQuoteLine
							}
						}
					}
				case "h2", "h3": // start of new episode or end of quote section
					episodeQuotes = append(episodeQuotes, quote)
					break findNextQuote
				}

			case html.SelfClosingTagToken:
				if "hr" == tokenizer.Token().Data { // line break between quotes
					episodeQuotes = append(episodeQuotes, quote)
					break getQuoteLines
				}
			}
		}
	}

	for _, x := range episodeQuotes {
		unique.Sort(unique.StringSlice{&x.characters})
		unique.Strings(&x.characters)
		fmt.Println(x.characters)
		fmt.Println(x.lines)
	}
	return episodeQuotes
}

func getSeasonFiveQuotes(resp *http.Response) Season {
	var season = Season{name: "Season " + strconv.Itoa(QuoteSeason)}
	var ep = Episode{name: QuoteEpisode}

	// tokenize WikiQuote response
	tokenizer := html.NewTokenizer(resp.Body)

episodeLoop:
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break episodeLoop //end of the file, break out of the loop
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		case html.StartTagToken:
			token := tokenizer.Token()
			if "span" == token.Data {
				for _, attr := range token.Attr {
					if attr.Val == "Dialogue" { // start parsing quotes
						ep.quotes = getEpisodeQuotes(tokenizer)
						break episodeLoop
					}
				}
			}
		}
	}

	season.episodes = append(season.episodes, ep)
	return season
}

func printQuotes(season Season) {
	fmt.Print("Season: ")
	fmt.Println(QuoteSeason)

	if AllQuotes {
		if QuoteEpisode != "" {
			fmt.Print("Episode: ")
			fmt.Println(QuoteEpisode)
			if QuoteCharacter != "" {
				// print all quotes from character in episode
			} else {
				// print all quotes from episode
			}
		} else {
			if QuoteCharacter != "" {
				// print all quotes from character in season
			} else {
				// print all quotes from season
			}
		}
	} else {
		fmt.Print("Episode: ")
		fmt.Println(QuoteEpisode)
	}
}
