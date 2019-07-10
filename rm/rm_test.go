package nexusrm

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRM(handler http.Handler) (rm RM, mock *httptest.Server, err error) {
	mock = httptest.NewServer(handler)

	rm, err = New(mock.URL, "dummy_user", "dummy_pass")
	if err != nil {
		return
	}

	return
}

func getTestRM(t *testing.T) RM {
	rm, err := New("http://localhost:8081", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}

	return rm
}
