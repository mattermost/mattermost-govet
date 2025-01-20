// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package license_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-govet/v2/license"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLicense(t *testing.T) {
	testCases := []struct {
		Description string
		Path        string
		Analyzer    *analysis.Analyzer
	}{
		{
			"Standard",
			"standard",
			license.Analyzer,
		},
		{
			"Enterprise",
			"enterprise",
			license.EEAnalyzer,
		},
		{
			"Standard",
			"source_available/enterprise",
			license.Analyzer,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			t.Run("ignored", func(t *testing.T) {
				t.Run("specified file", func(t *testing.T) {
					testdata := analysistest.TestData()
					require.NoError(t, testCase.Analyzer.Flags.Set("ignore", "ignored1.go,ignored2.go"))
					t.Cleanup(func() {
						require.NoError(t, testCase.Analyzer.Flags.Set("ignore", ""))
					})

					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "ignored_files"))
				})

				t.Run("mockery", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "ignored/mockery"))
				})

				t.Run("mockgen", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "ignored/mockgen"))
				})

				t.Run("go-bindata", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "ignored/go-bin-data"))
				})

				t.Run("manifest", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "ignored/manifest"))
				})
			})

			t.Run("missing license", func(t *testing.T) {
				testdata := analysistest.TestData()
				analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "missing"))
			})

			t.Run("valid license", func(t *testing.T) {
				testdata := analysistest.TestData()
				analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "valid"))
			})

			t.Run("invalid copyright on line 1", func(t *testing.T) {
				testdata := analysistest.TestData()
				analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "invalid_copyright"))
			})

			t.Run("invalid reference on line 2", func(t *testing.T) {
				testdata := analysistest.TestData()
				analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "invalid_reference"))
			})

			t.Run("build directives", func(t *testing.T) {
				t.Run("go:generate with valid license", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "build/gogenerate"))
				})

				t.Run("go:generate with invalid license", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "build/gogenerate_invalid"))
				})

				t.Run("go:build with valid license", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "build/buildtag"))
				})

				t.Run("directive with valid license but without newline", func(t *testing.T) {
					testdata := analysistest.TestData()
					analysistest.Run(t, testdata, testCase.Analyzer, filepath.Join(testCase.Path, "build/withoutnewline"))
				})
			})
		})
	}
}
