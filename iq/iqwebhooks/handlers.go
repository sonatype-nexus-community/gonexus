package iqwebhooks

import (
	"sync"
)

var mu sync.Mutex

var (
	appEvalEvents          = make(map[chan<- WebhookApplicationEvaluation]struct{})
	violationAlertEvents   = make(map[chan<- WebhookViolationAlert]struct{})
	policyManagementEvents = make(map[chan<- WebhookPolicyManagement]struct{})
	licenseOverrideEvents  = make(map[chan<- WebhookLicenseOverride]struct{})
	securityOverrideEvents = make(map[chan<- WebhookSecurityOverride]struct{})
)

// ApplicationEvaluationEvents returns a channel (and closer) where new Application Evaluation events are sent
func ApplicationEvaluationEvents() (<-chan WebhookApplicationEvaluation, func()) {
	events := make(chan WebhookApplicationEvaluation, 1)

	mu.Lock()
	appEvalEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(appEvalEvents, events)
	}
}

// ViolationAlertEvents returns a channel (and closer) where new Violation Alert events are sent
func ViolationAlertEvents() (<-chan WebhookViolationAlert, func()) {
	events := make(chan WebhookViolationAlert, 1)

	mu.Lock()
	violationAlertEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(violationAlertEvents, events)
	}
}

// PolicyManagementEvents returns a channel (and closer) where new Violation Alert events are sent
func PolicyManagementEvents() (<-chan WebhookPolicyManagement, func()) {
	events := make(chan WebhookPolicyManagement, 1)

	mu.Lock()
	policyManagementEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(policyManagementEvents, events)
	}
}

// LicenseOverrideEvents returns a channel (and closer) where new Violation Alert events are sent
func LicenseOverrideEvents() (<-chan WebhookLicenseOverride, func()) {
	events := make(chan WebhookLicenseOverride, 1)

	mu.Lock()
	licenseOverrideEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(licenseOverrideEvents, events)
	}
}

// SecurityOverrideEvents returns a channel (and closer) where new Violation Alert events are sent
func SecurityOverrideEvents() (<-chan WebhookSecurityOverride, func()) {
	events := make(chan WebhookSecurityOverride, 1)

	mu.Lock()
	securityOverrideEvents[events] = struct{}{}
	mu.Unlock()

	return events, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(securityOverrideEvents, events)
	}
}

func sendApplicationEvaluationEvent(event WebhookApplicationEvaluation) {
	mu.Lock()
	defer mu.Unlock()
	for c := range appEvalEvents {
		select {
		case c <- event:
		default:
		}
	}
}

func sendViolationAlertEvent(event WebhookViolationAlert) {
	mu.Lock()
	defer mu.Unlock()
	for c := range violationAlertEvents {
		select {
		case c <- event:
		default:
		}
	}
}

func sendPolicyManagementEvent(event WebhookPolicyManagement) {
	mu.Lock()
	defer mu.Unlock()
	for c := range policyManagementEvents {
		select {
		case c <- event:
		default:
		}
	}
}

func sendLicenseOverrideEvent(event WebhookLicenseOverride) {
	mu.Lock()
	defer mu.Unlock()
	for c := range licenseOverrideEvents {
		select {
		case c <- event:
		default:
		}
	}
}

func sendSecurityOverrideEvent(event WebhookSecurityOverride) {
	mu.Lock()
	defer mu.Unlock()
	for c := range securityOverrideEvents {
		select {
		case c <- event:
		default:
		}
	}
}
