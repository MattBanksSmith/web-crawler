package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// startCrawledLinksProcessor receives URLs from web crawlers and stores them
// if a url is not in the store, it forwards it on to sitesStream to be crawled
// has ownership of sitesStream as it's the sole sender
func startCrawledLinksProcessor(done chan struct{}, crawledLinksStream chan string, sitesStream chan<- string, urls map[string]struct{}, wg *sync.WaitGroup) {
	crawledURLs := urls
	wg.Add(1)
	defer func() {
		log.Printf("closing sitesStream")
		close(sitesStream)

		log.Printf("waiting for crawledLinksStream to empty")
		for _ = range crawledLinksStream {
		}

		//write results to file
		saveToFile(crawledURLs)
		wg.Done()
	}()

	for {
		select {
		case <-done:
			return
		case url, ok := <-crawledLinksStream:
			if !ok {
				log.Printf("crawledLinksStream closed")
				return
			}
			if _, exists := crawledURLs[url]; !exists {
				crawledURLs[url] = struct{}{}
				sitesStream <- url
			}
		}
	}
}

func saveToFile(crawledURLs map[string]struct{}) {
	sb := strings.Builder{}
	for u, _ := range crawledURLs {
		sb.WriteString(u)
		sb.WriteString("\n")
	}
	err := os.WriteFile("urls", []byte(sb.String()), 0)
	if err != nil {
		fmt.Println(err)
	}
}
