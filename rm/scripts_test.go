package nexusrm

import (
	// "encoding/json"
	// "fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

// const groovyEcho = `return args`
var dummyScripts = []Script{
	Script{Name: "script1", Content: "", Type: "groovy"},
	Script{Name: "script2", Content: "", Type: "groovy"},
	Script{Name: "script3", Content: "", Type: "groovy"},
	Script{Name: "script4", Content: "", Type: "groovy"},
}

func scriptsTestRM(t *testing.T) (rm RM, mock *httptest.Server, err error) {
	return newTestRM(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		t.Logf("%q\n", dump)

		switch {
		/*
			case r.Method == http.MethodGet:
				repo := r.URL.Query().Get("repository")

				lastComponentIdx := len(dummyComponents[repo]) - 1
				var components listComponentsResponse
				if r.URL.Query().Get("continuationToken") == "" {
					components.Items = dummyComponents[repo][:lastComponentIdx]
					components.ContinuationToken = dummyContinuationToken
				} else {
					components.Items = dummyComponents[repo][lastComponentIdx:]
				}

				resp, err := json.Marshal(components)
				if err != nil {
					t.Fatal(err)
				}

				fmt.Fprintln(w, string(resp))
		*/
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestScriptList(t *testing.T) {
	t.Skip("Needs new framework")
	// rm := getTestRM(t)
	rm, mock, err := scriptsTestRM(t)
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	scripts, err := ScriptList(rm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", scripts)
}

func TestScriptGet(t *testing.T) {
	t.Skip("TODO")

	// ScriptGet(rm RM, name string)
}

func TestScriptUpload(t *testing.T) {
	t.Skip("TODO")
	// ScriptUpload(rm RM, script Script)
}

func TestScriptUpdate(t *testing.T) {
	t.Skip("TODO")
	// ScriptUpdate(rm RM, script Script)
}

func TestScriptRun(t *testing.T) {
	t.Skip("TODO")
	// ScriptRun(rm RM, name string, arguments []byte)
}

func TestScriptRunOnce(t *testing.T) {
	t.Skip("TODO")
	// ScriptRunOnce(rm RM, script Script, arguments []byte) (err error)
}

func TestScriptDelete(t *testing.T) {
	t.Skip("TODO")
	// ScriptDelete(rm RM, name string)
}
