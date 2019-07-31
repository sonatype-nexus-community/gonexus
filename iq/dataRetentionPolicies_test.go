package nexusiq

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var dummyRetentionPolicies = DataRetentionPolicies{}

func dataRetentionPoliciesTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
	case r.Method == http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func dataRetentionPoliciesTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restOrganization):
			organizationTestFunc(t, w, r)
		default:
			dataRetentionPoliciesTestFunc(t, w, r)
		}
	})
}

func TestGetRetentionPolicies(t *testing.T) {
	t.Skip("TODO")
	// iq, mock := dataRetentionPoliciesTestIQ(t)
	// defer mock.Close()
}

func TestSetRetentionPolicies(t *testing.T) {
	t.Skip("TODO")
}
