// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package openApiSync

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	"github.com/sajari/fuzzy"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "openApiSync",
	Doc:  "check for inconsistencies between OpenAPI spec and the source code",
	Run:  run,
}

var specFile string

func init() {
	Analyzer.Flags.StringVar(&specFile, "spec", "", "Path to the OpenAPI 3 YAML spec file")
}

func formatNode(fset *token.FileSet, node interface{}) string {
	var typeNameBuf bytes.Buffer
	printer.Fprint(&typeNameBuf, fset, node)
	return typeNameBuf.String()
}

func processRouterInit(pass *analysis.Pass, name string, routerPrefixes map[string]string, swagger *openapi3.Swagger, cm *fuzzy.Model) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			decl, ok := n.(*ast.FuncDecl)
			if !ok || decl.Name.Name != name {
				return true
			}
			for _, stmt := range decl.Body.List {
				expr, ok := stmt.(*ast.ExprStmt)
				if !ok {
					continue
				}
				aexpr, ok := expr.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.CallExpr)
				if !ok || len(aexpr.Args) != 2 {
					continue
				}
				prefix, _ := strconv.Unquote(aexpr.Args[0].(*ast.BasicLit).Value)
				method, _ := strconv.Unquote(expr.X.(*ast.CallExpr).Args[0].(*ast.BasicLit).Value)
				name := aexpr.Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr).Sel.Name
				prefix = strings.Replace(prefix, ":[A-Za-z0-9]+", "", -1)
				routerPrefix := strings.Replace(routerPrefixes[name], ":[A-Za-z0-9]+", "", -1)
				handler := routerPrefix + prefix
				// TODO: make this more generic
				if strings.Contains(handler, "websocket:websocket") { // ignore special cases
					continue
				}
				if strings.HasPrefix(handler, "api/v4") {
					handler = handler[7:]
				}
				if !strings.HasPrefix(handler, "/") {
					handler = "/" + handler
				}
				handlers := []string{handler}
				// TODO: make this more generic
				sync_param_group := "{syncable_type:teams|channels}"
				if strings.Contains(handler, sync_param_group) {
					handlers = []string{strings.Replace(handler, sync_param_group, "teams", 1), strings.Replace(handler, sync_param_group, "channels", 1)}
				}
				for _, h := range handlers {
					if path := swagger.Paths.Find(h); path == nil {
						pass.Reportf(aexpr.Pos(), "Cannot find %v method: %v in OpenAPI 3 spec. (maybe you meant: %v)", h, method, cm.Suggestions(h, false))
					} else if path.GetOperation(method) == nil {
						pass.Reportf(aexpr.Pos(), "Handler %v is defined with method %s, but in not in the spec", h, method)
					}
				}
			}

			return true
		})
	}
}

func parseRoutesStruct(pass *analysis.Pass, decl *ast.GenDecl, routerPrefixes map[string]string) {
	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok || spec.Name.String() != "Routes" {
		return
	}
	for _, f := range spec.Type.(*ast.StructType).Fields.List {
		typeName := formatNode(pass.Fset, f.Type)
		if typeName != "*mux.Router" {
			continue
		}
		routerName := f.Names[0].Name
		if routerName == "ApiRoot" || routerName == "Root" {
			continue
		}
		if f.Comment != nil && len(f.Comment.List) > 0 {
			comment := f.Comment.List[0].Text
			if strings.HasPrefix(comment, "// '") && strings.HasSuffix(comment, "'") {
				routerPrefixes[routerName] = comment[4 : len(comment)-1]
			} else {
				pass.Reportf(f.Comment.List[0].Pos(), "Comment for field %s is not formatted correctly\n", routerName)
			}
		} else {
			pass.Reportf(f.Comment.Pos(), "Router field %s in Router struct is not commented properly\n", routerName)
		}
	}
}

func parseInitFunction(pass *analysis.Pass, decl *ast.FuncDecl, routerPrefixes map[string]string, initFunctions []string) []string {
	for _, stmt := range decl.Body.List {
		if expr, ok := stmt.(*ast.ExprStmt); ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				sel := call.Fun.(*ast.SelectorExpr)
				if sel.X.(*ast.Ident).Name == "api" {
					initFunctions = append(initFunctions, sel.Sel.Name)
				}
			}
		} else if assgn, ok := stmt.(*ast.AssignStmt); ok {
			if len(assgn.Lhs) == 1 && strings.HasPrefix(formatNode(pass.Fset, assgn.Lhs[0]), "api.BaseRoutes") {
				subRouterName := formatNode(pass.Fset, assgn.Lhs[0])[15:]
				if subRouterName == "ApiRoot" || subRouterName == "Root" {
					continue
				}
				rhs := formatNode(pass.Fset, assgn.Rhs[0])[15:]
				router := rhs[:strings.Index(rhs, ".")]
				path := rhs[strings.Index(rhs, ".")+13 : strings.LastIndex(rhs, ".")-2]
				prefix := ""
				switch router {
				case "ApiRoot":
					prefix = "api/v4"
				case "Root":
					prefix = ""
				default:
					if s, ok := routerPrefixes[router]; ok {
						prefix = s
					} else {
						pass.Reportf(assgn.Rhs[0].Pos(), "cannot find prefix for %s\n", router)

					}
				}
				s := fmt.Sprintf("%v%v", prefix, path)
				s2 := routerPrefixes[subRouterName]
				if s2 != s {
					pass.Reportf(assgn.Rhs[0].Pos(), "PathPrefix doesn't match field comment for field '%s': '%s' vs '%s'\n", subRouterName, s, s2)
				}

			}
		}
	}
	return initFunctions
}

func validateComments(pass *analysis.Pass) ([]string, map[string]string) {
	initFunctions := []string{}
	routerPrefixes := map[string]string{}
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if decl, ok := n.(*ast.GenDecl); ok && len(decl.Specs) == 1 {
				parseRoutesStruct(pass, decl, routerPrefixes)
			} else if decl, ok := n.(*ast.FuncDecl); ok && decl.Name.Name == "Init" {
				initFunctions = parseInitFunction(pass, decl, routerPrefixes, initFunctions)
			}
			return true
		})
	}
	return initFunctions, routerPrefixes
}

func run(pass *analysis.Pass) (interface{}, error) {
	if specFile == "" {
		return nil, errors.New("Please supply a path to OpenAPI spec yaml file via -openApiSync.spec")
	}
	if _, err := os.Stat(specFile); err != nil {
		return nil, errors.Wrapf(err, "spec file does not exist")
	}
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(specFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse swagger")
	}

	initFunctions, routerPrefixes := validateComments(pass)

	var swaggerPaths []string
	for p := range swagger.Paths {
		swaggerPaths = append(swaggerPaths, p)
	}
	model := fuzzy.NewModel()
	model.Train(swaggerPaths)

	for _, initFunc := range initFunctions {
		processRouterInit(pass, initFunc, routerPrefixes, swagger, model)
	}

	return nil, nil
}
