package nexusiq

import (
	"testing"
)

func TestCreateSourceControlEntry(t *testing.T) {
	iq := getTestIQ(t)

	err := CreateSourceControlEntry(iq, "WebGoat", "https://github.com/HokieGeek/WebGoat", "c564532366dc0fdd1abece4cd1d6f9fd4abe4840")
	if err != nil {
		t.Error(err)
	}
}

func TestGetSourceControlEntry(t *testing.T) {
	iq := getTestIQ(t)

	entries, err := GetSourceControlEntry(iq, "WebGoat")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", entries)
}

func TestGetAllSourceControlEntries(t *testing.T) {
	iq := getTestIQ(t)

	entries, err := GetAllSourceControlEntries(iq)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", entries)
}
