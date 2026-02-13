package cmd_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestNoDirectPtermColorHelpersOutsideThemeRegistry(t *testing.T) {
	t.Parallel()

	violations, err := collectPtermColorHelperViolations(repoRootFromThisFile(t), []string{"cmd/theme.go"})
	if err != nil {
		t.Fatalf("collect pterm color helper violations: %v", err)
	}

	if len(violations) > 0 {
		slices.Sort(violations)
		t.Fatalf(
			"found direct pterm color helper usage outside theme registry:\n%s\n\nuse pterm.ThemeDefault.<Style>.Sprint(...) instead",
			strings.Join(violations, "\n"),
		)
	}
}

func collectPtermColorHelperViolations(repoRoot string, allowList []string) ([]string, error) {
	fset := token.NewFileSet()
	violations := make([]string, 0)

	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if shouldSkipPath(path, d) {
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if slices.Contains(allowList, rel) {
			return nil
		}

		fileViolations, err := fileColorHelperViolations(path, rel, fset)
		if err != nil {
			return err
		}
		violations = append(violations, fileViolations...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return violations, nil
}

func shouldSkipPath(_ string, d fs.DirEntry) bool {
	if d == nil || !d.IsDir() {
		return false
	}
	return d.Name() == ".git"
}

func fileColorHelperViolations(path string, rel string, fset *token.FileSet) ([]string, error) {
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	violations := make([]string, 0)
	ast.Inspect(file, func(node ast.Node) bool {
		sel, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "pterm" {
			return true
		}
		if !isDirectColorHelper(sel.Sel.Name) {
			return true
		}

		pos := fset.Position(sel.Sel.Pos())
		violations = append(violations, rel+":"+strconv.Itoa(pos.Line)+" uses pterm."+sel.Sel.Name)
		return true
	})

	return violations, nil
}

func repoRootFromThisFile(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve caller path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func isDirectColorHelper(name string) bool {
	if strings.HasPrefix(name, "Light") || strings.HasPrefix(name, "Dark") {
		return true
	}
	switch name {
	case "Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White":
		return true
	default:
		return false
	}
}
