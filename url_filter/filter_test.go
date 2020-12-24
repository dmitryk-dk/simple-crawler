package url_filter

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFilter_CollectLinks(t *testing.T) {
	tests := []struct {
		name   string
		hrefs  []string
		want   []*url.URL
		filter *Filter
	}{
		{
			name:   "empty hrefs slice",
			hrefs:  []string{},
			want:   []*url.URL{},
			filter: New("https://www.brightlocal.com"),
		},
		{
			name:   "one empty hrefs slice",
			hrefs:  []string{""},
			want:   []*url.URL{},
			filter: New("https://www.brightlocal.com"),
		},
		{
			name:   "error parse href",
			hrefs:  []string{"%a"},
			want:   []*url.URL{},
			filter: New("https://www.brightlocal.com"),
		},
		{
			name:  "valid hrefs",
			hrefs: []string{"http://www.google.com/?q=go+language#foo&bar", "http://j@ne:password@google.com/p@th?q=@go", "/foo?query=http://bad"},
			want: []*url.URL{
				{
					Scheme:   "http",
					Host:     "www.google.com",
					Path:     "/",
					RawQuery: "q=go+language",
					Fragment: "foo&bar",
				},
				{
					Scheme:   "http",
					User:     url.UserPassword("j@ne", "password"),
					Host:     "google.com",
					Path:     "/p@th",
					RawQuery: "q=@go",
				},
				{
					Path:     "/foo",
					RawQuery: "query=http://bad",
				},
			},
			filter: New("https://www.brightlocal.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.filter.CollectLinks(tt.hrefs)
			if !reflect.DeepEqual(tt.filter.links, tt.want) {
				t.Errorf("CollectLinks() = %#v, want %#v", tt.filter.links, tt.want)
			}
		})
	}
}

func TestFilter_ExtractLinks(t *testing.T) {
	tests := []struct {
		name   string
		html   []byte
		want   []string
		filter *Filter
	}{
		{
			name:   "empty html",
			html:   nil,
			want:   nil,
			filter: New("https://google.com"),
		},
		{
			name:   "html without links",
			html:   []byte("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Title</title>\n</head>\n<body>\n\n</body>\n</html>"),
			want:   []string(nil),
			filter: New("https://google.com"),
		},
		{
			name:   "has many link",
			html:   []byte("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Title</title>\n</head>\n<body>\n    <a href=\"/some/href\">Some href</a>\n    <div>\n        <div>\n            <div>\n                <a href=\"https://google.com\">Some href</a>\n            </div>\n        </div>\n    </div>\n</body>\n</html>"),
			want:   []string{"/some/href", "https://google.com"},
			filter: New("https://google.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.ExtractLinks(tt.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractLinks() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFilter_FilterLinks(t *testing.T) {
	tests := []struct {
		name   string
		filter *Filter
		hrefs  []string
		want   []string
	}{
		{
			name:   "empty links",
			filter: New("https://google.com"),
			hrefs:  []string{},
			want:   []string(nil),
		},
		{
			name:   "every link not in same domain",
			filter: New("https://google.com"),
			hrefs:  []string{"http://foo.com/blah_blah", "http://foo.com/blah_blah_(wikipedia)"},
			want:   []string(nil),
		},
		{
			name:   "link with wrong schema",
			filter: New("https://google.com"),
			hrefs:  []string{"http://10.1.1.255", "http://google.com", "https://www.google.com"},
			want:   []string{"https://www.google.com"},
		},
		{
			name:   "has relative links",
			filter: New("https://google.com"),
			hrefs:  []string{"/some/relative", "http://google.com", "https://google.com"},
			want:   []string{"https://google.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.filter.CollectLinks(tt.hrefs)
			if got := tt.filter.FilterLinks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterLinks() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFilter_hasSameDomain(t *testing.T) {
	tests := []struct {
		name   string
		filter *Filter
		link   string
		want   bool
	}{
		{
			name:   "empty link",
			filter: New("https://google.com"),
			link:   "",
			want:   false,
		},
		{
			name:   "with not same schema",
			filter: New("https://google.com"),
			link:   "http://google.com",
			want:   true,
		},
		{
			name:   "not same domain",
			filter: New("https://google.com"),
			link:   "http://googles.com",
			want:   false,
		},
		{
			name:   "hase subdomain",
			filter: New("https://google.com"),
			link:   "https://tools.googles.com",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, _ := url.Parse(tt.link)
			if got := tt.filter.hasSameDomain(link); got != tt.want {
				t.Errorf("hasSameDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
