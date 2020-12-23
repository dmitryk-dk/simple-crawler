package url_cache

import (
	"sync"
)

type URLCache struct {
	visited map[string]struct{}
	mux     sync.Mutex
}

func New() *URLCache {
	return &URLCache{
		visited: make(map[string]struct{}),
	}
}

func (uc *URLCache) GetLenVisitedLinks() int {
	return len(uc.visited)
}

func (uc *URLCache) GetVisitedLinks() []string {
	uc.mux.Lock()
	defer uc.mux.Unlock()
	var visited []string
	for link := range uc.visited {
		visited = append(visited, link)
	}
	return visited
}

func (uc *URLCache) FilterLinks(links []string) []string {
	var filteredLinks []string
	for _, link := range links {
		if !uc.isVisited(link) {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func (uc *URLCache) isVisited(url string) bool {
	uc.mux.Lock()
	defer uc.mux.Unlock()
	if _, ok := uc.visited[url]; ok {
		return ok
	}
	uc.visited[url] = struct{}{}
	return false
}
