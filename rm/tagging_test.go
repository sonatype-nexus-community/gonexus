package nexusrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"strings"
	"testing"
)

var dummyTags = []Tag{
	{
		Name: "dummyTag1",
		/*
			Attributes: struct {
				Foo:  "bar",
				Ping: "pong",
			},
		*/
		FirstCreated: "2017-06-12T22:42:55.019+0000",
		LastUpdated:  "2017-06-12T22:42:55.019+0000",
	},
	{
		Name:         "dummyTag2",
		FirstCreated: "2017-06-12T22:42:55.019+0000",
		LastUpdated:  "2017-06-12T22:42:55.019+0000",
	},
}

func taggingTestFunc(t *testing.T, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path[1:] == restTagging:
		resp, err := json.Marshal(tagsResponse{Items: dummyTags})
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(resp))
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path[1:], restTagging):
		tagName := path.Base(r.URL.Path)

		if tagName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, tag := range dummyTags {
			if tag.Name == tagName {
				resp, err := json.Marshal(tag)
				if err != nil {
					t.Fatal(err)
				}

				fmt.Fprintln(w, string(resp))
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	case r.Method == http.MethodPost && r.URL.Path[1:] == restTagging:
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var tag Tag
		if err = json.Unmarshal(body, &tag); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tag.FirstCreated = "2017-06-12T22:42:55.019+0000"
		tag.LastUpdated = "2017-06-12T22:42:55.019+0000"

		dummyTags = append(dummyTags, tag)

		resp, err := json.Marshal(tag)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprintln(w, string(resp))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func taggingTestRM(t *testing.T) (rm RM, mock *httptest.Server) {
	return newTestRM(t, taggingTestFunc)
}

func TestTagsList(t *testing.T) {
	rm, mock := taggingTestRM(t)
	defer mock.Close()

	tags, err := TagsList(rm)
	if err != nil {
		t.Error(err)
	}

	if len(tags) == len(dummyTags) {
		t.Errorf("received %d tags instead of the expected %d\n", len(tags), len(dummyTags))
	}

	for i, tag := range tags {
		if !reflect.DeepEqual(tag, dummyTags[i]) {
			t.Fatal("Did not receive expected tag")
		}
	}
}

func TestGetTag(t *testing.T) {
	rm, mock := taggingTestRM(t)
	defer mock.Close()

	want := dummyTags[0]

	got, err := GetTag(rm, want.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("got", got)
		t.Error("want", want)
		t.Fatal("Did not receive expected tag")
	}
}

func TestAddTag(t *testing.T) {
	rm, mock := taggingTestRM(t)
	defer mock.Close()

	newName := "newTestTag"

	got, err := AddTag(rm, newName, nil)
	if err != nil {
		t.Error(err)
	}

	if got.Name != newName {
		t.Error("Did not get tag with expected name")
	}

	gotAgain, err := GetTag(rm, newName)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, gotAgain) {
		t.Error("got", got)
		t.Error("want", gotAgain)
		t.Fatal("Did not receive expected tag")
	}
}
