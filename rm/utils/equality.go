package utils

import (
	"github.com/go-test/deep"
)

// IsEqual checks equality between two variables
func IsEqual(a, b interface{}) bool {
	if diff := deep.Equal(a, b); diff != nil {
		return false
	}
	return true
}
