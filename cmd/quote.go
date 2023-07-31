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
  - a random episode in a random season from a user-defined character
  
  Or get all quotes from a user-defined episode. `,
	Example: `  futurama get-quote (no flags = randomized season and episode)
  futurama get quote --season 2
  futurama get quote --episode "Space Pilot 3000"
  futurama get quote --character "Fry"
  futurama get quote --all --episode "The Series Has Landed"`,
	Run: func(cmd *cobra.Command, args []string) {
		err := validateInput(cmd.Flags())
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			cmd.Help()
		} else {
			randomize()
			season := getQuotes()
			printQuotes(season)
		}
	},
}

func init() {
	getCmd.AddCommand(quoteCmd)
	quoteCmd.Flags().IntVarP(&QuoteSeason, "season", "s", 0, "Season number (1-7)")
	quoteCmd.Flags().StringVarP(&QuoteEpisode, "episode", "e", "", "Episode name (use 'futurama get episodes' command for assistance)")
	quoteCmd.Flags().StringVarP(&QuoteCharacter, "character", "c", "", "Character name (e.g. 'Fry', 'Bender')")
	quoteCmd.Flags().BoolVarP(&AllQuotes, "all", "a", false, "Toggle for returning all quotes from an episode")
	// limit flag combos
	quoteCmd.MarkFlagsMutuallyExclusive("season", "episode")
	quoteCmd.MarkFlagsMutuallyExclusive("season", "all")
	quoteCmd.MarkFlagsMutuallyExclusive("season", "character")
	quoteCmd.MarkFlagsMutuallyExclusive("character", "all")
	quoteCmd.MarkFlagsMutuallyExclusive("character", "episode")
}

func validateInput(flags *pflag.FlagSet) error {
	var err error

	// validate season number
	invalidSeason := true
	for i := 0; i < 8; i++ {
		if i == QuoteSeason {
			invalidSeason = false
		}
	}

	if QuoteSeason == 8 {
		return errors.New("Season 8 compatibility coming soon! Please select a value from 1-7.")
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

	// validate --all is set with --episode
	if AllQuotes && QuoteEpisode == "" {
		return errors.New("The --all flag must be set with the --episode flag.")
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
	// if character is specified, this will be re-randomized later
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

	return season
}

func getHttpResponse(url string) *http.Response {
	var resp *http.Response
	var err error

	for i := 0; i < 5; i++ { // retry in case of bad response (mainly 404)
		resp, err = http.Get(url)
		ctype := resp.Header.Get("Content-Type")

		if err == nil && resp.StatusCode == http.StatusOK && strings.HasPrefix(ctype, "text/html") {
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
	}

	if err != nil {
		//.Fatalf() prints the error and exits the process
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
	// fmt.Println("\n" + ep.name)
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
								character := normalizeName(token.Data)
								quote.characters = append(quote.characters, character)
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
		// fmt.Println(x.characters)
		// fmt.Println(x.lines)
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
	var ep Episode

	fmt.Print("Season: ")
	fmt.Println(QuoteSeason)

	// find and print quote from character
	if QuoteCharacter != "" {
		// get subset of episodes with character present
		subset := getCharacterEpisodes(season)

		// re-randomize episode
		epIndex := randomIndex(len(subset.episodes) - 1)
		QuoteEpisode = subset.episodes[epIndex].name

		fmt.Print("Episode: ")
		fmt.Println(QuoteEpisode)
		fmt.Println()

		// get random quote
		qIndex := randomIndex(len(subset.episodes[epIndex].quotes) - 1)
		for _, line := range subset.episodes[epIndex].quotes[qIndex].lines {
			fmt.Println(line)
		}

	} else if AllQuotes { // print all quotes from an episode
		fmt.Print("Episode: ")
		fmt.Println(QuoteEpisode)
		fmt.Println()

		ep = getEpisodeObject(season)
		for _, q := range ep.quotes {
			for _, line := range q.lines {
				fmt.Println(line)
			}
			fmt.Println("----")
		}

	} else { // season/episode have either been set or randomized
		fmt.Print("Episode: ")
		fmt.Println(QuoteEpisode)
		fmt.Println()

		ep = getEpisodeObject(season)
		qIndex := randomIndex(len(ep.quotes) - 1)
		for _, line := range ep.quotes[qIndex].lines {
			fmt.Println(line)
		}
	}

}

func getCharacterEpisodes(season Season) Season {
	subset := Season{
		name: season.name,
	}

	// loop through episodes to find onew with QuoteCharacter
	for _, ep := range season.episodes {
		subsetEp := Episode{name: ep.name}
		for _, q := range ep.quotes {
			for _, n := range q.characters {
				if n == QuoteCharacter {
					subsetEp.quotes = append(subsetEp.quotes, q)
				}
			}
		}
		if len(subsetEp.quotes) > 0 {
			subset.episodes = append(subset.episodes, subsetEp)
		}
	}
	// fmt.Println(subset)
	return subset
}

func getEpisodeObject(season Season) Episode {
	for _, ep := range season.episodes {
		if ep.name == QuoteEpisode {
			return ep
		}
	}

	return Episode{
		name:   "Error",
		quotes: []Quote{{characters: []string{}, lines: []string{"error: no quotes found"}}},
	}
}

func randomIndex(max int) int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	randIndex := rand.Intn(max-min+1) + min

	return randIndex
}
