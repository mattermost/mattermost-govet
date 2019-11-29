package main

import (
	"github.com/mattermost/mattermost-govet/license"
	"github.com/mattermost/mattermost-govet/structuredLogging"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		license.Analyzer,
		license.EEAnalyzer,
		structuredLogging.Analyzer,
		// appErrorWhere.Analyzer,
	)
}
