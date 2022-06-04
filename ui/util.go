package ui

import (
	"fmt"
	"github.com/gilliek/go-opml/opml"
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

func renderPreviewWindow(item *gofeed.Item) string {
	var author string
	if item.Author != nil {
		author = item.Author.Name
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
		"■ %s\n\n  %s\n\n  %s %s\n\n%s\n",
		item.Title,
		sprintfIfNotEmpty("by %s", author),
		sprintfIfNotEmpty("published at %s", publishedAt),
		sprintfIfNotEmpty("updated at %s", updatedAt),
		item.Description,
	)
}

func renderContent(item *gofeed.Item) string {
	var author string
	if item.Author != nil {
		sprintfIfNotEmpty("by %s ", item.Author.Name)
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
		`%s%s %s
──────
%s
%s
──────
%s
`,
		author,
		sprintfIfNotEmpty("published at %s", publishedAt),
		sprintfIfNotEmpty("updated at %s", updatedAt),
		sprintfIfNotEmpty("%s", item.Description),
		sprintfIfNotEmpty("%s", item.Content),
		sprintfIfNotEmpty("%s", strings.Join(item.Links, "\n")),
	)
}

func sprintfIfNotEmpty(format string, str string) string {
	if str == "" || format == "" {
		return ""
	}
	return fmt.Sprintf(format, str)
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

func ParseOPML(path string) (*opml.OPML, error) {
	doc, err := opml.NewOPMLFromFile(path)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func ExtractFeedURL(items []opml.Outline) []string {
	arr := make([]string, 0, 0)

	for _, category := range items {
		for _, feed := range category.Outlines {
			arr = append(arr, feed.XMLURL)
		}
	}

	return arr
}
