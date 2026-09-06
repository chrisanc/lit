package tests

import (
	"CLI_App/internal/adapters/exporter"
	"CLI_App/internal/domain"
	"fmt"
	"testing"
)

func createMockResults(count int) map[string][]*domain.FunctionData {
	res := make(map[string][]*domain.FunctionData)
	for i := 0; i < count; i++ {
		file := fmt.Sprintf("internal/pkg%d/file%d.go", i/10, i)
		fn := &domain.FunctionData{
			Name:         fmt.Sprintf("ProcessData%d", i),
			TotalParams:  uint(i % 6),
			Complexity:   uint((i % 15) + 1),
			InvalidNames: uint(i % 3),
			Size:         uint((i % 100) + 10),
			StartPosition: domain.Point{
				Row:    uint(i * 5),
				Column: 4,
			},
			Feedback: "INFO: Method parameters or complexity exceed standard threshold.\n",
		}
		res[file] = append(res[file], fn)
	}
	return res
}

func BenchmarkSARIFExporter(b *testing.B) {
	mockData := createMockResults(500)
	exp := &exporter.SarifExporter{}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := exp.Export(mockData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONExporter(b *testing.B) {
	mockData := createMockResults(500)
	exp := &exporter.JSONExporter{}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := exp.Export(mockData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarkdownExporter(b *testing.B) {
	mockData := createMockResults(500)
	exp := &exporter.MarkdownExporter{}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := exp.Export(mockData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
