package nexusiq

import (
	"net/http"
	"net/http/httptest"
)

func newTestIQ(handler http.Handler) (iq IQ, mock *httptest.Server, err error) {
	mock = httptest.NewServer(handler)

	iq, err = New(mock.URL, "dummy_user", "dummy_pass")
	if err != nil {
		return
	}

	return
}
