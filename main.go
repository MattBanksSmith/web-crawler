package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.SetOutput(os.Stderr)
	log.Println("starting application")
	urls := map[string]struct{}{
		"https://theuselessweb.com/": {},
	}
	crawledLinksStream := make(chan string, 10)
	sitesStream := make(chan string)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	go startCrawledLinksProcessor(done, crawledLinksStream, sitesStream, urls, &wg)
	go startOrchestrator(done, sitesStream, crawledLinksStream, &wg, 50, urls)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("sig received, waiting for modules to close")
	close(done)
	wg.Wait()
	log.Println("ending application")
}
