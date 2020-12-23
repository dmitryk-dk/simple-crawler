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

// TODO continue add tests
// func TestFilter_FilterLinks(t *testing.T) {
// 	type fields struct {
// 		baseURL *url.URL
// 		links   []*url.URL
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   []string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Filter{
// 				baseURL: tt.fields.baseURL,
// 				links:   tt.fields.links,
// 			}
// 			if got := f.FilterLinks(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("FilterLinks() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestFilter_hasSameDomain(t *testing.T) {
// 	type fields struct {
// 		baseURL *url.URL
// 		links   []*url.URL
// 	}
// 	type args struct {
// 		link *url.URL
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Filter{
// 				baseURL: tt.fields.baseURL,
// 				links:   tt.fields.links,
// 			}
// 			if got := f.hasSameDomain(tt.args.link); got != tt.want {
// 				t.Errorf("hasSameDomain() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestFilter_hasSameSchema(t *testing.T) {
// 	type fields struct {
// 		baseURL *url.URL
// 		links   []*url.URL
// 	}
// 	type args struct {
// 		link *url.URL
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Filter{
// 				baseURL: tt.fields.baseURL,
// 				links:   tt.fields.links,
// 			}
// 			if got := f.hasSameSchema(tt.args.link); got != tt.want {
// 				t.Errorf("hasSameSchema() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestFilter_hasSubdomains(t *testing.T) {
// 	type fields struct {
// 		baseURL *url.URL
// 		links   []*url.URL
// 	}
// 	type args struct {
// 		link *url.URL
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Filter{
// 				baseURL: tt.fields.baseURL,
// 				links:   tt.fields.links,
// 			}
// 			if got := f.hasSubdomains(tt.args.link); got != tt.want {
// 				t.Errorf("hasSubdomains() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestFilter_notEmptyLink(t *testing.T) {
// 	type fields struct {
// 		baseURL *url.URL
// 		links   []*url.URL
// 	}
// 	type args struct {
// 		link *url.URL
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Filter{
// 				baseURL: tt.fields.baseURL,
// 				links:   tt.fields.links,
// 			}
// 			if got := f.notEmptyLink(tt.args.link); got != tt.want {
// 				t.Errorf("notEmptyLink() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		baseURL string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *Filter
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := New(tt.args.baseURL); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("New() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
