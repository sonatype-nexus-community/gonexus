package nexusrm

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

const dummyContinuationToken = "go_on..."

func newTestRM(t *testing.T, handler func(t *testing.T, w http.ResponseWriter, r *http.Request)) (rm RM, mock *httptest.Server) {
	mock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		handler(t, w, r)
	}))

	rm, err := New(mock.URL, "dummy_user", "dummy_pass")
	if err != nil {
		t.Fatal(err)
	}

	return
}
