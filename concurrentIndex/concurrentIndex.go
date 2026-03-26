package concurrentIndex

import (
	"go/ast"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	createIndexRegex      = regexp.MustCompile(`(?i)\bCREATE\s+(UNIQUE\s+)?INDEX\b`)
	concurrentlyRegex     = regexp.MustCompile(`(?i)\bCREATE\s+(UNIQUE\s+)?INDEX\s+CONCURRENTLY\b`)
	dropIndexRegex        = regexp.MustCompile(`(?i)\bDROP\s+INDEX\b`)
	dropConcurrentlyRegex = regexp.MustCompile(`(?i)\bDROP\s+INDEX\s+CONCURRENTLY\b`)
)

var sqlPath string

var Analyzer = &analysis.Analyzer{
	Name: "concurrentIndex",
	Doc:  "checks that CREATE INDEX and DROP INDEX use CONCURRENTLY to avoid blocking DML",
	Run:  run,
}

func init() {
	Analyzer.Flags.StringVar(&sqlPath, "path", "", "Relative path to a directory of .sql files to scan recursively")
}

type diagnostic struct {
	message string
}

func checkLine(line string) *diagnostic {
	if createIndexRegex.MatchString(line) && !concurrentlyRegex.MatchString(line) {
		return &diagnostic{message: "use CREATE INDEX CONCURRENTLY instead of CREATE INDEX to avoid blocking DML"}
	}
	if dropIndexRegex.MatchString(line) && !dropConcurrentlyRegex.MatchString(line) {
		return &diagnostic{message: "use DROP INDEX CONCURRENTLY instead of DROP INDEX to avoid blocking DML"}
	}
	return nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			lit, ok := node.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}

			val, err := strconv.Unquote(lit.Value)
			if err != nil {
				val = strings.Trim(lit.Value, "`")
			}

			if d := checkLine(val); d != nil {
				pass.Reportf(lit.Pos(), "%s", d.message)
			}
			return true
		})
	}

	if sqlPath != "" {
		scanSQLDir(pass, sqlPath)
	}

	return nil, nil
}

func scanSQLDir(pass *analysis.Pass, root string) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		checkSQLFile(pass, path)
		return nil
	})
}

func checkSQLFile(pass *analysis.Pass, name string) {
	content, err := os.ReadFile(name)
	if err != nil {
		return
	}

	tf := pass.Fset.AddFile(name, -1, len(content))
	tf.SetLinesForContent(content)

	for i, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		if d := checkLine(line); d != nil {
			pass.Reportf(tf.LineStart(i+1), "%s", d.message)
		}
	}
}
