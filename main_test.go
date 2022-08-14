package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

//nolint:paralleltest
func TestFetchFeed(t *testing.T) {
	validURL := "https://zenn.dev/topics/go/feed"

	feed, err := gofeed.NewParser().ParseURL(validURL)
	if feed == nil {
		t.Errorf("Feed is nil: %s", feed)
	}

	if err != nil {
		t.Errorf("Failed to fetch feed: %s", err)
	}

	invalidURL := "hoge"

	feed, err = gofeed.NewParser().ParseURL(invalidURL)
	if feed != nil {
		t.Errorf("Feed is not nil: %s", feed)
	}

	if err == nil {
		t.Errorf("Got invalid URL but no error occurred: (feed: %s)", feed)
	}
}
