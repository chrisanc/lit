package tests

import (
	"CLI_App/internal/adapters/exporter"
	"CLI_App/internal/domain"
	"encoding/json"
	"strings"
	"testing"
)

func createMockFunctions() map[string][]*domain.FunctionData {
	return map[string][]*domain.FunctionData{
		"internal/service/scanner.go": {
			{
				Name:          "traverseFiles",
				TotalParams:   3,
				Complexity:    12,
				InvalidNames:  1,
				StartPosition: domain.Point{Row: 120, Column: 0},
				Size:          75,
				Feedback:      "INFO: High complexity detected.\n",
			},
		},
	}
}

func TestGetExporter(t *testing.T) {
	formats := []string{"sarif", "json", "markdown", "text", ""}
	for _, fmtName := range formats {
		exp, err := exporter.GetExporter(fmtName)
		if err != nil {
			t.Fatalf("expected exporter for format '%s', got error: %v", fmtName, err)
		}
		if exp == nil {
			t.Fatalf("expected non-nil exporter for format '%s'", fmtName)
		}
	}

	_, err := exporter.GetExporter("invalid_format")
	if err == nil {
		t.Fatalf("expected error for unsupported format, got nil")
	}
}

func TestSarifExporter(t *testing.T) {
	exp := &exporter.SarifExporter{}
	data, err := exp.Export(createMockFunctions())
	if err != nil {
		t.Fatalf("SarifExporter.Export failed: %v", err)
	}

	var report exporter.SarifReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("failed to unmarshal generated SARIF report: %v", err)
	}

	if report.Version != "2.1.0" {
		t.Errorf("expected SARIF version '2.1.0', got '%s'", report.Version)
	}

	if len(report.Runs) == 0 || len(report.Runs[0].Results) == 0 {
		t.Fatalf("expected SARIF results in run, got 0")
	}
}

func TestJSONExporter(t *testing.T) {
	exp := &exporter.JSONExporter{}
	data, err := exp.Export(createMockFunctions())
	if err != nil {
		t.Fatalf("JSONExporter.Export failed: %v", err)
	}

	var report exporter.JSONReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("failed to unmarshal JSON report: %v", err)
	}

	if report.TotalImprovements != 1 {
		t.Errorf("expected total improvements 1, got %d", report.TotalImprovements)
	}
}

func TestMarkdownExporter(t *testing.T) {
	exp := &exporter.MarkdownExporter{}
	data, err := exp.Export(createMockFunctions())
	if err != nil {
		t.Fatalf("MarkdownExporter.Export failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "# Lit Static Analysis Report") {
		t.Errorf("expected markdown title in report, got:\n%s", content)
	}
	if !strings.Contains(content, "internal/service/scanner.go") {
		t.Errorf("expected file path in markdown table, got:\n%s", content)
	}
}

func TestTextExporter(t *testing.T) {
	exp := &exporter.TextExporter{}
	data, err := exp.Export(createMockFunctions())
	if err != nil {
		t.Fatalf("TextExporter.Export failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Found 1 possible improvements") {
		t.Errorf("expected text summary header, got:\n%s", content)
	}
}
