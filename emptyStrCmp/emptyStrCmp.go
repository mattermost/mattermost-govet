// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package emptystrcmp

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "emptyStrCmp",
	Doc:  "check for len(s) == 0 and len(s) != 0 where s is string",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.BinaryExpr:
				b := node.(*ast.BinaryExpr)
				if callExpr, ok := b.X.(*ast.CallExpr); ok {
					if idt, ok := callExpr.Fun.(*ast.Ident); ok && idt.Name == "len" &&
						(b.Op == token.EQL || b.Op == token.NEQ) {
						arg0 := callExpr.Args[0]
						typ, ok := pass.TypesInfo.Types[arg0]
						if ok && typ.Type.String() == "string" {
							if bLit, ok := b.Y.(*ast.BasicLit); ok && bLit.Value == "0" {
								pass.Reportf(callExpr.Pos(), "calling len for string and comparing to 0, please compare it to empty string(\"\") instead")
								return false
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
