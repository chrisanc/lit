package tests

import (
	"CLI_App/internal/adapters/analysis/languages"
	"CLI_App/internal/domain"
	"CLI_App/internal/service"
	"context"
	"testing"
)

func createBenchAnalyzer() domain.Analyzer {
	cfg := &domain.Config{
		NamingConventionIndex: 1,
		Alerts: domain.Alerts{
			Parameters: domain.FeedbackValues{Info: 3, Warning: 5, Error: 8},
			Complexity: domain.FeedbackValues{Info: 5, Warning: 10, Error: 15},
			MethodSize: domain.FeedbackValues{Info: 30, Warning: 50, Error: 100},
		},
	}
	convIdx := cfg.NamingConventionIndex
	if convIdx < 1 || int(convIdx) > len(domain.Conventions) {
		convIdx = 1
	}
	return languages.NewFileAnalyzer(
		domain.Conventions[convIdx-1],
		domain.NewFeedback(cfg),
		convIdx,
	)
}

func BenchmarkScanner_ScanFiles(b *testing.B) {
	analyzer := createBenchAnalyzer()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := service.NewScannerService(analyzer)
		scanner.ScanFiles(ctx)
	}
}

func BenchmarkScanner_ExecuteLOC(b *testing.B) {
	analyzer := createBenchAnalyzer()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := service.NewScannerService(analyzer)
		scanner.ExecuteLOC(ctx)
	}
}

func BenchmarkScanner_FixFileDryRun(b *testing.B) {
	analyzer := createBenchAnalyzer()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner := service.NewScannerService(analyzer)
		scanner.FixFile(ctx, true)
	}
}
