package url_filter

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func (f *Filter) ExtractLinks(htmlDoc []byte) []string {
	var hrefs []string
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(htmlDoc))
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if s != nil {
				href, _ := s.Attr("href")
				hrefs = append(hrefs, href)
			}
		})
		return hrefs
	}
	return hrefs
}

func (f *Filter) CollectLinks(hrefs []string) {
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

func (f *Filter) GetBaseURL() string {
	return f.baseURL.String()
}

func (f *Filter) hasSubdomains(link *url.URL) bool {
	return !strings.HasSuffix(link.Host, "."+f.baseURL.String()) || strings.HasSuffix(link.Host, "www."+f.baseURL.String())
}

func (f *Filter) hasSameDomain(link *url.URL) bool {
	return link.Host == f.baseURL.Host || link.Host == "www."+f.baseURL.Host
}

func (f *Filter) hasSameSchema(link *url.URL) bool {
	return link.Scheme == f.baseURL.Scheme
}

func (f *Filter) notEmptyLink(link *url.URL) bool {
	return link.String() != ""
}
