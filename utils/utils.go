package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/go-pdf/fpdf"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

	return err
}

func getFileName(url string) string {
	parts := strings.Split(url, "/")
	if !strings.Contains(parts[len(parts)-1], "?") {
		return parts[len(parts)-1]
	}

	parts = strings.Split(parts[len(parts)-1], "?")
	return parts[0]
}

func ConvertToPng(dirName string) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	svgPaths, _ := os.ReadDir(path.Join(cwd, dirName))

	pngPaths := make([]string, len(svgPaths))

	var wg sync.WaitGroup

	for i, svgPath := range svgPaths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fileData, err := os.ReadFile(path.Join(cwd, dirName, svgPath.Name()))
			fileName := svgPath.Name()

			if strings.Contains(fileName, ".png") {
				pngPaths[i] = path.Join(cwd, dirName, fileName)
				return
			}

			if err != nil {
				return
			}

			html := fmt.Sprintf(`
		<html>
			<body style="margin:0">
				<img src="data:image/svg+xml;base64,%s" />
			</body>
		</html>`, base64.StdEncoding.EncodeToString(fileData))

			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			var output []byte
			err = chromedp.Run(ctx,
				chromedp.Navigate("data:text/html,"+html),
				chromedp.WaitVisible("img"),
				chromedp.Screenshot("img", &output, chromedp.NodeVisible),
			)
			if err != nil {
				log.Fatal(err)
			}

			pngPath := strings.Replace(fileName, ".svg", ".png", 1)

			os.WriteFile(path.Join(cwd, dirName, pngPath), output, 0644)
			pngPaths[i] = path.Join(cwd, dirName, pngPath)

			os.Remove(path.Join(cwd, dirName, fileName))
		}()

	}

	wg.Wait()

	return pngPaths, nil
}

func ConvertToPdf(dirName string, pngs []string) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	pdf := fpdf.New("P", "mm", "A4", "")

	for _, imgPath := range pngs {
		pdf.AddPage()

		fmt.Println("Processing", imgPath)

		pageWidth, _ := pdf.GetPageSize()

		pdf.ImageOptions(
			imgPath,
			0, 0,
			pageWidth, 0,
			false,
			fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
			0, "",
		)

		os.Remove(imgPath)
	}

	err = pdf.OutputFileAndClose(path.Join(cwd, dirName, "output.pdf"))

	return err
}
