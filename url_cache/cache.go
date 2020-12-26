package url_cache

import (
	"sync"
	"sync/atomic"
)

type URLCache struct {
	visited sync.Map
	count   int32
}

func New() *URLCache {
	return &URLCache{}
}

func (uc *URLCache) GetLenVisitedLinks() int {
	return int(atomic.LoadInt32(&uc.count))
}

func (uc *URLCache) GetVisitedLinks() []string {
	var visited []string
	uc.visited.Range(func(key, val interface{}) bool {
		if link, ok := key.(string); ok {
			visited = append(visited, link)
			return ok
		}
		return true
	})
	return visited
}

func (uc *URLCache) FilterLinks(links []string) []string {
	var filteredLinks []string
	for _, link := range links {
		if !uc.isVisited(link) {
			uc.setVisited(link)
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func (uc *URLCache) isVisited(url string) bool {
	if _, ok := uc.visited.Load(url); ok {
		return ok
	}
	return false
}

func (uc *URLCache) setVisited(url string) {
	atomic.AddInt32(&uc.count, 1)
	uc.visited.Store(url, struct{}{})
}
