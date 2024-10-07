package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type WebPage struct {
	Url  string
	Data []byte
}

func fetchURL(url string) (WebPage, error) {
	resp, err := http.Get(url)
	if err != nil {
		return WebPage{}, fmt.Errorf("error fetching %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WebPage{}, fmt.Errorf("error reading body from %s: %v", url, err)
	}
	return WebPage{Url: url, Data: body}, nil
}

func scrapeURLS(urls []string) ([]WebPage, []error) {
	var wg sync.WaitGroup
	results := make([]WebPage, 0, len(urls))
	errors := make([]error, 0, len(urls))
	resultsChan := make(chan WebPage, len(urls))
	errorsChan := make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			webpage, err := fetchURL(url)
			if err != nil {
				errorsChan <- err
			} else {
				resultsChan <- webpage
			}
		}(url)
	}
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	for result := range resultsChan {
		results = append(results, result)
	}

	for err := range errorsChan {
		errors = append(errors, err)
	}

	return results, errors
}

func saveToFile(webpage *WebPage, fileDir string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", fileDir, err)
	}

	parsedUrl, err := url.Parse(webpage.Url)
	if err != nil {
		return fmt.Errorf("error parsing URL %s: %w", webpage.Url, err)
	}

	filename := fileDir + "/" + parsedUrl.Hostname() + ".html"
	err = os.WriteFile(filename, webpage.Data, 0644)
	if err != nil {
		return fmt.Errorf("error saving %s: %w", webpage.Url, err)
	}
	return nil
}

func main() {
	urls := []string{
		"https://example.com",
		"https://golang.com",
		"https://openai.com",
	}

	folderDir := "scrapedWebsite"
	results, errors := scrapeURLS(urls)

	for _, result := range results {
		fmt.Printf("Saving %s\n", result.Url)
		saveToFile(&result, folderDir)
	}
	for _, err := range errors {
		log.Println(err)
	}
}
