package nexusrm

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	restStatusReadable = "service/rest/v1/status"
	restStatusWritable = "service/rest/v1/status/writable"
	restStatusCheck = "service/rest/v1/status/check"
)

type StatusStateErrorTrace struct {
	MethodName string `json:"methodName"`
	FileName string `json:"fileName"`
	LineNumber int `json:"lineNumber"`
	ClassName string `json:"className"`
	NativeMethod bool `json:"nativeMethod"`
}

type StatusStateError struct {
	StackTrace []StatusStateErrorTrace `json:"stackTrace"`
	Message string `json:"message"`
	LocalizedMessage string `json:"localizedMessage"`
	Suppressed []interface{} `json:"suppressed"`
}

type StatusStatePart struct {
	Healthy bool `json:"healthy"`
	Message string `json:"message"`
	Error *StatusStateError `json:"error"`
	Details interface{} `json:"details"`
	Time int `json:"time"`
	Duration int `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

type StatusState struct {
	AvalaibleCPUs StatusStatePart `json:"Available CPUs"`
	BlobStores StatusStatePart `json:"Blob Stores"`
	DefaultAdlinCredentials StatusStatePart `json:"Default Admin Credentials"`
	DefaultRoleAdmin StatusStatePart `json:"DefaultRoleRealm"`
	FileBlobStoresPath StatusStatePart `json:"File Blob Stores Path"`
	FileDescriptors StatusStatePart `json:"File Descriptors"`
	LifecyclePhase StatusStatePart `json:"Lifecycle Phase"`
	ReadOnlyDetector StatusStatePart `json:"Read-Only Detector"`
	Scheduler StatusStatePart `json:"Scheduler"`
	ThreadDeadlockDetector StatusStatePart `json:"Thread Deadlock Detector"`
	Transactions StatusStatePart `json:"Transactions"`
}

// StatusReadable returns true if the RM instance can serve read requests
func StatusReadable(rm RM) (_ bool) {
	_, resp, err := rm.Get(restStatusReadable)
	return err == nil && resp.StatusCode == http.StatusOK
}

// StatusWritable returns true if the RM instance can serve read requests
func StatusWritable(rm RM) (_ bool) {
	_, resp, err := rm.Get(restStatusWritable)
	return err == nil && resp.StatusCode == http.StatusOK
}

// StatusCheck return a lot of controls to determine Nexus health
func StatusCheck(rm RM) (StatusState, error) {

	var state StatusState

	body, resp, err := rm.Get(restStatusCheck)
	if err != nil || resp.StatusCode != http.StatusOK {
		return state, err
	}

	if err := json.Unmarshal(body, &state); err != nil {
		return state, err
	}

	return state, nil
}
