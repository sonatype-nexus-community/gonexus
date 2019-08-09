/*
Package nexusrm provides a number of functions that interact with the Nexus Repository Manager REST API.
All functions require a new RM instance which can be instantiated as such:

	rm, err := nexusrm.New("http://localhost:8081", "username", "password")
	if err != nil {
	    panic(err)
	}
*/
package nexusrm
