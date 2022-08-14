package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

//nolint:gochecknoglobals
var (
	cacheDir  = filepath.Join(configdir.LocalCache(), "srss")
	cacheFile = filepath.Join(cacheDir, "cache.gob")
)

func Export(items []*gofeed.Item) (err error) {
	if err := mkdir(cacheDir); err != nil {
		return fmt.Errorf("failed to create cache parent directory(%s): %w", cacheDir, err)
	}

	file, err := open(cacheFile)
	if err != nil {
		return fmt.Errorf("failed to open the cache file(%s): %w", cacheFile, err)
	}

	defer func() {
		if e := file.Close(); e != nil {
			err = fmt.Errorf("failed to close the cache file(%s): %w", cacheFile, e)
		}
	}()

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err = enc.Encode(&items); err != nil {
		return fmt.Errorf("failed to encode cache data: %w", err)
	}

	fmt.Fprint(file, buf)

	return nil
}

//nolint:nonamedreturns
func Import() (items []*gofeed.Item, err error) {
	if err := mkdir(cacheDir); err != nil {
		return nil, fmt.Errorf("failed to create cache parent directory(%s): %w", cacheDir, err)
	}

	file, err := open(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open the cache file(%s): %w", cacheFile, err)
	}

	defer func() {
		if e := file.Close(); e != nil {
			err = fmt.Errorf("failed to close the cache file(%s): %w", cacheFile, e)
		}
	}()

	err = gob.NewDecoder(file).Decode(&items)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, io.EOF
		}

		return nil, fmt.Errorf("failed to load cache from the file (%s): %w", cacheFile, err)
	}

	return items, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)

	// return !os.IsNotExist(err)
	return err == nil
}

func mkdir(path string) error {
	if !exists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory(%s), %w", path, err)
		}
	}

	return nil
}

func open(path string) (*os.File, error) {
	if !exists(path) {
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create the file(%s): %w", path, err)
		}

		return file, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file(%s): %w", path, err)
	}

	return file, nil
}
