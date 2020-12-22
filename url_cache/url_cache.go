package url_cache

import (
	"net/url"
	"strings"
	"sync"
)

type URLProcessor struct {
	visited map[string]struct{}
	mux     sync.Mutex
}

func New() *URLProcessor {
	return &URLProcessor{
		visited: make(map[string]struct{}),
	}
}

func (uc *URLProcessor) FilterLinks(baseURL string, hrefs []string) []string {
	var internalURLs []string
	for _, href := range hrefs {
		if href == "" {
			continue
		}
		link, err := url.Parse(href)
		if err != nil {
			continue
		}
		if link.IsAbs() &&
			!strings.HasSuffix(link.Host, "."+baseURL) &&
			// TODO parse base url and compare hosts
			strings.HasSuffix(baseURL, link.Host) &&
			!uc.isVisited(link.String()) {
			internalURLs  = append(internalURLs, link.String())
		}
	}
	return internalURLs
}

func (uc *URLProcessor) GetLenVisitedLinks() int {
	return len(uc.visited)
}

func (uc *URLProcessor) GetVisitedLinks() []string {
	uc.mux.Lock()
	defer uc.mux.Unlock()
	var visited []string
	for link := range uc.visited {
		visited = append(visited, link)
	}
	return visited
}

func (uc *URLProcessor) isVisited(url string) bool {
	uc.mux.Lock()
	defer uc.mux.Unlock()
	if _, ok := uc.visited[url]; ok {
		return ok
	}
	uc.visited[url] = struct{}{}
	return false
}
