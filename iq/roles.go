package nexusiq

import (
	"encoding/json"
	"fmt"
)

const restRoles = "api/v2/applications/roles"

type rolesResponse struct {
	Roles []Role `json:"roles"`
}

// Role describes an IQ role
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

const (
	MemberTypeUser  = "USER"
	MemberTypeGroup = "GROUP"
)

type MemberMappings struct {
	MemberMappings []MemberMapping `json:"memberMappings"`
}

type MemberMapping struct {
	RoleID  string   `json:"roleId"`
	Members []Member `json:"members"`
}

type Member struct {
	Type            string `json:"type"`
	UserOrGroupName string `json:"userOrGroupName"`
}

// Roles returns a slice of all the roles in the IQ instance
func Roles(iq IQ) ([]Role, error) {
	body, _, err := iq.Get(restRoles)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve roles: %v", err)
	}

	var resp rolesResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("could not marshal roles response: %v", err)
	}

	return resp.Roles, nil
}
