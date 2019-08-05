package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
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
	case r.Method == http.MethodHead:
		if !strings.HasPrefix(r.URL.Path[1:], "api/v2/organizations") || !strings.HasSuffix(r.URL.Path, "roleMembers") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodGet:
		pathParts := strings.Split(r.URL.Path, "/")

		authType := pathParts[3]
		id := pathParts[4]

		var found bool
		var mappings []MemberMapping
		switch authType {
		case "organizations":
			mappings, found = dummyRoleMappingsOrgs[id]
		case "applications":
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
	case r.Method == http.MethodPut:
		pathParts := strings.Split(r.URL.Path, "/")

		authType := pathParts[3]
		id := pathParts[4]

		var dummyMappings map[string][]MemberMapping
		switch authType {
		case "organizations":
			dummyMappings = dummyRoleMappingsOrgs
		case "applications":
			dummyMappings = dummyRoleMappingsApps
		}
		if _, ok := dummyMappings[id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		var mapping memberMappings
		if err = json.Unmarshal(body, &mapping); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		dummyMappings[id] = mapping.MemberMappings
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func roleMembershipsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodHead:
		if !strings.HasPrefix(r.URL.Path[1:], "api/v2/roleMemberships") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodGet:
		id := path.Base(r.URL.Path) // Could be better?
		var found bool
		var mappings []MemberMapping
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restRoleMembersOrgGet[:len(restRoleMembersOrgGet)-2]):
			mappings, found = dummyRoleMappingsOrgs[id]
		case strings.HasPrefix(r.URL.Path[1:], restRoleMembersAppGet[:len(restRoleMembersAppGet)-2]):
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
	case r.Method == http.MethodPut:
		pathParts := strings.Split(r.URL.Path[1:], "/")

		//* api/v2/roleMemberships/[organization|application]/{organizationId}/role/{roleId}/[user|group]/{userName}
		authType := pathParts[3]
		id := pathParts[4]
		roleID := pathParts[6]
		memberType := pathParts[7]
		memberName := pathParts[8]

		var dummyMappings map[string][]MemberMapping
		switch authType {
		case "organization":
			dummyMappings = dummyRoleMappingsOrgs
		case "application":
			dummyMappings = dummyRoleMappingsApps
		}
		if _, ok := dummyMappings[id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		newTestMapping := MemberMapping{
			RoleID: roleID,
			Members: []Member{
				{
					Type:            memberType,
					UserOrGroupName: memberName,
				},
			},
		}

		// TODO deal with duplicates?

		dummyMappings[id] = append(dummyMappings[id], newTestMapping)
	case r.Method == http.MethodDelete:
		// TODO
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func roleMembershipsTestIQ(t *testing.T, useDeprecated bool) (IQ, *httptest.Server) {
	return newTestIQ(t, func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path[1:] == restOrganization:
			organizationTestFunc(t, w, r)
		case r.URL.Path[1:] == restApplication:
			applicationTestFunc(t, w, r)
		case r.URL.Path[1:] == restRoles:
			rolesTestFunc(t, w, r)
		default:
			if useDeprecated {
				roleMembershipsDeprecatedTestFunc(t, w, r)
			} else {
				roleMembershipsTestFunc(t, w, r)
			}
		}
	})
}

func testWithDeprecated(t *testing.T, subtest func(*testing.T, IQ)) {
	t.Helper()
	t.Run("role memberships < r70", func(t *testing.T) {
		iq, mock := roleMembershipsTestIQ(t, true)
		subtest(t, iq)
		defer mock.Close()
	})
	t.Run("role memberships >= r70", func(t *testing.T) {
		iq, mock := roleMembershipsTestIQ(t, false)
		subtest(t, iq)
		defer mock.Close()
	})
}

func testGetOrganizationAuthorizations(t *testing.T, iq IQ) {
	t.Helper()
	dummyIdx := 0

	got, err := OrganizationAuthorizations(iq, dummyOrgs[dummyIdx].Name)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(got)

	want := dummyRoleMappingsOrgs[dummyOrgs[dummyIdx].ID]
	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected organization mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func testGetApplicationAuthorizations(t *testing.T, iq IQ) {
	t.Helper()
	dummyIdx := 0

	got, err := ApplicationAuthorizations(iq, dummyApps[dummyIdx].PublicID)
	if err != nil {
		t.Error(err)
	}

	want := dummyRoleMappingsApps[dummyApps[dummyIdx].ID]
	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected application mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func testSetAuth(t *testing.T, iq IQ, authTarget string, memberType string) {
	t.Helper()
	dummyIdx := 0
	role := 0
	memberName := "dummyDumDum"

	want := MemberMapping{
		RoleID: dummyRoles[role].ID,
		Members: []Member{
			{
				Type:            memberType,
				UserOrGroupName: memberName,
			},
		},
	}

	var got []MemberMapping
	var err error

	switch authTarget {
	case "organization":
		switch memberType {
		case MemberTypeUser:
			err = SetOrganizationUser(iq, dummyOrgs[dummyIdx].Name, dummyRoles[role].Name, memberName)
		case MemberTypeGroup:
			err = SetOrganizationGroup(iq, dummyOrgs[dummyIdx].Name, dummyRoles[role].Name, memberName)
		}
		if err == nil {
			got, err = OrganizationAuthorizations(iq, dummyOrgs[dummyIdx].Name)
		}
	case "application":
		switch memberType {
		case MemberTypeUser:
			err = SetApplicationUser(iq, dummyApps[dummyIdx].PublicID, dummyRoles[role].Name, memberName)
		case MemberTypeGroup:
			err = SetApplicationGroup(iq, dummyApps[dummyIdx].PublicID, dummyRoles[role].Name, memberName)
		}
		if err == nil {
			got, err = ApplicationAuthorizations(iq, dummyApps[dummyIdx].PublicID)
		}
	}
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(got)

	var found bool
	for _, mapping := range got {
		if reflect.DeepEqual(mapping, want) {
			found = true
			break
		}
	}

	if !found {
		t.Error("User role mapping not updated")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func testSetOrganizationUser(t *testing.T, iq IQ) {
	t.Helper()
	testSetAuth(t, iq, "organization", MemberTypeUser)
}

func testSetOrganizationGroup(t *testing.T, iq IQ) {
	t.Helper()
	testSetAuth(t, iq, "organization", MemberTypeGroup)
}

func testSetApplicationUser(t *testing.T, iq IQ) {
	t.Helper()
	testSetAuth(t, iq, "organization", MemberTypeUser)
}

func testSetApplicationGroup(t *testing.T, iq IQ) {
	t.Helper()
	testSetAuth(t, iq, "organization", MemberTypeGroup)
}

func TestOrganizationAuthorizations(t *testing.T) {
	testWithDeprecated(t, testGetOrganizationAuthorizations)
}

func TestSetOrganizationUser(t *testing.T) {
	testWithDeprecated(t, testSetOrganizationUser)
}

func TestSetOrganizationGroup(t *testing.T) {
	testWithDeprecated(t, testSetOrganizationGroup)
}

func TestApplicationAuthorizations(t *testing.T) {
	testWithDeprecated(t, testGetApplicationAuthorizations)
}

func TestSetApplicationUser(t *testing.T) {
	testWithDeprecated(t, testSetApplicationUser)
}

func TestSetApplicationGroup(t *testing.T) {
	testWithDeprecated(t, testSetApplicationGroup)
}
