// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package equalLenAsserts

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "equalLenAsserts",
	Doc:  "check for (require/assert).Equal(t, X, len(Y))",
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

				if fun.Sel.Name != "Equal" || (module.Name != "require" && module.Name != "assert") {
					return true
				}

				for _, arg := range callExpr.Args {
					call, ok := arg.(*ast.CallExpr)
					if ok {
						callFun, ok := call.Fun.(*ast.Ident)
						if ok {
							if callFun.Name == "len" {
								pass.Reportf(callFun.Pos(), "calling len inside require/assert.Equal, please use require/assert.Len instead")
								return false
							}
						}
					}
				}
				return true

				// 				ast.Inspect(fun, func(subnode ast.Node) bool {
				// 					switch y := subnode.(type) {
				// 					case *ast.CallExpr:
				// 						subfun, ok := y.Fun.(*ast.Ident)
				// 						if !ok {
				// 							return true
				// 						}

				// 						if subfun.Name != "len" {
				// 							return true
				// 						}

				// 						pass.Reportf(subnode.Pos(), "calling len inside assert")
				// 					}
				// 					return true
				// 				})
			}
			return true
		})
	}
	return nil, nil
}
