package iqwebhooks

import (
	"net/http"
	"strings"

	"github.com/sonatype-nexus-community/gonexus/iq"
)

// WebhookEventType identifies a webhook event typu
type WebhookEventType string

// Enumeration of every Webhook event type
const (
	WebhookEventApplicationEvaluation WebhookEventType = "iq:applicationEvaluation"
	WebhookEventViolationAlert        WebhookEventType = "iq:policyAlert"
	WebhookEventPolicyManagement      WebhookEventType = "iq:policyManagement"
	WebhookEventLicenseOverride       WebhookEventType = "iq:licenseOverrideManagement"
	WebhookEventSecurityOverride      WebhookEventType = "iq:securityVulnerabilityOverrideManagement"
)

// WebhookEvent identifies a webhook event
type WebhookEvent interface{}

// WebhookApplicationEvaluation is the payload for an Application Evaluation webhook event
type WebhookApplicationEvaluation struct {
	Timestamp             string         `json:"timestamp"`
	Initiator             string         `json:"initiator"`
	ID                    string         `json:"id"`
	ApplicationEvaluation webhookAppEval `json:"applicationEvaluation"`
}

type webhookAppEval struct {
	PolicyEvaluationID     string `json:"policyEvaluationId"`
	Stage                  string `json:"stage,omitempty"`
	OwnerID                string `json:"ownerId,omitempty"`
	EvaluationDate         string `json:"evaluationDate,omitempty"`
	AffectedComponentCount int64  `json:"affectedComponentCount,omitempty"`
	CriticalComponentCount int64  `json:"criticalComponentCount,omitempty"`
	SevereComponentCount   int64  `json:"severeComponentCount,omitempty"`
	ModerateComponentCount int64  `json:"moderateComponentCount,omitempty"`
	Outcome                string `json:"outcome,omitempty"`
	ReportID               string `json:"reportId,omitempty"`
}

// WebhookViolationAlert is the payload for a Violation Alert webhook event
type WebhookViolationAlert struct {
	Initiator             string              `json:"initiator"`
	ApplicationEvaluation webhookAppEval      `json:"applicationEvaluation"`
	Application           nexusiq.Application `json:"application"`
	PolicyAlerts          []policyAlert       `json:"policyAlerts"`
}

type policyAlert struct {
	PolicyID          string          `json:"policyId"`
	PolicyName        string          `json:"policyName"`
	ThreatLevel       int64           `json:"threatLevel"`
	ComponentFacts    []componentFact `json:"componentFacts"`
	PolicyViolationID string          `json:"policyViolationId"`
}

type componentFact struct {
	Hash                string                      `json:"hash"`
	DisplayName         string                      `json:"displayName"`
	ComponentIdentifier nexusiq.ComponentIdentifier `json:"componentIdentifier"`
	PathNames           []string                    `json:"pathNames"`
	ConstraintFacts     []constraintFact            `json:"constraintFacts"`
}

type constraintFact struct {
	ConstraintName      string               `json:"constraintName"`
	SatisfiedConditions []satisfiedCondition `json:"satisfiedConditions"`
}

type satisfiedCondition struct {
	Summary string `json:"summary"`
	Reason  string `json:"reason"`
}

// WebhookPolicyManagement is the payload for a Policy Management webhook event
type WebhookPolicyManagement struct {
	Owner policyOwner `json:"owner"`
}

type policyOwner struct {
	ID                  string          `json:"id,omitempty"`
	PublicID            string          `json:"publicId,omitempty"`
	Name                string          `json:"name,omitempty"`
	ParentOwnerID       string          `json:"parentOwnerId,omitempty"`
	Type                string          `json:"type,omitempty"`
	Tags                []tag           `json:"tags,omitempty"`
	Labels              []tag           `json:"labels,omitempty"`
	LicenseThreatGroups []policyDetails `json:"licenseThreatGroups,omitempty"`
	Policies            []policyDetails `json:"policies,omitempty"`
	Access              []access        `json:"access,omitempty"`
}

type access struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Members []struct {
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"members,omitempty"`
}

type tag struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type policyDetails struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ThreatLevel int64  `json:"threatLevel,omitempty"`
}

// WebhookLicenseOverride is the payload for a License Override webhook event
type WebhookLicenseOverride struct {
	LicenseOverride licenseOverride `json:"licenseOverride"`
}

type licenseOverride struct {
	ID                  string                      `json:"id"`
	OwnerID             string                      `json:"ownerId"`
	Status              string                      `json:"status"`
	Comment             string                      `json:"comment"`
	LicenseIDS          []string                    `json:"licenseIds"`
	ComponentIdentifier nexusiq.ComponentIdentifier `json:"componentIdentifier"`
}

// WebhookSecurityOverride is the payload for a Security Vulnerability Override webhook event
type WebhookSecurityOverride struct {
	SecurityVulnerabilityOverride securityVulnerabilityOverride `json:"securityVulnerabilityOverride"`
}

type securityVulnerabilityOverride struct {
	ID          string `json:"id"`
	OwnerID     string `json:"ownerId"`
	Hash        string `json:"hash"`
	Source      string `json:"source"`
	ReferenceID string `json:"referenceId"`
	Status      string `json:"status"`
	Comment     string `json:"comment"`
}

// IsWebhookEvent determines if HTTP request is an IQ Webhook payload and identifies the type
func IsWebhookEvent(r *http.Request) (ok bool, whtype WebhookEventType) {
	for k, v := range r.Header {
		switch k {
		case "User-Agent":
			if !strings.HasPrefix(v[0], "Sonatype_CLM_Server") {
				break
			}
		case "X-Nexus-Webhook-Id":
			return true, WebhookEventType(v[0])
		}
	}
	return
}
