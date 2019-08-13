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
			RoleID:  dummyRoles[0].ID,
			Members: []Member{},
		},
		{
			RoleID: dummyRoles[1].ID,
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
			RoleID:  dummyRoles[2].ID,
			Members: []Member{},
		},
		{
			RoleID:  dummyRoles[3].ID,
			Members: []Member{},
		},
		{
			RoleID: dummyRoles[0].ID,
			Members: []Member{
				{
					Type:            MemberTypeUser,
					UserOrGroupName: "testrina",
				},
			},
		},
	},
}

var dummyRoleMappingsApps = map[string][]MemberMapping{
	dummyApps[0].ID: []MemberMapping{
		{
			RoleID:  dummyRoles[0].ID,
			Members: []Member{},
		},
		{
			RoleID: dummyRoles[1].ID,
			Members: []Member{
				{
					Type:            MemberTypeUser,
					UserOrGroupName: "foo",
				},
			},
		},
		{
			RoleID: dummyRoles[0].ID,
			Members: []Member{
				{
					Type:            MemberTypeUser,
					UserOrGroupName: "le test",
				},
			},
		},
	},
}

var dummyRoleMappingsRepos = []MemberMapping{
	{
		RoleID: dummyRoles[0].ID,
		Members: []Member{
			{
				Type:            MemberTypeGroup,
				UserOrGroupName: "oof",
			},
		},
	},
	{
		RoleID: dummyRoles[1].ID,
		Members: []Member{
			{
				Type:            MemberTypeUser,
				UserOrGroupName: "foo",
			},
		},
	},
}

