// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package nilErrors

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

var Analyzer = &analysis.Analyzer{
	Name: "nilErrors",
	Doc:  "finds nil error dereferences",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.IfStmt:
				con, ok := x.Cond.(*ast.BinaryExpr)
				if !ok {
					return false
				}
				if con.Op != token.NEQ {
					return false
				}
				xIdent, ok := con.X.(*ast.Ident)
				if !ok {
					return false
				}

				if obj, ok := pass.TypesInfo.Defs[xIdent]; ok && obj.Type().Underlying().(*types.Interface) != errType {
					return false
				}

				identifiers := map[string]bool{
					xIdent.Name: true,
				}
				_ = errorChecked(pass, x.Body, identifiers)
				return false

			}
			return true
		})
	}
	return nil, nil
}

func errorChecked(pass *analysis.Pass, body *ast.BlockStmt, objects map[string]bool) bool {
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
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

			if strings.Contains(xIdent.Name, "log") || xIdent.Name == "t" || xIdent.Name == "http" {
				return false
			}

			if _, ok := objects[xIdent.Name]; !ok {
				pass.Reportf(x.Pos(), ".Error() call found of a possible nil error interface.")
			}

			return true
		case *ast.SwitchStmt:
			if x.Body == nil {
				return false
			}

			return false
			// for _, stmt := range x.Body.List {
			// 	cc, ok := stmt.(*ast.CaseClause)
			// 	if !ok {
			// 		continue
			// 	}
			// 	identifiers := map[string]bool{}
			// 	for _, exp := range cc.List {
			// 		ast.Inspect(exp, func(n ast.Node) bool {
			// 			ident, ok := n.(*ast.Ident)
			// 			if !ok {
			// 				return true
			// 			}

			// 			if obj, ok := pass.TypesInfo.Defs[ident]; ok && obj.Type().String() != "error" {
			// 				return false
			// 			}
			// 			identifiers[ident.Name] = true
			// 			return true
			// 		})
			// 	}
			// 	for _, stmt := range cc.Body {
			// 		ast.Inspect(stmt, func(n ast.Node) bool {
			// 			x, ok := n.(*ast.CallExpr)
			// 			if !ok {
			// 				return true
			// 			}
			// 			fn := x.Fun
			// 			exp, ok := fn.(*ast.SelectorExpr)
			// 			if !ok {
			// 				return false
			// 			}

			// 			if exp.Sel.Name != "Error" {
			// 				return false
			// 			}

			// 			xIdent, ok := exp.X.(*ast.Ident)
			// 			if !ok {
			// 				return false
			// 			}

			// 			if _, ok := objects[xIdent.Name]; !ok {
			// 				pass.Reportf(x.Pos(), ".Error() call found of a possible nil error interface.")
			// 			}

			// 			return false
			// 		})
			// 	}

			// }
		}

		return true
	})

	return false
}
