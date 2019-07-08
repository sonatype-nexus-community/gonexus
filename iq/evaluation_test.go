package nexusiq

import (
	"testing"
)

func TestEvaluateComponents(t *testing.T) {
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

	report, err := EvaluateComponents(iq, []Component{dummy}, appID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", report)
}