var dummyRoleMappingsGlobal = []MemberMapping{
	{
		RoleID: dummyRoles[0].ID,
		Members: []Member{
			{
				OwnerID:         "global",
				OwnerType:       "GLOBAL",
				Type:            MemberTypeGroup,
				UserOrGroupName: "oof",
			},
		},
	},
	{
		RoleID: dummyRoles[1].ID,
		Members: []Member{
			{
				OwnerID:         "global",
				OwnerType:       "GLOBAL",
				Type:            MemberTypeUser,
				UserOrGroupName: "foo",
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
		var found bool
		var mappings []MemberMapping
		switch {
		case strings.HasPrefix(r.URL.Path[1:], restRoleMembersOrgGet[:len(restRoleMembersOrgGet)-2]):
			id := path.Base(r.URL.Path)
			mappings, found = dummyRoleMappingsOrgs[id]
		case strings.HasPrefix(r.URL.Path[1:], restRoleMembersAppGet[:len(restRoleMembersAppGet)-2]):
			id := path.Base(r.URL.Path)
			mappings, found = dummyRoleMappingsApps[id]
		case strings.HasPrefix(r.URL.Path[1:], restRoleMembersReposGet[:len(restRoleMembersReposGet)-2]):
			found = true
			mappings = dummyRoleMappingsRepos
		case strings.HasPrefix(r.URL.Path[1:], restRoleMembersGlobalGet[:len(restRoleMembersGlobalGet)-2]):
			found = true
			mappings = dummyRoleMappingsGlobal
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
	case r.Method == http.MethodPut && strings.Contains(r.URL.Path, "global"):
		fallthrough
	case r.Method == http.MethodPut && strings.Contains(r.URL.Path, "repository_container"):
		pathParts := strings.Split(r.URL.Path[1:], "/")

		//* api/v2/roleMemberships/[repository_container|global]/role/{roleId}/[user|group]/{name}
		authType := pathParts[3]
		roleID := pathParts[5]
		memberType := pathParts[6]
		memberName := pathParts[7]

		newTestMapping := MemberMapping{
			RoleID: roleID,
			Members: []Member{
				{
					Type:            strings.ToUpper(memberType),
					UserOrGroupName: memberName,
				},
			},
		}

		// TODO deal with duplicates?

		switch authType {
		case "global":
			newTestMapping.Members[0].OwnerID = "global"
			newTestMapping.Members[0].OwnerType = "GLOBAL"
			dummyRoleMappingsGlobal = append(dummyRoleMappingsGlobal, newTestMapping)
		case "repository_container":
			dummyRoleMappingsRepos = append(dummyRoleMappingsRepos, newTestMapping)
		}

	case r.Method == http.MethodPut:
		pathParts := strings.Split(r.URL.Path[1:], "/")

		//* api/v2/roleMemberships/[organization|application]/{id}/role/{roleId}/[user|group]/{name}
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
					Type:            strings.ToUpper(memberType),
					UserOrGroupName: memberName,
				},
			},
		}

		// TODO deal with duplicates?

		dummyMappings[id] = append(dummyMappings[id], newTestMapping)
	case r.Method == http.MethodDelete:
		pathParts := strings.Split(r.URL.Path[1:], "/")

		var (
			roleID, memberType, memberName string
			mappings                       []MemberMapping
		)
		switch {
		case strings.Contains(r.URL.Path, "global"):
			fallthrough
		case strings.Contains(r.URL.Path, "repository_container"):
			//* api/v2/roleMemberships/[repository_container|global]/role/{roleId}/[user|group]/{name}
			authType := pathParts[3]
			roleID = pathParts[5]
			memberType = strings.ToUpper(pathParts[6])
			memberName = pathParts[7]
			switch authType {
			case "repository_container":
				mappings = dummyRoleMappingsRepos
			case "global":
				mappings = dummyRoleMappingsGlobal
			}
		default:
			//* api/v2/roleMemberships/[organization|application]/{id}/role/{roleId}/[user|group]/{name}
			authType := pathParts[3]
			id := pathParts[4]
			roleID = pathParts[6]
			memberType = strings.ToUpper(pathParts[7])
			memberName = pathParts[8]

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
			mappings = dummyMappings[id]
		}

		for _, mapping := range mappings {
			if mapping.RoleID == roleID {
				for i, member := range mapping.Members {
					if member.Type == memberType && member.UserOrGroupName == memberName {
						copy(mapping.Members[i:], mapping.Members[i+1:])
						mapping.Members[len(mapping.Members)-1] = Member{}
						mapping.Members = mapping.Members[:len(mapping.Members)-1]

						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
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

func testGetOrganizationAuthorizationsByRole(t *testing.T, iq IQ) {
	t.Helper()
	role := dummyRoles[0]

	want := make([]MemberMapping, 0)
	for _, v := range dummyRoleMappingsOrgs {
		for _, m := range v {
			if m.RoleID == role.ID {
				want = append(want, m)
			}
		}
	}

	got, err := OrganizationAuthorizationsByRole(iq, role.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected organization mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func TestOrganizationAuthorizations(t *testing.T) {
	testWithDeprecated(t, testGetOrganizationAuthorizations)
}

func TestOrganizationAuthorizationsByRole(t *testing.T) {
	testWithDeprecated(t, testGetOrganizationAuthorizationsByRole)
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

func testGetApplicationAuthorizationsByRole(t *testing.T, iq IQ) {
	t.Helper()
	role := dummyRoles[0]

	want := make([]MemberMapping, 0)
	for _, v := range dummyRoleMappingsApps {
		for _, m := range v {
			if m.RoleID == role.ID {
				want = append(want, m)
			}
		}
	}

	got, err := ApplicationAuthorizationsByRole(iq, role.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected application mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func TestApplicationAuthorizations(t *testing.T) {
	testWithDeprecated(t, testGetApplicationAuthorizations)
}

func TestApplicationAuthorizationsByRole(t *testing.T) {
	testWithDeprecated(t, testGetApplicationAuthorizationsByRole)
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

func TestSetOrganizationUser(t *testing.T) {
	testWithDeprecated(t, testSetOrganizationUser)
}

func TestSetOrganizationGroup(t *testing.T) {
	testWithDeprecated(t, testSetOrganizationGroup)
}

func TestSetApplicationUser(t *testing.T) {
	testWithDeprecated(t, testSetApplicationUser)
}

func TestSetApplicationGroup(t *testing.T) {
	testWithDeprecated(t, testSetApplicationGroup)
}

func testRevoke(t *testing.T, iq IQ, authType, memberType string) {
	t.Helper()
	role := dummyRoles[0]
	name := "deleteMe"
	/*
		want := MemberMapping{
			RoleID: role.ID,
			Members: []Member{
				{
					Type:            memberType,
					UserOrGroupName: name,
				},
			},
		}
	*/

	var mappings []MemberMapping
	var err error
	switch authType {
	case "organization":
		dummyOrgName := dummyOrgs[0].Name
		switch memberType {
		case MemberTypeUser:
			err = SetOrganizationUser(iq, dummyOrgName, role.Name, name)
			if err == nil {
				err = RevokeOrganizationUser(iq, dummyOrgName, role.Name, name)
			}
		case MemberTypeGroup:
			err = SetOrganizationGroup(iq, dummyOrgName, role.Name, name)
			if err == nil {
				t.Log("HERE1")
				err = RevokeOrganizationGroup(iq, dummyOrgName, role.Name, name)
			}
		}
		if err == nil {
			mappings, err = OrganizationAuthorizations(iq, dummyOrgName)
		}
	case "application":
		dummyAppName := dummyApps[0].PublicID
		switch memberType {
		case MemberTypeUser:
			err = SetApplicationUser(iq, dummyAppName, role.Name, name)
			if err == nil {
				err = RevokeApplicationUser(iq, dummyAppName, role.Name, name)
			}
		case MemberTypeGroup:
			err = SetApplicationGroup(iq, dummyAppName, role.Name, name)
			if err == nil {
				err = RevokeApplicationGroup(iq, dummyAppName, role.Name, name)
			}
		}
		if err == nil {
			mappings, err = ApplicationAuthorizations(iq, dummyAppName)
		}
	case "repository_container":
		switch memberType {
		case MemberTypeUser:
			err = SetRepositoriesUser(iq, role.Name, name)
			if err == nil {
				err = RevokeRepositoriesUser(iq, role.Name, name)
			}
		case MemberTypeGroup:
			err = SetRepositoriesGroup(iq, role.Name, name)
			if err == nil {
				err = RevokeRepositoriesGroup(iq, role.Name, name)
			}
		}
		if err == nil {
			mappings, err = RepositoriesAuthorizations(iq)
		}
	case "global":
		switch memberType {
		case MemberTypeUser:
			err = SetGlobalUser(iq, role.Name, name)
			if err == nil {
				err = RevokeGlobalUser(iq, role.Name, name)
			}
		case MemberTypeGroup:
			err = SetGlobalGroup(iq, role.Name, name)
			if err == nil {
				err = RevokeGlobalGroup(iq, role.Name, name)
			}
		}
		if err == nil {
			mappings, err = GlobalAuthorizations(iq)
		}
	}
	if err != nil {
		t.Error(err)
	}

	for _, mapping := range mappings {
		if mapping.RoleID == role.ID {
			for _, member := range mapping.Members {
				if member.Type == memberType && member.UserOrGroupName == name {
					t.Error("found mapping which should have been revoked")
				}
			}
		}
	}
}

func testRevokeOrganizationUser(t *testing.T, iq IQ) {
	t.Helper()
	testRevoke(t, iq, "organization", MemberTypeUser)
}

func testRevokeOrganizationGroup(t *testing.T, iq IQ) {
	t.Helper()
	testRevoke(t, iq, "organization", MemberTypeGroup)
}

func testRevokeApplicationUser(t *testing.T, iq IQ) {
	t.Helper()
	testRevoke(t, iq, "application", MemberTypeUser)
}

func testRevokeApplicationGroup(t *testing.T, iq IQ) {
	t.Helper()
	testRevoke(t, iq, "application", MemberTypeGroup)
}

func TestRevokeOrganizationUser(t *testing.T) {
	testWithDeprecated(t, testRevokeOrganizationUser)
}

func TestRevokeOrganizationGroup(t *testing.T) {
	testWithDeprecated(t, testRevokeOrganizationGroup)
}

func TestRevokeApplicationUser(t *testing.T) {
	testWithDeprecated(t, testRevokeApplicationUser)
}

func TestRevokeApplicationGroup(t *testing.T) {
	testWithDeprecated(t, testRevokeApplicationGroup)
}

func TestRepositoriesAuthorizations(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	got, err := RepositoriesAuthorizations(iq)
	if err != nil {
		t.Error(err)
	}

	want := dummyRoleMappingsRepos
	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected repositories mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func TestGetApplicationAuthorizationsByRole(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	role := dummyRoles[0]

	want := make([]MemberMapping, 0)
	for _, m := range dummyRoleMappingsRepos {
		if m.RoleID == role.ID {
			want = append(want, m)
		}
	}

	got, err := RepositoriesAuthorizationsByRole(iq, role.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected repositories mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func testSetRepositories(t *testing.T, memberType string) {
	t.Helper()
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	role := dummyRoles[0]
	memberName := "dummyDumDum"

	want := MemberMapping{
		RoleID: role.ID,
		Members: []Member{
			{
				Type:            memberType,
				UserOrGroupName: memberName,
			},
		},
	}

	var err error
	switch memberType {
	case MemberTypeUser:
		err = SetRepositoriesUser(iq, role.Name, memberName)
	case MemberTypeGroup:
		err = SetRepositoriesGroup(iq, role.Name, memberName)
	}
	if err != nil {
		t.Error(err)
	}

	got, err := RepositoriesAuthorizations(iq)
	if err != nil {
		t.Error(err)
	}

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

func TestSetRepositoriesUser(t *testing.T) {
	testSetRepositories(t, MemberTypeUser)
}

func TestSetRepositoriesGroup(t *testing.T) {
	testSetRepositories(t, MemberTypeGroup)
}

func TestRevokeRepositoriesUser(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	testRevoke(t, iq, "repository_container", MemberTypeUser)
}

func TestRevokeRepositoriesGroup(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	testRevoke(t, iq, "repository_container", MemberTypeGroup)
}

func testMembersByRole(t *testing.T, iq IQ) {
	role := dummyRoles[0]

	// Build the want slice
	want := make([]MemberMapping, 0)
	for _, v := range dummyRoleMappingsOrgs {
		for _, m := range v {
			if m.RoleID == role.ID {
				want = append(want, m)
			}
		}
	}
	for _, v := range dummyRoleMappingsApps {
		for _, m := range v {
			if m.RoleID == role.ID {
				want = append(want, m)
			}
		}
	}
	if hasRev70API(iq) {
		for _, m := range dummyRoleMappingsRepos {
			if m.RoleID == role.ID {
				want = append(want, m)
			}
		}
	}

	got, err := MembersByRole(iq, role.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected members for a role")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func TestMembersByRole(t *testing.T) {
	testWithDeprecated(t, testMembersByRole)
}

func TestGlobalAuthorizations(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	got, err := GlobalAuthorizations(iq)
	if err != nil {
		t.Error(err)
	}

	want := dummyRoleMappingsGlobal
	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected global mapping")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func testSetGlobal(t *testing.T, memberType string) {
	t.Helper()
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	role := dummyRoles[0]
	memberName := "dummyDumDum"

	want := MemberMapping{
		RoleID: role.ID,
		Members: []Member{
			{
				OwnerID:         "global",
				OwnerType:       "GLOBAL",
				Type:            memberType,
				UserOrGroupName: memberName,
			},
		},
	}

	var err error
	switch memberType {
	case MemberTypeUser:
		err = SetGlobalUser(iq, role.Name, memberName)
	case MemberTypeGroup:
		err = SetGlobalGroup(iq, role.Name, memberName)
	}
	if err != nil {
		t.Error(err)
	}

	got, err := GlobalAuthorizations(iq)
	if err != nil {
		t.Error(err)
	}

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

func TestSetGlobalUser(t *testing.T) {
	testSetGlobal(t, MemberTypeUser)
}

func TestSetGlobalGroup(t *testing.T) {
	testSetGlobal(t, MemberTypeGroup)
}

func TestRevokeGlobalUser(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	testRevoke(t, iq, "global", MemberTypeUser)
}

func TestRevokeGlobalGroup(t *testing.T) {
	iq, mock := roleMembershipsTestIQ(t, false)
	defer mock.Close()

	testRevoke(t, iq, "global", MemberTypeGroup)
}
