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

// Equals compares two Script objects
func (a *Script) Equals(b *Script) (_ bool) {
	if a == b {
		return true
	}

	if a.Name != b.Name {
		return
	}

	if a.Content != b.Content {
		return
	}

	if a.Type != b.Type {
		return
	}

	return true
}

type runResponse struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

// ScriptList lists all of the uploaded scripts in Repository Manager
func ScriptList(rm RM) (scripts []Script, err error) {
	body, _, err := rm.Get(restScript)
	if err != nil {
		return
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

	_, resp, err := rm.Post(restScript, bytes.NewBuffer(json))
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
	resp, err := rm.Put(endpoint, bytes.NewBuffer(json))
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

// ScriptRun executes the named Script
func ScriptRun(rm RM, name string, arguments []byte) (ret string, err error) {
	endpoint := fmt.Sprintf(restScriptRun, name)
	body, _, err := rm.Post(endpoint, bytes.NewBuffer(arguments)) // TODO: Better response handling
	if err != nil {
		return
	}

	var resp runResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return
	}
	ret = resp.Result

	return
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
		return err
	}
	return nil
}
