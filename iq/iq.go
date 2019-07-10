package nexusiq

import (
	"net/http"

	"github.com/hokiegeek/gonexus"
)

// IQ is the interface which any IQ Server implementation would need to satisfy
type IQ interface {
	Get(endpoint string) ([]byte, *http.Response, error)
	Post(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Put(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Del(endpoint string) (resp *http.Response, err error)
}

type iqServer struct {
	nexus.DefaultServer
}

// New creates a new IQ instance
func New(host, username, password string) (IQ, error) {
	iq := new(iqServer)
	iq.Host = host
	iq.Username = username
	iq.Password = password
	return iq, nil
}
