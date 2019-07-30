/*
Package nexusiq provides a number of functions that interact with the Nexus IQ REST API.
All functions require a new RM instance which can be instantiated as such:

	iq, err := nexusiq.New("http://localhost:8070", "username", "password")
	if err != nil {
	    panic(err)
	}
*/
package nexusiq
