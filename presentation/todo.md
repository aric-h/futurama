# TODO

## Week of July 17

- [ ] add universal flag for generating Lambda output
- [ ] add print function
- [ ] revise randomization function
- [ ] update Episode struct
  - [ ] add new [Quote struct](#quote-struct)
  - [ ] add function for getting characters from quote lines
- [ ] revise getSeasonQuotes function
  - [ ] clean up comments
  - [ ] move basic token-parsing logic to function
  - [ ] add function for getting quote characters
- [ ] add `version` command to display ascii art
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
- [ ] build basic GHA automation
