// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package auditable

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

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
		commentMap := make(map[token.Pos]*ast.CommentGroup)
		for i := range file.Comments {
			commentMap[file.Comments[i].Pos()] = file.Comments[i]
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

			if fn.Sel != nil && fn.Sel.Name != "AddEventParameter" {
				return true
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

			// we check for the comment next to the whole expression, so we need to add 1 to the end position
			// if we find a comment with the auditable:ignore tag, we skip the check for this node
			if comment, ok := commentMap[expr.End()+1]; ok {
				if strings.Contains(comment.Text(), "auditable:ignore") {
					return false
				}
			}

			if auditable, err := isAuditable(typ.Type); err == nil && !auditable {
				pass.Reportf(param.Pos(), "%s is not auditable, but it is added to the audit record", typ.Type.String())
			} else if err != nil {
				pass.Reportf(param.Pos(), "error checking if %s is auditable: %v", typ.Type.String(), err)
			}

			return true
		})
	}

	return nil, nil
}

func isAuditable(typ types.Type) (bool, error) {
	var nt *types.Named
	if _, ok := typ.(*types.Basic); ok {
		return true, nil
	}

	switch v := typ; v.(type) {
	case *types.Interface:
		return false, nil
	case *types.Map:
		return isAuditable(v.(*types.Map).Elem())
	case *types.Slice:
		return isAuditable(v.(*types.Slice).Elem())
	case *types.Pointer:
		return isAuditable(v.(*types.Pointer).Elem())
	case *types.Named:
		nt = v.(*types.Named)
	default:
		return false, fmt.Errorf("unexpected type %T", v)
	}

	// check if the type implements the Auditable interface
	for i := 0; i < nt.NumMethods(); i++ {
		if nt.Method(i).Name() == "Auditable" && nt.Method(i).Type().(*types.Signature).String() == "func() map[string]interface{}" {
			return true, nil
		}
	}

	return false, nil
}
