package nexusiq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	policyWaivers = "/api/v2/policyWaivers/%s/%s"
)

// OwnerType type defines the organizational scope of the waiver
type OwnerType string

// Provides guide rails for the owner types
const (
	OwnerApplication         OwnerType = "application"
	OwnerOrganization        OwnerType = "organization"
	OwnerRepository          OwnerType = "repository"
	OwnerRepositoryContainer OwnerType = "repository_container"
)

// ComponentMatcher type defines the component scope of the waiver
type ComponentMatcher string

// Provides guide rails for the matcher strategies
const (
	MatchExactComponent ComponentMatcher = "EXACT_COMPONENT"
	MatchAllComponents  ComponentMatcher = "ALL_COMPONENTS"
	MatchAllVersions    ComponentMatcher = "ALL_VERSIONS"
)

type PolicyWaiverProperties struct {
	Comment              string    `json:"comment,omitempty"`
	ExpiryTime           time.Time `json:"expiryTime,omitempty"` // time.RFC3339Nano
	MatcherStrategy      string    `json:"matcherStrategy"`
	ApplyToAllComponents bool      `json:"applytoAllComponents"`
}

func CreateWaiverByViolationId(iq IQ, properties PolicyWaiverProperties, ownerType, ownerId string) (err error) {
	_, err = validateComponentMatcher(properties.MatcherStrategy)
	if err != nil {
		return fmt.Errorf("invalid matcher strategy: %v", properties.MatcherStrategy)
	}

	_, err = validateOwnerType(ownerType)
	if err != nil {
		return fmt.Errorf("invalid owner type: %v", ownerType)
	}
	request, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("could not marshal policy waiver properties: %v", err)
	}

	endpoint := fmt.Sprintf(policyWaivers, ownerType, ownerId)
	_, resp, err := iq.Post(endpoint, bytes.NewBuffer(request))
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("did not succeeed in creating waiver: %v", err)
	}
	defer resp.Body.Close()

	return
}

func validateOwnerType(s string) (OwnerType, error) {
	switch OwnerType(s) {
	case OwnerApplication, OwnerOrganization, OwnerRepository, OwnerRepositoryContainer:
		return OwnerType(s), nil
	default:
		return "", errors.New("invalid OwnerType value")
	}
}

func validateComponentMatcher(s string) (ComponentMatcher, error) {
	switch ComponentMatcher(s) {
	case MatchExactComponent, MatchAllComponents, MatchAllVersions:
		return ComponentMatcher(s), nil
	default:
		return "", errors.New("invalid ComponentMatcher value")
	}
}
