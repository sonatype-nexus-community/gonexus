package nexusrm

import (
	"bytes"
	"text/template"

	"github.com/hokiegeek/gonexus"
)

/*
type blobStoreDelete struct {
	Name string
}

const groovyDeleteBlobStore = `blobStore.delete('{{.Name}}')`
*/

type BlobStoreFile struct {
	Name, Path string
}

const groovyCreateFileBlobStore = `blobStore.createFileBlobStore('{{.Name}}', '{{.Path}}')`

type BlobStoreS3 struct {
	Name, BucketName, AwsAccessKey, AwsSecret, AwsIamRole, AwsRegion string
}

const groovyCreateS3BlobStore = `def config = [:]
config['bucket'] = '{{.BucketName}}'
config['accessKeyId'] = '{{.AwsAccessKey}}'
config['secretAccessKey'] = '{{.AwsSecret}}'
config['assumeRole'] = '{{.AwsIamRole}}'
config['region'] = '{{.AwsRegion}}'
blobStore.createS3BlobStore('{{.Name}}', config)`

type blobStoreGroup struct {
	Name       string
	BlobStores []string
}

const groovyCreateBlobStoreGroup = `blobStore.createBlobStoreGroup('{{.Name}}', [{{range .BlobStores}}'{{.}}',{{end}}], 'writeToFirst')`

/*
// DeleteBlobStore creates a blobstore
func DeleteBlobStore(rm nexus.Server, name string) error {
	tmpl, err := template.New("dbs").Parse(groovyDeleteBlobStore)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, blobStoreDelete{name})
	if err != nil {
		return err
	}

	return ScriptRunOnce(rm, newAnonGroovyScript(buf.String()), nil)
}
*/

// CreateFileBlobStore creates a blobstore
func CreateFileBlobStore(rm nexus.Server, name, path string) error {
	tmpl, err := template.New("fbs").Parse(groovyCreateFileBlobStore)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, BlobStoreFile{name, path})
	if err != nil {
		return err
	}

	return ScriptRunOnce(rm, newAnonGroovyScript(buf.String()), nil)
}

// CreateBlobStoreGroup creates a blobstore
func CreateBlobStoreGroup(rm nexus.Server, name string, blobStores []string) error {
	tmpl, err := template.New("group").Parse(groovyCreateBlobStoreGroup)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, blobStoreGroup{name, blobStores})
	if err != nil {
		return err
	}

	return ScriptRunOnce(rm, newAnonGroovyScript(buf.String()), nil)
}
