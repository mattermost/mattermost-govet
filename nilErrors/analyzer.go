// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package nilErrors

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nilErrors",
	Doc:  "",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		errors := make(map[token.Pos]bool)
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.AssignStmt:
				for _, lhsExpr := range x.Lhs {
					ident, ok := lhsExpr.(*ast.Ident)
					if !ok {
						return false
					}
					if obj, ok := pass.TypesInfo.Defs[ident]; ok && obj.Type().String() == "error" {
						// add errors to declaration list
						errors[ident.Pos()] = true
					}
				}
			case *ast.BinaryExpr:
				if x.Op != token.NEQ {
					return true
				}
				xIdent, ok := x.X.(*ast.Ident)
				if !ok {
					return true
				}
				if _, ok := errors[xIdent.Obj.Pos()]; ok {
					delete(errors, xIdent.Obj.Pos())
				}
			case *ast.CallExpr:
				fn := x.Fun
				exp, ok := fn.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if exp.Sel.Name != "Error" {
					return true
				}

				xIdent, ok := exp.X.(*ast.Ident)
				if !ok {
					return true
				}

				if _, ok := errors[xIdent.Obj.Pos()]; ok {
					pass.Reportf(x.Pos(), ".Error() call found of an error interface.")
				}

				return true
			}
			return true
		})
	}
	return nil, nil
}
