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

const (
	diagCreateIndex = "use CREATE INDEX CONCURRENTLY instead of CREATE INDEX to avoid blocking DML"
	diagDropIndex   = "use DROP INDEX CONCURRENTLY instead of DROP INDEX to avoid blocking DML"
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

func checkStatement(stmt string) string {
	if createIndexRegex.MatchString(stmt) && !concurrentlyRegex.MatchString(stmt) {
		return diagCreateIndex
	}
	if dropIndexRegex.MatchString(stmt) && !dropConcurrentlyRegex.MatchString(stmt) {
		return diagDropIndex
	}
	return ""
}

func hasNolintDirective(s string) bool {
	if !strings.Contains(s, "nolint:concurrentIndex") {
		return false
	}
	for _, raw := range strings.Split(s, "\n") {
		t := strings.TrimSpace(raw)
		if strings.HasPrefix(t, "--") && strings.Contains(t, "nolint:concurrentIndex") {
			return true
		}
	}
	return false
}

func containsIndexKeyword(s string) bool {
	return strings.Contains(strings.ToLower(s), "index")
}

func checkLine(line string) string {
	if !containsIndexKeyword(line) {
		return ""
	}
	if hasNolintDirective(line) {
		return ""
	}
	for _, stmt := range strings.Split(line, ";") {
		stmt = strings.TrimSpace(stripSQLComments(stmt))
		if stmt == "" {
			continue
		}
		if msg := checkStatement(stmt); msg != "" {
			return msg
		}
	}
	return ""
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

			if msg := checkLine(val); msg != "" {
				pass.Reportf(lit.Pos(), "%s", msg)
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
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
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
	nolint := false

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
			if strings.Contains(trimmed, "nolint:concurrentIndex") {
				nolint = true
			}
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
				if nolint {
					nolint = false
				} else if msg := checkStatement(stmt); msg != "" {
					pass.Reportf(tf.LineStart(startLine), "%s", msg)
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

	if stmtBuf.Len() > 0 && !nolint {
		if msg := checkStatement(stmtBuf.String()); msg != "" {
			pass.Reportf(tf.LineStart(startLine), "%s", msg)
		}
	}

	return nil
}
