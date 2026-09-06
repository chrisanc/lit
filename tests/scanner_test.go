package tests

import (
	"CLI_App/internal/adapters/analysis/languages"
	"CLI_App/internal/domain"
	"CLI_App/internal/service"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestScannerExportResults(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "TestScannerExportResults")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &domain.Config{
		NamingConventionIndex: 1,
		Alerts: domain.Alerts{
			Parameters: domain.FeedbackValues{Info: 1, Warning: 2, Error: 3},
			Complexity: domain.FeedbackValues{Info: 1, Warning: 2, Error: 3},
			MethodSize: domain.FeedbackValues{Info: 1, Warning: 2, Error: 3},
		},
	}

	analyzer := languages.NewFileAnalyzer(cfg.GetVariableConvention(), cfg.GetFunctionConvention(), domain.NewFeedback(cfg), cfg.GetVariableConventionIndex(), cfg.GetFunctionConventionIndex(), cfg.IgnoredSymbols)
	scanner := service.NewScannerService(analyzer)

	ctx := context.Background()
	scanner.ScanFiles(ctx)

	formats := []string{"text", "sarif", "json", "markdown"}

	for _, fmtName := range formats {
		outPath := filepath.Join(tmpDir, "report."+fmtName)
		if err := scanner.ExportResults(fmtName, outPath); err != nil {
			t.Errorf("ExportResults failed for format '%s': %v", fmtName, err)
		}

		info, err := os.Stat(outPath)
		if err != nil {
			t.Errorf("expected generated file for format '%s', got stat error: %v", fmtName, err)
		} else if info.Size() == 0 {
			t.Errorf("expected non-empty generated file for format '%s'", fmtName)
		}
	}
}
