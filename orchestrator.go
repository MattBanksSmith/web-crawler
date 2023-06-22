package main

import (
	"log"
	"net/url"
	"sync"
)

func startOrchestrator(done chan struct{}, siteStream <-chan string, crawledLinksStream chan string, wg *sync.WaitGroup, crawlerCount int, urls map[string]struct{}) {
	wg.Add(1)
	defer func() { wg.Done() }()

	urlCrawlerStream := make(chan string)
	crawlerWg := sync.WaitGroup{}
	for i := 0; i < crawlerCount; i++ {
		go startWebReader(crawledLinksStream, urlCrawlerStream, &crawlerWg)
	}

	for u, _ := range urls {
		urlCrawlerStream <- u
	}

	defer func() {
		log.Printf("closing urlCrawlerStream")
		close(urlCrawlerStream)
		log.Printf("waiting on crawler wg")
		crawlerWg.Wait()

		log.Printf("closing crawledLinksStream")
		close(crawledLinksStream)
	}()
	for {
		select {
		case data, ok := <-siteStream:
			if !ok {
				log.Printf("siteStream closed\n")
				return
			}
			_, err := url.Parse(data)
			if err != nil {
				continue
			}
			urlCrawlerStream <- data
		}
	}

}
