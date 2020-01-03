// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package inconsistentStructName

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "inconsistentStructName",
	Doc:  "check for inconsistent struct name in the methods definition",
	Run:  run,
}

type firstVarName = struct {
	Name string
	Node ast.Node
}

func run(pass *analysis.Pass) (interface{}, error) {
	firstVarNameForType := make(map[string]firstVarName)
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch funDecl := node.(type) {
			case *ast.FuncDecl:
				if funDecl.Recv != nil && len(funDecl.Recv.List) > 0 && len(funDecl.Recv.List[0].Names) > 0 {
					field := funDecl.Recv.List[0]
					typeIdent, ok := field.Type.(*ast.Ident)
					if !ok {
						var starExpr *ast.StarExpr
						starExpr, ok = field.Type.(*ast.StarExpr)
						if ok {
							typeIdent, ok = starExpr.X.(*ast.Ident)
						}
					}

					if ok {
						varName := field.Names[0].Name
						if _, exists := firstVarNameForType[typeIdent.Name]; !exists {
							firstVarNameForType[typeIdent.Name] = firstVarName{
								Name: varName,
								Node: node,
							}
						} else if firstVarNameForType[typeIdent.Name].Name != varName {
							otherVarName := firstVarNameForType[typeIdent.Name].Name
							otherVarPosition := pass.Fset.Position(firstVarNameForType[typeIdent.Name].Node.Pos())
							pass.Reportf(node.Pos(), "Different variable name used for the struct %s in different methods, using %s here but was named as %s before at %s", typeIdent.Name, varName, otherVarName, otherVarPosition.String())
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
