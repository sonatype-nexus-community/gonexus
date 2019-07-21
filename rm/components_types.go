package nexusrm

import (
	"fmt"
	"io"
)

// uploadComponentFormMapper defines the interface which describes a component to upload
type uploadComponentFormMapper interface {
	formData() (fields map[string]string, files map[string]io.Reader)
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

func (a UploadComponentMaven) formData() (map[string]string, map[string]io.Reader) {
	fields := make(map[string]string)
	files := make(map[string]io.Reader)

	fields["maven2.groupId"] = a.GroupID
	fields["maven2.artifactId"] = a.ArtifactID
	fields["maven2.version"] = a.Version
	fields["maven2.packaging"] = a.Packaging
	fields["maven2.tag"] = a.Tag
	fields["maven2.generate-pom"] = fmt.Sprintf("%v", a.GeneratePom)

	for i, a := range a.Assets {
		if a.File != nil {
			fieldName := fmt.Sprintf("maven2.asset%d", i+1)
			files[fieldName] = a.File
			fields[fieldName+".classifier"] = a.Classifier
			fields[fieldName+".extension"] = a.Extension
		}
	}

	return fields, files
}

// UploadComponentNpm encapsulates data needed to upload an NPM component
type UploadComponentNpm struct {
	File io.Reader
	Tag  string
}

func (a UploadComponentNpm) formData() (map[string]string, map[string]io.Reader) {
	return map[string]string{"npm.tag": a.Tag}, map[string]io.Reader{"npm.asset": a.File}
}
