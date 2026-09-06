package service_test

import (
	"CLI_App/internal/domain"
	"CLI_App/internal/service"
	"os"
	"path/filepath"
	"testing"
)

type mockAnalyzer struct{}

func (m *mockAnalyzer) AnalyzeFile(filePath string, code *[]string) []*domain.FunctionData {
	return []*domain.FunctionData{
		{
			Name:          "MockFunc",
			TotalParams:   5,
			Complexity:    3,
			InvalidNames:  1,
			StartPosition: domain.Point{Row: 1, Column: 0},
			Size:          55,
		},
	}
}

func (m *mockAnalyzer) FixFile(filePath string, code *[]string) int {
	return 0
}

func TestScannerExportResults(t *testing.T) {
	scanner := service.NewScannerService(&mockAnalyzer{})
	
	// Access GetDangerousFunctions
	funcs := scanner.GetDangerousFunctions()
	if funcs == nil {
		t.Fatalf("GetDangerousFunctions returned nil")
	}

	tempDir := t.TempDir()

	for _, format := range []string{"text", "sarif", "json", "markdown"} {
		outputPath := filepath.Join(tempDir, "report."+format)
		err := scanner.ExportResults(format, outputPath)
		if err != nil {
			t.Errorf("ExportResults(%q, %q) failed: %v", format, outputPath, err)
		}

		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Errorf("Expected output file %s to exist for format %s", outputPath, format)
		}
	}
}
