package nexusiq

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httputil"
	"path"
	"strings"
	"time"
)

const (
	restReports       = "api/v2/reports/applications"
	restReportsRaw    = "api/v2/applications/%s/reports/%s/raw"
	restReportsPolicy = "api/v2/applications/%s/reports/%s/policy"
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
	EvaluationDateStr       string `json:"evaluationDate"`
	ReportDataURL           string `json:"reportDataUrl"`
	ReportHTMLURL           string `json:"reportHtmlUrl"`
	ReportPdfURL            string `json:"reportPdfUrl"`
	Stage                   string `json:"stage"`
	evaluationDate          time.Time
}

// ReportID compares two ReportInfo objects
func (r *ReportInfo) ReportID() string {
	return path.Base(r.ReportHTMLURL)
}

// EvaluationDate returns a time object of the report's EvaluationDate
func (r *ReportInfo) EvaluationDate() time.Time {
	if r.evaluationDate.IsZero() {
		t, err := time.Parse(time.RFC3339, r.EvaluationDateStr)
		if err != nil {
			r.evaluationDate = time.Now()
		}
		r.evaluationDate = t
	}
	return r.evaluationDate
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

// ReportRaw describes the raw data of an application report
type ReportRaw struct {
	Components   []rawReportComponent  `json:"components"`
	MatchSummary rawReportMatchSummary `json:"matchSummary"`
	ReportInfo   ReportInfo            `json:"reportInfo,omitempty"`
}

type policyReportViolation struct {
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
}

// PolicyReportComponent encapsulates a component which violates a policy
type PolicyReportComponent struct {
	Component
	Violations []policyReportViolation `json:"violations"`
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
	Components  []PolicyReportComponent `json:"components"`
	Counts      policyReportCounts      `json:"counts"`
	ReportTime  int64                   `json:"reportTime"`
	ReportTitle string                  `json:"reportTitle"`
	ReportInfo  ReportInfo              `json:"reportInfo,omitempty"`
}

// Report encapsulates the policy and raw report of an application
type Report struct {
	Policy ReportPolicy `json:"policyReport"`
	Raw    ReportRaw    `json:"rawReport"`
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

// GetAllReports returns all policy and raw reports
func GetAllReports(iq IQ) ([]Report, error) {
	infos, err := GetAllReportInfos(iq)
	if err != nil {
		return nil, fmt.Errorf("could not get report infos: %v", err)
	}

	reports := make([]Report, 0)

	for _, info := range infos {
		raw, _ := getRawReportByURL(iq, info.ReportDataURL)
		policy, _ := getPolicyReportByURL(iq, strings.Replace(info.ReportDataURL, "/raw", "/policy", 1))

		raw.ReportInfo = info
		policy.ReportInfo = info

		reports = append(reports,
			Report{
				Raw:    raw,
				Policy: policy,
			})
	}

	return reports, err
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
	body, resp, err := iq.Get(URL)
	if err != nil {
		log.Printf("error: could not retrieve raw report: %v\n", err)
		dump, _ := httputil.DumpRequest(resp.Request, true)
		log.Printf("error: policy raw request: %s\n", string(dump))
		return ReportRaw{}, fmt.Errorf("could not get raw report at URL %s: %v", URL, err)
	}

	var report ReportRaw
	if err = json.Unmarshal(body, &report); err != nil {
		return report, fmt.Errorf("could not unmarshal raw report: %v", err)
	}
	return report, nil
}

func getRawReportByAppReportID(iq IQ, appID, reportID string) (ReportRaw, error) {
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
			report, err := getRawReportByURL(iq, info.ReportDataURL)
			report.ReportInfo = info
			return report, err
		}
	}

	return ReportRaw{}, fmt.Errorf("could not find raw report for stage %s", stage)
}

func getPolicyReportByURL(iq IQ, URL string) (ReportPolicy, error) {
	body, _, err := iq.Get(URL)
	if err != nil {
		return ReportPolicy{}, fmt.Errorf("could not get policy report at URL %s: %v", URL, err)
	}

	var report ReportPolicy
	if err = json.Unmarshal(body, &report); err != nil {
		return report, fmt.Errorf("could not unmarshal policy report: %v", err)
	}
	return report, nil
}

