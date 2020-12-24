package fetcher

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type Client struct {
	client *http.Client
}

func New(timeout time.Duration) *Client {
	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:     true,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper), // Disable HTTP/2
				ResponseHeaderTimeout: timeout,
			},
			Jar:     cookieJar,
			Timeout: timeout,
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
	defer func() { _ = resp.Body.Close() }()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}
