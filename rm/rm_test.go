package nexusrm

import (
	"fmt"
	"testing"
)

func TestGetRepositories(t *testing.T) {
	rm, err := New("http://localhost:8081", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}

	repos, err := rm.GetRepositories()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q\n", repos)
}

func ExampleRM_GetComponents() {
	rm, err := New("http://localhost:8081", "user", "password")
	if err != nil {
		panic(err)
	}

	items, err := rm.GetComponents("maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
