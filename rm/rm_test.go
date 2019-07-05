package nexusrm

import (
	"fmt"
)

func ExampleRM_ListComponents() {
	rm, err := New("http://localhost:8081", "user", "password")
	if err != nil {
		panic(err)
	}

	items, err := rm.ListComponents("maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
