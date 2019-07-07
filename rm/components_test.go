package nexusrm

import (
	"fmt"
)

func ExampleGetComponents() {
	rm, err := New("http://localhost:8081", "user", "password")
	if err != nil {
		panic(err)
	}

	items, err := GetComponents(rm, "maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
