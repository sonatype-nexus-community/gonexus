package nexusrm

import (
	"testing"
)

func TestGetRepositories(t *testing.T) {
	rm := getTestRM(t)

	repos, err := GetRepositories(rm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%q\n", repos)
}
