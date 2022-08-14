package urlentry

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kirsle/configdir"
)

//nolint:gochecknoglobals
var urlFile = filepath.Join(configdir.LocalConfig(), "srss", "urls.txt")

func Add(url string) error {
	if !isValidURL(url) {
		//nolint:goerr113
		return fmt.Errorf("invalid URL(%s)", url)
	}

	if err := createConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	//nolint:gomnd,nosnakecase
	file, err := os.OpenFile(urlFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open URL entry file (%s): %w", urlFile, err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, url)
	if err != nil {
		return fmt.Errorf("writing failed to the URL entry file (%s): %w", urlFile, err)
	}

	return nil
}

//nolint:wsl
func Load() ([]string, error) {
	if err := createConfigDir(); err != nil {
		return []string{}, fmt.Errorf("failed to create config directory: %w", err)
	}

	//nolint:gomnd,nosnakecase
	file, err := os.OpenFile(urlFile, os.O_RDONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to open URL entry file (%s): %w", urlFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var urls []string
	for scanner.Scan() {
		url := scanner.Text()
		if !isValidURL(url) {
			//nolint:goerr113
			return nil, fmt.Errorf("invalid URL(%s)", url)
		}
		urls = append(urls, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("failed to scan from URL entry file (%s): %w", urlFile, err)
	}

	return urls, nil
}

func IsUniqueURL(url string) bool {
	list, err := Load()
	if err != nil {
		return true
	}

	return isUnique(list, url)
}

func OpenEditor(editor string) error {
	err := execEditor(editor, urlFile)

	return err
}

func isUnique(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return false
		}
	}

	return true
}

func createConfigDir() error {
	if err := configdir.MakePath(filepath.Dir(urlFile)); err != nil {
		return fmt.Errorf("failed to create the directory: %w", err)
	}

	return nil
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)

	return err == nil
}

// https://doloopwhile.hatenablog.com/entry/2014/08/05/213819
func execEditor(editor string, args ...string) error {
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run editor (%s) %w", editor, err)
	}

	return nil
}
