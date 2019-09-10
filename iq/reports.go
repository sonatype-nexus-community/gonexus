package nexusiq

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

const (
	restReports    = "api/v2/reports/applications"
	restReportsRaw = "api/v2/applications/%s/reports/%s/raw"
)

// Stage type describes a pipeline stage
type Stage string

// Provides a constants for the IQ stages
const (
	StageProxy                = "proxy"
	StageDevelop              = "develop"
	StageBuild                = "build"
	StageStageRelease         = "stage-release"
	StageRelease              = "release"
	StageOperate              = "operate"
	StageContinuousMonitoring = "continuous-monitoring"
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

// ReportID compares two ReportInfo objects
func (a *ReportInfo) ReportID() string {
	return path.Base(a.ReportHTMLURL)
}

type rawReportComponent struct {
	Component
	LicensesData LicenseData `json:"licenseData"`
	SecurityData struct {
		SecurityIssues []SecurityIssue `json:"securityIssues"`
	} `json:"securityData"`
}

type rawReportMatchSummary struct {
	KnownComponentCount int64 `json:"knownComponentCount"`
	TotalComponentCount int64 `json:"totalComponentCount"`
}

// ReportRaw descrpibes the raw data of an application report
type ReportRaw struct {
	Components   []rawReportComponent  `json:"components"`
	MatchSummary rawReportMatchSummary `json:"matchSummary"`
	ReportInfo   ReportInfo            `json:"reportInfo,omitempty"`
}

type policyReportComponent struct {
	Component
	Violations []struct {
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
}

type policyReportCounts struct {
	ExactlyMatchedComponentCount      int64 `json:"exactlyMatchedComponentCount"`
	GrandfatheredPolicyViolationCount int64 `json:"grandfatheredPolicyViolationCount"`
	PartiallyMatchedComponentCount    int64 `json:"partiallyMatchedComponentCount"`
	TotalComponentCount               int64 `json:"totalComponentCount"`
}

// ReportPolicy descrpibes the policies violated by the components in an application report
type ReportPolicy struct {
	Application Application             `json:"application"`
	Components  []policyReportComponent `json:"components"`
	Counts      policyReportCounts      `json:"counts"`
	ReportTime  int64                   `json:"reportTime"`
	ReportTitle string                  `json:"reportTitle"`
	ReportInfo  ReportInfo              `json:"reportInfo,omitempty"`
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
		return nil, fmt.Errorf("could not get info for application: %v", err)
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

// GetReportInfoByAppIDStage returns report information by application public ID and stage
func GetReportInfoByAppIDStage(iq IQ, appID, stage string) (ReportInfo, error) {
	if infos, err := GetReportInfosByAppID(iq, appID); err == nil {
		for _, info := range infos {
			if info.Stage == stage {
				return info, nil
			}
		}
	}

	return ReportInfo{}, fmt.Errorf("did not find report for '%s'", appID)
}

func getRawReportByURL(iq IQ, URL string) (ReportRaw, error) {
	body, _, err := iq.Get(URL)
	if err != nil {
		return ReportRaw{}, fmt.Errorf("could not get raw report: %v", err)
	}

	var report ReportRaw
	if err = json.Unmarshal(body, &report); err != nil {
		return report, fmt.Errorf("could not unmarshal raw report: %v", err)
	}
	return report, nil
}

// GetRawReportByAppReportID returns raw report information by application and application public ID
func GetRawReportByAppReportID(iq IQ, appID, reportID string) (ReportRaw, error) {
	return getRawReportByURL(iq, fmt.Sprintf(restReportsRaw, appID, reportID))
}

// GetRawReportByAppID returns report information by application public ID
func GetRawReportByAppID(iq IQ, appID, stage string) (ReportRaw, error) {
	infos, err := GetReportInfosByAppID(iq, appID)
	if err != nil {
		return ReportRaw{}, fmt.Errorf("could not get report info for app '%s': %v", appID, err)
	}

	for _, info := range infos {
		if info.Stage == stage {
			return getRawReportByURL(iq, info.ReportDataURL)
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
