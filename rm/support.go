package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

const restSupportZip = "service/rest/v1/support/supportzip"

// SupportZipOptions encapsulates the various information you can toggle for inclusion in a support zip
type SupportZipOptions struct {
	SystemInformation bool `json:"systemInformation"`
	ThreadDump        bool `json:"threadDump"`
	Metrics           bool `json:"metrics"`
	Configuration     bool `json:"configuration"`
	Security          bool `json:"security"`
	Log               bool `json:"log"`
	TaskLog           bool `json:"taskLog"`
	AuditLog          bool `json:"auditLog"`
	Jmx               bool `json:"jmx"`
	LimitFileSizes    bool `json:"limitFileSizes"`
	LimitZipSize      bool `json:"limitZipSize"`
}

// NewSupportZipOptions creates a SupportZipOptions intance with all options enabled
func NewSupportZipOptions() (o SupportZipOptions) {
	o.SystemInformation = true
	o.ThreadDump = true
	o.Metrics = true
	o.Configuration = true
	o.Security = true
	o.Log = true
	o.TaskLog = true
	o.AuditLog = true
	o.Jmx = true
	o.LimitFileSizes = true
	o.LimitZipSize = true
	return
}

// GetSupportZip generates a support zip with the given options
func GetSupportZip(rm RM, options SupportZipOptions) ([]byte, string, error) {
	request, err := json.Marshal(options)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving support zip: %v", err)
	}

	body, resp, err := rm.Post(restSupportZip, bytes.NewBuffer(request))
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving support zip: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("error retrieving support zip: %s", resp.Status)
	}

	_, params, err := mime.ParseMediaType(resp.Header["Content-Disposition"][0])
	if err != nil {
		return nil, "", fmt.Errorf("error determining name of support zip: %v", err)
	}

	return body, params["filename"], nil
}
