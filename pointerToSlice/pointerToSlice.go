// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package pointerToSlice

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "pointerToSlice",
	Doc:  "check for usage of pointer to slice in function definitions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	checkNode := func(expr ast.Expr) {
		sexpr, ok := expr.(*ast.StarExpr)
		if !ok {
			return
		}
		if _, ok := sexpr.X.(*ast.ArrayType); ok {
			pass.Reportf(sexpr.Pos(), "use of pointer to slice in function definition")
			return
		}
		if tv, ok := pass.TypesInfo.Types[sexpr.X]; ok {
			if strings.HasPrefix(tv.Type.String(), "[]") || strings.HasPrefix(tv.Type.Underlying().String(), "[]") {
				pass.Reportf(sexpr.Pos(), "use of pointer to slice in function definition")
				return
			}
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			f, ok := node.(*ast.FuncDecl)
			if !ok || f.Type == nil {
				return true
			}

			if params := f.Type.Params; params != nil {
				for _, p := range params.List {
					checkNode(p.Type)
				}
			}

			if results := f.Type.Results; results != nil {
				for _, r := range results.List {
					checkNode(r.Type)
				}
			}

			return true
		})
	}

	return nil, nil
}
