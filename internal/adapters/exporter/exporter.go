package exporter

import (
	"CLI_App/internal/domain"
	"fmt"
)

type ReportExporter interface {
	Export(results map[string][]*domain.FunctionData) ([]byte, error)
}

func GetExporter(format string) (ReportExporter, error) {
	switch format {
	case "sarif":
		return &SarifExporter{}, nil
	case "json":
		return &JSONExporter{}, nil
	case "markdown", "md":
		return &MarkdownExporter{}, nil
	case "text", "":
		return &TextExporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s (supported: text, sarif, json, markdown)", format)
	}
}
