<div align="right">

![CI](https://github.com/sheepla/srss/actions/workflows/ci.yml/badge.svg)

</div>

<div align="center">

# 📘 srss

</div>

<div align="center">

A fast & simple command line RSS/ATOM/JSON feed reader

![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)

*This repository is still under development!. Specifications are subject to change without notice.*

</div>

## Features

- Fast, efficient, and easy-to-use interface for CLI lovers
- Supports multiple feed format: RSS, Atom and JSON

## Usage

### Commands and Options

```
NAME:
   srss - A simple command line RSS feed reader

USAGE:
   srss [global options] command [command options] [arguments...]

VERSION:
   unknown

COMMANDS:
   open, o  Open feed URL on your browser
   add, a   Add URL entry
   edit, e  Edit URL entry file
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Register or edit the feeds URL

Use the `add` command to register the feed URL.

The feeds URL is saved in a plain text file and you can edit it using the `edit` command.
You can specify the command name of the editor in the argument of the `-e`, `--editor` option.

```bash
srss add https://zenn.dev/topics/go/feed
srss edit --editor nvim
```

### View feed items on the terminal

*TODO*

### Open links on feed items in the browser

Use the `open`, `o` command to filter from feed items and open the URL of the selected item in your browser.
## Installation

To build from source, clone this repository then run `go install`. 
Developing with Go `v1.18.2 linux/amd64`

## Contributing

Welcome any bug reports, requests, typo fixes, etc.

## LICENSE

MIT

## Author

[Sheepla](https://github.com/sheepla)

