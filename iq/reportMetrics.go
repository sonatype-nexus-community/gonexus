package nexusiq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const restMetrics = "api/v2/reports/metrics"

// MetricTimePeriod is the time period to use when analyzing the data
type metricTimePeriod string

const (
	metricsMonthly metricTimePeriod = "MONTH"
	metricsWeekly  metricTimePeriod = "WEEK"
)

type metricRequest struct {
	// "MONTH" or "WEEK"
	TimePeriod metricTimePeriod `json:"timePeriod,omitempty"`

	// If timePeriod is MONTH - an ISO 8601 year and month without timezone.
	// If timePeriod is WEEK  - an ISO 8601 week year and week (e.g. week of 29 December 2008 is "2009-W01")
	FirstTimePeriod string `json:"firstTimePeriod,omitempty"`

	// Same rules as above. Must be equal to or after firstTimePeriod. Can be omitted,
	// in which case data for all successive time periods is provided including partial data for the current one.
	LastTimePeriod string `json:"lastTimePeriod,omitempty"`

	// If both of these are null or empty, data for all applications (that the user has access to) is returned.
	// applicationIds are Internal ids.
	ApplicationIDS  []string `json:"applicationIds,omitempty"`
	OrganizationIDS []string `json:"organizationIds,omitempty"`
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
	ApplicationID       string        `json:"applicationId,omitempty"`
	ApplicationPublicID string        `json:"applicationPublicId,omitempty"`
	ApplicationName     string        `json:"applicationName,omitempty"`
	OrganizationID      string        `json:"organizationId,omitempty"`
	OrganizationName    string        `json:"organizationName,omitempty"`
	Aggregations        []aggregation `json:"aggregations,omitempty"`
}

type aggregation struct {
	TimePeriodStart           string             `json:"timePeriodStart,omitempty"`
	EvaluationCount           int64              `json:"evaluationCount,omitempty"`
	MttrLowThreat             int64              `json:"mttrLowThreat,omitempty"`
	MttrModerateThreat        int64              `json:"mttrModerateThreat,omitempty"`
	MttrSevereThreat          int64              `json:"mttrSevereThreat,omitempty"`
	MttrCriticalThreat        int64              `json:"mttrCriticalThreat,omitempty"`
	DiscoveredCounts          violationCountsMap `json:"discoveredCounts,omitempty"`
	FixedCounts               violationCountsMap `json:"fixedCounts,omitempty"`
	WaivedCounts              violationCountsMap `json:"waivedCounts,omitempty"`
	OpenCountsAtTimePeriodEnd violationCountsMap `json:"openCountsAtTimePeriodEnd,omitempty"`
}

// MetricsRequestBuilder builds a request to retrieve metrics data from IQ
type MetricsRequestBuilder struct {
	timePeriod                      metricTimePeriod
	firstTimePeriod, lastTimePeriod time.Time
	apps, orgs                      []string
}

// Monthly sets the request to use a monthly time period
func (b *MetricsRequestBuilder) Monthly() *MetricsRequestBuilder {
	b.timePeriod = metricsMonthly
	return b
}

// Weekly sets the request to use a weekly time period
func (b *MetricsRequestBuilder) Weekly() *MetricsRequestBuilder {
	b.timePeriod = metricsWeekly
	return b
}

// StartingOn allows you to set the starting time period for the data gathering
func (b *MetricsRequestBuilder) StartingOn(v time.Time) *MetricsRequestBuilder {
	b.firstTimePeriod = v
	return b
}

// EndingOn allows you to set the ending time period for the data gathering. Optional
func (b *MetricsRequestBuilder) EndingOn(v time.Time) *MetricsRequestBuilder {
	b.lastTimePeriod = v
	return b
}

// WithApplication adds an application (by Public ID) whose data to include
func (b *MetricsRequestBuilder) WithApplication(v string) *MetricsRequestBuilder {
	if b.apps == nil {
		b.apps = make([]string, 0)
	}
	b.apps = append(b.apps, v)
	return b
}

// WithOrganization adds an organization whose data to include in†he report
func (b *MetricsRequestBuilder) WithOrganization(v string) *MetricsRequestBuilder {
	if b.orgs == nil {
		b.orgs = make([]string, 0)
	}
	b.orgs = append(b.orgs, v)
	return b
}

func (b *MetricsRequestBuilder) build(iq IQ) (req metricRequest, err error) {
	// If timePeriod is MONTH - an ISO 8601 year and month without timezone.
	// If timePeriod is WEEK  - an ISO 8601 week year and week (e.g. week of 29 December 2008 is "2009-W01")
	formatTime := func(t time.Time) string {
		switch req.TimePeriod {
		case metricsWeekly:
			y, w := t.ISOWeek()
			return fmt.Sprintf("%d-W%02d", y, w)
		case metricsMonthly:
			fallthrough
		default:
			return t.Format("2006-01")
		}
	}

	req.TimePeriod = metricsMonthly // TODO: validate
	if b.timePeriod != "" {
		req.TimePeriod = b.timePeriod
	}

	// TODO: would be nice to not have to require this
	if b.firstTimePeriod.IsZero() {
		return req, errors.New("a starting time period is required")
	}
	req.FirstTimePeriod = formatTime(b.firstTimePeriod)

	if !b.lastTimePeriod.IsZero() {
		// Must be equal to or after firstTimePeriod
		if b.lastTimePeriod.Before(b.firstTimePeriod) {
			return req, errors.New("ending time period must be equal to, or after, the starting time period")
		}
		req.LastTimePeriod = formatTime(b.lastTimePeriod)
	}

	if b.apps != nil {
		req.ApplicationIDS = make([]string, len(b.apps))
		for i, a := range b.apps {
			app, er := GetApplicationByPublicID(iq, a)
			if er != nil {
				return req, fmt.Errorf("could not find application with public id %s: %v", a, er)
			}
			req.ApplicationIDS[i] = app.ID
		}
	}

	if b.orgs != nil {
		req.OrganizationIDS = make([]string, len(b.orgs))
		for i, o := range b.orgs {
			org, er := GetOrganizationByName(iq, o)
			if er != nil {
				return req, fmt.Errorf("could not find organization with name %s: %v", o, er)
			}
			req.OrganizationIDS[i] = org.ID
		}
	}

	return
}

// NewMetricsRequestBuilder returns a new builder instance
func NewMetricsRequestBuilder() *MetricsRequestBuilder {
	return new(MetricsRequestBuilder)
}

// GenerateMetrics creates metrics from the given qualifiers
func GenerateMetrics(iq IQ, builder *MetricsRequestBuilder) ([]Metrics, error) {
	// TODO: Accept header: application/json or text/csv

	req, err := builder.build(iq)
	if err != nil {
		return nil, fmt.Errorf("could not build request: %v", err)
	}

	buf, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	body, _, err := iq.Post(restMetrics, bytes.NewBuffer(buf))
	if err != nil {
		return nil, fmt.Errorf("could not issue request to IQ: %v", err)
	}

	var metrics []Metrics
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		return nil, fmt.Errorf("could not read response from IQ: %v", err)
	}

	return metrics, nil
}
