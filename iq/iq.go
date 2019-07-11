package nexusiq

import (
	"net/http"

	"github.com/hokiegeek/gonexus"
)

// IQ is the interface which any IQ Server implementation would need to satisfy
type IQ interface {
	Get(endpoint string) ([]byte, *http.Response, error)
	Post(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Put(endpoint string, payload []byte) (*http.Response, error)
	Del(endpoint string) (*http.Response, error)
}

type iqClient struct {
	nexus.DefaultClient
}

// New creates a new IQ instance
func New(host, username, password string) (IQ, error) {
	iq := new(iqClient)
	iq.Host = host
	iq.Username = username
	iq.Password = password
	return iq, nil
}
