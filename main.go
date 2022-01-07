package main

import (
	"github.com/mattermost/mattermost-govet/v2/apiAuditLogs"
	"github.com/mattermost/mattermost-govet/v2/appErrorReturn"
	"github.com/mattermost/mattermost-govet/v2/configtelemetry"
	"github.com/mattermost/mattermost-govet/v2/emptyStrCmp"
	"github.com/mattermost/mattermost-govet/v2/equalLenAsserts"
	"github.com/mattermost/mattermost-govet/v2/errorAssertions"
	"github.com/mattermost/mattermost-govet/v2/errorVars"
	"github.com/mattermost/mattermost-govet/v2/errorVarsName"
	"github.com/mattermost/mattermost-govet/v2/immut"
	"github.com/mattermost/mattermost-govet/v2/inconsistentReceiverName"
	"github.com/mattermost/mattermost-govet/v2/license"
	"github.com/mattermost/mattermost-govet/v2/mutexLock"
	"github.com/mattermost/mattermost-govet/v2/openApiSync"
	"github.com/mattermost/mattermost-govet/v2/pointerToSlice"
	"github.com/mattermost/mattermost-govet/v2/rawSql"
	"github.com/mattermost/mattermost-govet/v2/structuredLogging"
	"github.com/mattermost/mattermost-govet/v2/tFatal"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		license.Analyzer,
		license.EEAnalyzer,
		structuredLogging.Analyzer,
		// appErrorWhere.Analyzer,
		appErrorReturn.Analyzer,
		tFatal.Analyzer,
		equalLenAsserts.Analyzer,
		openApiSync.Analyzer,
		rawSql.Analyzer,
		inconsistentReceiverName.Analyzer,
		apiAuditLogs.Analyzer,
		immut.Analyzer,
		emptyStrCmp.Analyzer,
		configtelemetry.Analyzer,
		errorAssertions.Analyzer,
		errorVarsName.Analyzer,
		errorVars.Analyzer,
		pointerToSlice.Analyzer,
		mutexLock.Analyzer,
	)
}
