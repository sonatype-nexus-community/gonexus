package nexusiq

import (
	"time"
)

const restMetrics = "api/v2/reports/metrics"

type metricTimePeriod string

const (
	timePeriodMonthly metricTimePeriod = "MONTH"
	timePeriodWeekly  metricTimePeriod = "WEEK"
)

type metricRequest struct {
	// "MONTH" or "WEEK"
	TimePeriod metricTimePeriod `json:"timePeriod"`

	// If timePeriod is MONTH - an ISO 8601 year and month without timezone.
	// If timePeriod is WEEK  - an ISO 8601 week year and week (e.g. week of 29 December 2008 is "2009-W01")
	FirstTimePeriod string `json:"firstTimePeriod"`

	// Same rules as above. Must be equal to or after firstTimePeriod. Can be omitted,
	// in which case data for all successive time periods is provided including partial data for the current one.
	LastTimePeriod string `json:"lastTimePeriod"`

	// If both of these are null or empty, data for all applications (that the user has access to) is returned.
	// applicationIds are Internal ids.
	ApplicationIDS  []string `json:"applicationIds"`
	OrganizationIDS []string `json:"organizationIds"`
}

type violationCountType string

const (
	securityCount violationCountType = "SECURITY"
	licenseCount  violationCountType = "LICENSE"
	qualityCount  violationCountType = "QUALITY"
	otherCount    violationCountType = "OTHER"
)

type violationCountsMap map[violationCountType]violationCounts

type violationCounts struct {
	Low      int64 `json:"LOW"`
	Moderate int64 `json:"MODERATE"`
	Severe   int64 `json:"SEVERE"`
	Critical int64 `json:"CRITICAL"`
}

type Metrics struct {
	ApplicationID       string        `json:"applicationId"`
	ApplicationPublicID string        `json:"applicationPublicId"`
	ApplicationName     string        `json:"applicationName"`
	OrganizationID      string        `json:"organizationId"`
	OrganizationName    string        `json:"organizationName"`
	Aggregations        []aggregation `json:"aggregations"`
}

type aggregation struct {
	TimePeriodStart           string             `json:"timePeriodStart"`
	EvaluationCount           int64              `json:"evaluationCount"`
	MttrLowThreat             int64              `json:"mttrLowThreat"`
	MttrModerateThreat        int64              `json:"mttrModerateThreat"`
	MttrSevereThreat          int64              `json:"mttrSevereThreat"`
	MttrCriticalThreat        int64              `json:"mttrCriticalThreat"`
	DiscoveredCounts          violationCountsMap `json:"discoveredCounts"`
	FixedCounts               violationCountsMap `json:"fixedCounts"`
	WaivedCounts              violationCountsMap `json:"waivedCounts"`
	OpenCountsAtTimePeriodEnd violationCountsMap `json:"openCountsAtTimePeriodEnd"`
}

// TODO: Accept header: application/json or text/csv

func generateMetrics(iq IQ, req metricRequest) Metrics {
	// time.ISOWeek() // year, week int
	return Metrics{}
}

// GenerateMetrics TODO
func GenerateMetrics(iq IQ, timePeriod string, firstTimePeriod, lastTimePeriod time.Time) Metrics {
	var req metricRequest
	return generateMetrics(iq, req)
}

// GenerateApplicationMetrics TODO
func GenerateApplicationMetrics(iq IQ, timePeriod string, firstTimePeriod, lastTimePeriod time.Time, appPublicID string) Metrics {
	var req metricRequest
	return generateMetrics(iq, req)
}

// GenerateOrganizationMetrics TODO
func GenerateOrganizationMetrics(iq IQ, timePeriod string, firstTimePeriod, lastTimePeriod time.Time, orgName string) Metrics {
	var req metricRequest
	return generateMetrics(iq, req)
}
