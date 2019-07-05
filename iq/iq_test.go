package nexusiq

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"testing"
)

func getTestIQ(t *testing.T) *IQ {
	iq, err := New("http://localhost:8070", "admin", "admin123")
	if err != nil {
		t.Fatal(err)
	}

	return iq
}

func createTempApplication(t *testing.T) (orgID string, appName string, appID string, err error) {
	rand.Seed(time.Now().UnixNano())
	name := strconv.Itoa(rand.Int())

	iq := getTestIQ(t)

	orgID, err = iq.CreateOrganization(name)
	if err != nil {
		return
	}

	appName = fmt.Sprintf("%s_app", name)

	appID, err = iq.CreateApplication(appName, orgID)
	if err != nil {
		return
	}

	return
}

func deleteTempApplication(t *testing.T, applicationName string) error {
	iq := getTestIQ(t)

	appInfo, err := iq.GetApplicationDetailsByName(applicationName)
	if err != nil {
		return err
	}

	if err := iq.DeleteApplication(appInfo.ID); err != nil {
		return err
	}

	// if err := iq.DeleteOrganization(appInfo.OrganizationID); err != nil {
	// 	return err
	// }

	return nil
}

func TestIQ_EvaluateComponents(t *testing.T) {
	iq := getTestIQ(t)
	iq.Debug = true

	var dummy Component
	dummy.Hash = "045c37a03be19f3e0db8"
	dummy.ComponentID.Format = "maven"
	dummy.ComponentID.Coordinates.ArtifactID = "jackson-databind"
	dummy.ComponentID.Coordinates.GroupID = "com.fasterxml.jackson.core"
	dummy.ComponentID.Coordinates.Version = "2.6.1"
	dummy.ComponentID.Coordinates.Extension = "jar"

	_, appName, appID, err := createTempApplication(t)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTempApplication(t, appName)

	report, err := iq.EvaluateComponents([]Component{dummy}, appID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", report)
}

func ExampleIQ_CreateOrganization() {
	iq, err := New("http://localhost:8070", "user", "password")
	if err != nil {
		panic(err)
	}

	orgID, err := iq.CreateOrganization("DatabaseTeam")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Organization ID: %s\n", orgID)
}
