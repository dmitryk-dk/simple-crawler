package crawler

import (
	"log"
	"time"

	"github.com/dmitryk-dk/simlpe-crawler/fetcher"
	"github.com/dmitryk-dk/simlpe-crawler/url_cache"
	"github.com/dmitryk-dk/simlpe-crawler/url_filter"
)

type SimpleCrawler struct {
	fetcher           fetcher.Fetcher
	urlCache          *url_cache.URLCache
	urlParser         url_filter.URLFilter
	doneC             chan struct{}
	linksC            chan []string
	visitedLinksC     chan []string
	reqLimit          chan string
	baseURL           string
	numberOfViewLinks int
}

func New(baseURL string, numberOfViewLinks, reqLimit int, timeout time.Duration) *SimpleCrawler {
	c := &SimpleCrawler{
		fetcher:           fetcher.New(timeout),
		urlCache:          url_cache.New(),
		urlParser:         url_filter.New(baseURL),
		linksC:            make(chan []string),
		doneC:             make(chan struct{}),
		visitedLinksC:     make(chan []string),
		reqLimit:          make(chan string, reqLimit),
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
				c.reqLimit <- link
				go c.getLinks(link)
			}
		case <-c.doneC:
			c.visitedLinksC <- c.urlCache.GetVisitedLinks()
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
		if c.urlCache.GetLenVisitedLinks() == c.numberOfViewLinks {
			close(c.doneC)
			return
		}
	}
}

func (c *SimpleCrawler) getLinks(link string) {
	<-c.reqLimit
	res, err := c.fetcher.Fetch(link)
	if err != nil {
		log.Printf("err fetch url %s : %s", link, err)
	} else {
		c.linksC <- c.urlCache.FilterLinks(c.urlParser.CollectLinks(c.urlParser.ExtractLinks(res)).FilterLinks())
	}
}
