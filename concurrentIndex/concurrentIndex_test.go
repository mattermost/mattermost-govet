package concurrentIndex

import (
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGoFiles(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}

func TestSQLDir(t *testing.T) {
	testdata := analysistest.TestData()
	migrationsDir := filepath.Join(testdata, "migrations")

	fset := token.NewFileSet()
	var diags []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer,
		Fset:     fset,
		Report: func(d analysis.Diagnostic) {
			diags = append(diags, d)
		},
	}

	scanSQLDir(pass, migrationsDir)

	var files []string
	for _, d := range diags {
		pos := fset.Position(d.Pos)
		files = append(files, filepath.Base(pos.Filename)+":"+d.Message)
	}

	createCount := 0
	dropCount := 0
	for _, f := range files {
		switch {
		case strings.Contains(f, "CREATE INDEX CONCURRENTLY"):
			createCount++
		case strings.Contains(f, "DROP INDEX CONCURRENTLY"):
			dropCount++
		}
	}
	assert.Equal(t, 2, createCount, "expected 2 CREATE INDEX diagnostics")
	assert.Equal(t, 2, dropCount, "expected 2 DROP INDEX diagnostics")

	// Verify correct files were flagged
	badCount := 0
	nestedCount := 0
	for _, f := range files {
		switch {
		case strings.HasPrefix(f, "000002_bad"):
			badCount++
		case strings.HasPrefix(f, "000003_nested"):
			nestedCount++
		}
	}
	assert.Equal(t, 3, badCount, "expected 3 diagnostics from 000002_bad.up.sql")
	assert.Equal(t, 1, nestedCount, "expected 1 diagnostic from sub/000003_nested.up.sql")
}

func TestSQLDirSkipsComments(t *testing.T) {
	testdata := analysistest.TestData()
	migrationsDir := filepath.Join(testdata, "migrations")

	fset := token.NewFileSet()
	var diags []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer,
		Fset:     fset,
		Report: func(d analysis.Diagnostic) {
			diags = append(diags, d)
		},
	}

	scanSQLDir(pass, migrationsDir)

	for _, d := range diags {
		pos := fset.Position(d.Pos)
		require.NotContains(t, filepath.Base(pos.Filename), "000001_good",
			"good migration file should not have diagnostics")
	}
}

func TestCheckLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantMsg string
	}{
		{
			name:    "CREATE INDEX without CONCURRENTLY",
			line:    "CREATE INDEX IF NOT EXISTS idx_foo_bar ON foo (bar);",
			wantMsg: "use CREATE INDEX CONCURRENTLY instead of CREATE INDEX to avoid blocking DML",
		},
		{
			name:    "CREATE UNIQUE INDEX without CONCURRENTLY",
			line:    "CREATE UNIQUE INDEX IF NOT EXISTS idx_foo_bar ON foo (bar);",
			wantMsg: "use CREATE INDEX CONCURRENTLY instead of CREATE INDEX to avoid blocking DML",
		},
		{
			name:    "DROP INDEX without CONCURRENTLY",
			line:    "DROP INDEX IF EXISTS idx_foo_bar;",
			wantMsg: "use DROP INDEX CONCURRENTLY instead of DROP INDEX to avoid blocking DML",
		},
		{
			name:    "lowercase create index",
			line:    "create index if not exists idx_foo on foo (bar);",
			wantMsg: "use CREATE INDEX CONCURRENTLY instead of CREATE INDEX to avoid blocking DML",
		},
		{
			name:    "lowercase drop index",
			line:    "drop index if exists idx_foo;",
			wantMsg: "use DROP INDEX CONCURRENTLY instead of DROP INDEX to avoid blocking DML",
		},
		{
			name: "CREATE INDEX CONCURRENTLY is fine",
			line: "CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_foo_bar ON foo (bar);",
		},
		{
			name: "CREATE UNIQUE INDEX CONCURRENTLY is fine",
			line: "CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_foo_bar ON foo (bar);",
		},
		{
			name: "DROP INDEX CONCURRENTLY is fine",
			line: "DROP INDEX CONCURRENTLY IF EXISTS idx_foo_bar;",
		},
		{
			name: "unrelated SQL",
			line: "CREATE TABLE foo (id int);",
		},
		{
			name: "SELECT statement",
			line: "SELECT * FROM indexes;",
		},
		{
			name: "plain string",
			line: "just a regular string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkLine(tt.line)
			if tt.wantMsg == "" {
				assert.Nil(t, d)
			} else {
				assert.NotNil(t, d)
				assert.Equal(t, tt.wantMsg, d.message)
			}
		})
	}
}
