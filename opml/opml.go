package opml

import (
	"fmt"

	"github.com/gilliek/go-opml/opml"
)

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
