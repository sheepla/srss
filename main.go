package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mmcdole/gofeed"
	"github.com/sheepla/srss/cache"
	"github.com/sheepla/srss/opml"
	"github.com/sheepla/srss/ui"
	"github.com/sheepla/srss/urlentry"
	"github.com/urfave/cli/v2"
)

//nolint:gochecknoglobals
var (
	appName    = "srss"
	appVersion = "unknown"
	appUsage   = "A simple command line RSS feed reader"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrFetchFeeds
	exitCodeErrURLEntry
	exitCodeErrFuzzyFinder
	exitCodeErrPager
	exitCodeErrOPML
	exitCodeErrEditor
	exitCodeErrBrowser
	exitCodeErrCache
)

const asciiArt = `
   ─────────────┐   ┌────────────────────────────────┐
                │   │                                │
   ─────────┐   │   │ ┌─────  ─┬──── ┌─────  ┌─────  │
            │   │   │ │        │     │       │       │
   ─────┐   │   │   │ └─────┐  │     └─────┐ └─────┐ │
        │   │   │   │       │  │           │       │ │
┌───┐   │   │   │   │  ─────┘ ─┴─     ─────┘  ─────┘ │
│   │   │   │   │   │                                │
└───┘               └────────────────────────────────┘
`

func main() {
	app := initApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

//nolint:funlen,exhaustruct,exhaustivestruct
func initApp() *cli.App {
	return &cli.App{
		Name:                 appName,
		HelpName:             appName,
		Version:              appVersion,
		Usage:                appUsage,
		Suggest:              false,
		EnableBashCompletion: true,
		Before: func(ctx *cli.Context) error {
			return nil
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				//nolint:forbidigo
				fmt.Print(asciiArt)

				return cli.Exit("must require arguments", int(exitCodeErrArgs))
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add url entry",
				Action:  runAddCommand,
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "Edit URL entry file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "editor",
						Aliases: []string{"e"},
						Usage:   "Editor command to edit URL entry file",
						Value:   "vim",
						EnvVars: []string{"EDITOR"},
					},
				},
				Action: runEditCommand,
			},
			{
				Name:    "tui",
				Aliases: []string{"t"},
				Usage:   "View items in the feed with built-in pager",
				Action:  runTUICommand,
			},
			{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open feed URL on your browser",
				Action:  runOpenCommand,
			},
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "Import Feed URL from OPML file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Usage:   "OPML file path",
					},
				},
				Action: runImportCommand,
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Fetch the latest feeds and update the cache",
				Action:  runUpdateCommand,
			},
		},
	}
}

func runAddCommand(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return cli.Exit(
			"requires URL as an argument",
			int(exitCodeErrArgs),
		)
	}

	url := strings.TrimSpace(ctx.Args().Get(0))
	if !urlentry.IsUniqueURL(url) {
		return cli.Exit(
			fmt.Sprintf("the URL(%s) has already registered", url),
			int(exitCodeErrURLEntry),
		)
	}

	if err := urlentry.Add(url); err != nil {
		return cli.Exit(
			fmt.Sprintf("failed to add URL(%s) to entry file: %s", url, err),
			int(exitCodeErrURLEntry),
		)
	}

	return cli.Exit("", int(exitCodeOK))
}

func runEditCommand(ctx *cli.Context) error {
	if ctx.NArg() != 0 {
		return cli.Exit(
			fmt.Sprintf("extra arguments (%s)", ctx.Args().Slice()),
			int(exitCodeErrArgs),
		)
	}

	editor := strings.TrimSpace(ctx.String("editor"))
	if editor == "" {
		return cli.Exit(
			"requires editor command name",
			int(exitCodeErrArgs),
		)
	}

	if err := urlentry.OpenEditor(editor); err != nil {
		return cli.Exit(
			fmt.Sprintf("failed to launch editor: %s", err),
			int(exitCodeErrEditor),
		)
	}

	return cli.Exit("", int(exitCodeOK))
}

