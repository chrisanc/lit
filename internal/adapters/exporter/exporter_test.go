package exporter_test

import (
	"CLI_App/internal/adapters/exporter"
	"CLI_App/internal/domain"
	"encoding/json"
	"strings"
	"testing"
)

func sampleData() map[string][]*domain.FunctionData {
	return map[string][]*domain.FunctionData{
		"main.go": {
			{
				Name:          "ProcessData",
				TotalParams:   6,
				Complexity:    12,
				InvalidNames:  2,
				StartPosition: domain.Point{Row: 10, Column: 1},
				Size:          65,
				Feedback:      "High complexity",
			},
		},
	}
}

func TestGetExporter(t *testing.T) {
	formats := []string{"sarif", "json", "markdown", "md", "text", ""}
	for _, f := range formats {
		exp, err := exporter.GetExporter(f)
		if err != nil {
			t.Errorf("GetExporter(%q) returned unexpected error: %v", f, err)
		}
		if exp == nil {
			t.Errorf("GetExporter(%q) returned nil exporter", f)
		}
	}

	_, err := exporter.GetExporter("invalid")
	if err == nil {
		t.Errorf("GetExporter(\"invalid\") expected error, got nil")
	}
}

func TestSarifExporter(t *testing.T) {
	exp, _ := exporter.GetExporter("sarif")
	data, err := exp.Export(sampleData())
	if err != nil {
		t.Fatalf("Sarif export failed: %v", err)
	}

	var report exporter.SarifReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse SARIF output as JSON: %v", err)
	}

	if report.Version != "2.1.0" {
		t.Errorf("Expected version 2.1.0, got %s", report.Version)
	}
	if len(report.Runs) == 0 {
		t.Fatalf("Expected at least one run in SARIF report")
	}
	if len(report.Runs[0].Results) < 4 {
		t.Errorf("Expected at least 4 SARIF results for LIT001-LIT004, got %d", len(report.Runs[0].Results))
	}

	ruleIDs := make(map[string]bool)
	for _, r := range report.Runs[0].Results {
		ruleIDs[r.RuleID] = true
	}

	for _, expectedID := range []string{"LIT001", "LIT002", "LIT003", "LIT004"} {
		if !ruleIDs[expectedID] {
			t.Errorf("Expected rule %s in SARIF results, but not found", expectedID)
		}
	}
}

func TestJSONExporter(t *testing.T) {
	exp, _ := exporter.GetExporter("json")
	data, err := exp.Export(sampleData())
	if err != nil {
		t.Fatalf("JSON export failed: %v", err)
	}

	var report exporter.JSONReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if report.TotalImprovements != 1 {
		t.Errorf("Expected TotalImprovements = 1, got %d", report.TotalImprovements)
	}
	if len(report.Files) != 1 {
		t.Fatalf("Expected 1 file in report, got %d", len(report.Files))
	}
	if report.Files[0].FilePath != "main.go" {
		t.Errorf("Expected file path main.go, got %s", report.Files[0].FilePath)
	}
}

func TestMarkdownExporter(t *testing.T) {
	exp, _ := exporter.GetExporter("markdown")
	data, err := exp.Export(sampleData())
	if err != nil {
		t.Fatalf("Markdown export failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "# Lit Static Analysis Report") {
		t.Errorf("Markdown report missing title header")
	}
	if !strings.Contains(content, "| `main.go` | `ProcessData` |") {
		t.Errorf("Markdown report missing expected table row")
	}
}

func TestTextExporter(t *testing.T) {
	exp, _ := exporter.GetExporter("text")
	data, err := exp.Export(sampleData())
	if err != nil {
		t.Fatalf("Text export failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Found 1 possible improvements") {
		t.Errorf("Text report missing expected summary line")
	}
	if !strings.Contains(content, "ProcessData") {
		t.Errorf("Text report missing function name")
	}
}
