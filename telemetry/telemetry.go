// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package telemetry

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	telemetryPkgPath = "github.com/mattermost/mattermost-server/v5/services/telemetry"
	modelPkgPath     = "github.com/mattermost/mattermost-server/v5/model"
)

// Analyzer describes telemetry status analysis function for a config setting
var Analyzer = &analysis.Analyzer{
	Name: "telemetry",
	Doc:  "reports if a config field is not used in services/telemetry package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Path() != telemetryPkgPath {
		return nil, nil
	}

	// we need to find model.Config definition to parse its Type Spec
	var configObject *types.TypeName
	for _, obj := range pass.TypesInfo.Uses {
		if tn, ok := obj.(*types.TypeName); ok && obj.Name() == "Config" && obj.Pkg().Path() == modelPkgPath {
			configObject = tn
			break
		}
	}

	if configObject == nil {
		return nil, errors.New("could not find model.Config type")
	}

	configMap, err := typeFieldMap(pass.Fset, configObject)
	if err != nil {
		return nil, fmt.Errorf("could generate fields map: %w", err)
	}

	var allExpressions []string

	// We look every expression in the file so that we can
	// check if the expected references exist in this package.

	// NOTE: this approach does not check if the reference (config field) is added
	// to telemetry queue, it assumes that if a config field is used in the file,
	// it probably added to the telemetry. This is a naive approach but should work.
	// Also, we can narrow down this for specific function calls, declarations etc.
	// so that we can be sure if the config field is added to telemetry but all these
	// fine tuning will require extra maintenance cost.
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			expr, ok := n.(ast.Expr)
			if !ok {
				return true
			}

			sel, ok := expr.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			allExpressions = append(allExpressions, getFieldAccessSequence(sel))

			return true
		})
	}

	// now remove the expresions from what's expected
	for _, ref := range allExpressions {
		delete(configMap, ref)
	}

	// report remaining fields from the map
	for k, v := range configMap {
		pass.Reportf(v, "%s is not used in telemetry\n", k)
	}

	return nil, nil
}

// getFieldAccessSequence is used recursively to get string representation of
// a selector expression sequence. e.g. config.ServiceSettings.SiteURL
func getFieldAccessSequence(sel *ast.SelectorExpr) string {
	switch v := sel.X.(type) {
	case *ast.SelectorExpr:
		return strings.Join([]string{getFieldAccessSequence(v), sel.Sel.Name}, ".")
	default:
		return sel.Sel.Name
	}
}
