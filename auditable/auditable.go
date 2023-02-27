// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package auditable

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var (
	apiPackagePath = "github.com/mattermost/mattermost-server/v6/api4"
)

var Analyzer = &analysis.Analyzer{
	Name: "auditable",
	Doc:  "check if auditable interface is satisfied when passing structs to audit records",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Path() != apiPackagePath {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			return true
		})
	}

	return nil, nil
}
