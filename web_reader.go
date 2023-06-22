package main

import (
	"log"
	"sync"
)

func startWebReader(crawledLinksStream chan<- string, urlCrawlerStream <-chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case url, ok := <-urlCrawlerStream:
			if !ok {
				log.Printf("closing web reader")
				return
			}
			urls := readAndExtract(url)
			//log.Printf("url [%v] parsed urls [%v]\n", url, urls)

			if urls == nil {
				continue
			}

			for extractedUrl, _ := range urls {
				//log.Printf("url [%v] parsed sub url [%v]\n", url, extractedUrl)

				go func(u string, group *sync.WaitGroup) {
					group.Add(1)
					defer group.Done()
					crawledLinksStream <- u
				}(extractedUrl, wg)
			}
		}
	}

}

func readAndExtract(url string) map[string]struct{} {
	data, err := GetWebPage(url)
	if err != nil {
		log.Printf("unable to fetch data for URL [%v] due to err [%v] skipping\n", url, err)
	}
	//pass into extractor
	urls, err := extractURLs(data)
	if err != nil {
		log.Printf("error extracting urls for data returned from URL [%v] due to err [%v] skipping\n", url, err)
	}
	return urls
}
