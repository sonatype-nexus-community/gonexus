package nexusiq

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var dummySearchResults = []SearchResult{}

func searchTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func searchTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, searchTestFunc)
}

func TestSearchComponent(t *testing.T) {
	t.Skip("TODO")
	// iq, mock := searchTestIQ(t)
	// defer mock.Close()

}
