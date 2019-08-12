package nexusiq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"testing"
)

var dummyUsers = []User{
	{
		Username:  "dummy1",
		FirstName: "Dummy",
		LastName:  "One",
		Email:     "dummy@dumdum.one",
		Password:  "onetwothree",
	},
	{
		Username:  "dummy2",
		FirstName: "Dummy",
		LastName:  "Two",
		Email:     "dummy@dumdum.two",
		Password:  "fourfivesix",
	},
}

func usersTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	findUserByUsername := func(username string) (User, int, bool) {
		for i, u := range dummyUsers {
			if u.Username == username {
				return u, i, true
			}
		}
		return User{}, -1, false
	}

	switch {
	case r.Method == http.MethodGet:
		username := path.Base(r.URL.Path)

		user, _, ok := findUserByUsername(username)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		buf, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
			return
		}

		fmt.Fprintln(w, string(buf))
	case r.Method == http.MethodPost:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user User
		if err = json.Unmarshal(body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := user.Username
		if _, _, ok := findUserByUsername(username); ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dummyUsers = append(dummyUsers, user)

		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodPut:
		username := path.Base(r.URL.Path)
		_, idx, ok := findUserByUsername(username)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var userUpdate User
		if err = json.Unmarshal(body, &userUpdate); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if userUpdate.Username != "" && username != userUpdate.Username {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if dummyUsers[idx].FirstName != userUpdate.FirstName {
			dummyUsers[idx].FirstName = userUpdate.FirstName
		}

		if dummyUsers[idx].LastName != userUpdate.LastName {
			dummyUsers[idx].LastName = userUpdate.LastName
		}

		if dummyUsers[idx].Email != userUpdate.Email {
			dummyUsers[idx].Email = userUpdate.Email
		}

		buf, err := json.Marshal(dummyApps[idx])
		if err != nil {
			t.Fatal(err)
			return
		}

		fmt.Fprintln(w, string(buf))
	case r.Method == http.MethodDelete:
		username := path.Base(r.URL.Path)
		if _, i, ok := findUserByUsername(username); ok {
			copy(dummyUsers[i:], dummyUsers[i+1:])
			dummyUsers[len(dummyUsers)-1] = User{}
			dummyUsers = dummyUsers[:len(dummyUsers)-1]

			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func usersTestIQ(t *testing.T, useDeprecated bool) (IQ, *httptest.Server) {
	return newTestIQ(t, usersTestFunc)
}

func checkExists(t *testing.T, iq IQ, want User) {
	t.Helper()
	got, err := GetUser(iq, want.Username)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Did not get expected user")
		t.Error(" got", got)
		t.Error("want", want)
	}
}

func TestGetUser(t *testing.T) {
	iq, mock := usersTestIQ(t, true)
	defer mock.Close()

	checkExists(t, iq, dummyUsers[0])
}

func setUser(t *testing.T, iq IQ, want User) {
	err := SetUser(iq, want)
	if err != nil {
		t.Error(err)
	}

	checkExists(t, iq, want)
}

func testCreateUser(t *testing.T, iq IQ) {
	want := User{
		Username:  "newUser",
		FirstName: "New",
		LastName:  "Dummy",
		Email:     "dummy@dumdum.new",
		Password:  "spankin",
	}

	setUser(t, iq, want)
}

func testUpdateUser(t *testing.T, iq IQ) {
	want := User{
		Username:  "anotherNewUser",
		FirstName: "Newer",
		LastName:  "Dummy",
		Email:     "dummy@dumdum.new",
		Password:  "spankin",
	}

	// Create new dummy user
	setUser(t, iq, want)

	// Update the dummy
	want.FirstName = "updatedNewer"
	want.LastName = "updatedDummy"
	want.Email = "updatedEmail"

	setUser(t, iq, want)
}

func TestSetUser(t *testing.T) {
	t.Run("creating new user", func(t *testing.T) {
		iq, mock := usersTestIQ(t, true)
		defer mock.Close()
		testCreateUser(t, iq)
	})

	t.Run("updating existing user", func(t *testing.T) {
		iq, mock := usersTestIQ(t, true)
		defer mock.Close()
		testUpdateUser(t, iq)
	})
}

func TestDeleteUser(t *testing.T) {
	iq, mock := usersTestIQ(t, true)
	defer mock.Close()

	want := User{
		Username:  "deleteMe",
		FirstName: "Delete",
		LastName:  "Me",
		Email:     "dummy@delete.me",
		Password:  "delete",
	}

	// Create new dummy user
	setUser(t, iq, want)

	err := DeleteUser(iq, want.Username)
	if err != nil {
		t.Error(err)
	}

	if _, err := GetUser(iq, want.Username); err == nil {
		t.Error("Found user which I tried to delete")
	}
}
