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

// Server is the interface which all Nexus server objects need to adhere to
type Server interface {
	// NewRequest(method, endpoint string, payload io.Reader) (*http.Request, error)
	// Do(request *http.Request) ([]byte, *http.Response, error)
	Get(endpoint string) ([]byte, *http.Response, error)
	Post(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Put(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Del(endpoint string) (resp *http.Response, err error)
}

// DefaultServer provides an HTTP wrapper with optimized for communicating with a Nexus server
type DefaultServer struct {
	Host, Username, Password string
	Debug                    bool
}

// NewRequest created an http.Request object based on an endpoint and fills in basic auth
func (s DefaultServer) NewRequest(method, endpoint string, payload io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", s.Host, endpoint)
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(s.Username, s.Password)
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

// Do performs an http.Request and reads the body if StatusOK
func (s DefaultServer) Do(request *http.Request) ([]byte, *http.Response, error) {
	if s.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Printf("%q\n", dump)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return body, resp, err
	}

	return nil, resp, errors.New(resp.Status)
}

func (s DefaultServer) http(method, endpoint string, payload io.Reader) ([]byte, *http.Response, error) {
	request, err := s.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, nil, err
	}

	return s.Do(request)
}

// Get performs an HTTP GET against the indicated endpoint
func (s DefaultServer) Get(endpoint string) ([]byte, *http.Response, error) {
	return s.http("GET", endpoint, nil)
}

// Post performs an HTTP POST against the indicated endpoint
func (s DefaultServer) Post(endpoint string, payload []byte) ([]byte, *http.Response, error) {
	return s.http("POST", endpoint, bytes.NewBuffer(payload))
}

// Put performs an HTTP PUT against the indicated endpoint
func (s DefaultServer) Put(endpoint string, payload []byte) ([]byte, *http.Response, error) {
	return s.http("PUT", endpoint, bytes.NewBuffer(payload))
}

// Del performs an HTTP DELETE against the indicated endpoint
func (s DefaultServer) Del(endpoint string) (resp *http.Response, err error) {
	_, resp, err = s.http("DELETE", endpoint, nil)
	return
}

/*
func iqResultToComponent(r nexusiq.ComponentEvaluationResult) component {
	var c component
	c.ComponentIdentifier.Format = r.Component.ComponentIdentifier.Format
	c.ComponentIdentifier.Coordinates.ArtifactID = r.Component.ComponentIdentifier.Coordinates.ArtifactID
	c.ComponentIdentifier.Coordinates.GroupID = r.Component.ComponentIdentifier.Coordinates.GroupID
	c.ComponentIdentifier.Coordinates.Version = r.Component.ComponentIdentifier.Coordinates.Version
	c.ComponentIdentifier.Coordinates.Extension = r.Component.ComponentIdentifier.Coordinates.Extension
	// c.Quarantined = false
	if highestViolation := r.HighestThreatPolicy(); highestViolation != nil {
		// c.HighestThreatLevel = true
		c.ThreatLevel = highestViolation.ThreatLevel
		c.PolicyName = highestViolation.PolicyName
	}
	return c
}

// RmItemToIQComponent converts a repo item to an IQ component
func RmItemToIQComponent(rm nexusrm.RepositoryItem) nexusiq.Component {
	var iqc nexusiq.Component
	switch rm.Format {
	case "maven2":
		iqc.ComponentIdentifier.Format = "maven"
		iqc.ComponentIdentifier.Coordinates.Extension = "jar"
	case "rubygems":
		iqc.ComponentIdentifier.Format = "gem"
	case "npm":
		iqc.ComponentIdentifier.Format = "npm"
		iqc.ComponentIdentifier.Coordinates.Extension = "tgz"
	case "pipy":
		iqc.ComponentIdentifier.Format = "pypi"
	default:
		iqc.ComponentIdentifier.Format = rm.Format
	}
	iqc.ComponentIdentifier.Coordinates.ArtifactID = rm.Name
	iqc.ComponentIdentifier.Coordinates.GroupID = rm.Group
	iqc.ComponentIdentifier.Coordinates.Version = rm.Version
	iqc.Hash = rm.Hash()
	return iqc
}
*/
