package utils

import (
	"github.com/playwright-community/playwright-go"
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
