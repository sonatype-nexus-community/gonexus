package nexusiq

const restDataRetentionPolicies = "api/v2/dataRetentionPolicies/organizations/%s"

// DataRetentionPolicies encapsulates an organization's retention policies
type DataRetentionPolicies struct {
	ApplicationReports ApplicationReports  `json:"applicationReports"`
	SuccessMetrics     DataRetentionPolicy `json:"successMetrics"`
}

// ApplicationReports captures the policies related to application reports
type ApplicationReports struct {
	Stages map[Stage]DataRetentionPolicy `json:"stages"`
}

// DataRetentionPolicy describes the retention policies for a pipeline stage
type DataRetentionPolicy struct {
	InheritPolicy bool   `json:"inheritPolicy"`
	EnablePurging bool   `json:"enablePurging"`
	MaxAge        string `json:"maxAge"`
}

// GetRetentionPolicies returns the current retention policies
func GetRetentionPolicies(iq IQ) (DataRetentionPolicies, error) {
	//GET
	return DataRetentionPolicies{}, nil
}

// SetRetentionPolicies updates the retention policies
func SetRetentionPolicies(iq IQ, policies DataRetentionPolicies) error {
	//PUT
	return nil
}
