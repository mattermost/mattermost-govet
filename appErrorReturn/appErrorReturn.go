// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package appErrorReturn

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const mattermostPackagePath = "github.com/mattermost/mattermost-server/v6/"

var Analyzer = &analysis.Analyzer{
	Name:     "appErrorReturn",
	Doc:      "check that when an error occurs in the handlers, we call return. Can also skip checks on api handler by adding 'skip appErrReturn check' to handler function doc",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Path() != mattermostPackagePath+"api4" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if isApiHandler(x) && !skipCheckForApiHandler(x) {
					stmts := x.Body.List
					for _, s := range stmts {
						switch stmt := s.(type) {
						case *ast.IfStmt:
							if errIfCheck(stmt) {
								if !hasReturnStmtAtEnd(stmt.Body) {
									pass.Report(analysis.Diagnostic{
										Pos:     stmt.Pos(),
										Message: "Did not find a return statement in if block",
									})
								}
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

func errIfCheck(st *ast.IfStmt) bool {
	cond, ok := st.Cond.(*ast.BinaryExpr)
	if ok {
		return checkCond(cond)
	}
	return false
}

func checkCond(biE *ast.BinaryExpr) bool {
	x, okXIdent := biE.X.(*ast.Ident)
	y, okYIdent := biE.Y.(*ast.Ident)

	s, okXParen := biE.X.(*ast.ParenExpr)
	t, okYParen := biE.Y.(*ast.ParenExpr)

	_, ghOk := biE.X.(*ast.SelectorExpr)
	if ghOk {
		return false
	}

	_, bhOk := biE.X.(*ast.CallExpr)
	if bhOk {
		return false
	}

	_, shOk := biE.X.(*ast.StarExpr)
	if shOk {
		return false
	}

	_, uhOk := biE.X.(*ast.UnaryExpr)
	if uhOk {
		return false
	}

	_, lithOk := biE.X.(*ast.BasicLit)
	if lithOk {
		return false
	}

	_, inDhOk := biE.X.(*ast.IndexExpr)
	if inDhOk {
		return false
	}

	// Y

	_, ghOY := biE.Y.(*ast.SelectorExpr)
	if ghOY {
		return false
	}

	_, bhOY := biE.Y.(*ast.CallExpr)
	if bhOY {
		return false
	}

	_, shOY := biE.Y.(*ast.StarExpr)
	if shOY {
		return false
	}

	_, uhOY := biE.Y.(*ast.UnaryExpr)
	if uhOY {
		return false
	}

	_, lithOY := biE.Y.(*ast.BasicLit)
	if lithOY {
		return false
	}

	_, inDhOY := biE.Y.(*ast.IndexExpr)
	if inDhOY {
		return false
	}

	if okXIdent && okYIdent {
		aName := x.Name
		bName := y.Name
		a := aName == "err" || aName == "e" || aName == "appErr"
		b := bName == "nil"
		return a && b
	} else if okXIdent && !okYIdent {
		biExprY := biE.Y.(*ast.BinaryExpr)
		return checkCond(biExprY)
	} else if !okXIdent && okYIdent {
		biExprX := biE.X.(*ast.BinaryExpr)
		return checkCond(biExprX)
	} else if okXParen && okYParen {
		xParenBiExpr := s.X.(*ast.BinaryExpr)
		yParenBiExpr := t.X.(*ast.BinaryExpr)
		return checkCond(xParenBiExpr) || checkCond(yParenBiExpr)
	} else if okXParen && !okYParen {
		xParenBiExpr := s.X.(*ast.BinaryExpr)
		biExprY := biE.Y.(*ast.BinaryExpr)
		return checkCond(xParenBiExpr) || checkCond(biExprY)
	} else if !okXParen && okYParen {
		yParenBiExpr := t.X.(*ast.BinaryExpr)
		biExprX := biE.X.(*ast.BinaryExpr)
		return checkCond(biExprX) || checkCond(yParenBiExpr)
	}

	biExprX := biE.X.(*ast.BinaryExpr)
	biExprY := biE.Y.(*ast.BinaryExpr)
	return checkCond(biExprX) || checkCond(biExprY)
}

func hasReturnStmtAtEnd(b *ast.BlockStmt) bool {
	l := b.List
	lastItem := l[(len(l) - 1)]

	_, ok := lastItem.(*ast.ReturnStmt)
	return ok
}

func skipCheckForApiHandler(node *ast.FuncDecl) bool {
	comments := node.Doc.Text()
	return strings.Contains(comments, "skip appErrReturn check")
}

func isApiHandler(funDecl *ast.FuncDecl) bool {
	funcType := funDecl.Type
	if len(funcType.Params.List) < 3 {
		return false
	}
	arg0Type, ok := funcType.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	arg0X, ok := arg0Type.X.(*ast.Ident)
	if !ok {
		return false
	}
	if arg0X.Name != "Context" {
		return false
	}

	arg1Type, ok := funcType.Params.List[1].Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	arg1X, ok := arg1Type.X.(*ast.Ident)
	if !ok {
		return false
	}
	if arg1X.Name != "http" || arg1Type.Sel.Name != "ResponseWriter" {
		return false
	}

	arg2Type, ok := funcType.Params.List[2].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	arg2X, ok := arg2Type.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	arg2XX, ok := arg2X.X.(*ast.Ident)
	if !ok {
		return false
	}

	if arg2XX.Name != "http" || arg2X.Sel.Name != "Request" {
		return false
	}
	return true
}
