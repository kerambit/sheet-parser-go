package utils

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func WaitLink(parentElem playwright.Locator, selector string) (string, error) {
	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elem := parentElem.Locator(selector)
			count, err := elem.Count()
			if err != nil {
				return "", err
			}
			if count == 0 {
				continue
			}

			link, err := elem.First().GetAttribute("src")
			if err != nil {
				return "", err
			}
			if link != "" {
				return link, nil
			}
		}
	}
}

func FormatTitle(str string) string {
	re := regexp.MustCompile(`[\s\-~()]+`)
	return re.ReplaceAllString(str, "_")
}

func DownloadImage(dirName, url string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := getFileName(url)
	filePath := filepath.Join(cwd, dirName, fileName)
	fmt.Println("Downloading", filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	return nil
}

func getFileName(url string) string {
	parts := strings.Split(url, "/")
	if !strings.Contains(parts[len(parts)-1], "?") {
		return parts[len(parts)-1]
	}

	parts = strings.Split(parts[len(parts)-1], "?")
	return parts[0]
}
