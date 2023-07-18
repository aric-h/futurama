package cmd

import "fmt"

type SeasonEpisodes struct {
	name     string
	episodes []string
}

type Season struct {
	name     string
	episodes []Episode
}

type Episode struct {
	name   string
	quotes [][]string
}

func getSeries() [7]SeasonEpisodes {
	series := [7]SeasonEpisodes{
		{
			name: "Season 1",
			episodes: []string{
				"Space Pilot 3000",
				"The Series Has Landed",
				"I, Roommate",
				"Love's Labors Lost in Space",
				"Fear of a Bot Planet",
				"A Fishful of Dollars",
				"My Three Suns",
				"A Big Piece of Garbage",
				"Hell Is Other Robots",
				"A Flight to Remember",
				"Mars University",
				"When Aliens Attack",
				"Fry and the Slurm Factory",
			},
		},
		{
			name: "Season 2",
			episodes: []string{
				"I Second That Emotion",
				"Brannigan, Begin Again",
				"A Head in the Polls",
				"Xmas Story",
				"Why Must I Be a Crustacean in Love?",
				"Lesser of Two Evils",
				"Put Your Head on my Shoulders",
				"Raging Bender",
				"A Bicyclops Built For Two",
				"A Clone of My Own",
				"How Hermes Requisitioned His Groove Back",
				"The Deep South",
				"Bender Gets Made",
				"Mother's Day",
				"The Problem With Popplers",
				"Anthology of Interest I",
				"War Is the H-Word",
				"The Honking",
				"The Cryonic Woman",
			},
		},
		{
			name: "Season 3",
			episodes: []string{
				"Amazon Women in the Mood",
				"Parasites Lost",
				"A Tale of Two Santas",
				"The Luck of the Fryrish",
				"The Birdbot of Ice-Catraz",
				"Bendless Love",
				"The Day the Earth Stood Stupid",
				"That's Lobstertainment",
				"The Cyber House Rules",
				"Where the Buggalo Roam",
				"Insane in the Mainframe",
				"The Route of All Evil",
				"Bendin' in the Wind",
				"Time Keeps on Slippin'",
				"I Dated a Robot",
				"A Leela of Her Own",
				"A Pharaoh to Remember",
				"Anthology of Interest II",
				"Roswell That Ends Well",
				"Godfellas",
				"Future Stock",
				"The 30% Iron Chef",
			},
		},
		{
			name: "Season 4",
			episodes: []string{
				"Kif Gets Knocked Up A Notch",
				"Leela's Homeworld",
				"Love and Rocket",
				"Less Than Hero",
				"A Taste of Freedom",
				"Bender Should Not Be Allowed On TV",
				"Jurassic Bark",
				"Crimes of the Hot",
				"Teenage Mutant Leela's Hurdles",
				"The Why of Fry",
				"Where No Fan Has Gone Before",
				"The Sting",
				"Bend Her",
				"Obsoletely Fabulous",
				"The Farnsworth Parabox",
				"Three Hundred Big Boys",
				"Spanish Fry",
				"The Devil's Hands are Idle Playthings",
			},
		},
		{
			name: "Season 5",
			episodes: []string{
				"Bender's Big Score",
				"The Beast with a Billion Backs",
				"Bender's Game",
				"Into the Wild Green Yonder",
			},
		},
		{
			name: "Season 6",
			episodes: []string{
				"Rebirth",
				"In-A-Gadda-Da-Leela",
				"Attack of the Killer App",
				"Proposition Infinity",
				"The Duh-Vinci Code",
				"Lethal Inspection",
				"The Late Philip J. Fry",
				"That Darn Katz!",
				"A Clockwork Origin",
				"The Prisoner of Benda",
				"Lrrreconcilable Ndndifferences",
				"The Mutants Are Revolting",
				"The Futurama Holiday Spectacular",
				"Neutopia",
				"Benderama",
				"Ghost in the Machines",
				"Law and Oracle",
				"The Silence of the Clamps",
				"Yo Leela Leela",
				"All the Presidents' Heads",
				"Möbius Dick",
				"Fry Am the Egg Man",
				"The Tip of the Zoidberg",
				"Cold Warriors",
				"Overclockwise",
				"Reincarnation",
			},
		},
		{
			name: "Season 7",
			episodes: []string{
				"The Bots and the Bees",
				"A Farewell to Arms",
				"Decision 3012",
				"The Thief of Baghead",
				"Zapp Dingbat",
				"The Butterjunk Effect",
				"The Six Million Dollar Mon",
				"Fun on a Bun",
				"Free Will Hunting",
				"Near-Death Wish",
				"31st Century Fox",
				"Viva Mars Vegas",
				"Naturama",
				"2-D Blacktop",
				"Fry and Leela's Big Fling",
				"T.: The Terrestrial",
				"Forty Percent Leadbelly",
				"The Inhuman Torch",
				"Saturday Morning Fun Pit",
				"Calculon 2.0",
				"Assie Come Home",
				"Leela and the Genestalk",
				"Game of Tones",
				"Murder on the Planet Express",
				"Stench and Stenchibility",
				"Meanwhile",
				"Simpsons Crossover: Simpsorama",
			},
		},
	}

	return series
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
