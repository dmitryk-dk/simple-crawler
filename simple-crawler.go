package main

import (
	"flag"
	"log"
	"time"

	"github.com/dmitryk-dk/simlpe-crawler/crawler"
	"github.com/dmitryk-dk/simlpe-crawler/fetcher"
	_ "github.com/dmitryk-dk/simlpe-crawler/fetcher"
	"github.com/dmitryk-dk/simlpe-crawler/url_cache"
	"github.com/dmitryk-dk/simlpe-crawler/url_filter"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var baseURL string
	var numLinks int
	var limit int
	var timeout time.Duration
	flag.StringVar(&baseURL, "url", "", "start crawling url")
	flag.IntVar(&numLinks, "links", 40, "number of collected links")
	flag.IntVar(&limit, "limit", 10, "limit of requests")
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout of response waiting")
	flag.Parse()
	if baseURL == "" {
		log.Fatal("Error run crawler on empty url")
	}
	c := crawler.New(
		fetcher.New(timeout),
		url_filter.New(baseURL),
		url_cache.New(),
		numLinks,
		limit)
	go c.Process()
	log.Printf("visited links -> %+v", <-c.VisitedLinks())
}
