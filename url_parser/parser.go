package url_parser

import (
	"net/url"
	"strings"
)

type LinkFilter interface {
	CollectLinks(hrefs []string) *Filter
	FilterLinks() []string
}

type Filter struct {
	baseURL *url.URL
	links   []*url.URL
}

func New(baseURL string) *Filter {
	link, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}
	return &Filter{
		baseURL: link,
		links:   make([]*url.URL, 0),
	}
}

func (f *Filter) CollectLinks(hrefs []string) *Filter {
	for _, href := range hrefs {
		if href == "" {
			continue
		}
		link, err := url.Parse(href)
		if err != nil {
			continue
		}
		f.links = append(f.links, link)
	}
	return f
}

func (f *Filter) FilterLinks() []string {
	var internalLinks []string
	for _, link := range f.links {
		if link.IsAbs() &&
			f.hasSameDomain(link) &&
			f.hasSubdomains(link) &&
			f.hasSameSchema(link) &&
			f.notEmptyLink(link) {
			internalLinks = append(internalLinks, link.String())
		}
	}
	return internalLinks
}

func (f *Filter) hasSubdomains(link *url.URL) bool {
	return !strings.HasSuffix(link.Host, "."+f.baseURL.String())
}

func (f *Filter) hasSameDomain(link *url.URL) bool {
	return link.Host == f.baseURL.Host
}

func (f *Filter) hasSameSchema(link *url.URL) bool {
	return link.Scheme == f.baseURL.Scheme
}

func (f *Filter) notEmptyLink(link *url.URL) bool {
	return link.String() != ""
}
