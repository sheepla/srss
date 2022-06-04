package ui

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/mmcdole/gofeed"
)

const padding = 5

// nolint:wrapcheck
func FindItem(items []*gofeed.Item) (int, error) {
	return fuzzyfinder.Find(
		items,
		func(i int) string {
			return fmt.Sprintf("%s [%s]", items[i].Title, humanizeTime(items[i].PublishedParsed))
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			return runewidth.Wrap(renderPreviewWindow(items[i]), width/2-padding)
		}),
	)
}

// nolint:wrapcheck
func FindItemMulti(items []*gofeed.Item) ([]int, error) {
	return fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			return fmt.Sprintf("%s [%s]", items[i].Title, humanizeTime(items[i].PublishedParsed))
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}

			return runewidth.Wrap(renderPreviewWindow(items[i]), width/2-padding)
		}),
	)
}
