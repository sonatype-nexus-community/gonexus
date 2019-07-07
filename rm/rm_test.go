package nexusrm

import (
	"testing"
)

func getTestRM(t *testing.T) *RM {
	rm, err := New("http://localhost:8081", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}

	return rm
}
