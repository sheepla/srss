package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kirsle/configdir"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mmcdole/gofeed"
	"github.com/sheepla/srss/ui"
	"github.com/toqueteos/webbrowser"
	"github.com/urfave/cli/v2"
)

// nolint:gochecknoglobals
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
	exitCodeErrEditor
	exitCodeErrBrowser
)

// nolint:gochecknoglobals
var urlFile = filepath.Join(configdir.LocalConfig(), appName, "urls.txt")

// var cacheDBFile = path.Join(configdir.LocalCache(), appName, "cache.db")

func main() {
	app := initApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// nolint:funlen,gocognit,cyclop,exhaustruct,exhaustivestruct,maintidx
func initApp() *cli.App {
	return &cli.App{
		Name:                 appName,
		HelpName:             appName,
		Version:              appVersion,
		Usage:                appUsage,
		Suggest:              false,
		EnableBashCompletion: true,
		Before: func(ctx *cli.Context) error {
			if err := configdir.MakePath(filepath.Dir(urlFile)); err != nil {
				return fmt.Errorf("failed to create URL entry file: %w", err)
			}

			return nil
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return cli.Exit("must require arguments", int(exitCodeErrArgs))
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add url entry",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() != 1 {
						return cli.Exit(
							"requires URL as an argument",
							int(exitCodeErrArgs),
						)
					}
					url := ctx.Args().Get(0)
					if !isValidURL(url) {
						return cli.Exit(
							fmt.Sprintf("invalid URL (%s)", url),
							int(exitCodeErrURLEntry),
						)
					}
					urls, err := readURLEntry()
					if err != nil {
						return cli.Exit(
							fmt.Sprintf("failed to read URL entry file (%s)", url),
							int(exitCodeErrURLEntry),
						)
					}
					if !isUniqueURL(urls, url) {
						return cli.Exit(
							fmt.Sprintf("the URL is already registered (%s)", url),
							int(exitCodeErrURLEntry),
						)
					}
					if err := addURLEntry(url); err != nil {
						return cli.Exit(
							fmt.Sprintf("failed to add URL entry (%s): %s", url, err),
							int(exitCodeErrURLEntry),
						)
					}

					return cli.Exit("", int(exitCodeOK))
				},
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
				Action: func(ctx *cli.Context) error {
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
					err := execEditor(editor, urlFile)
					if err != nil {
						return cli.Exit(
							fmt.Sprintf("failed to launch editor: %s", err),
							int(exitCodeErrEditor),
						)
					}

					return cli.Exit("", int(exitCodeOK))
				},
			},
			{
				Name:    "tui",
				Aliases: []string{"t"},
				Usage:   "View items in the feed with built-in pager",
				Action: func(ctx *cli.Context) error {
					urls, err := readURLEntry()
					if err != nil {
						return err
					}
					var feeds []gofeed.Feed
					for _, v := range urls {
						f, err := fetchFeed(v)
						if err != nil {
							return err
						}
						feeds = append(feeds, *f)
					}

					var items []*gofeed.Item
					for i := 0; i < len(feeds); i++ {
						items = append(items, feeds[i].Items...)
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

						// nolint:staticcheck
						return cli.Exit("", int(exitCodeOK))
					}
				},
			},
			{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open feed URL on your browser",
				Action: func(ctx *cli.Context) error {
					urls, err := readURLEntry()
					if err != nil {
						return err
					}
					var feeds []gofeed.Feed
					for _, v := range urls {
						feed, err := fetchFeed(v)
						if err != nil {
							return cli.Exit(
								fmt.Sprintf("failed to fetch feeds: %s", err),
								int(exitCodeErrFetchFeeds),
							)
						}
						feeds = append(feeds, *feed)
					}

					var items []*gofeed.Item
					for i := 0; i < len(feeds); i++ {
						items = append(items, feeds[i].Items...)
					}

					choises, err := ui.FindItemMulti(items)
					if err != nil {
						return cli.Exit(
							fmt.Sprintf("an error occurred on fuzzyfinder: %s", err),
							int(exitCodeErrFuzzyFinder),
						)
					}
					for _, idx := range choises {
						if err := openURL(items[idx].Link); err != nil {
							return cli.Exit(
								fmt.Sprintf("failed to open URL in browser: %s", err),
								int(exitCodeErrBrowser),
							)
						}
					}

					return cli.Exit("", int(exitCodeOK))
				},
			},
		},
	}
}

func fetchFeed(url string) (*gofeed.Feed, error) {
	feed, err := gofeed.NewParser().ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch or parse feed at %s: %w", url, err)
	}

	return feed, nil
}

func isUniqueURL(urls []string, u string) bool {
	for _, v := range urls {
		if v == u {
			return false
		}
	}

	return true
}

// nolint:wsl
func addURLEntry(url string) error {
	// nolint:gomnd
	file, err := os.OpenFile(urlFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open URL entry file (%s): %w", urlFile, err)
	}
	defer file.Close()
	if err != nil {
		return fmt.Errorf("writing failed to the URL entry file (%s): %w", urlFile, err)
	}

	return nil
}

// nolint:wsl
func readURLEntry() ([]string, error) {
	var urls []string
	// nolint:gomnd
	file, err := os.OpenFile(urlFile, os.O_RDONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to open URL entry file (%s): %w", urlFile, err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		urls = append(urls, s.Text())
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("failed to scan from URL entry file (%s): %w", urlFile, err)
	}

	return urls, nil
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)

	return err == nil
}

func openURL(url string) error {
	if err := webbrowser.Open(url); err != nil {
		return fmt.Errorf("failed to open the URL (%s): %w", url, err)
	}

	return nil
}

// nolint:wsl
// https://doloopwhile.hatenablog.com/entry/2014/08/05/213819
func execEditor(editor string, args ...string) error {
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run editor (%s) %w", editor, err)
	}

	return nil
}
