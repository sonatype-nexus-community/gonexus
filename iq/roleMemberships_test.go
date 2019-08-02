package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var dummyRoleMappingsOrgs = map[string][]MemberMapping{
	dummyOrgs[0].ID: []MemberMapping{
		{
			RoleID:  "2cb71b3468d649789163ea2e212b541e",
			Members: []Member{},
		},
		{
			RoleID: "90c7c98683b4471cb77a916744540bcc",
			Members: []Member{
				{
					Type:            MemberTypeUser,
					UserOrGroupName: "foo",
				},
				{
					Type:            MemberTypeGroup,
					UserOrGroupName: "bar",
				},
			},
		},
		{
			RoleID:  "1da70fae1fd54d6cb7999871ebdb9a36",
			Members: []Member{},
		},
		{
			RoleID:  "1cddabf7fdaa47d6833454af10e0a3ef",
			Members: []Member{},
		},
	},
}

var dummyRoleMappingsApps = map[string][]MemberMapping{
	dummyApps[0].ID: []MemberMapping{
		{
			RoleID:  "2cb71b3468d649789163ea2e212b541e",
			Members: []Member{},
		},
		{
			RoleID: "90c7c98683b4471cb77a916744540bcc",
			Members: []Member{
				{
					Type:            MemberTypeUser,
					UserOrGroupName: "foo",
				},
			},
		},
	},
}

func roleMembershipsDeprecatedTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		id := strings.Split(r.URL.Path, "/")[4] // Could be better
		var found bool
		var mappings []MemberMapping
		switch {
		case strings.HasPrefix(r.URL.Path[1:], "api/v2/organizations"):
			mappings, found = dummyRoleMappingsOrgs[id]
		case strings.HasPrefix(r.URL.Path[1:], "api/v2/applications"):
			mappings, found = dummyRoleMappingsApps[id]
		}
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		buf, err := json.Marshal(memberMappings{mappings})
		if err != nil {
			t.Fatal(err)
			return
		}

		fmt.Fprintln(w, string(buf))
	// case r.Method == http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func roleMembershipsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func roleMembershipsTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path[1:] == restOrganization:
			organizationTestFunc(t, w, r)
		case strings.HasPrefix(r.URL.Path[1:], "api/v2/organizations"), strings.HasPrefix(r.URL.Path[1:], "api/v2/applications"):
			roleMembershipsDeprecatedTestFunc(t, w, r)
		default:
			roleMembershipsTestFunc(t, w, r)
		}
	})
}

func TestOrganizationAuthorizations(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t)
	defer mock.Close()

	org := 0

	got, err := OrganizationAuthorizations(iq, dummyOrgs[org].Name)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)

	if !reflect.DeepEqual(got, dummyRoleMappingsOrgs[dummyOrgs[0].ID]) {
		t.Error(err)
	}
}
