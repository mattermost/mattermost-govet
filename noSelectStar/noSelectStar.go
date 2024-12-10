package noSelectStar

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noSelectStar",
	Doc:  "checks for SQL queries containing SELECT * which breaks forwards compatibility",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		lit, ok := node.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		if strings.Contains(strings.ToUpper(lit.Value), "SELECT") &&
			strings.Contains(lit.Value, "*") {
			pass.Reportf(lit.Pos(), "do not use SELECT *: explicitly select the needed columns instead")
		}
		return true
	}

	inspectCalls := func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		if fun, ok := call.Fun.(*ast.Ident); ok {
			if fun.Name == "Select" || fun.Name == "Columns" {
				// Check all arguments for "*"
				for _, arg := range call.Args {
					if lit, ok := arg.(*ast.BasicLit); ok && 
						lit.Kind == token.STRING && 
						strings.Contains(lit.Value, "*") {
						pass.Reportf(lit.Pos(), "do not use %s with *: explicitly select the needed columns instead", fun.Name)
					}
				}
			}
		}
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
		ast.Inspect(f, inspectCalls)
	}
	return nil, nil
}
