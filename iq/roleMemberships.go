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
	restRoleMembersOrgGet    = "api/v2/roleMemberships/organization/%s"
	restRoleMembersAppGet    = "api/v2/roleMemberships/application/%s"
	restRoleMembersReposGet  = "api/v2/roleMemberships/repository_container"
	restRoleMembersGlobalGet = "api/v2/roleMemberships/global"

	restRoleMembersOrgUser         = "api/v2/roleMemberships/organization/%s/role/%s/user/%s"
	restRoleMembersOrgGroup        = "api/v2/roleMemberships/organization/%s/role/%s/group/%s"
	restRoleMembersAppUser         = "api/v2/roleMemberships/application/%s/role/%s/user/%s"
	restRoleMembersAppGroup        = "api/v2/roleMemberships/application/%s/role/%s/group/%s"
	restRoleMembersRepositoryUser  = "api/v2/roleMemberships/repository_container/role/%s/user/%s"
	restRoleMembersRepositoryGroup = "api/v2/roleMemberships/repository_container/role/%s/group/%s"
	restRoleMembersGlobalUser      = "api/v2/roleMemberships/global/role/%s/user/%s"
	restRoleMembersGlobalGroup     = "api/v2/roleMemberships/global/role/%s/group/%s"
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
	OwnerID         string `json:"ownerId,omitempty"`
	OwnerType       string `json:"ownerType,omitempty"`
	Type            string `json:"type"`
	UserOrGroupName string `json:"userOrGroupName"`
}

func hasRev70API(iq IQ) bool {
	api := fmt.Sprintf(restRoleMembersOrgGet, RootOrganization)
	request, _ := iq.NewRequest("HEAD", api, nil)
	_, resp, _ := iq.Do(request)
	return resp.StatusCode != http.StatusNotFound
}

func newMapping(roleID, memberType, memberName string) MemberMapping {
	return MemberMapping{
		RoleID: roleID,
		Members: []Member{
			{
				Type:            memberType,
				UserOrGroupName: memberName,
			},
		},
	}
}

func newMappings(roleID, memberType, memberName string) memberMappings {
	return memberMappings{
		MemberMappings: []MemberMapping{newMapping(roleID, memberType, memberName)},
	}
}

func organizationAuthorizationsByID(iq IQ, orgID string) ([]MemberMapping, error) {
	var endpoint string
	if hasRev70API(iq) {
		endpoint = fmt.Sprintf(restRoleMembersOrgGet, orgID)
	} else {
		endpoint = fmt.Sprintf(restRoleMembersOrgDeprecated, orgID)
	}

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve role mapping for organization %s: %v", orgID, err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)

	return mappings.MemberMappings, err
}

func organizationAuthorizationsByRoleID(iq IQ, roleID string) ([]MemberMapping, error) {
	orgs, err := GetAllOrganizations(iq)
	if err != nil {
		return nil, fmt.Errorf("could not find organizations: %v", err)
	}

	mappings := make([]MemberMapping, 0)
	for _, org := range orgs {
		orgMaps, _ := organizationAuthorizationsByID(iq, org.ID)
		for _, m := range orgMaps {
			if m.RoleID == roleID {
				mappings = append(mappings, m)
			}
		}
	}

	return mappings, nil
}

// OrganizationAuthorizations returns the member mappings of an organization
func OrganizationAuthorizations(iq IQ, name string) ([]MemberMapping, error) {
	org, err := GetOrganizationByName(iq, name)
	if err != nil {
		return nil, fmt.Errorf("could not find organization with name %s: %v", name, err)
	}

	return organizationAuthorizationsByID(iq, org.ID)
}

