package auditable

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	apiPackagePath = "api4"
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "api4")
}
