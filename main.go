package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"sheet-parser/utils"
)

const UrlToParse = "https://example.com/some/sheeet/1231"

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch()

	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}

	page, err := browser.NewPage()

	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err = page.Goto(UrlToParse, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	entities, err := page.Locator(".EEnGW").All()

	if err != nil {
		log.Fatalf("could not get entities: %v", err)
	}

	var elems []string

	for _, ent := range entities {
		ent.ScrollIntoViewIfNeeded()
		link, err := utils.WaitLink(ent, "img")
		if err != nil {
			continue
		}

		elems = append(elems, link)
	}

	fmt.Println(elems, len(elems))

	pageTitle, _ := page.Title()
	fmt.Println(pageTitle)

	browser.Close()
}
