package crawler

import (
	"log"

	"github.com/dmitryk-dk/simlpe-crawler/fetcher"
	"github.com/dmitryk-dk/simlpe-crawler/url_cache"
)

type SimpleCrawler struct {
	fetcher           fetcher.Fetcher
	urlProcessor      *url_cache.URLProcessor
	doneC             chan struct{}
	linksC            chan []string
	visitedLinksC     chan []string
	baseURL           string
	numberOfViewLinks int
}

func New(baseURL string, numberOfViewLinks int) *SimpleCrawler {
	c := &SimpleCrawler{
		fetcher:           fetcher.New(),
		urlProcessor:      url_cache.New(),
		linksC:            make(chan []string),
		doneC:             make(chan struct{}),
		visitedLinksC:     make(chan []string),
		baseURL:           baseURL,
		numberOfViewLinks: numberOfViewLinks,
	}
	go c.start()
	go c.checkLimit()
	return c
}

func (c *SimpleCrawler) Process() {
	for {
		select {
		case links := <-c.linksC:
			for _, link := range links {
				go c.getLinks(link)
			}
		case <-c.doneC:
			c.visitedLinksC <- c.urlProcessor.GetVisitedLinks()
			return
		}
	}
}

func (c *SimpleCrawler) VisitedLinks() chan []string {
	return c.visitedLinksC
}

func (c *SimpleCrawler) start() {
	c.linksC <- []string{c.baseURL}
}

func (c *SimpleCrawler) checkLimit() {
	for {
		if c.urlProcessor.GetLenVisitedLinks() == c.numberOfViewLinks {
			close(c.doneC)
			return
		}
	}
}

func (c *SimpleCrawler) getLinks(link string) {
	res, err := c.fetcher.Fetch(link)
	if err != nil {
		log.Printf("err fetch url %s : %s", link, err)
	} else {
		c.linksC <- c.urlProcessor.FilterLinks(c.baseURL, fetcher.ExtractLinks(res))
	}
}
