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
	Policy ReportPolicy
	Raw    ReportRaw
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

// GetReportInfosByAppID returns report information by application public ID
func GetReportInfosByAppID(iq IQ, appID string) (infos []ReportInfo, err error) {
	app, err := GetApplicationByPublicID(iq, appID)
	if err != nil {
		return nil, fmt.Errorf("could not get application: %v", err)
	}

	endpoint := fmt.Sprintf("%s/%s", restReports, app.ID)
	body, _, err := iq.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not get report infos: %v", err)
	}

	infos = make([]ReportInfo, 0)
	if err = json.Unmarshal(body, &infos); err != nil {
		return infos, fmt.Errorf("could not get report infos: %v", err)
	}

	return
}

// GetRawReportByAppID returns report information by application public ID
func GetRawReportByAppID(iq IQ, appID, stage string) (ReportRaw, error) {
	infos, err := GetReportInfosByAppID(iq, appID)
	if err != nil {
		return ReportRaw{}, fmt.Errorf("could not get report info for app '%s': %v", appID, err)
	}

	for _, info := range infos {
		if info.Stage == stage {
			body, _, err := iq.Get(info.ReportDataURL)
			if err != nil {
				return ReportRaw{}, fmt.Errorf("could not get raw report: %v", err)
			}

			var report ReportRaw
			if err = json.Unmarshal(body, &report); err != nil {
				return report, fmt.Errorf("could not unmarshal raw report: %v", err)
			}
			return report, nil
		}
	}

	return ReportRaw{}, fmt.Errorf("could not find raw report for stage %s", stage)
}

// GetPolicyReportByAppID returns report information by application public ID
func GetPolicyReportByAppID(iq IQ, appID, stage string) (ReportPolicy, error) {
	infos, err := GetReportInfosByAppID(iq, appID)
	if err != nil {
		return ReportPolicy{}, fmt.Errorf("could not get report info for app '%s': %v", appID, err)
	}

	for _, info := range infos {
		if info.Stage == stage {
			body, _, err := iq.Get(strings.Replace(infos[0].ReportDataURL, "/raw", "/policy", 1))
			if err != nil {
				return ReportPolicy{}, fmt.Errorf("could not get policy report: %v", err)
			}

			var report ReportPolicy
			if err = json.Unmarshal(body, &report); err != nil {
				return report, fmt.Errorf("could not unmarshal policy report: %v", err)
			}
			return report, nil
		}
	}

	return ReportPolicy{}, fmt.Errorf("could not find policy report for stage %s", stage)
}

// GetReportByAppID returns report information by application public ID
func GetReportByAppID(iq IQ, appID, stage string) (report Report, err error) {
	report.Policy, err = GetPolicyReportByAppID(iq, appID, stage)
	if err != nil {
		return report, fmt.Errorf("could not retrieve policy report: %v", err)
	}

	report.Raw, err = GetRawReportByAppID(iq, appID, stage)
	if err != nil {
		return report, fmt.Errorf("could not retrieve raw report: %v", err)
	}

	return
}
