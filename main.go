package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/kirsle/configdir"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mmcdole/gofeed"
	"github.com/toqueteos/webbrowser"
	"github.com/urfave/cli/v2"
)

const (
	appName = "srss"
)

var urlFile = path.Join(configdir.LocalConfig(), appName, "urls.txt")

// var cacheDBFile = path.Join(configdir.LocalCache(), appName, "cache.db")

func main() {
	app := initApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func initApp() cli.App {
	return cli.App{
		Name:  "srss",
		Usage: "A simple command line RSS feed reader",
		Action: func(ctx *cli.Context) error {
			return nil
		},
		// Flags: []cli.Flag{
		// 	&cli.BoolFlag{
		// 		Name:    "version",
		// 		Aliases: []string{"V", "v"},
		// 		Usage:   "Show version",
		// 	},
		// },
		Commands: []*cli.Command{
			{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open feed URL on your browser",
				Action: func(ctx *cli.Context) error {
					urls, err := readURLsFromEntry()
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

					choises, err := findItem(items)
					if err != nil {
						return err
					}
					for _, idx := range choises {
						if err := openURL(items[idx].Link); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add url entry",
				Action: func(ctx *cli.Context) error {
					url := ctx.Args().Get(0)
					return addURLEntry(url)
				},
			},
			// {
			// 	Name:    "edit",
			// 	Aliases: []string{"o"},
			// 	Usage:   "Edit URL resource file",
			// },
		},
	}
}

func fetchFeed(url string) (*gofeed.Feed, error) {
	feed, err := gofeed.NewParser().ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch or parse feed at %s: %w", url, err)
	}
	return feed, nil
}

func addURLEntry(url string) error {
	file, err := os.OpenFile(urlFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("Failed to open URL entry file (%s): %w", urlFile, err)
	}
	defer file.Close()
	_, err = fmt.Fprintln(file, url)
	if err != nil {
		return fmt.Errorf("Writing failed to the URL entry file (%s): %w", urlFile, err)
	}
	return nil
}

func readURLsFromEntry() ([]string, error) {
	var urls []string
	file, err := os.OpenFile(urlFile, os.O_RDONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("Failed to open URL entry file (%s): %w", urlFile, err)
	}
	s := bufio.NewScanner(file)
	for s.Scan() {
		urls = append(urls, s.Text())
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("Failed to scan from URL entry file (%s): %w", urlFile, err)
	}
	return urls, nil
}

func findItem(items []*gofeed.Item) ([]int, error) {
	return fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			return items[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf(
				"%s\n\n%s\n\n%s\n",
				items[i].Title,
				items[i].Description,
				items[i].Content,
			)
		}),
	)
}

func openURL(url string) error {
	if err := webbrowser.Open(url); err != nil {
		return fmt.Errorf("Failed to open the URL (%s): %w", url, err)
	}
	return nil
}
