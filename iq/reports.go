package nexusiq

import (
	"encoding/json"
	"fmt"
	"strings"
)

const restReports = "api/v2/reports/applications"

const (
	StageBuild       = "build"
	StageStageRelase = "stage-release"
	StageRelease     = "release"
)

// ReportInfo encapsulates the summary information on a given report
type ReportInfo struct {
	ApplicationID           string `json:"applicationId"`
	EmbeddableReportHTMLURL string `json:"embeddableReportHtmlUrl"`
	EvaluationDate          string `json:"evaluationDate"`
	ReportDataURL           string `json:"reportDataUrl"`
	ReportHTMLURL           string `json:"reportHtmlUrl"`
	ReportPdfURL            string `json:"reportPdfUrl"`
	Stage                   string `json:"stage"`
}

// ReportRaw descrpibes the raw data of an application report
type ReportRaw struct {
	Components []struct {
		ComponentID  ComponentIdentifier `json:"componentIdentifier,omitempty"`
		Hash         string              `json:"hash"`
		MatchState   string              `json:"matchState"`
		PackageURL   string              `json:"packageUrl"`
		Pathnames    []string            `json:"pathnames"`
		Proprietary  bool                `json:"proprietary"`
		LicensesData LicenseData         `json:"licenseData"`
		SecurityData struct {
			SecurityIssues []SecurityIssue `json:"securityIssues"`
		} `json:"securityData"`
	} `json:"components"`
	MatchSummary struct {
		KnownComponentCount int64 `json:"knownComponentCount"`
		TotalComponentCount int64 `json:"totalComponentCount"`
	} `json:"matchSummary"`
}

// ReportPolicy descrpibes the policies violated by the components in an application report
type ReportPolicy struct {
	Application Application `json:"application"`
	Components  []struct {
		ComponentID ComponentIdentifier `json:"componentIdentifier,omitempty"`
		Hash        string              `json:"hash"`
		MatchState  string              `json:"matchState"`
		PackageURL  string              `json:"packageUrl"`
		Pathnames   []string            `json:"pathnames"`
		Proprietary bool                `json:"proprietary"`
		Violations  []struct {
			Constraints []struct {
				Conditions []struct {
					ConditionReason  string `json:"conditionReason"`
					ConditionSummary string `json:"conditionSummary"`
				} `json:"conditions"`
				ConstraintID   string `json:"constraintId"`
				ConstraintName string `json:"constraintName"`
			} `json:"constraints"`
			Grandfathered        bool   `json:"grandfathered"`
			PolicyID             string `json:"policyId"`
			PolicyName           string `json:"policyName"`
			PolicyThreatCategory string `json:"policyThreatCategory"`
			PolicyThreatLevel    int64  `json:"policyThreatLevel"`
			Waived               bool   `json:"waived"`
		} `json:"violations"`
	} `json:"components"`
	Counts struct {
		ExactlyMatchedComponentCount      int64 `json:"exactlyMatchedComponentCount"`
		GrandfatheredPolicyViolationCount int64 `json:"grandfatheredPolicyViolationCount"`
		PartiallyMatchedComponentCount    int64 `json:"partiallyMatchedComponentCount"`
		TotalComponentCount               int64 `json:"totalComponentCount"`
	} `json:"counts"`
	ReportTime  int64  `json:"reportTime"`
	ReportTitle string `json:"reportTitle"`
}

// Report encapsulates the policy and raw report of an application
type Report struct {
	PolicyReport ReportPolicy
	RawReport    ReportRaw
}

// GetAllReportInfos returns all report infos
func GetAllReportInfos(iq IQ) ([]ReportInfo, error) {
	body, _, err := iq.Get(restReports)
	if err != nil {
		return nil, fmt.Errorf("could not get report info: %v", err)
	}

	infos := make([]ReportInfo, 0)
	err = json.Unmarshal(body, &infos)

	return infos, err
}

// GetReportInfoByAppID returns report information by application public ID
func GetReportInfoByAppID(iq IQ, appID, stage string) (info ReportInfo, err error) {
	app, err := GetApplicationByPublicID(iq, appID)
	if err != nil {
		return info, fmt.Errorf("could not get application: %v", err)
	}

	endpoint := fmt.Sprintf("%s/%s", restReports, app.ID)
	body, _, err := iq.Get(endpoint)
	if err != nil {
		return info, fmt.Errorf("could not get report info: %v", err)
	}

	err = json.Unmarshal(body, &info)

	return
}

// GetRawReportByAppID returns report information by application public ID
func GetRawReportByAppID(iq IQ, appID, stage string) (ReportRaw, error) {
	info, err := GetReportInfoByAppID(iq, appID, stage)
	if err != nil {
		return ReportRaw{}, fmt.Errorf("could not report info for app '%s': %v", appID, err)
	}

	fmt.Println(info.ReportDataURL)
	// body, resp, err := iq.Get(info.ReportDataURL)

	return ReportRaw{}, nil
}

// GetPolicyReportByAppID returns report information by application public ID
func GetPolicyReportByAppID(iq IQ, appID, stage string) (ReportPolicy, error) {
	info, err := GetReportInfoByAppID(iq, appID, stage)
	if err != nil {
		return ReportPolicy{}, fmt.Errorf("could not report info for app '%s': %v", appID, err)
	}

	fmt.Println(strings.Replace(info.ReportDataURL, "/raw", "/policy", 1))
	// body, resp, err := iq.Get(policyURL)

	return ReportPolicy{}, nil
}

// GetReportByAppID returns report information by application public ID
func GetReportByAppID(iq IQ, appID, stage string) (report Report, err error) {
	report.PolicyReport, err = GetPolicyReportByAppID(iq, appID, stage)
	if err != nil {
		return report, fmt.Errorf("could not retrieve policy report: %v", err)
	}

	report.RawReport, err = GetRawReportByAppID(iq, appID, stage)
	if err != nil {
		return report, fmt.Errorf("could not retrieve raw report: %v", err)
	}

	return
}
