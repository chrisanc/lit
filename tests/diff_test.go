package tests

import (
	"CLI_App/internal/service"
	"strings"
	"testing"
)

func TestGenerateUnifiedDiff_ModifiedLine(t *testing.T) {
	orig := []string{"func foo_bar() {"}
	mod := []string{"func fooBar() {"}
	filename := "test.go"

	diff := service.GenerateUnifiedDiff(filename, orig, mod)

	if !strings.Contains(diff, "--- a/test.go") {
		t.Errorf("expected header '--- a/test.go' in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "+++ b/test.go") {
		t.Errorf("expected header '+++ b/test.go' in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "\033[36m@@ -1,1 +1,1 @@\033[0m") {
		t.Errorf("expected hunk header '@@ -1,1 +1,1 @@' with cyan ANSI in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "\033[31m-func foo_bar() {\033[0m") {
		t.Errorf("expected red deleted line in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "\033[32m+func fooBar() {\033[0m") {
		t.Errorf("expected green added line in diff, got:\n%s", diff)
	}
}

func TestGenerateUnifiedDiff_Unchanged(t *testing.T) {
	orig := []string{"line1", "line2", "line3"}
	mod := []string{"line1", "line2", "line3"}
	filename := "test.go"

	diff := service.GenerateUnifiedDiff(filename, orig, mod)

	if diff != "" {
		t.Errorf("expected empty diff for unchanged slices, got:\n%s", diff)
	}
}

func TestGenerateUnifiedDiff_InsertionAndDeletion(t *testing.T) {
	orig := []string{"line1", "line2", "line3"}
	mod := []string{"line1", "line2_modified", "line3", "line4"}
	filename := "sample.py"

	diff := service.GenerateUnifiedDiff(filename, orig, mod)

	if !strings.Contains(diff, "--- a/sample.py") {
		t.Errorf("expected header '--- a/sample.py' in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "+++ b/sample.py") {
		t.Errorf("expected header '+++ b/sample.py' in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "\033[31m-line2\033[0m") {
		t.Errorf("expected red deleted '-line2', got:\n%s", diff)
	}
	if !strings.Contains(diff, "\033[32m+line2_modified\033[0m") {
		t.Errorf("expected green added '+line2_modified', got:\n%s", diff)
	}
	if !strings.Contains(diff, "\033[32m+line4\033[0m") {
		t.Errorf("expected green added '+line4', got:\n%s", diff)
	}
}
