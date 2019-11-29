package main

import (
	"github.com/mattermost/mmgovet/license"
	"github.com/mattermost/mmgovet/structuredLogging"
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
