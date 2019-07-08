package nexusiq

import (
	"testing"
)

func getTestIQ(t *testing.T) *IQ {
	iq, err := New("http://localhost:8070", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}
	iq.Debug = true

	return iq
}
