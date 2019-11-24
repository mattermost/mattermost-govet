package main

import (
	"github.com/mattermost/mmgovet/license"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		license.Analyzer,
		// structuredLogging.Analyzer,
		// appErrorWhere.Analyzer,
	)
}
