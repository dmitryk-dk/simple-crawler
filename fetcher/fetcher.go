package fetcher

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type Client struct {
	client *http.Client
}

func New() *Client {
	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:     true,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper), // Disable HTTP/2
				ResponseHeaderTimeout: time.Second * 10,
			},
			Jar:     cookieJar,
			Timeout: time.Second * 10,
		},
	}
}

func (f *Client) Fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(){ _ = resp.Body.Close() }()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ExtractLinks(htmlDoc []byte) []string {
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
