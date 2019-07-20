package nexusrm

import (
	"github.com/sonatype-nexus-community/gonexus"
)

// RM is the interface which any Repository Manager implementation would need to satisfy
type RM interface {
	nexus.Client
}

type rmClient struct {
	nexus.DefaultClient
}

// New creates a new Repository Manager instance
func New(host, username, password string) (RM, error) {
	rm := new(rmClient)
	rm.Host = host
	rm.Username = username
	rm.Password = password
	return rm, nil
}
