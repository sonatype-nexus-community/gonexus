package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"strings"
	"testing"
)

var dummyDatabaseStates = map[string]DatabaseState{
	AccessLogDB: DatabaseState{true, 11},
	ComponentDB: DatabaseState{false, 22},
	ConfigDB:    DatabaseState{true, 33},
	SecurityDB:  DatabaseState{false, 44},
}

func maintenanceTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPut && path.Base(r.URL.Path[1:]) == "check" && strings.HasPrefix(r.URL.Path[1:], path.Dir(path.Dir(restMaintenanceDBCheck))):
		// path.Dir(path.Dir(r.URL.Path[1:])))
		dbName := path.Base(path.Dir(r.URL.Path[1:]))
		t.Log(dbName)
		if state, ok := dummyDatabaseStates[dbName]; ok {
			resp, err := json.Marshal(state)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(resp))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func maintenanceTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, maintenanceTestFunc)
}

func TestCheckDatabase(t *testing.T) {
	rm, mock := maintenanceTestRM(t)
	defer mock.Close()

	db := ComponentDB

	state, err := CheckDatabase(rm, db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", state)

	expectedState := dummyDatabaseStates[db]

	if !reflect.DeepEqual(state, expectedState) {
		t.Fatal("Did not receive expected database state")
	}
}

func TestCheckAllDatabases(t *testing.T) {
	rm, mock := maintenanceTestRM(t)
	defer mock.Close()

	states, err := CheckAllDatabases(rm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", states)

	if len(states) != len(dummyDatabaseStates) {
		t.Fatalf("Received %d states instead of %d\n", len(states), len(dummyDatabaseStates))
	}

	for k, v := range dummyDatabaseStates {
		state, ok := states[k]
		if !ok {
			t.Fatal("Received state does not exist in expected")
		}
		if !reflect.DeepEqual(state, v) {
			t.Fatal("Database state does not match expected")
		}
	}
}
