package nexusiq

import (
	nexus "github.com/sonatype-nexus-community/gonexus"
)

// IQ is the interface which allows interacting with an IQ server
type IQ interface {
	nexus.Client
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
