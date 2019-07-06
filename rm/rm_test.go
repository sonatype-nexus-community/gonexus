package nexusrm

import (
	"fmt"
	"testing"
)

func getTestRM(t *testing.T) *RM {
	rm, err := New("http://localhost:8081", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}

	return rm
}

func TestGetRepositories(t *testing.T) {
	rm := getTestRM(t)

	repos, err := GetRepositories(rm)
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

	items, err := GetComponents(rm, "maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
