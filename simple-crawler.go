package main

import (
	"flag"
	"log"
	"time"

	"github.com/dmitryk-dk/simlpe-crawler/crawler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var baseURL string
	var numLinks int
	var limit int
	var duration time.Duration
	flag.StringVar(&baseURL, "url", "", "start crawling url")
	flag.IntVar(&numLinks, "links", 40, "number of collected links")
	flag.IntVar(&limit, "limit", 10, "limit of requests")
	flag.DurationVar(&duration, "timeout", time.Second*10, "timeout of response waiting")
	flag.Parse()
	if baseURL == "" {
		log.Fatal("Error run crawler on empty url")
	}

	c := crawler.New(baseURL, numLinks, limit, duration)
	go c.Process()
	log.Printf("visited links -> %+v", <-c.VisitedLinks())
}
