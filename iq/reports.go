package nexusiq

import (
	"encoding/json"
	"fmt"
	"strings"
)

const restReports = "api/v2/reports/applications"

// Provides a constants for the IQ stages
const (
	StageProxy       = "proxy"
	StageDevelop     = "develop"
	StageBuild       = "build"
	StageStageRelase = "stage-release"
	StageRelease     = "release"
	StageOperate     = "operate"
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

// Equals compares two ReportInfo objects
func (a *ReportInfo) Equals(b *ReportInfo) (_ bool) {
	if a == b {
		return true
	}

	if a.ApplicationID != b.ApplicationID {
		return
	}

	if a.EmbeddableReportHTMLURL != b.EmbeddableReportHTMLURL {
		return
	}

	if a.EvaluationDate != b.EvaluationDate {
		return
	}

	if a.ReportDataURL != b.ReportDataURL {
		return
	}

	if a.ReportHTMLURL != b.ReportHTMLURL {
		return
	}

	if a.ReportPdfURL != b.ReportPdfURL {
		return
	}

	if a.Stage != b.Stage {
		return
	}

	return true
}

// ReportRaw descrpibes the raw data of an application report
type ReportRaw struct {
	Components []struct {
		Component
		LicensesData LicenseData `json:"licenseData"`
		SecurityData struct {
			SecurityIssues []SecurityIssue `json:"securityIssues"`
		} `json:"securityData"`
	} `json:"components"`
	MatchSummary struct {
		KnownComponentCount int64 `json:"knownComponentCount"`
		TotalComponentCount int64 `json:"totalComponentCount"`
	} `json:"matchSummary"`
}

// Equals compares two ReportRaw objects
func (a *ReportRaw) Equals(b *ReportRaw) (_ bool) {
	if a == b {
		return true
	}

	if a.MatchSummary.KnownComponentCount != b.MatchSummary.KnownComponentCount {
		return
	}

	if a.MatchSummary.TotalComponentCount != b.MatchSummary.TotalComponentCount {
		return
	}

	if len(a.Components) != len(b.Components) {
		return
	}

	for i, c := range a.Components {
		// TODO: Component ??

		if !c.LicensesData.Equals(&b.Components[i].LicensesData) {
			return
		}

		if len(c.SecurityData.SecurityIssues) != len(b.Components[i].SecurityData.SecurityIssues) {
			return
		}

		for j, s := range c.SecurityData.SecurityIssues {
			if !s.Equals(&b.Components[i].SecurityData.SecurityIssues[j]) {
				return
			}
		}
	}

	return true
}

// ReportPolicy descrpibes the policies violated by the components in an application report
type ReportPolicy struct {
	Application Application `json:"application"`
	Components  []struct {
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

// Equals compares two ReportPolicy objects
func (a *ReportPolicy) Equals(b *ReportPolicy) (_ bool) {
	if a == b {
		return true
	}

	if !a.Application.Equals(&b.Application) {
		return
	}

	if a.ReportTime != b.ReportTime {
		return
	}

	if a.ReportTitle != b.ReportTitle {
		return
	}

	if a.Counts.ExactlyMatchedComponentCount != b.Counts.ExactlyMatchedComponentCount {
		return
	}

	if a.Counts.GrandfatheredPolicyViolationCount != b.Counts.GrandfatheredPolicyViolationCount {
		return
	}

	if a.Counts.PartiallyMatchedComponentCount != b.Counts.PartiallyMatchedComponentCount {
		return
	}

	if a.Counts.TotalComponentCount != b.Counts.TotalComponentCount {
		return
	}

	if len(a.Components) != len(b.Components) {
		return
	}

	for i, c := range a.Components {
		// TODO: Component.Equals??

		if len(c.Violations) != len(b.Components[i].Violations) {
			return
		}

		for j, v := range c.Violations {
			if v.Grandfathered != b.Components[i].Violations[j].Grandfathered {
				return
			}

			if v.PolicyID != b.Components[i].Violations[j].PolicyID {
				return
			}

			if v.PolicyName != b.Components[i].Violations[j].PolicyName {
				return
			}

			if v.PolicyThreatCategory != b.Components[i].Violations[j].PolicyThreatCategory {
				return
			}

			if v.PolicyThreatLevel != b.Components[i].Violations[j].PolicyThreatLevel {
				return
			}

			if v.Waived != b.Components[i].Violations[j].Waived {
				return
			}

			if len(v.Constraints) != len(b.Components[i].Violations[j].Constraints) {
				return
			}

			for k, d := range v.Constraints {
				if d.ConstraintID != b.Components[i].Violations[j].Constraints[k].ConstraintID {
					return
				}

				if d.ConstraintName != b.Components[i].Violations[j].Constraints[k].ConstraintName {
					return
				}

				if len(d.Conditions) != len(b.Components[i].Violations[j].Constraints[k].Conditions) {
					return
				}

				for l, t := range d.Conditions {
					if t.ConditionReason != b.Components[i].Violations[j].Constraints[k].Conditions[l].ConditionReason {
						return
					}

					if t.ConditionSummary != b.Components[i].Violations[j].Constraints[k].Conditions[l].ConditionSummary {
						return
					}
				}
			}

		}
	}

	return true
}

// Report encapsulates the policy and raw report of an application
type Report struct {
	Policy ReportPolicy
	Raw    ReportRaw
}

// Equals compares two Report objects
func (a *Report) Equals(b *Report) (_ bool) {
	if a == b {
		return true
	}

	if a.Policy.Equals(&b.Policy) {
		return
	}

	if a.Raw.Equals(&b.Raw) {
		return
	}

	return true
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
