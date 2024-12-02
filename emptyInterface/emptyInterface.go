// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package emptyInterface

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	Analyzer = &analysis.Analyzer{
		Name: "emptyInterface",
		Doc:  "check for usage of interface{} instead of any",
		Run:  run,
	}
	ignoreFilesPattern string
)

func init() {
	Analyzer.Flags.StringVar(&ignoreFilesPattern, "ignore", "", "Comma separated list of files to ignore")
}

func run(pass *analysis.Pass) (interface{}, error) {
	var ignoreFiles []string
	if ignoreFilesPattern != "" {
		ignoreFiles = strings.Split(ignoreFilesPattern, ",")
	}

	checkList := []string{
		"Regenerate this file using `make plugin-mocks`.",
		"Regenerate this file using `make store-mocks`.",
		"Regenerate this file using `make sharedchannel-mocks`.",
	}

	for _, file := range pass.Files {
		comments := ""
		for _, commentGroup := range file.Comments {
			for _, comment := range commentGroup.List {
				comments = comments + comment.Text
			}
		}
		skip := false
		for _, check := range checkList {
			if strings.Contains(comments, check) {
				skip = true
			}
		}
		if skip {
			continue
		}

		if node := file; node != nil {
			f := pass.Fset.File(node.Pos())
			if isIgnore(f.Name(), ignoreFiles) {
				continue
			}
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.InterfaceType:
				if len(x.Methods.List) == 0 {
					pass.Reportf(x.Pos(), "using 'interface{}', please replace it with 'any'.")
				}
				return true
			}
			return true
		})
	}
	return nil, nil
}

func isIgnore(file string, ignoreFiles []string) bool {
	for _, f := range ignoreFiles {
		if strings.HasSuffix(file, f) {
			return true
		}
	}
	return false
}
