package nexusiq

import (
	"fmt"
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

// Metrics encapsulates the data generate
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

// MetricsRequestBuilder builds a request to retrieve metrics data from IQ
type MetricsRequestBuilder struct {
	timePeriod                      string
	firstTimePeriod, lastTimePeriod time.Time
	apps, orgs                      []string
}

// TimePeriod allows you to set the time period type. Defaults to MONTH
func (b *MetricsRequestBuilder) TimePeriod(v string) *MetricsRequestBuilder {
	b.timePeriod = v
	return b
}

// FirstTimePeriod allows you to set the starting time period for the data gathering
func (b *MetricsRequestBuilder) FirstTimePeriod(v time.Time) *MetricsRequestBuilder {
	b.firstTimePeriod = v
	return b
}

// LastTimePeriod allows you to set the ending time period for the data gathering. Optional
func (b *MetricsRequestBuilder) LastTimePeriod(v time.Time) *MetricsRequestBuilder {
	b.lastTimePeriod = v
	return b
}

// Application adds an application whose data to include
func (b *MetricsRequestBuilder) Application(v string) *MetricsRequestBuilder {
	if b.apps == nil {
		b.apps = make([]string, 0)
	}
	b.apps = append(b.apps, v)
	return b
}

// Organization adds an application whose data to include
func (b *MetricsRequestBuilder) Organization(v string) *MetricsRequestBuilder {
	if b.orgs == nil {
		b.orgs = make([]string, 0)
	}
	b.orgs = append(b.orgs, v)
	return b
}

func (b *MetricsRequestBuilder) build() (req metricRequest) {
	// If timePeriod is MONTH - an ISO 8601 year and month without timezone.
	// If timePeriod is WEEK  - an ISO 8601 week year and week (e.g. week of 29 December 2008 is "2009-W01")
	formatTime := func(t time.Time) string {
		switch req.TimePeriod {
		case timePeriodWeekly:
			y, w := t.ISOWeek()
			return fmt.Sprintf("%d-W%02d", y, w)
		case timePeriodMonthly:
			fallthrough
		default:
			return t.Format("2006-01")
		}
	}

	req.TimePeriod = timePeriodMonthly
	if b.timePeriod != "" {
		req.TimePeriod = metricTimePeriod(b.timePeriod)
	}
	req.FirstTimePeriod = formatTime(b.firstTimePeriod)
	if !b.lastTimePeriod.IsZero() {
		req.LastTimePeriod = formatTime(b.lastTimePeriod)
	}

	req.ApplicationIDS = b.apps
	req.OrganizationIDS = b.orgs

	return
}

// GenerateMetrics creates metrics from the given qualifiers
func GenerateMetrics(iq IQ, builder MetricsRequestBuilder) Metrics {
	return Metrics{}
}
