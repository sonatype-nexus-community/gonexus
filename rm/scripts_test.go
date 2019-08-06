package nexusrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var dummyScripts = []Script{
	Script{Name: "script1", Content: "log.info('script1')", Type: "groovy"},
	Script{Name: "script2", Content: "log.info('script2')", Type: "groovy"},
	Script{Name: "script3", Content: "log.info('script3')", Type: "groovy"},
	Script{Name: "script4", Content: "log.info('script4')", Type: "groovy"},
}

func scriptsTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	getDummyScriptByName := func(scriptName string) (int, Script, bool) {
		for i, s := range dummyScripts {
			if s.Name == scriptName {
				return i, s, true
			}
		}
		return 0, Script{}, false
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path[1:] == restScript:
		resp, err := json.Marshal(dummyScripts)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(resp))
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path[1:], restScript):
		scriptName := strings.Replace(r.URL.Path[1:], restScript+"/", "", 1)
		if _, s, ok := getDummyScriptByName(scriptName); ok {
			resp, err := json.Marshal(s)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodPost && r.URL.Path[1:] == restScript:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		var uploadedScript Script
		if err = json.Unmarshal(body, &uploadedScript); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		dummyScripts = append(dummyScripts, uploadedScript)

	case r.Method == http.MethodPut:
		scriptName := strings.Replace(r.URL.Path[1:], restScript+"/", "", 1)
		if i, _, ok := getDummyScriptByName(scriptName); ok {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			var updatedScript Script
			if err = json.Unmarshal(body, &updatedScript); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			dummyScripts[i] = updatedScript
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path[1:], restScript) && strings.HasSuffix(r.URL.Path, "/run"):
		scriptName := strings.Replace(r.URL.Path[1:], restScript+"/", "", 1)
		scriptName = strings.Replace(scriptName, "/run", "", 1)

		if _, _, ok := getDummyScriptByName(scriptName); ok {
			defer r.Body.Close()
			args, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			resp, err := json.Marshal(runResponse{Name: scriptName, Result: string(args)})
			if err != nil {
				w.WriteHeader(http.StatusTeapot)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case r.Method == http.MethodDelete:
		scriptName := strings.Replace(r.URL.Path[1:], restScript+"/", "", 1)
		var found bool
		for i, s := range dummyScripts {
			if s.Name == scriptName {
				found = true
				copy(dummyScripts[i:], dummyScripts[i+1:])
				dummyScripts[len(dummyScripts)-1] = Script{}
				dummyScripts = dummyScripts[:len(dummyScripts)-1]

				w.WriteHeader(http.StatusNoContent)
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func scriptsTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, scriptsTestFunc)
}

func TestScriptList(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	scripts, err := ScriptList(rm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", scripts)

	if len(scripts) != len(dummyScripts) {
		t.Errorf("Received %d scripts but expected %d\n", len(scripts), len(dummyScripts))
	}

	for i, s := range scripts {
		if !reflect.DeepEqual(s, dummyScripts[i]) {
			t.Fatal("Did not receive the expected scripts")
		}
	}
}

func TestScriptGet(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	dummyScriptsIdx := 1

	script, err := ScriptGet(rm, dummyScripts[dummyScriptsIdx].Name)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", script)

	if !reflect.DeepEqual(script, dummyScripts[dummyScriptsIdx]) {
		t.Fatal("Did not receive the expected script")
	}
}

func TestScriptUpload(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	newScript := Script{Name: "newScript", Content: "log.info('I am new!')", Type: "groovy"}

	if err := ScriptUpload(rm, newScript); err != nil {
		t.Error(err)
	}

	script, err := ScriptGet(rm, newScript.Name)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", script)

	if !reflect.DeepEqual(script, newScript) {
		t.Fatal("Did not receive the expected script")
	}
}

func TestScriptUpdate(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	updatedScript := Script{
		Name:    dummyScripts[0].Name,
		Content: "log.info('I have been updated!')",
		Type:    dummyScripts[0].Type,
	}

	if reflect.DeepEqual(updatedScript, dummyScripts[0]) {
		t.Fatal("I am an idiot")
	}

	if err := ScriptUpdate(rm, updatedScript); err != nil {
		t.Error(err)
	}

	script, err := ScriptGet(rm, updatedScript.Name)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", script)

	if !reflect.DeepEqual(script, updatedScript) {
		t.Fatal("Did not receive the expected script")
	}
}

func TestScriptDelete(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	deleteMe := Script{Name: "deleteMe", Content: "log.info('Existence is pain!')", Type: "groovy"}

	if err := ScriptUpload(rm, deleteMe); err != nil {
		t.Error(err)
	}

	if err := ScriptDelete(rm, deleteMe.Name); err != nil {
		t.Error(err)
	}

	if _, err := ScriptGet(rm, deleteMe.Name); err == nil {
		t.Error("Found script which should have been deleted")
	}
}

func TestScriptRun(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	script := Script{Name: "scriptArgsTest", Content: "return args", Type: "groovy"}
	input := "this is a test"

	if err := ScriptUpload(rm, script); err != nil {
		t.Error(err)
	}

	ret, err := ScriptRun(rm, script.Name, []byte(input))
	if err != nil {
		t.Error(err)
	}

	if ret != input {
		t.Errorf("Did not get expected script output: %s\n", ret)
	}

	if err = ScriptDelete(rm, script.Name); err != nil {
		t.Error(err)
	}
}

func TestScriptRunOnce(t *testing.T) {
	rm, mock := scriptsTestRM(t)
	defer mock.Close()

	script := Script{Name: "scriptArgsTest", Content: "return args", Type: "groovy"}
	input := "this is a test"

	ret, err := ScriptRunOnce(rm, script, []byte(input))
	if err != nil {
		t.Error(err)
	}

	if ret != input {
		t.Errorf("Did not get expected script output: %s\n", ret)
	}

	if _, err = ScriptGet(rm, script.Name); err == nil {
		t.Error("Found script which should have been deleted")
	}
}
