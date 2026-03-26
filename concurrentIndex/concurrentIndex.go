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

func checkStatement(stmt string) *diagnostic {
	if createIndexRegex.MatchString(stmt) && !concurrentlyRegex.MatchString(stmt) {
		return &diagnostic{message: "use CREATE INDEX CONCURRENTLY instead of CREATE INDEX to avoid blocking DML"}
	}
	if dropIndexRegex.MatchString(stmt) && !dropConcurrentlyRegex.MatchString(stmt) {
		return &diagnostic{message: "use DROP INDEX CONCURRENTLY instead of DROP INDEX to avoid blocking DML"}
	}
	return nil
}

func checkLine(line string) *diagnostic {
	for _, stmt := range strings.Split(line, ";") {
		stmt = strings.TrimSpace(stripSQLComments(stmt))
		if stmt == "" {
			continue
		}
		if d := checkStatement(stmt); d != nil {
			return d
		}
	}
	return nil
}

func stripSQLComments(s string) string {
	var b strings.Builder
	for _, raw := range strings.Split(s, "\n") {
		t := strings.TrimSpace(raw)
		if t == "" || strings.HasPrefix(t, "--") {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(t)
	}
	return b.String()
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
		if err := scanSQLDir(pass, sqlPath); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func scanSQLDir(pass *analysis.Pass, root string) error {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		return checkSQLFile(pass, path)
	})
}

func checkSQLFile(pass *analysis.Pass, name string) error {
	content, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	tf := pass.Fset.AddFile(name, -1, len(content))
	tf.SetLinesForContent(content)

	lines := strings.Split(string(content), "\n")
	var stmtBuf strings.Builder
	startLine := 1
	inBlockComment := false

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if inBlockComment {
			if idx := strings.Index(trimmed, "*/"); idx >= 0 {
				inBlockComment = false
				trimmed = strings.TrimSpace(trimmed[idx+2:])
			} else {
				continue
			}
		}

		if strings.HasPrefix(trimmed, "/*") {
			if idx := strings.Index(trimmed, "*/"); idx >= 0 {
				trimmed = strings.TrimSpace(trimmed[idx+2:])
			} else {
				inBlockComment = true
				continue
			}
		}

		if strings.HasPrefix(trimmed, "--") {
			continue
		}

		if trimmed == "" {
			continue
		}

		if stmtBuf.Len() == 0 {
			startLine = lineNum
		}
		if stmtBuf.Len() > 0 {
			stmtBuf.WriteByte(' ')
		}
		stmtBuf.WriteString(trimmed)

		for {
			current := stmtBuf.String()
			semi := strings.Index(current, ";")
			if semi < 0 {
				break
			}

			stmt := strings.TrimSpace(current[:semi])
			if stmt != "" {
				if d := checkStatement(stmt); d != nil {
					pass.Reportf(tf.LineStart(startLine), "%s", d.message)
				}
			}

			rest := strings.TrimSpace(current[semi+1:])
			stmtBuf.Reset()
			if rest != "" {
				stmtBuf.WriteString(rest)
				startLine = lineNum
			}
		}
	}

	if stmtBuf.Len() > 0 {
		if d := checkStatement(stmtBuf.String()); d != nil {
			pass.Reportf(tf.LineStart(startLine), "%s", d.message)
		}
	}

	return nil
}
