package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restScript    = "service/rest/v1/script"
	restScriptRun = "service/rest/v1/script/%s/run"
)

// Script encapsulates a Repository Manager script
type Script struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type runResponse struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

// ScriptList lists all of the uploaded scripts in Repository Manager
func ScriptList(rm RM) ([]Script, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not list scripts: %v", err)
	}

	body, _, err := rm.Get(restScript)
	if err != nil {
		return nil, doError(err)
	}

	scripts := make([]Script, 0)
	if err = json.Unmarshal(body, &scripts); err != nil {
		return nil, doError(err)
	}

	return scripts, nil
}

// ScriptGet returns the named script
func ScriptGet(rm RM, name string) (Script, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not find script '%s': %v", name, err)
	}

	var script Script

	endpoint := fmt.Sprintf("%s/%s", restScript, name)
	body, _, err := rm.Get(endpoint)
	if err != nil {
		return script, doError(err)
	}

	if err = json.Unmarshal(body, &script); err != nil {
		return script, doError(err)
	}

	return script, nil
}

// ScriptUpload uploads the given Script to Repository Manager
func ScriptUpload(rm RM, script Script) error {
	doError := func(err error) error {
		return fmt.Errorf("could not upload script '%s': %v", script.Name, err)
	}

	json, err := json.Marshal(script)
	if err != nil {
		return doError(err)
	}

	_, resp, err := rm.Post(restScript, bytes.NewBuffer(json))
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return doError(err)
	}

	return nil
}

// ScriptUpdate update the contents of the given script
func ScriptUpdate(rm RM, script Script) error {
	doError := func(err error) error {
		return fmt.Errorf("could not update script '%s': %v", script.Name, err)
	}

	json, err := json.Marshal(script)
	if err != nil {
		return doError(err)
	}

	endpoint := fmt.Sprintf("%s/%s", restScript, script.Name)
	_, resp, err := rm.Put(endpoint, bytes.NewBuffer(json))
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return doError(err)
	}

	return nil
}

// ScriptRun executes the named Script
func ScriptRun(rm RM, name string, arguments []byte) (string, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not run script '%s': %v", name, err)
	}

	endpoint := fmt.Sprintf(restScriptRun, name)
	body, _, err := rm.Post(endpoint, bytes.NewBuffer(arguments)) // TODO: Better response handling
	if err != nil {
		return "", doError(err)
	}

	var resp runResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", doError(err)
	}

	return resp.Result, nil
}

// ScriptRunOnce takes the given Script, uploads it, executes it, and deletes it
func ScriptRunOnce(rm RM, script Script, arguments []byte) (string, error) {
	if err := ScriptUpload(rm, script); err != nil {
		return "", err
	}
	defer ScriptDelete(rm, script.Name)

	return ScriptRun(rm, script.Name, arguments)
}

// ScriptDelete removes the name, uploaded script
func ScriptDelete(rm RM, name string) error {
	endpoint := fmt.Sprintf("%s/%s", restScript, name)
	resp, err := rm.Del(endpoint)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete '%s': %v", name, err)
	}
	return nil
}
