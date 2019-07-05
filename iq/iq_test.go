package nexusiq

import (
	"fmt"
)

func ExampleIQ_CreateOrganization() {
	iq, err := New("http://localhost:8070", "user", "password")
	if err != nil {
		panic(err)
	}

	orgID, err := iq.CreateOrganization("DatabaseTeam")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Organization ID: %s\n", orgID)
}
