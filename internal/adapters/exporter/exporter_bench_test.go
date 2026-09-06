package exporter_test

import (
	"CLI_App/internal/adapters/exporter"
	"CLI_App/internal/domain"
	"testing"
)

func benchData() map[string][]*domain.FunctionData {
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
			{
				Name:          "HandleRequest",
				TotalParams:   3,
				Complexity:    5,
				InvalidNames:  0,
				StartPosition: domain.Point{Row: 80, Column: 1},
				Size:          30,
				Feedback:      "OK",
			},
		},
		"utils.go": {
			{
				Name:          "ParseConfig",
				TotalParams:   2,
				Complexity:    8,
				InvalidNames:  1,
				StartPosition: domain.Point{Row: 15, Column: 1},
				Size:          45,
				Feedback:      "Moderate complexity",
			},
		},
	}
}

func BenchmarkSARIFExporter(b *testing.B) {
	exp, err := exporter.GetExporter("sarif")
	if err != nil {
		b.Fatalf("failed to get exporter: %v", err)
	}
	data := benchData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exp.Export(data)
	}
}

func BenchmarkJSONExporter(b *testing.B) {
	exp, err := exporter.GetExporter("json")
	if err != nil {
		b.Fatalf("failed to get exporter: %v", err)
	}
	data := benchData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exp.Export(data)
	}
}

func BenchmarkMarkdownExporter(b *testing.B) {
	exp, err := exporter.GetExporter("markdown")
	if err != nil {
		b.Fatalf("failed to get exporter: %v", err)
	}
	data := benchData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exp.Export(data)
	}
}
