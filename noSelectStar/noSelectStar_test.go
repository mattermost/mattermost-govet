package noSelectStar_test

import (
	"testing"

	"github.com/mattermost/mattermost-govet/v2/noSelectStar"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noSelectStar.Analyzer, "a")
}
