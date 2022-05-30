package ui

import (
	"fmt"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/mmcdole/gofeed"
)

func bodyFromItem(item *gofeed.Item) string {
	var author string
	if item.Author != nil {
		author = fmt.Sprintf("  by %s\n", item.Author.Name)
	}

	return fmt.Sprintf(
		"■ %s\n\n%s  published at %s, updated at %s\n\n%s\n",
		item.Title,
		author,
		item.Published,
		item.Updated,
		item.Description,
	)
}

func FindItemMulti(items []*gofeed.Item) ([]int, error) {
	return fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			return fmt.Sprintf("%s [%s ago]", items[i].Title, CompTimeDiff(items[i].PublishedParsed))
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}
			return runewidth.Wrap(bodyFromItem(items[i]), width/2-5)
		}),
	)
}

func FindItem(items []*gofeed.Item) (int, error) {
	return fuzzyfinder.Find(
		items,
		func(i int) string {
			return fmt.Sprintf("%s [%s ago]", items[i].Title, CompTimeDiff(items[i].PublishedParsed))
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}
			return runewidth.Wrap(bodyFromItem(items[i]), width/2-5)
		}),
	)
}

func CompTimeDiff(t *time.Time) string {
	now := time.Now()
	diff := int(now.Sub(*t).Hours())

	day := diff / 24
	if day > 30 {
		month := day / 30
		return fmt.Sprintf("%dmon", month)
	}
	if day == 0 {
		hours := diff % 24
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dd", day)
}
