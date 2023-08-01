# Futurama CLI

## What is it?

The futurama CLI is a utility for retrieving random quotes and plot descriptions from the critically-acclaimed animated series *Futurama*.

## Available commands

### `get quote`

Get random Futurama quote from WikiQuote

Available flags:

- `--season`, `-s` - int - Season number (1-7)
- `--episode`, `-e` - string - Episode name (use 'futurama get episodes' command for assistance)
- `--all`, `a` - Toggle for returning all quotes from an episode
- `--character`, `-c` - string - Character name (e.g. 'Fry', 'Bender')

### `get episodes`

Get list of episode names

Available flags:

- `--season`, `-s` - int - Season number (1-7)
- `--all`, `a` - Toggle for returning all episodes from the entire series
  
### `get characters`

Get list of supported characters for the `get quote --character` flag

### `describe episode`

Describe plot of a Futurama episode

Required flag:

- `--name`, `-n` - string - Episode name (use 'futurama get episodes' command for assistance)

## Installation

If you have Go installed:

```bash
go install github.com/aric-h/futurama@latest
```

If you *do not* have Go installed:

Visit the [release page](https://github.com/aric-h/futurama/releases) and download the appropriate binary for your OS.
