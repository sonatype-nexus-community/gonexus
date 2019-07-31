package nexusiq

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const restComponentVersions = "api/v2/components/versions"

// ComponentVersions returns all known versions of a given component
func ComponentVersions(iq IQ, comp Component) (versions []string, err error) {
	str, err := json.Marshal(comp)
	if err != nil {
		return nil, fmt.Errorf("could not process component: %v", err)
	}

	body, _, err := iq.Post(restComponentVersions, bytes.NewBuffer(str))
	if err != nil {
		return nil, fmt.Errorf("could not request component: %v", err)
	}

	if err := json.Unmarshal(body, &versions); err != nil {
		return nil, fmt.Errorf("could not process versions list: %v", err)
	}

	return
}
