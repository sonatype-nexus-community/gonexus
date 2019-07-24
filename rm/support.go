package nexusrm

const restSupportZip = "service/rest/v1/support/supportzip"

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

func GetSupportZip(rm RM, options SupportZipOptions) (io.Reader, error) {
	request, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("error retrieving support zip: %v", err)
	}

	body, resp, err := iq.Post(restSupportZip, bytes.NewBuffer(request))
	if err != nil {
		return nil, fmt.Errorf("error retrieving support zip: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error retrieving support zip: %s", resp.Status)
	}

	return bytes.NewBuffer(body), nil
}
