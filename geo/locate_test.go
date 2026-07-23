package geo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocate(t *testing.T) {
	tests := []struct {
		name    string   // test-case ID + intent
		body    string   // body the fake ip-api.com server returns
		wantErr bool     // whether Locate should return a non-nil error
		want    Location // expected location on the happy path
	}{
		{
			name:    "UT-010 success returns decoded location",
			body:    `{"status":"success","lat":38.72,"lon":-9.14,"city":"Lisbon","country":"Portugal"}`,
			wantErr: false,
			want:    Location{Lat: 38.72, Lon: -9.14, City: "Lisbon", Country: "Portugal"},
		},
		{
			name:    "UT-011 fail status returns error",
			body:    `{"status":"fail"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/json" {
					t.Errorf("unexpected request path: got %q, want %q", r.URL.Path, "/json")
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			// Point the package at the fake server, restoring globals afterwards.
			origBase, origClient := BaseURL, HTTPClient
			BaseURL, HTTPClient = srv.URL, srv.Client()
			defer func() { BaseURL, HTTPClient = origBase, origClient }()

			got, err := Locate(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Locate() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Locate() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Locate() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
