package nexusiq

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestIQ(handler http.Handler) (iq IQ, mock *httptest.Server, err error) {
	mock = httptest.NewServer(handler)

	iq, err = New(mock.URL, "dummy_user", "dummy_pass")
	if err != nil {
		return
	}

	return
}

func getTestIQ(t *testing.T) IQ {
	iq, err := New("http://localhost:8070", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}
	// iq.Debug = true

	return iq
}
