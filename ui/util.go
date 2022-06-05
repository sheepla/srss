package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gilliek/go-opml/opml"
	"github.com/mmcdole/gofeed"
)

func renderPreviewWindow(item *gofeed.Item) string {
	author := func() string {
		if item.Author != nil {
			return item.Author.Name
		}

		return ""
	}()
	publishedAt := func() string {
		if item.PublishedParsed != nil {
			return humanizeTime(item.PublishedParsed)
		}

		return item.Published
	}()
	updatedAt := func() string {
		if item.UpdatedParsed != nil {
			return humanizeTime(item.UpdatedParsed)
		}

		return item.Updated
	}()

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
	author := func() string {
		if item.Author != nil {
			return item.Author.Name
		}

		return ""
	}()
	publishedAt := func() string {
		if item.PublishedParsed != nil {
			return humanizeTime(item.PublishedParsed)
		}

		return item.Published
	}()
	updatedAt := func() string {
		if item.UpdatedParsed != nil {
			return humanizeTime(item.UpdatedParsed)
		}

		return item.Updated
	}()

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

// nolint:gomnd
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
		return nil, fmt.Errorf("failed to parse OPML file: %w", err)
	}

	return doc, nil
}

func ExtractFeedURL(items []opml.Outline) []string {
	arr := make([]string, 0)

	for _, category := range items {
		for _, feed := range category.Outlines {
			arr = append(arr, feed.XMLURL)
		}
	}

	return arr
}
