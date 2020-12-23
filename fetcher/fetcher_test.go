package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestClient_Fetch(t *testing.T) {
	tests := []struct {
		name    string
		server  *httptest.Server
		want    []byte
		wantErr bool
	}{
		{
			name: "return Hello, client",
			server: httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprintln(w, "Hello, client")
			})),
			want:    []byte("Hello, client\n"),
			wantErr: false,
		},
		{
			name: "has timeout for answer",
			server: httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Second * 11)
				_, _ = fmt.Fprintln(w, "Hello, client")
			})),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			f := New()
			got, err := f.Fetch(tt.server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
