// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package license

import (
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "license",
	Doc:  "check the license header.",
	Run:  run,
}

const licenseLine1 = "// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved."
const licenseLine2 = "// See LICENSE.txt for license information."

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if len(file.Comments) == 0 {
			pass.Reportf(file.Pos(), "License not found")
			continue
		}
		if len(file.Comments[0].List) < 2 {
			pass.Reportf(file.Pos(), "License not found or wrong")
			continue
		}

		if file.Comments[0].List[0].Text != licenseLine1 {
			pass.Reportf(file.Comments[0].List[0].Pos(), "License wrong:\n\tseen:\n\t%s\n\t%s\n\n\texpected:\n\t%s\n\t%s", file.Comments[0].List[0].Text, file.Comments[0].List[1].Text, licenseLine1, licenseLine2)
			continue
		}

		if file.Comments[0].List[1].Text != licenseLine2 {
			pass.Reportf(file.Comments[0].List[1].Pos(), "License wrong:\n\tseen:\n\t%s\n\t%s\n\n\texpected:\n\t%s\n\t%s", file.Comments[0].List[0].Text, file.Comments[0].List[1].Text, licenseLine1, licenseLine2)
			continue
		}
	}
	return nil, nil
}
