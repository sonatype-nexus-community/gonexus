package nexusiq

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testIQ struct {
	IQ
	server *httptest.Server
}

// NewRequest creates an http.Request object with private session
func (iq testIQ) NewRequest(method, endpoint string, payload io.Reader) (*http.Request, error) {
	// req, err := iq.defaultServer.NewRequest(method, endpoint, payload)
	/*
		req, err := iq.pub.NewRequest(method, endpoint, payload)
		if err != nil {
			return nil, err
		}

		// _, resp, err := iq.defaultServer.Get(iqRestSessionPrivate)
		_, resp, err := iq.pub.Get(iqRestSessionPrivate)
		if err != nil {
			return nil, err
		}

		for _, cookie := range resp.Cookies() {
			req.AddCookie(cookie)
			if cookie.Name == "CLM-CSRF-TOKEN" {
				req.Header.Add("X-CSRF-TOKEN", cookie.Value)
			}
		}

	*/
	return nil, nil
}

func (iq testIQ) http(method, endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	request, err := iq.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, nil, err
	}

	return iq.Do(request)
}

// Get performs an HTTP GET against the indicated endpoint
func (iq testIQ) Get(endpoint string) ([]byte, *http.Response, error) {
	return iq.http("GET", endpoint, nil)
}

// Post performs an HTTP POST against the indicated endpoint
func (iq testIQ) Post(endpoint string, payload []byte) ([]byte, *http.Response, error) {
	return iq.http("POST", endpoint, bytes.NewBuffer(payload))
}

// Put performs an HTTP PUT against the indicated endpoint
func (iq testIQ) Put(endpoint string, payload []byte) ([]byte, *http.Response, error) {
	return iq.http("PUT", endpoint, bytes.NewBuffer(payload))
}

// Del performs an HTTP DELETE against the indicated endpoint
func (iq testIQ) Del(endpoint string) (resp *http.Response, err error) {
	_, resp, err = iq.http("DELETE", endpoint, nil)
	return
}

func getTestIQ(t *testing.T) *IQ {
	iq, err := New("http://localhost:8070", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}
	iq.Debug = true

	return iq
}
