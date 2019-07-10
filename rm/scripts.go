package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const restScript = "service/rest/v1/script"
const restScriptRun = "service/rest/v1/script/%s/run"

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
func ScriptList(rm RM) (scripts []Script, err error) {
	body, _, err := rm.Get(restScript)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &scripts)

	return
}

// ScriptGet returns the named script
func ScriptGet(rm RM, name string) (script Script, err error) {
	endpoint := fmt.Sprintf("%s/%s", restScript, name)
	body, _, err := rm.Get(endpoint)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &script)

	return
}

// ScriptUpload uploads the given Script to Repository Manager
func ScriptUpload(rm RM, script Script) error {
	json, err := json.Marshal(script)
	if err != nil {
		return err
	}

	_, resp, err := rm.Post(restScript, json)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

// ScriptUpdate update the contents of the given script
func ScriptUpdate(rm RM, script Script) error {
	json, err := json.Marshal(script)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/%s", restScript, script.Name)
	_, resp, err := rm.Put(endpoint, json)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

// ScriptRun executes the named Script
func ScriptRun(rm RM, name string, arguments []byte) error {
	endpoint := fmt.Sprintf(restScriptRun, name)
	_, _, err := rm.Post(endpoint, arguments) // TODO: Better response handling
	if err != nil {
		return err
	}

	return nil
}

// ScriptRunOnce takes the given Script, uploads it, executes it, and deletes it
func ScriptRunOnce(rm RM, script Script, arguments []byte) (err error) {
	if err = ScriptUpload(rm, script); err != nil {
		return
	}
	defer ScriptDelete(rm, script.Name)

	return ScriptRun(rm, script.Name, arguments)
}

// ScriptDelete removes the name, uploaded script
func ScriptDelete(rm RM, name string) error {
	endpoint := fmt.Sprintf("%s/%s", restScript, name)
	resp, err := rm.Del(endpoint) // TODO handle output
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}
	return nil
}
