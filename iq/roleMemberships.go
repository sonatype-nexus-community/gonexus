package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// Before 70
	restRoleMembersOrgDeprecated = "api/v2/organizations/%s/roleMembers"
	restRoleMembersAppDeprecated = "api/v2/applications/%s/roleMembers"

	// After 70
	restRoleMembersOrgGet = "api/v2/roleMemberships/organization/%s"
	restRoleMembersAppGet = "api/v2/roleMemberships/application/%s"
	// restRoleMembersReposGet  = "api/v2/roleMemberships/repository_container"
	// restRoleMembersGlobalGet = "api/v2/roleMemberships/global"

	restRoleMembersOrgUser  = "api/v2/roleMemberships/organization/%s/role/%s/user/%s"
	restRoleMembersOrgGroup = "api/v2/roleMemberships/organization/%s/role/%s/group/%s"
	restRoleMembersAppUser  = "api/v2/roleMemberships/application/%s/role/%s/user/%s"
	restRoleMembersAppGroup = "api/v2/roleMemberships/application/%s/role/%s/group/%s"
	// restRoleMembersRepositoryUser  = "api/v2/roleMemberships/repository_container/role/{roleId}/user/{userName}"
	// restRoleMembersRepositoryGroup = "api/v2/roleMemberships/repository_container/role/{roleId}/group/{groupName}"
	// restRoleMembersGlobalUser      = "api/v2/roleMemberships/global/role/{roleId}/user/{userName}"
	// restRoleMembersGlobalGroup     = "api/v2/roleMemberships/global/role/{roleId}/group/{groupName}"
)

// Constants to describe a Member Type
const (
	MemberTypeUser  = "USER"
	MemberTypeGroup = "GROUP"
)

type memberMappings struct {
	MemberMappings []MemberMapping `json:"memberMappings"`
}

// MemberMapping describes a list of Members against a Role
type MemberMapping struct {
	RoleID  string   `json:"roleId"`
	Members []Member `json:"members"`
}

// Member describes a member to map with a role
type Member struct {
	Type            string `json:"type"`
	UserOrGroupName string `json:"userOrGroupName"`
}

func hasDeprecatedAPI(iq IQ) bool {
	api := fmt.Sprintf(restRoleMembersOrgDeprecated, rootOrganizationID)
	request, _ := iq.NewRequest("HEAD", api, nil)
	_, resp, _ := iq.Do(request)
	return resp.StatusCode != http.StatusNotFound
}

func newMapping(roleID, memberType, memberName string) memberMappings {
	return memberMappings{
		MemberMappings: []MemberMapping{
			{
				RoleID: roleID,
				Members: []Member{
					{
						Type:            memberType,
						UserOrGroupName: memberName,
					},
				},
			},
		},
	}
}

// OrganizationAuthorizations returns the member mappings of an organization
func OrganizationAuthorizations(iq IQ, name string) ([]MemberMapping, error) {
	org, err := GetOrganizationByName(iq, name)
	if err != nil {
		return nil, fmt.Errorf("could not find organization with name %s: %v", name, err)
	}

	var endpoint string
	if hasDeprecatedAPI(iq) {
		endpoint = fmt.Sprintf(restRoleMembersOrgDeprecated, org.ID)
	} else {
		endpoint = fmt.Sprintf(restRoleMembersOrgGet, org.ID)
	}

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve role mapping for organization %s: %v", name, err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)

	return mappings.MemberMappings, err
}

func setOrganizationAuth(iq IQ, name, roleName, member, memberType string) error {
	org, err := GetOrganizationByName(iq, name)
	if err != nil {
		return fmt.Errorf("could not find organization with name %s: %v", name, err)
	}

	role, err := RoleByName(iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	var payload io.Reader
	if hasDeprecatedAPI(iq) {
		endpoint = fmt.Sprintf(restRoleMembersOrgDeprecated, org.ID)
		buf, err := json.Marshal(newMapping(role.ID, memberType, member))
		if err != nil {
			return fmt.Errorf("could not create mapping: %v", err)
		}
		payload = bytes.NewBuffer(buf)
	} else {
		switch memberType {
		case MemberTypeUser:
			endpoint = fmt.Sprintf(restRoleMembersOrgUser, org.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersOrgGroup, org.ID, role.ID, member)
		}
	}

	_, _, err = iq.Put(endpoint, payload)
	if err != nil {
		return fmt.Errorf("could not update organization role mapping: %v", err)
	}

	return nil
}

// SetOrganizationUser sets the role and user that can have access to an organization
func SetOrganizationUser(iq IQ, name, roleName, user string) error {
	return setOrganizationAuth(iq, name, roleName, user, MemberTypeUser)
}

// SetOrganizationGroup sets the role and group that can have access to an organization
func SetOrganizationGroup(iq IQ, name, roleName, group string) error {
	return setOrganizationAuth(iq, name, roleName, group, MemberTypeGroup)
}

// ApplicationAuthorizations returns the member mappings of an application
func ApplicationAuthorizations(iq IQ, name string) ([]MemberMapping, error) {
	app, err := GetApplicationByPublicID(iq, name)
	if err != nil {
		return nil, fmt.Errorf("could not find application with name %s: %v", name, err)
	}

	var endpoint string
	if hasDeprecatedAPI(iq) {
		endpoint = fmt.Sprintf(restRoleMembersAppDeprecated, app.ID)
	} else {
		endpoint = fmt.Sprintf(restRoleMembersAppGet, app.ID)
	}

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve role mapping for application %s: %v", name, err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)

	return mappings.MemberMappings, err
}

func setApplicationAuth(iq IQ, name, roleName, member, memberType string) error {
	app, err := GetApplicationByPublicID(iq, name)
	if err != nil {
		return fmt.Errorf("could not find application with name %s: %v", name, err)
	}

	role, err := RoleByName(iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	var payload io.Reader
	if hasDeprecatedAPI(iq) {
		endpoint = fmt.Sprintf(restRoleMembersAppDeprecated, app.ID)
		buf, err := json.Marshal(newMapping(role.ID, memberType, member))
		if err != nil {
			return fmt.Errorf("could not create mapping: %v", err)
		}
		payload = bytes.NewBuffer(buf)
	} else {
		switch memberType {
		case MemberTypeUser:
			endpoint = fmt.Sprintf(restRoleMembersAppUser, app.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersAppGroup, app.ID, role.ID, member)
		}
	}

	_, _, err = iq.Put(endpoint, payload)
	if err != nil {
		return fmt.Errorf("could not update organization role mapping: %v", err)
	}

	return nil
}

// SetApplicationUser sets the role and user that can have access to an application
func SetApplicationUser(iq IQ, name, roleName, user string) error {
	return setApplicationAuth(iq, name, roleName, user, MemberTypeUser)
}

// SetApplicationGroup sets the role and group that can have access to an application
func SetApplicationGroup(iq IQ, name, roleName, group string) error {
	return setApplicationAuth(iq, name, roleName, group, MemberTypeGroup)
}
