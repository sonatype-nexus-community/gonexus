package nexusrm

import (
	nexus "github.com/hokiegeek/gonexus"
)

// RM holds basic and state info of the Repository Manager server we will connect to
type RM struct {
	nexus.DefaultServer
}

// New creates a new Repository Manager instance
func New(host, username, password string) (rm *RM, err error) {
	rm = new(RM)
	rm.Host = host
	rm.Username = username
	rm.Password = password
	return
}
