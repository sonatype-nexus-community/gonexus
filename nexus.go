package nexus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

/*
// ClientResponse provides an HTTP wrapper with optimized for communicating with a Nexus server
type ClientResponse struct {
	Response *http.Response
	Err      error
}

// Body will read the http.Response body
func (c *ClientResponse) Body() ([]byte, error) {
	if c.Response.StatusCode != http.StatusOK {
		return nil
	}

	defer c.Response.Body.Close()
	return ioutil.ReadAll(c.Response.Body)
}
*/

// DefaultClient provides an HTTP wrapper with optimized for communicating with a Nexus server
type DefaultClient struct {
	Host, Username, Password string
	Debug                    bool
}

// NewRequest created an http.Request object based on an endpoint and fills in basic auth
func (s DefaultClient) NewRequest(method, endpoint string, payload io.Reader) (request *http.Request, err error) {
	url := fmt.Sprintf("%s/%s", s.Host, endpoint)
	request, err = http.NewRequest(method, url, payload)
	if err != nil {
		return
	}

	request.SetBasicAuth(s.Username, s.Password)
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return
}

// Do performs an http.Request and reads the body if StatusOK
func (s DefaultClient) Do(request *http.Request) (body []byte, resp *http.Response, err error) {
	if s.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Printf("%q\n", dump)
	}

	client := &http.Client{}
	resp, err = client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Trying to decide if this is a horrible idea or just kinda bad
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		return
	}

	err = errors.New(resp.Status)
	return
}

func (s DefaultClient) http(method, endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	request, err := s.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, nil, err
	}

	return s.Do(request)
}

// Get performs an HTTP GET against the indicated endpoint
func (s DefaultClient) Get(endpoint string) ([]byte, *http.Response, error) {
	return s.http(http.MethodGet, endpoint, nil)
}

// Post performs an HTTP POST against the indicated endpoint
func (s DefaultClient) Post(endpoint string, payload []byte) ([]byte, *http.Response, error) {
	return s.http(http.MethodPost, endpoint, bytes.NewBuffer(payload))
}

// Put performs an HTTP PUT against the indicated endpoint
func (s DefaultClient) Put(endpoint string, payload []byte) (resp *http.Response, err error) {
	_, resp, err = s.http(http.MethodPut, endpoint, bytes.NewBuffer(payload))
	return
}

// Del performs an HTTP DELETE against the indicated endpoint
func (s DefaultClient) Del(endpoint string) (resp *http.Response, err error) {
	_, resp, err = s.http(http.MethodDelete, endpoint, nil)
	return
}
