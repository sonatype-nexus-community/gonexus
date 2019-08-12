package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var dummyRoles = []Role{
	{
		ID:          "2cb71b3468d649789163ea2e212b541e",
		Name:        "Application Evaluator",
		Description: "Evaluates applications and views policy violation summary results.",
	},
	{
		ID:          "90c7c98683b4471cb77a916744540bcc",
		Name:        "Component Evaluator",
		Description: "Evaluates individual components and views policy violation results for a specified application.",
	},
	{
		ID:          "1b92fae3e55a411793a091fb821c422d",
		Name:        "System Administrator",
		Description: "Manages system configuration and users.",
	},
	{
		ID:          "1da70fae1fd54d6cb7999871ebdb9a36",
		Name:        "Developer",
		Description: "Views all information for their assigned organization or application.",
	},
	{
		ID:          "1cddabf7fdaa47d6833454af10e0a3ef",
		Name:        "Owner",
		Description: "Manages assigned organizations, applications, policies, and policy violations.",
	},
	{
		ID:          "a4fe90b9655643e7a4b1bb488f86a627",
		Name:        "SubContractor",
		Description: "SubContractor",
	},
	{
		ID:          "1556f306b2424ef5b447bfdd28249960",
		Name:        "UXDev",
		Description: "UXDev",
	},
	{
		ID:          "94d2b05ebeef4207bb55cbfc9de539b9",
		Name:        "app creator",
		Description: "test",
	},
	{
		ID:          "2c20beff71884f6a9fb7dbd6b07d3728",
		Name:        "arstar",
		Description: "arst",
	},
}

func rolesTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		resp, err := json.Marshal(rolesResponse{dummyRoles})
		if err != nil {
			t.Fatal(err)
			return
		}

		fmt.Fprintln(w, string(resp))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func rolesTestIQ(t *testing.T) (iq IQ, mock *httptest.Server) {
	return newTestIQ(t, rolesTestFunc)
}

func TestRoles(t *testing.T) {
	iq, mock := rolesTestIQ(t)
	defer mock.Close()

	got, err := Roles(iq)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, dummyRoles) {
		t.Error("Did not get expected roles")
	}
}

func TestRoleByName(t *testing.T) {
	iq, mock := rolesTestIQ(t)
	defer mock.Close()

	want := dummyRoles[0]

	got, err := RoleByName(iq, want.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected role")
	}
}
