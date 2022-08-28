package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
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
	description := func() string {
		content, err := renderHTML(item.Description)
		if err != nil {
			return item.Description
		}

		return content
	}()

	return fmt.Sprintf(
		"■ %s\n\n  %s\n\n  %s %s\n\n%s\n",
		item.Title,
		sprintfIfNotEmpty("by %s", author),
		sprintfIfNotEmpty("published at %s", publishedAt),
		sprintfIfNotEmpty("updated at %s", updatedAt),
		description,
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
	content := func() string {
		c, err := renderHTML(item.Content)
		if err != nil {
			return item.Content
		}

		return c
	}()
	description := func() string {
		c, err := renderHTML(item.Description)
		if err != nil {
			return item.Content
		}

		return c
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
		sprintfIfNotEmpty("%s", description),
		sprintfIfNotEmpty("%s", content),
		sprintfIfNotEmpty("%s", strings.Join(item.Links, "\n")),
	)
}

func sprintfIfNotEmpty(format string, str string) string {
	if str == "" || format == "" {
		return ""
	}

	return fmt.Sprintf(format, str)
}

//nolint:gomnd
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

func renderHTML(content string) (string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse content as HTML: %w", err)
	}

	var buf bytes.Buffer

	removeHTMLTags(doc, &buf)

	return buf.String(), nil
}

//nolint:interfacer
func removeHTMLTags(node *html.Node, buf *bytes.Buffer) {
	if node.Type == html.TextNode {
		buf.WriteString(node.Data)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		removeHTMLTags(child, buf)
	}
}
