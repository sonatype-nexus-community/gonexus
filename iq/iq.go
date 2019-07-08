package nexusiq

import (
	"github.com/hokiegeek/gonexus"
)

// IQ holds basic and state info on the IQ Server we will connect to
type IQ struct {
	nexus.DefaultServer
}

// New creates a new IQ instance
func New(host, username, password string) (*IQ, error) {
	iq := new(IQ)
	iq.Host = host
	iq.Username = username
	iq.Password = password
	return iq, nil
}
