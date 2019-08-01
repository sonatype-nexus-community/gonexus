package nexusiq

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restRolesDeprecated = "api/v2/applications/roles" // Before r70
	restRoles           = "api/v2/roles"
)

type rolesResponse struct {
	Roles []Role `json:"roles"`
}

// Role describes an IQ role
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Roles returns a slice of all the roles in the IQ instance
func Roles(iq IQ) ([]Role, error) {
	body, resp, err := iq.Get(restRoles)
	if resp.StatusCode == http.StatusNotFound {
		body, _, err = iq.Get(restRolesDeprecated)
	}
	if err != nil {
		return nil, fmt.Errorf("could not retrieve roles: %v", err)
	}

	var list rolesResponse
	if err = json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("could not marshal roles response: %v", err)
	}

	return list.Roles, nil
}

// RoleByName returns the named role
func RoleByName(iq IQ, name string) (Role, error) {
	roles, err := Roles(iq)
	if err != nil {
		return Role{}, err
	}

	for _, r := range roles {
		if r.Name == name {
			return r, nil
		}
	}

	return Role{}, fmt.Errorf("did not find role with name %s", name)
}
