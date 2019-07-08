package nexusiq

import (
	"fmt"
	// "testing"
)

func ExampleCreateOrganization() {
	iq, err := New("http://localhost:8070", "user", "password")
	if err != nil {
		panic(err)
	}

	orgID, err := CreateOrganization(iq, "DatabaseTeam")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Organization ID: %s\n", orgID)
}