// GetPolicyReportByAppID returns report information by application public ID
func GetPolicyReportByAppID(iq IQ, appID, stage string) (ReportPolicy, error) {
	infos, err := GetReportInfosByAppID(iq, appID)
	if err != nil {
		return ReportPolicy{}, fmt.Errorf("could not get report info for app '%s': %v", appID, err)
	}

	for _, info := range infos {
		if info.Stage == stage {
			report, err := getPolicyReportByURL(iq, strings.Replace(infos[0].ReportDataURL, "/raw", "/policy", 1))
			report.ReportInfo = info
			return report, err
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

	return report, nil
}

// GetReportByAppReportID returns raw and policy report information for a given report ID
func GetReportByAppReportID(iq IQ, appID, reportID string) (report Report, err error) {
	report.Policy, err = getPolicyReportByURL(iq, fmt.Sprintf(restReportsPolicy, appID, reportID))
	if err != nil {
		return report, fmt.Errorf("could not retrieve policy report: %v", err)
	}

	report.Raw, err = getRawReportByURL(iq, fmt.Sprintf(restReportsRaw, appID, reportID))
	if err != nil {
		return report, fmt.Errorf("could not retrieve raw report: %v", err)
	}

	infos, err := GetReportInfosByAppID(iq, appID)
	if err != nil {
		return report, fmt.Errorf("could not retrieve report infos: %v", err)
	}
	for _, info := range infos {
		if info.ReportID() == reportID {
			report.Policy.ReportInfo = info
			report.Raw.ReportInfo = info
		}
	}

	return report, nil
}

// GetReportInfosByOrganization returns report information by organization name
func GetReportInfosByOrganization(iq IQ, organizationName string) (infos []ReportInfo, err error) {
	apps, err := GetApplicationsByOrganization(iq, organizationName)
	if err != nil {
		return nil, fmt.Errorf("could not get applications for organization '%s': %v", organizationName, err)
	}

	infos = make([]ReportInfo, 0)
	for _, app := range apps {
		if appInfos, err := GetReportInfosByAppID(iq, app.PublicID); err == nil {
			infos = append(infos, appInfos...)
		}
	}

	return infos, nil
}

// GetReportsByOrganization returns all reports for an given organization
func GetReportsByOrganization(iq IQ, organizationName string) (reports []Report, err error) {
	apps, err := GetApplicationsByOrganization(iq, organizationName)
	if err != nil {
		return nil, fmt.Errorf("could not get applications for organization '%s': %v", organizationName, err)
	}

	stages := []Stage{StageBuild, StageStageRelease, StageRelease, StageOperate}

	reports = make([]Report, 0)
	for _, app := range apps {
		for _, s := range stages {
			if appReport, err := GetReportByAppID(iq, app.PublicID, string(s)); err == nil {
				reports = append(reports, appReport)
			}
		}
	}

	return reports, nil
}

// ReportDiff encapsulates the differences between reports
type ReportDiff struct {
	Reports []Report                `json:"reports"`
	Waived  []PolicyReportComponent `json:"waived,omitempty"`
	Fixed   []PolicyReportComponent `json:"fixed,omitempty"`
}

// ReportsDiff returns a structure describing various differences between two reports
func ReportsDiff(iq IQ, appID, report1ID, report2ID string) (ReportDiff, error) {
	var (
		report1, report2 Report
		err              error
	)

	report1, err = GetReportByAppReportID(iq, appID, report1ID)
	if err == nil {
		report2, err = GetReportByAppReportID(iq, appID, report2ID)
	}
	if err != nil {
		return ReportDiff{}, fmt.Errorf("could not retrieve raw reports: %v", err)
	}

	diff := func(iq IQ, report1, report2 Report) (ReportDiff, error) {
		var d ReportDiff
		d.Reports = make([]Report, 2)
		d.Reports[0] = report1
		d.Reports[1] = report2

		// TODO
		report2Components := make(map[string]PolicyReportComponent)
		for _, c := range report2.Policy.Components {
			report2Components[c.Hash] = c
		}

		for _, comp1 := range report1.Policy.Components {
			comp2, ok := report2Components[comp1.Hash]
			// If the component is no longer listed in report2, then it has been fixed
			if !ok {
				d.Fixed = append(d.Fixed, comp1)
				continue
			}
			for _, pol1 := range comp1.Violations {
				var found bool
				for _, pol2 := range comp2.Violations {
					if pol1.PolicyID != pol2.PolicyID {
						continue
					}
					// Marking as waived
					if pol2.Waived {
						d.Waived = append(d.Waived, comp2)
					}
				}
				if !found {
					// If the component in report1 has a policy that is not in report2, then it was fixed
					d.Fixed = append(d.Fixed, comp1)
				}
			}
		}

		return d, nil
	}

	// determine report ordering
	if report2.Raw.ReportInfo.EvaluationDate().After(report1.Raw.ReportInfo.EvaluationDate()) {
		return diff(iq, report1, report2)
	}

	return diff(iq, report2, report1)
}
