package iqwebhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func parseRequest(r *http.Request) (whtype WebhookEventType, err error) {
	ok, whtype := IsWebhookEvent(r)
	if !ok {
		return whtype, errors.New("not a valid IQ webhook")
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return whtype, errors.New("could not read payload")
	}

	switch whtype {
	case WebhookEventApplicationEvaluation:
		var event WebhookApplicationEvaluation
		log.Println("app eval", string(body))
		if err = json.Unmarshal(body, &event); err == nil {
			sendApplicationEvaluationEvent(event)
		}
	case WebhookEventViolationAlert:
		var event WebhookViolationAlert
		if err = json.Unmarshal(body, &event); err == nil {
			sendViolationAlertEvent(event)
		}
	case WebhookEventPolicyManagement:
		var event WebhookPolicyManagement
		if err = json.Unmarshal(body, &event); err == nil {
			sendPolicyManagementEvent(event)
		}
	case WebhookEventLicenseOverride:
		var event WebhookLicenseOverride
		if err = json.Unmarshal(body, &event); err == nil {
			sendLicenseOverrideEvent(event)
		}
	case WebhookEventSecurityOverride:
		var event WebhookSecurityOverride
		if err = json.Unmarshal(body, &event); err == nil {
			sendSecurityOverrideEvent(event)
		}
	default:
		return whtype, fmt.Errorf("IQ webhook type '%s' not supported", whtype)
	}

	return whtype, err
}

// Listen will handle any HTTP requests which are genuine Nexus IQ Webhooks
func Listen(w http.ResponseWriter, r *http.Request) {
	if _, err := parseRequest(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
