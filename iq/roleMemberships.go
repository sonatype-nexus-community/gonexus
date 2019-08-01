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

	restRoleMembersOrgGet = "api/v2/roleMemberships/organization/%s"
	restRoleMembersAppGet = "api/v2/roleMemberships/application/%s"
	// restRoleMembersRepos  = "api/v2/roleMemberships/repository_container"
	// restRoleMembersGlobal = "api/v2/roleMemberships/global"

	restRoleMembersOrgUserPut  = "api/v2/roleMemberships/organization/%s/role/%s/user/%s"
	restRoleMembersOrgGroupPut = "api/v2/roleMemberships/organization/%s/role/%s/group/%s"

//* GET api/v2/organizations/{organizationInternalId}/roleMembers
// PUT api/v2/organizations/{organizationInternalId}/roleMembers
// GET api/v2/applications/{applicationInternalId}/roleMembers
// PUT api/v2/applications/{applicationInternalId}/roleMembers

// After 70
//* GET api/v2/roleMemberships/application/{applicationInternalId}
//* GET api/v2/roleMemberships/organization/{organizationId}
// GET api/v2/roleMemberships/repository_container
// GET api/v2/roleMemberships/global

// PUT api/v2/roleMemberships/organization/{organizationId}/role/{roleId}/user/{userName}
// PUT api/v2/roleMemberships/organization/{organizationId}/role/{roleId}/group/{groupName}
// PUT api/v2/roleMemberships/application/{applicationInternalId}/role/{roleId}/user/{userName}
// PUT api/v2/roleMemberships/application/{applicationInternalId}/role/{roleId}/group/{groupName}
// PUT api/v2/roleMemberships/repository_container/role/{roleId}/user/{userName}
// PUT api/v2/roleMemberships/repository_container/role/{roleId}/group/{groupName}
// PUT api/v2/roleMemberships/global/role/{roleId}/user/{userName}
// PUT api/v2/roleMemberships/global/role/{roleId}/group/{groupName}
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

func newUserMapping(roleID, name string) memberMappings {
	return newMapping(roleID, MemberTypeUser, name)
}

func newGroupMapping(roleID, name string) memberMappings {
	return newMapping(roleID, MemberTypeGroup, name)
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
		return nil, fmt.Errorf("could not retrieve role mapping: %v", err)
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
			endpoint = fmt.Sprintf(restRoleMembersOrgUserPut, org.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersOrgGroupPut, org.ID, role.ID, member)
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
	return nil, nil
}

// SetApplicationUser sets the role and user that can have access to an application
func SetApplicationUser(iq IQ, name, role, user string) error {
	return nil
}

// SetApplicationGroup sets the role and group that can have access to an application
func SetApplicationGroup(iq IQ, name, role, group string) error {
	return nil
}

// RoleMemberships returns the member mappings of all the things
