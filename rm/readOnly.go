package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restReadOnly             = "service/rest/v1/read-only"
	restReadOnlyFreeze       = "service/rest/v1/read-only/freeze"
	restReadOnlyRelease      = "service/rest/v1/read-only/release"
	restReadOnlyForceRelease = "service/rest/v1/read-only/force-release"
)

// ReadOnlyState returns information about the read-only state of an RM instance
type ReadOnlyState struct {
	SystemInitiated bool   `json:"systemInitiated"`
	SummaryReason   string `json:"summaryReason"`
	Frozen          bool   `json:"frozen"`
}

func (s ReadOnlyState) String() string {
	var buf bytes.Buffer

	buf.WriteString("SystemInitiated: ")
	buf.WriteString(fmt.Sprintf("%v", s.SystemInitiated))
	buf.WriteString("\n")

	buf.WriteString("Frozen: ")
	buf.WriteString(fmt.Sprintf("%v", s.Frozen))
	buf.WriteString("\n")

	buf.WriteString("SummaryReason: ")
	buf.WriteString(s.SummaryReason)
	buf.WriteString("\n")

	return buf.String()
}

// GetReadOnlyState returns the read-only state of the RM instance
func GetReadOnlyState(rm RM) (state ReadOnlyState, err error) {
	body, resp, err := rm.Get(restReadOnly)
	if err != nil {
		return state, fmt.Errorf("could not get read-only state: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return state, fmt.Errorf("could not get read-only state: %v", resp.Status)
	}

	err = json.Unmarshal(body, &state)

	return
}

// ReadOnlyEnable enables read-only mode for the RM instance
func ReadOnlyEnable(rm RM) (state ReadOnlyState, err error) {
	body, resp, err := rm.Post(restReadOnlyFreeze, nil)
	if err != nil && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return
	}

	err = json.Unmarshal(body, &state)

	return
}

// ReadOnlyRelease disables read-only mode for the RM instance
func ReadOnlyRelease(rm RM, force bool) (state ReadOnlyState, err error) {
	endpoint := restReadOnlyRelease
	if force {
		endpoint = restReadOnlyForceRelease
	}

	body, resp, err := rm.Post(endpoint, nil)
	if err != nil && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return
	}

	err = json.Unmarshal(body, &state)

	return
}
