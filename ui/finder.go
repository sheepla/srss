package ui

import (
	"fmt"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/mmcdole/gofeed"
)

func renderPreviewWindow(item *gofeed.Item) string {
	var author string
	if item.Author != nil {
		author = fmt.Sprintf("  by %s\n", item.Author.Name)
	}
	var publishedAt string
	if item.PublishedParsed != nil {
		publishedAt = humanizeTime(item.PublishedParsed)
	} else {
		publishedAt = item.Published
	}
	var updatedAt string
	if item.UpdatedParsed != nil {
		updatedAt = humanizeTime(item.UpdatedParsed)
	} else {
		updatedAt = item.Updated
	}
	return fmt.Sprintf(
		"■ %s\n\n%s  published at %s, updated at %s\n\n%s\n",
		item.Title,
		author,
		publishedAt,
		updatedAt,
		item.Description,
	)
}

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
			return runewidth.Wrap(renderPreviewWindow(items[i]), width/2-5)
		}),
	)
}

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
			return runewidth.Wrap(renderPreviewWindow(items[i]), width/2-5)
		}),
	)
}

func humanizeTime(t *time.Time) string {
	now := time.Now()
	diff := int(now.Sub(*t).Hours())
	day := diff / 24
	if day >= 30 {
		month := day / 30
		return fmt.Sprintf("%dmon ago", month)
	}
	if day == 0 {
		hours := diff % 24
		return fmt.Sprintf("%dh ago", hours)
	}
	return fmt.Sprintf("%dd ago", day)
}
