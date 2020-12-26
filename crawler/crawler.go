package crawler

import (
	"log"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type URLFilter interface {
	ExtractLinks([]byte) []string
	CollectLinks([]string)
	FilterLinks() []string
	GetBaseURL() string
}

type URLCacher interface {
	GetLenVisitedLinks() int
	GetVisitedLinks() []string
	FilterLinks([]string) []string
}

type SimpleCrawler struct {
	fetcher           Fetcher
	urlFilter         URLFilter
	urlCache          URLCacher
	doneC             chan struct{}
	linksC            chan []string
	visitedLinksC     chan []string
	reqLimit          chan string
	numberOfViewLinks int
}

func New(fetcher Fetcher, urlFilter URLFilter, urlCacher URLCacher, numberOfViewLinks, reqLimit int) *SimpleCrawler {
	c := &SimpleCrawler{
		fetcher:           fetcher,
		urlFilter:         urlFilter,
		urlCache:          urlCacher,
		linksC:            make(chan []string),
		doneC:             make(chan struct{}),
		visitedLinksC:     make(chan []string),
		reqLimit:          make(chan string, reqLimit),
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
	c.linksC <- []string{c.urlFilter.GetBaseURL()}
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
		c.urlFilter.CollectLinks(c.urlFilter.ExtractLinks(res))
		c.linksC <- c.urlCache.FilterLinks(c.urlFilter.FilterLinks())
	}
}
