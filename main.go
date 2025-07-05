package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"path"
	"sheet-parser/utils"
	"sync"
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

	browser.Close()

	formattedTitle := utils.FormatTitle(pageTitle)

	dirPath := path.Join("./", formattedTitle)

	_, err = os.ReadDir(dirPath)

	if err != nil {
		os.Mkdir(dirPath, os.ModePerm)
	}

	var wg sync.WaitGroup

	for _, link := range elems {
		wg.Add(1)
		go func() {
			err := utils.DownloadImage(dirPath, link)
			if err != nil {
				fmt.Printf("could not download image: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	pngPaths, err := utils.ConvertToPng(dirPath)

	if err != nil {
		log.Fatalf("could not convert to png: %v", err)
	}

	err = utils.ConvertToPdf(dirPath, pngPaths)

	if err != nil {
		log.Fatalf("could not convert to pdf: %v", err)
	}

	fmt.Printf("Converted to PDF successfully to %s", dirPath)
}
