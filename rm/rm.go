package nexusrm

import (
	"net/http"

	"github.com/hokiegeek/gonexus"
)

// RM is the interface which any Repository Manager implementation would need to satisfy
type RM interface {
	Get(endpoint string) ([]byte, *http.Response, error)
	Post(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Put(endpoint string, payload []byte) ([]byte, *http.Response, error)
	Del(endpoint string) (resp *http.Response, err error)
}

type rmServer struct {
	nexus.DefaultServer
}

// New creates a new Repository Manager instance
func New(host, username, password string) (RM, error) {
	rm := new(rmServer)
	rm.Host = host
	rm.Username = username
	rm.Password = password
	return rm, nil
}
