package nexusiq

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func createTempApplication(t *testing.T) (orgID string, appName string, appID string, err error) {
	rand.Seed(time.Now().UnixNano())
	name := strconv.Itoa(rand.Int())

	iq := getTestIQ(t)

	orgID, err = CreateOrganization(iq, name)
	if err != nil {
		return
	}

	appName = fmt.Sprintf("%s_app", name)

	appID, err = CreateApplication(iq, appName, orgID)
	if err != nil {
		return
	}

	return
}

func deleteTempApplication(t *testing.T, applicationName string) error {
	iq := getTestIQ(t)

	appInfo, err := GetApplicationDetailsByPublicID(iq, applicationName)
	if err != nil {
		return err
	}

	if err := DeleteApplication(iq, appInfo.ID); err != nil {
		return err
	}

	// if err := DeleteOrganization(iq, appInfo.OrganizationID); err != nil {
	// 	return err
	// }

	return nil
}

func TestEvaluateComponents(t *testing.T) {
	t.Skip("Skipping until I figure out a better test")
	iq := getTestIQ(t)

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

	report, err := EvaluateComponents(iq, []Component{dummy}, appID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", report)
}
