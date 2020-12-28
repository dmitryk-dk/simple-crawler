package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	client *http.Client
}

func New(timeout time.Duration) *Client {
	// ignore err because cookiejar.New(nil) return always nil in err
	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{ResponseHeaderTimeout: timeout},
			Jar:       cookieJar,
			Timeout:   timeout,
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
