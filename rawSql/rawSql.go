// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package rawSql

import (
	"go/ast"
	"go/token"
	"strings"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

const sqlstorePackagePath = "github.com/mattermost/mattermost-server/v5/store/sqlstore"

var Analyzer = &analysis.Analyzer{
	Name: "rawSql",
	Doc: "check invalid usage of raw SQL queries instead of using the squirrel lib",
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !strings.HasPrefix(pass.Pkg.Path(), sqlstorePackagePath) {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.BasicLit:
				if n.Kind != token.STRING {
					return false
				}

				unquoted, err := strconv.Unquote(n.Value)
				if err != nil {
					return false
				}

				words := strings.Fields(strings.TrimSpace(unquoted))
				if len(words) == 0 {
					return false
				}

				switch lead := strings.ToLower(words[0]); lead {
				case "select":
					pass.Reportf(n.Pos(), "Found leading \"select\" in a string. Creating raw SQL queries is not allowed, please use the squirrel builder instead.")
				case "insert":
					pass.Reportf(n.Pos(), "Found leading \"insert\" in a string. Creating raw SQL queries is not allowed, please use the squirrel builder instead.")
				case "update":
					pass.Reportf(n.Pos(), "Found leading \"update\" in a string. Creating raw SQL queries is not allowed, please use the squirrel builder instead.")
				}
				return false
			}
			return true
		})
	}
	return nil, nil
}