// OrganizationAuthorizationsByRole returns the member mappings of all organizations which match the given role
func OrganizationAuthorizationsByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByName(iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	return organizationAuthorizationsByRoleID(iq, role.ID)
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
	if hasRev70API(iq) {
		switch memberType {
		case MemberTypeUser:
			endpoint = fmt.Sprintf(restRoleMembersOrgUser, org.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersOrgGroup, org.ID, role.ID, member)
		}
	} else {
		endpoint = fmt.Sprintf(restRoleMembersOrgDeprecated, org.ID)
		current, err := OrganizationAuthorizations(iq, name)
		if err != nil && current == nil {
			current = make([]MemberMapping, 0)
		}
		current = append(current, newMapping(role.ID, memberType, member))

		buf, err := json.Marshal(memberMappings{MemberMappings: current})
		if err != nil {
			return fmt.Errorf("could not create mapping: %v", err)
		}
		payload = bytes.NewBuffer(buf)
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

func applicationAuthorizationsByID(iq IQ, appID string) ([]MemberMapping, error) {
	var endpoint string
	if hasRev70API(iq) {
		endpoint = fmt.Sprintf(restRoleMembersAppGet, appID)
	} else {
		endpoint = fmt.Sprintf(restRoleMembersAppDeprecated, appID)
	}

	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve role mapping for application %s: %v", appID, err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)

	return mappings.MemberMappings, err
}

func applicationAuthorizationsByRoleID(iq IQ, roleID string) ([]MemberMapping, error) {
	apps, err := GetAllApplications(iq)
	if err != nil {
		return nil, fmt.Errorf("could not find applications: %v", err)
	}

	mappings := make([]MemberMapping, 0)
	for _, app := range apps {
		appMaps, _ := applicationAuthorizationsByID(iq, app.ID)
		for _, m := range appMaps {
			if m.RoleID == roleID {
				mappings = append(mappings, m)
			}
		}
	}

	return mappings, nil
}

// ApplicationAuthorizations returns the member mappings of an application
func ApplicationAuthorizations(iq IQ, name string) ([]MemberMapping, error) {
	app, err := GetApplicationByPublicID(iq, name)
	if err != nil {
		return nil, fmt.Errorf("could not find application with name %s: %v", name, err)
	}

	return applicationAuthorizationsByID(iq, app.ID)
}

// ApplicationAuthorizationsByRole returns the member mappings of all applications which match the given role
func ApplicationAuthorizationsByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByName(iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	return applicationAuthorizationsByRoleID(iq, role.ID)
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
	if hasRev70API(iq) {
		switch memberType {
		case MemberTypeUser:
			endpoint = fmt.Sprintf(restRoleMembersAppUser, app.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersAppGroup, app.ID, role.ID, member)
		}
	} else {
		endpoint = fmt.Sprintf(restRoleMembersAppDeprecated, app.ID)
		current, err := ApplicationAuthorizations(iq, name)
		if err != nil && current == nil {
			current = make([]MemberMapping, 0)
		}
		current = append(current, newMapping(role.ID, memberType, member))

		buf, err := json.Marshal(memberMappings{MemberMappings: current})
		if err != nil {
			return fmt.Errorf("could not create mapping: %v", err)
		}
		payload = bytes.NewBuffer(buf)
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

func revokeLT70(iq IQ, authType, authName, roleName, memberType, memberName string) error {
	var err error
	role, err := RoleByName(iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var (
		authID, baseEndpoint string
		mapping              []MemberMapping
	)
	switch authType {
	case "organization":
		org, err := GetOrganizationByName(iq, authName)
		if err == nil {
			authID = org.ID
			baseEndpoint = restRoleMembersOrgDeprecated
			mapping, err = OrganizationAuthorizations(iq, authName)
		}
	case "application":
		app, err := GetApplicationByPublicID(iq, authName)
		if err == nil {
			authID = app.ID
			baseEndpoint = restRoleMembersAppDeprecated
			mapping, err = ApplicationAuthorizations(iq, authName)
		}
	}
	if err != nil && mapping != nil {
		return fmt.Errorf("could not get current authorizations for %s: %v", authName, err)
	}

	for i, auth := range mapping {
		if auth.RoleID == role.ID {
			for j, member := range auth.Members {
				if member.Type == memberType && member.UserOrGroupName == memberName {
					copy(mapping[i].Members[j:], mapping[i].Members[j+1:])
					mapping[i].Members[len(mapping[i].Members)-1] = Member{}
					mapping[i].Members = mapping[i].Members[:len(mapping[i].Members)-1]
				}
			}
		}
	}

	buf, err := json.Marshal(memberMappings{MemberMappings: mapping})
	if err != nil {
		return fmt.Errorf("could not create mapping: %v", err)
	}

	endpoint := fmt.Sprintf(baseEndpoint, authID)
	_, _, err = iq.Put(endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("could not remove role mapping: %v", err)
	}

	return nil
}

func revoke(iq IQ, authType, authName, roleName, memberType, memberName string) error {
	role, err := RoleByName(iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var (
		authID, baseEndpoint string
	)
	switch authType {
	case "organization":
		org, err := GetOrganizationByName(iq, authName)
		if err == nil {
			authID = org.ID
			switch memberType {
			case MemberTypeUser:
				baseEndpoint = restRoleMembersOrgUser
			case MemberTypeGroup:
				baseEndpoint = restRoleMembersOrgGroup
			}
		}
	case "application":
		app, err := GetApplicationByPublicID(iq, authName)
		if err == nil {
			authID = app.ID
			switch memberType {
			case MemberTypeUser:
				baseEndpoint = restRoleMembersAppUser
			case MemberTypeGroup:
				baseEndpoint = restRoleMembersAppGroup
			}
		}
	}

	endpoint := fmt.Sprintf(baseEndpoint, authID, role.ID, memberName)
	_, err = iq.Del(endpoint)
	return err
}

// RevokeOrganizationUser removes a user and role from the named organization
func RevokeOrganizationUser(iq IQ, name, roleName, user string) error {
	if !hasRev70API(iq) {
		return revokeLT70(iq, "organization", name, roleName, MemberTypeUser, user)
	}
	return revoke(iq, "organization", name, roleName, MemberTypeUser, user)
}

// RevokeOrganizationGroup removes a group and role from the named organization
func RevokeOrganizationGroup(iq IQ, name, roleName, group string) error {
	if !hasRev70API(iq) {
		return revokeLT70(iq, "organization", name, roleName, MemberTypeGroup, group)
	}
	return revoke(iq, "organization", name, roleName, MemberTypeGroup, group)
}

// RevokeApplicationUser removes a user and role from the named application
func RevokeApplicationUser(iq IQ, name, roleName, user string) error {
	if !hasRev70API(iq) {
		return revokeLT70(iq, "application", name, roleName, MemberTypeUser, user)
	}
	return revoke(iq, "application", name, roleName, MemberTypeUser, user)
}

// RevokeApplicationGroup removes a group and role from the named application
func RevokeApplicationGroup(iq IQ, name, roleName, group string) error {
	if !hasRev70API(iq) {
		return revokeLT70(iq, "application", name, roleName, MemberTypeGroup, group)
	}
	return revoke(iq, "application", name, roleName, MemberTypeGroup, group)
}

func repositoriesAuth(iq IQ, method, roleName, memberType, member string) error {
	if !hasRev70API(iq) {
		return fmt.Errorf("did not find revision 70 API")
	}

	role, err := RoleByName(iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	switch memberType {
	case MemberTypeUser:
		endpoint = fmt.Sprintf(restRoleMembersRepositoryUser, role.ID, member)
	case MemberTypeGroup:
		endpoint = fmt.Sprintf(restRoleMembersRepositoryGroup, role.ID, member)
	}

	switch method {
	case http.MethodPut:
		_, _, err = iq.Put(endpoint, nil)
	case http.MethodDelete:
		_, err = iq.Del(endpoint)
	}
	if err != nil {
		return fmt.Errorf("could not affect repositories role mapping: %v", err)
	}

	return nil
}

func repositoriesAuthorizationsByRoleID(iq IQ, roleID string) ([]MemberMapping, error) {
	auths, err := RepositoriesAuthorizations(iq)
	if err != nil {
		return nil, fmt.Errorf("could not find authorization mappings for repositories: %v", err)
	}

	mappings := make([]MemberMapping, 0)
	for _, m := range auths {
		if m.RoleID == roleID {
			mappings = append(mappings, m)
		}
	}

	return mappings, nil
}

// RepositoriesAuthorizations returns the member mappings of all repositories
func RepositoriesAuthorizations(iq IQ) ([]MemberMapping, error) {
	body, _, err := iq.Get(restRoleMembersReposGet)
	if err != nil {
		return nil, fmt.Errorf("could not get repositories mappings: %v", err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal mapping: %v", err)
	}

	return mappings.MemberMappings, nil
}

// RepositoriesAuthorizationsByRole returns the member mappings of all repositories which match the given role
func RepositoriesAuthorizationsByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByName(iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	return repositoriesAuthorizationsByRoleID(iq, role.ID)
}

// SetRepositoriesUser sets the role and user that can have access to the repositories
func SetRepositoriesUser(iq IQ, roleName, user string) error {
	return repositoriesAuth(iq, http.MethodPut, roleName, MemberTypeUser, user)
}

// SetRepositoriesGroup sets the role and group that can have access to the repositories
func SetRepositoriesGroup(iq IQ, roleName, group string) error {
	return repositoriesAuth(iq, http.MethodPut, roleName, MemberTypeGroup, group)
}

// RevokeRepositoriesUser revoke the role and user that can have access to the repositories
func RevokeRepositoriesUser(iq IQ, roleName, user string) error {
	return repositoriesAuth(iq, http.MethodDelete, roleName, MemberTypeUser, user)
}

// RevokeRepositoriesGroup revoke the role and group that can have access to the repositories
func RevokeRepositoriesGroup(iq IQ, roleName, group string) error {
	return repositoriesAuth(iq, http.MethodDelete, roleName, MemberTypeGroup, group)
}

func membersByRoleID(iq IQ, roleID string) ([]MemberMapping, error) {
	members := make([]MemberMapping, 0)

	if m, err := organizationAuthorizationsByRoleID(iq, roleID); err == nil && len(m) > 0 {
		members = append(members, m...)
	}

	if m, err := applicationAuthorizationsByRoleID(iq, roleID); err == nil && len(m) > 0 {
		members = append(members, m...)
	}

	if hasRev70API(iq) {
		if m, err := repositoriesAuthorizationsByRoleID(iq, roleID); err == nil && len(m) > 0 {
			members = append(members, m...)
		}
	}

	return members, nil
}

// MembersByRole returns all users and groups by role name
func MembersByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByName(iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}
	return membersByRoleID(iq, role.ID)
}

// GlobalAuthorizations returns all of the users and roles who have the administrator role across all of IQ
func GlobalAuthorizations(iq IQ) ([]MemberMapping, error) {
	body, _, err := iq.Get(restRoleMembersGlobalGet)
	if err != nil {
		return nil, fmt.Errorf("could not get global members: %v", err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal mapping: %v", err)
	}

	return mappings.MemberMappings, nil
}

func globalAuth(iq IQ, method, roleName, memberType, member string) error {
	if !hasRev70API(iq) {
		return fmt.Errorf("did not find revision 70 API")
	}

	role, err := RoleByName(iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	switch memberType {
	case MemberTypeUser:
		endpoint = fmt.Sprintf(restRoleMembersGlobalUser, role.ID, member)
	case MemberTypeGroup:
		endpoint = fmt.Sprintf(restRoleMembersGlobalGroup, role.ID, member)
	}

	switch method {
	case http.MethodPut:
		_, _, err = iq.Put(endpoint, nil)
	case http.MethodDelete:
		_, err = iq.Del(endpoint)
	}
	if err != nil {
		return fmt.Errorf("could not affect global role mapping: %v", err)
	}

	return nil
}

// SetGlobalUser sets the role and user that can have access to the repositories
func SetGlobalUser(iq IQ, roleName, user string) error {
	return globalAuth(iq, http.MethodPut, roleName, MemberTypeUser, user)
}

// SetGlobalGroup sets the role and group that can have access to the global
func SetGlobalGroup(iq IQ, roleName, group string) error {
	return globalAuth(iq, http.MethodPut, roleName, MemberTypeGroup, group)
}

// RevokeGlobalUser revoke the role and user that can have access to the global
func RevokeGlobalUser(iq IQ, roleName, user string) error {
	return globalAuth(iq, http.MethodDelete, roleName, MemberTypeUser, user)
}

// RevokeGlobalGroup revoke the role and group that can have access to the global
func RevokeGlobalGroup(iq IQ, roleName, group string) error {
	return globalAuth(iq, http.MethodDelete, roleName, MemberTypeGroup, group)
}
