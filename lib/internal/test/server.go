package test

import (
	"net/http"
	"net/http/httptest"
)

func InitTestServer(patterns map[string]string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if resp, ok := patterns[r.RequestURI]; ok {
				w.Write([]byte(resp))
				return
			}
			http.Error(w, "not found", http.StatusNotFound)
		}),
	)
}
