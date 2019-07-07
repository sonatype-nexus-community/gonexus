package nexusrm

import (
	"testing"
)

const groovyEcho = `return args`

func TestScriptList(t *testing.T) {
	rm := getTestRM(t)

	scripts, err := ScriptList(rm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", scripts)
}
