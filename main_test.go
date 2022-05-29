package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestIsValidURL(t *testing.T) {
	validURL := "https://zenn.dev/topics/go/feed"
	have := isValidURL(validURL)
	want := true
	if have != want {
		t.Errorf("%s is valid URL", validURL)
	}

	invalidURL := "hoge"
	have = isValidURL(invalidURL)
	want = false
	if have != want {
		t.Errorf("%s is invalid URL", invalidURL)
	}
}

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
		t.Errorf("Got invalid URL but no error occured: (feed: %s)", feed)
	}
}

func TestIsUniqueURL(t *testing.T) {
	urls := []string{
		"https://zenn.dev/topics/go/feed",
		"https://www.archlinux.jp/feeds/news.xml",
	}

	uniqueURL := "https://zenn.dev/topics/python/feed"
	duplicateURL := "https://zenn.dev/topics/go/feed"
	have := isUniqueURL(urls, uniqueURL)
	want := true
	if have != want {
		t.Errorf("%s is unique URL", uniqueURL)
	}

	have = isUniqueURL(urls, duplicateURL)
	want = false
	if have != want {
		t.Errorf("%s is duplicate URL", duplicateURL)
	}
}
