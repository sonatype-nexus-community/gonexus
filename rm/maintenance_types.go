package nexusrm

const (
	AccessLogDB = "accesslog"
	ComponentDB = "component"
	ConfigDB    = "config"
	SecurityDB  = "security"
)

type DatabaseState struct {
	PageCorruption bool `json:"pageCorruption"`
	IndexErrors    int  `json:"indexErrors"`
}

// Equals compares two DatabaseState objects
func (a *DatabaseState) Equals(b *DatabaseState) (_ bool) {
	if a == b {
		return true
	}

	if a.PageCorruption != b.PageCorruption {
		return
	}

	if a.IndexErrors != b.IndexErrors {
		return
	}

	return true
}