func runTUICommand(ctx *cli.Context) error {
	items, err := cache.Import()
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("failed to load cache: %s", err),
			int(exitCodeErrCache),
		)
	}

	for {
		idx, err := ui.FindItem(items)
		if err != nil {
			if errors.Is(fuzzyfinder.ErrAbort, err) {
				return cli.Exit(
					"quit",
					int(exitCodeOK),
				)
			}

			return cli.Exit(
				fmt.Sprintf("an error occurred on fuzzyfinder: %s", err),
				int(exitCodeErrFuzzyFinder),
			)
		}

		pager, err := ui.NewPager(items[idx])
		if err != nil {
			return cli.Exit(
				fmt.Sprintf("failed to init pager: %s", err),
				int(exitCodeErrPager),
			)
		}

		if err := pager.Start(); err != nil {
			return cli.Exit(
				fmt.Sprintf("an error occurred on pager: %s", err),
				int(exitCodeErrPager),
			)
		}
	}
}

func runOpenCommand(ctx *cli.Context) error {
	items, err := cache.Import()
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("failed to load cache: %s", err),
			int(exitCodeErrCache),
		)
	}

	choises, err := ui.FindItemMulti(items)
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("an error occurred on fuzzyfinder: %s", err),
			int(exitCodeErrFuzzyFinder),
		)
	}

	for _, idx := range choises {
		if err := ui.OpenURL(items[idx].Link); err != nil {
			return cli.Exit(
				fmt.Sprintf("failed to open URL in browser: %s", err),
				int(exitCodeErrBrowser),
			)
		}
	}

	return cli.Exit("", int(exitCodeOK))
}

func runImportCommand(ctx *cli.Context) error {
	if ctx.NArg() != 0 {
		return cli.Exit(
			fmt.Sprintf("extra arguments (%s)", ctx.Args().Slice()),
			int(exitCodeErrArgs),
		)
	}

	path := strings.TrimSpace(ctx.String("path"))

	if path == "" {
		return cli.Exit(
			"requires OPML file path",
			int(exitCodeErrArgs),
		)
	}

	outlines, err := opml.ParseOPML(path)
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("failed to parse OPML file (%s) %s", path, err),
			int(exitCodeErrOPML),
		)
	}

	urls := opml.ExtractFeedURL(outlines.Outlines())

	for _, url := range urls {
		if err := urlentry.Add(url); err != nil {
			return cli.Exit(
				fmt.Sprintf("failed to regester the URL entry (%s): %s", url, err),
				int(exitCodeErrURLEntry),
			)
		}
	}

	return cli.Exit("", int(exitCodeOK))
}

func runUpdateCommand(ctx *cli.Context) error {
	if ctx.NArg() != 0 {
		return cli.Exit(
			fmt.Sprintf("extra arguments (%s)", ctx.Args().Slice()),
			int(exitCodeErrArgs),
		)
	}

	urls, err := urlentry.Load()
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("failed to load URL entry: %s", err),
			int(exitCodeErrURLEntry),
		)
	}

	if len(urls) == 0 {
		return cli.Exit(
			"URL entry not registered",
			int(exitCodeErrURLEntry),
		)
	}

	//nolint:prealloc
	var feeds []gofeed.Feed

	for _, url := range urls {
		feed, err := fetchFeed(url)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf("failed to fetch the feeds: %s", err),
				int(exitCodeErrFetchFeeds),
			)
		}

		//nolint:forbidigo
		fmt.Printf("Fetched the feed: %s\n", url)

		feeds = append(feeds, *feed)
	}

	var items []*gofeed.Item
	for i := 0; i < len(feeds); i++ {
		items = append(items, feeds[i].Items...)
	}

	if err := cache.Export(items); err != nil {
		return fmt.Errorf("failed save the cache: %w", err)
	}

	return nil
}

func fetchFeed(url string) (*gofeed.Feed, error) {
	feed, err := gofeed.NewParser().ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch or parse feed at %s: %w", url, err)
	}

	return feed, nil
}
