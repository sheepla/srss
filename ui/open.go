package ui

import (
	"fmt"

	"github.com/toqueteos/webbrowser"
)

func OpenURL(url string) error {
	if err := webbrowser.Open(url); err != nil {
		return fmt.Errorf("failed to open the URL (%s): %w", url, err)
	}

	return nil
}
