package ui

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/mmcdole/gofeed"
)

func FindItem(items []*gofeed.Item) ([]int, error) {
	return fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			return items[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}
			s := fmt.Sprintf(
				"■ %s\n\n  by %s\n  published at %s, updated at %s\n\n%s\n",
				items[i].Title,
				items[i].Author.Name,
				items[i].Published,
				items[i].Updated,
				items[i].Description,
			)
			return runewidth.Wrap(s, width/2-5)
		}),
	)
}
