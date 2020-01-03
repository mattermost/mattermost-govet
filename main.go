package main

import (
	"github.com/mattermost/mattermost-govet/equalLenAsserts"
	"github.com/mattermost/mattermost-govet/inconsistentStructName"
	"github.com/mattermost/mattermost-govet/license"
	"github.com/mattermost/mattermost-govet/structuredLogging"
	"github.com/mattermost/mattermost-govet/tFatal"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		license.Analyzer,
		license.EEAnalyzer,
		structuredLogging.Analyzer,
		// appErrorWhere.Analyzer,
		tFatal.Analyzer,
		equalLenAsserts.Analyzer,
		inconsistentStructName.Analyzer,
	)
}
