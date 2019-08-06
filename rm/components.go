package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	restComponents           = "service/rest/v1/components"
	restListComponentsByRepo = "service/rest/v1/components?repository=%s"
)

type listComponentsResponse struct {
	Items             []RepositoryItem `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

// RepositoryItem holds the data of a component in a repository
type RepositoryItem struct {
	ID         string                `json:"id"`
	Repository string                `json:"repository"`
	Format     string                `json:"format"`
	Group      string                `json:"group"`
	Name       string                `json:"name"`
	Version    string                `json:"version"`
	Assets     []RepositoryItemAsset `json:"assets"`
	Tags       []string              `json:"tags"`
}

const hashPart = 20

// Hash is a hack which returns the most appopriate IQable hash of a repo item
func (a *RepositoryItem) Hash() string {
	var hash string

	sumByExt := func(exts []string) string {
		ext := exts[0]
		for _, ass := range a.Assets {
			if strings.HasSuffix(ass.Path, ext) {
				return ass.Checksum.Sha1
			}
		}
		return ""
	}

	switch a.Format {
	case "maven2":
		hash = sumByExt([]string{"jar"})
	case "rubygems":
		hash = sumByExt([]string{"gem"})
	case "npm":
		hash = sumByExt([]string{"tar.gz"})
	case "pipy":
		hash = sumByExt([]string{"whl", "tar.gz"})
	default:
		hash = ""
	}
	if len(hash) < hashPart {
		return hash
	}
	return hash[0:hashPart]
}

// UploadComponentWriter defines the interface which describes a component to upload
type UploadComponentWriter interface {
	write(w *multipart.Writer) error
}

func writeMultipartAsset(w *multipart.Writer, key string, asset io.Reader) error {
	fw, err := w.CreateFormFile(key, "") // The name seems to not matter
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, asset)
	return err
}

// UploadAssetMaven encapsulates data needed to upload an maven2 asset
type UploadAssetMaven struct {
	File                  io.Reader
	Classifier, Extension string
}

// UploadComponentMaven encapsulates data needed to upload an maven2 component
type UploadComponentMaven struct {
	GroupID, ArtifactID, Version, Packaging, Tag string
	GeneratePom                                  bool
	Assets                                       []UploadAssetMaven
}

func (a UploadComponentMaven) write(w *multipart.Writer) error {
	w.WriteField("maven2.groupId", a.GroupID)
	w.WriteField("maven2.artifactId", a.ArtifactID)
	w.WriteField("maven2.version", a.Version)
	w.WriteField("maven2.packaging", a.Packaging)
	w.WriteField("maven2.tag", a.Tag)
	w.WriteField("maven2.generate-pom", fmt.Sprintf("%v", a.GeneratePom))

	for i, a := range a.Assets {
		if a.File != nil {
			fieldName := fmt.Sprintf("maven2.asset%d", i+1)

			w.WriteField(fieldName+".classifier", a.Classifier)
			w.WriteField(fieldName+".extension", a.Extension)

			if err := writeMultipartAsset(w, fieldName, a.File); err != nil {
				return fmt.Errorf("could not add asset: %v", err)
			}
		}
	}

	return nil
}

// NewUploadComponentMaven creates a new UploadComponentMaven struct with some defaults
func NewUploadComponentMaven(coordinate string, assets ...io.Reader) (comp UploadComponentMaven, err error) {
	coordSlice := strings.Split(coordinate, ":")

	if len(coordSlice) < 3 {
		return comp, fmt.Errorf("invalid coordinate for target maven2 repo")
	}

	comp = UploadComponentMaven{
		GroupID:    coordSlice[0],
		ArtifactID: coordSlice[1],
		Version:    coordSlice[2],
		Assets:     make([]UploadAssetMaven, len(assets)),
	}

	var havePom bool
	for i, a := range assets {
		comp.Assets[i] = UploadAssetMaven{Extension: "jar", File: a} // FIXME: highly assumed extension
	}

	if !havePom {
		comp.GeneratePom = true
	}

	return
}

// UploadAssetRaw encapsulates data needed to upload a raw asset
type UploadAssetRaw struct {
	File     io.Reader
	Filename string
}

// UploadComponentRaw encapsulates data needed to upload a raw component
type UploadComponentRaw struct {
	Directory, Tag string
	Assets         []UploadAssetRaw
}

func (a UploadComponentRaw) write(w *multipart.Writer) error {
	w.WriteField("raw.directory", a.Directory)
	w.WriteField("raw.tag", a.Tag)

	for i, a := range a.Assets {
		if a.File != nil {
			fieldName := fmt.Sprintf("raw.asset%d", i+1)

			w.WriteField(fieldName+".filename", a.Filename)

			if err := writeMultipartAsset(w, fieldName, a.File); err != nil {
				return fmt.Errorf("could not add asset: %v", err)
			}
		}
	}

	return nil
}

// UploadAssetYum encapsulates data needed to upload a raw asset
type UploadAssetYum struct {
	File     io.Reader
	Filename string
}

// UploadComponentYum encapsulates data needed to upload a raw component
type UploadComponentYum struct {
	Directory, Tag string
	Assets         []UploadAssetYum
}

func (a UploadComponentYum) write(w *multipart.Writer) error {
	w.WriteField("yum.directory", a.Directory)
	w.WriteField("yum.tag", a.Tag)

	for i, a := range a.Assets {
		if a.File != nil {
			fieldName := fmt.Sprintf("yum.asset%d", i+1)

			w.WriteField(fieldName+".filename", a.Filename)

			if err := writeMultipartAsset(w, fieldName, a.File); err != nil {
				return fmt.Errorf("could not add asset: %v", err)
			}
		}
	}

	return nil
}

// UploadComponentNpm encapsulates data needed to upload an NPM component
type UploadComponentNpm struct {
	File io.Reader
	Tag  string
}

func (a UploadComponentNpm) write(w *multipart.Writer) error {
	w.WriteField("npm.tag", a.Tag)

	if err := writeMultipartAsset(w, "npm.asset", a.File); err != nil {
		return fmt.Errorf("could not add asset: %v", err)
	}

	return nil
}

// UploadComponentPyPi encapsulates data needed to upload an PyPi component
type UploadComponentPyPi struct {
	File io.Reader
	Tag  string
}

func (a UploadComponentPyPi) write(w *multipart.Writer) error {
	w.WriteField("pypi.tag", a.Tag)

	if err := writeMultipartAsset(w, "pypi.asset", a.File); err != nil {
		return fmt.Errorf("could not add asset: %v", err)
	}

	return nil
}

// UploadComponentNuget encapsulates data needed to upload an NuGet component
type UploadComponentNuget struct {
	File io.Reader
	Tag  string
}

func (a UploadComponentNuget) write(w *multipart.Writer) error {
	w.WriteField("nuget.tag", a.Tag)

	if err := writeMultipartAsset(w, "nuget.asset", a.File); err != nil {
		return fmt.Errorf("could not add asset: %v", err)
	}

	return nil
}

// UploadComponentRubyGems encapsulates data needed to upload an RubyGems component
type UploadComponentRubyGems struct {
	File io.Reader
	Tag  string
}

func (a UploadComponentRubyGems) write(w *multipart.Writer) error {
	w.WriteField("rubygems.tag", a.Tag)

	if err := writeMultipartAsset(w, "rubygems.asset", a.File); err != nil {
		return fmt.Errorf("could not add asset: %v", err)
	}

	return nil
}

// UploadComponentApt encapsulates data needed to upload an Apt component
type UploadComponentApt struct {
	File io.Reader
	Tag  string
}

func (a UploadComponentApt) write(w *multipart.Writer) error {
	w.WriteField("apt.tag", a.Tag)

	if err := writeMultipartAsset(w, "apt.asset", a.File); err != nil {
		return fmt.Errorf("could not add asset: %v", err)
	}

	return nil
}

// GetComponents returns a list of components in the indicated repository
func GetComponents(rm RM, repo string) ([]RepositoryItem, error) {
	continuation := ""

	getComponents := func() (listResp listComponentsResponse, err error) {
		url := fmt.Sprintf(restListComponentsByRepo, repo)

		if continuation != "" {
			url += "&continuationToken=" + continuation
		}

		body, resp, err := rm.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			return
		}

		err = json.Unmarshal(body, &listResp)

		return
	}

	items := make([]RepositoryItem, 0)
	for {
		resp, err := getComponents()
		if err != nil {
			return items, fmt.Errorf("could not get components: %v", err)
		}

		items = append(items, resp.Items...)

		if resp.ContinuationToken == "" {
			break
		}

		continuation = resp.ContinuationToken
	}

	return items, nil
}

// GetComponentByID returns a component by ID
func GetComponentByID(rm RM, id string) (RepositoryItem, error) {
	doError := func(err error) error {
		return fmt.Errorf("no component with id '%s': %v", id, err)
	}

	var item RepositoryItem

	url := fmt.Sprintf("%s/%s", restComponents, id)
	body, resp, err := rm.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return item, doError(err)
	}

	if err := json.Unmarshal(body, &item); err != nil {
		return item, doError(err)
	}

	return item, nil
}

// DeleteComponentByID deletes the indicated component
func DeleteComponentByID(rm RM, id string) error {
	url := fmt.Sprintf("%s/%s", restComponents, id)

	if resp, err := rm.Del(url); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("component not deleted '%s': %v", id, err)
	}

	return nil
}

// UploadComponent uploads a component to repository manager
func UploadComponent(rm RM, repo string, component UploadComponentWriter) error {
	if _, err := GetRepositoryByName(rm, repo); err != nil {
		return fmt.Errorf("could not find repository: %v", err)
	}

	doError := func(err error) error {
		return fmt.Errorf("component not uploaded: %v", err)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	component.write(w)

	if err := w.Close(); err != nil {
		return doError(err)
	}

	url := fmt.Sprintf(restListComponentsByRepo, repo)
	req, err := rm.NewRequest("POST", url, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err != nil {
		return doError(err)
	}

	if _, resp, err := rm.Do(req); err != nil && resp.StatusCode != http.StatusNoContent {
		return doError(err)
	}

	return nil
}
