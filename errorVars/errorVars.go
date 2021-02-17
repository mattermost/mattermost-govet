// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package errorVars

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	appErrorString = "*github.com/mattermost/mattermost-server/v5/model.AppError"
)

var Analyzer = &analysis.Analyzer{
	Name: "errorVars",
	Doc:  "check for non valid type assignations to err and appErr prefixed variables",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.AssignStmt:
				for idx, lhsExpr := range x.Lhs {
					if typeAndValue, ok := pass.TypesInfo.Types[lhsExpr]; ok && typeAndValue.Type.String() == "error" {
						if len(x.Rhs) == 1 {
							if typeAndValue, ok := pass.TypesInfo.Types[x.Rhs[0]]; ok {
								returnTypes := strings.Split(strings.Trim(typeAndValue.Type.String(), "()"), ", ")
								if len(returnTypes) > idx && returnTypes[idx] == appErrorString {
									pass.Reportf(x.Pos(), "assigning a *model.AppError to a `error` type variable, please create a new variable to store this value.")
								}
							}
						} else if len(x.Rhs) == len(x.Lhs) {
							if typeAndValue, ok := pass.TypesInfo.Types[x.Rhs[idx]]; ok && typeAndValue.Type.String() == appErrorString {
								pass.Reportf(x.Pos(), "assigning a *model.AppError to a `error` type variable, please create a new variable to store this value.")
							}
						}
					}
				}
				return true
			}
			return true
		})
	}
	return nil, nil
}
