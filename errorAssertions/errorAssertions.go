// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package errorAssertions

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errorAssertoins",
	Doc:  "check for (require/assert).Nil/NotNil(t, error)",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				callExpr, _ := node.(*ast.CallExpr)
				fun, ok := x.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				module, ok := fun.X.(*ast.Ident)
				if !ok {
					return true
				}

				if fun.Sel.Name == "Nil" && (module.Name == "require" || module.Name == "assert") && len(callExpr.Args) > 1 {
					if typeAndValue, ok := pass.TypesInfo.Types[callExpr.Args[1]]; ok && typeAndValue.Type.String() == "error" {
						pass.Reportf(callExpr.Pos(), "calling require/assert.Nil on a regular error, please use require/assert.NoError instead")
					}
				}

				if fun.Sel.Name == "NotNil" && (module.Name == "require" || module.Name == "assert") && len(callExpr.Args) > 1 {
					if typeAndValue, ok := pass.TypesInfo.Types[callExpr.Args[1]]; ok && typeAndValue.Type.String() == "error" {
						pass.Reportf(callExpr.Pos(), "calling require/assert.NotNil on a regular error, please use require/assert.Error instead")
					}
				}

				if fun.Sel.Name == "Error" && (module.Name == "require" || module.Name == "assert") && len(callExpr.Args) > 1 {
					if typeAndValue, ok := pass.TypesInfo.Types[callExpr.Args[1]]; ok && typeAndValue.Type.String() == "*github.com/mattermost/mattermost-server/v5/model.AppError" {
						pass.Reportf(callExpr.Pos(), "calling require/assert.Error on a AppError, please use require/assert.NotNil instead")
					}
				}

				if fun.Sel.Name == "NoError" && (module.Name == "require" || module.Name == "assert") && len(callExpr.Args) > 1 {
					if typeAndValue, ok := pass.TypesInfo.Types[callExpr.Args[1]]; ok && typeAndValue.Type.String() == "*github.com/mattermost/mattermost-server/v5/model.AppError" {
						pass.Reportf(callExpr.Pos(), "calling require/assert.NoError on a AppError, please use require/assert.Nil instead")
					}
				}

				return true
			}
			return true
		})
	}
	return nil, nil
}
