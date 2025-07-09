package main

import (
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"path"
	"sheet-parser/utils"
	"sync"
)

const maxGoroutines = 4

func main() {
	UrlToParse := utils.ParseUrlFromArgs()

	var urls []string

	err := json.Unmarshal([]byte(UrlToParse), &urls)

	if err != nil {
		// processing as string
		parse(UrlToParse)
		return
	}

	semaphore := make(chan string, maxGoroutines)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		semaphore <- url

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
			parse(url)
		}()
	}

	wg.Wait()
}

func parse(urlToParse string) {
	log.Printf("Processing url %s", urlToParse)

	pw, err := playwright.Run()
	if err != nil {
		log.Printf("Could not start playwright: %v", err)
		return
	}

	browser, err := pw.Chromium.Launch()

	if err != nil {
		log.Printf("Could not launch Chromium: %v", err)
		return
	}

	page, err := browser.NewPage()

	if err != nil {
		log.Printf("Could not create page: %v", err)
		return
	}

	if _, err = page.Goto(urlToParse, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		log.Printf("Could not goto: %v", err)
		return
	}

	entities, err := page.Locator(".EEnGW").All()

	if err != nil {
		log.Printf("Could not get entities: %v", err)
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
		log.Printf("Could not convert to png: %v", err)
		return
	}

	err = utils.ConvertToPdf(dirPath, pngPaths)

	if err != nil {
		log.Printf("Could not convert to pdf: %v", err)
		return
	}

	log.Printf("Converted to PDF successfully to %s", dirPath)
}
