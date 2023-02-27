// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
