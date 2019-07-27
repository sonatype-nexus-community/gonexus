package nexusiq

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

func newTestIQ(t *testing.T, handler func(t *testing.T, w http.ResponseWriter, r *http.Request)) (iq IQ, mock *httptest.Server) {
	mock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		handler(t, w, r)
	}))

	iq, err := New(mock.URL, "dummy_user", "dummy_pass")
	if err != nil {
		t.Fatal(err)
	}

	return
}
