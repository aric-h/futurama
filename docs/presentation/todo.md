# TODO

## Week of July 17

- [x] ~~add universal flag for generating Lambda output~~
- [x] add print function
- [x] revise validation function
- [x] add validation rule for character input
- [x] revise randomization function
- [x] update Episode struct
  - [x] add new [Quote struct](#quote-struct)
  - [x] add function for getting characters from quote lines
- [x] revise getSeasonQuotes function
  - [x] clean up comments
  - [x] move basic token-parsing logic to function
  - [x] add function for getting quote characters
- [x] add `version` command to display ascii art
- [ ] Update `describe episode` to pull from Infosphere
  - [ ] consider adding other links and info

### quote struct

```go
type Quote struct {
    Characters []string
    Quote [string]
}
```

## Week of July 24

- [ ] finish slide deck
- [ ] prep live demo
- [x] build basic GHA automation
- [x] add screenshot of finished product to slide 7
- [ ] populate readme
- [x] ~~Update `describe episode` to pull from Infosphere~~
  - [x] consider adding other links and info
- [x] revert back to wikipedia for `describe episode`; add links to wikipedia, infosphere. and fandom
- [x] figure out logic for character validation
- [x] figure out logic for validating combinations
  - [x] consider making --all a hidden flag that can only be set with the episode flag
- [x] update `getEpisodeQuotes` function to normalize character names before appending to character list
- [x] test turning app into lambda
  - [x] or test separate, basic lambda for random quote generation
- [ ] add `describe character` command with links to wiki, infosphere, and fandom
- [x] add Zapp Brannigan to supported characters
- [ ] add colorization to main characters' names in quote output
- [x] overhaul repo structure
- [ ] add retry functionality if WikiQuote 404s
- [ ] update `get quote` to return season and call prinQuote separately
- [ ] create basic TF for API gateway and lambda

## July 31 & August 1

- [x] add retry functionality if WikiQuote 404s
- [ ] update `get quote` to return season and call prinQuote separately
- [ ] create basic TF for API gateway and lambda
- [ ] create slack app and request installation rights in BSC workspace
- [ ] add colorization to main characters' names in quote output
