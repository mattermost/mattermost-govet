package main

import (
	"github.com/mattermost/mattermost-govet/apiAuditLogs"
	"github.com/mattermost/mattermost-govet/equalLenAsserts"
	"github.com/mattermost/mattermost-govet/immut"
	"github.com/mattermost/mattermost-govet/inconsistentReceiverName"
	"github.com/mattermost/mattermost-govet/license"
	"github.com/mattermost/mattermost-govet/openApiSync"
	"github.com/mattermost/mattermost-govet/rawSql"
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
		openApiSync.Analyzer,
		rawSql.Analyzer,
		inconsistentReceiverName.Analyzer,
		apiAuditLogs.Analyzer,
		immut.Analyzer,
	)
}
