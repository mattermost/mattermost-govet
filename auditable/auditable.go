// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package auditable

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	apiPackagePath = "github.com/mattermost/mattermost-server/v6/api4"
	methodsToCheck = map[string]tester{"AddEventParameter": {
		fn:            isAuditable,
		vetMsg:        "%s is not auditable, but it is added to the audit record",
		errorMsg:      "error checking if %s is auditable: %v",
		ignoreComment: "auditable:ignore",
		name:          "auditable",
	}}
)

type tester struct {
	fn            func(types.Type) (bool, error)
	vetMsg        string
	errorMsg      string
	ignoreComment string
	name          string
}

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
		commentMap := make(map[string]map[int]struct{})
		for tester := range methodsToCheck {
			commentMap[tester] = make(map[int]struct{})
		}

		for method, tester := range methodsToCheck {
			for i := range file.Comments {
				if !strings.Contains(file.Comments[i].Text(), tester.ignoreComment) {
					continue
				}

				// since go/ast parse comments out of the node tree, it's tricky to get a comment
				// after an expression. You'll need to guess how many columns the comment is away
				// so we take a naive approach and just get the line number of the comment so that
				// we can check if it's the next line of the expression end position.
				line := pass.Fset.PositionFor(file.Comments[i].Pos(), false).Line
				commentMap[method][line] = struct{}{}
			}
		}

		ast.Inspect(file, func(n ast.Node) bool {
			expr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if expr.Fun == nil {
				return true
			}

			fn, ok := expr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if fn.Sel == nil {
				return true
			}

			for method, tester := range methodsToCheck {
				if fn.Sel.Name != method {
					continue
				}

				if len(expr.Args) < 2 {
					return false
				}

				param, ok := expr.Args[1].(*ast.Ident)
				if !ok {
					return false
				}

				typ, ok := pass.TypesInfo.Types[param]
				if !ok {
					return false
				}

				exprLine := pass.Fset.PositionFor(expr.End(), false).Line
				// we check for the comment next to the whole expression, so we need to add 1 to the end position
				// if we find a comment with the auditable:ignore tag, we skip the check for this node.
				if _, ok := commentMap[method][exprLine]; ok {
					return false
				}

				if auditable, err := tester.fn(typ.Type); err == nil && !auditable {
					pass.Reportf(param.Pos(), tester.vetMsg, typ.Type.String())
				} else if err != nil {
					pass.Reportf(param.Pos(), tester.errorMsg, typ.Type.String(), err)
				}
			}

			return true
		})
	}

	return nil, nil
}

func isAuditable(typ types.Type) (bool, error) {
	if _, ok := typ.(*types.Basic); ok {
		return true, nil
	}

	// this is the common interface that we need to check if the type implements
	// a method. We need to check for both *types.Interface and *types.Named
	var mt interface {
		NumMethods() int
		Method(int) *types.Func
	}

	switch v := typ; v.(type) {
	case *types.Interface:
		mt = v.(*types.Interface)
	case *types.Map:
		return isAuditable(v.(*types.Map).Elem())
	case *types.Slice:
		return isAuditable(v.(*types.Slice).Elem())
	case *types.Pointer:
		return isAuditable(v.(*types.Pointer).Elem())
	case *types.Named:
		mt = v.(*types.Named)
	default:
		return false, fmt.Errorf("unexpected type %T", v)
	}

	// check if the type implements the Auditable interface
	for i := 0; i < mt.NumMethods(); i++ {
		if mt.Method(i).Name() == "Auditable" && mt.Method(i).Type().(*types.Signature).String() == "func() map[string]interface{}" {
			return true, nil
		}
	}

	return false, nil
}
