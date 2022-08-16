<div align="right">

![CI](https://github.com/sheepla/srss/actions/workflows/ci.yml/badge.svg)
![Relase](https://github.com/sheepla/srss/actions/workflows/release.yml/badge.svg)

</div>

<div align="center">

# 📘 srss

</div>

<div align="center">

A fast & simple command line RSS/ATOM/JSON feed reader written in Go, inspired by [newsboat](https://github.com/newsboat/newsboat)

![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/sheepla/srss?style=flat-square)](https://github.com/sheepla/srss/releases/latest)

*This repository is still under development!. Specifications are subject to change without notice.*

</div>

## Features

- Fast, efficient, and easy-to-use interface for CLI lovers
- Supports multiple feed format: RSS, Atom and JSON
- You can import a file in OPML format and register URL entries

## Demo

<div align="center">

![demo](https://user-images.githubusercontent.com/62412884/184543394-b79df2de-e8ef-4812-a767-7b3a7d26e746.gif)

</div>

## Usage

### Commands and Options

```
NAME:
   srss - A simple command line RSS feed reader

USAGE:
   srss [global options] command [command options] [arguments...]

VERSION:
   0.0.3-alpha

COMMANDS:
   add, a     Add url entry
   edit, e    Edit URL entry file
   tui, t     View items in the feed with built-in pager
   open, o    Open feed URL on your browser
   import, i  Import Feed URL from OPML file
   update, u  Fetch the latest feeds and update the cache
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Register or edit the feeds URL

Use the `add` command to register the feed URL.

The feeds URL is saved in a plain text file and you can edit it using the `edit` command.
You can specify the command name of the editor in the argument of the `-e`, `--editor` option.
If the environment variable `$EDITOR` is set, will use it.

```bash
srss add https://zenn.dev/topics/go/feed
srss edit --editor nvim
```

*NOTE*

The location of the URL entry file depends on the OS. It is as follows:

|OS     |Path                                                                        |
|-------|----------------------------------------------------------------------------|
|Windows|`%APPDATA%\srss\urls.txt` or `C:\Users\%USER%\AppData\Roaming\srss\urls.txt`|
|Linux  |`$XDG_CONFIG_HOME/srss/urls.txt` or `$HOME/.config/srss/urls.txt`           |
|macOS  |`$HOME/Library/Application Support/srss/urls.txt`                           |

### Fetch the feeds and update cache

Run the `update`, `u` command before view the feeds.
The feed is fetched from the URL described in the URL entry file and the cache is updated.

```
srss update
```

*NOTE*

The location of the cache file depends on the OS. It is as follows:

|OS     |Path                                                                            |
|-------|--------------------------------------------------------------------------------|
|Windows|`%LOCALAPPDATA%\srss\cache.gob`or `C:\Users\%USER%\AppData\Local\srss\cache.gob`|
|Linux  |`$XDG_CACHE_HOME/srss/cache.gob`or `$HOME/.cache/srss/cache.gob`                |
|macOS  |`$HOME/Library/Caches/srss/cache.gob`                                           |
  
### View items in the feed on the terminal

Run the `tui`, `t` command then narrow down and select the items in the feed with a fuzzyfinder-like UI,
you can browse the items with a `less` like pager UI.

```
srss tui
```

The key bindings in fuzzyfinder UI are follows:

|Key        |Description     |
|-----------|----------------|
|`C-k` `C-p`|Move focus up   |
|`C-j` `C-n`|Move focus down |
|`Enter`    |Select the item |
|`q` `Esc`  |Quit fuzzyfinder|

The key bindings in pager UI are follows:

|Key       |Description                        |
|----------|-----------------------------------|
|`k` `Up`  |Scroll up                          |
|`j` `Down`|Scroll down                        |
|`g` `Home`|Scroll on top                      |
|`G` `End` |Scroll on bottom                   |
|`q` `Esc` |Quit pager then back to fuzzyfinder|

### Open links on items in the feed in the browser

Use the `open`, `o` command, you can open the link of the selected item in your browser.
You can select multiple items with `Tab` key.

```
srss open
```

### Import feeds URL from OPML file

Use the `import`, `i` command, you can import a file in [OPML](https://en.wikipedia.org/wiki/OPML) format and register feeds URL.

```
srss import --path path/to/file.opml
```

## Installation

Executable binaries are available from the latest release page.

> [![Latest Release](https://img.shields.io/github/v/release/sheepla/srss?style=flat-square)](https://github.com/sheepla/srss/releases/latest)

To build from source, clone this repository then run `go install`. 
Developing with Go `v1.18.2 linux/amd64`

## Contributing

Welcome any bug reports, requests, typo fixes, etc.

## LICENSE

[MIT](./LICENSE)

## Author

[Sheepla](https://github.com/sheepla)

