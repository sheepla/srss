package cache_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/sheepla/srss/cache"
)

//nolint:paralleltest
func TestExport(t *testing.T) {
	items := []*gofeed.Item{
		{
			Title:       "TITLE1",
			Description: "DESCRIPTION1",
			Content:     "CONTENT1",
		},
		{
			Title:       "TITLE2",
			Description: "DESCRIPTION2",
			Content:     "CONTENT2",
		},
		{
			Title:       "TITLE3",
			Description: "DESCRIPTION3",
			Content:     "CONTENT3",
		},
	}

	if err := cache.Export(items); err != nil {
		t.Errorf("an error occurred on `Export()`: %s", err)
	}
}

//nolint:paralleltest
func TestImport(t *testing.T) {
	// var items []gofeed.Item{}
	items, err := cache.Import()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			t.Errorf("an error occurred on `Import()`: %s", err)
		}

		//nolint:forbidigo
		fmt.Println("EOF")
	}

	//nolint:forbidigo
	fmt.Printf("items: %v\n", items)
}
