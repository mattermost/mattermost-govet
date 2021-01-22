package configtelemetry

import (
	"math/rand"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	rand.Seed(1)
	telemetryPkgPath = "telemetry"
	modelPkgPath = "model"
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "telemetry")
}
