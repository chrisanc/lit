package exporter

import (
	"CLI_App/internal/domain"
	"encoding/json"
)

type JSONExporter struct{}

type JSONResult struct {
	FilePath string                 `json:"file"`
	Methods  []*domain.FunctionData `json:"methods"`
}

type JSONReport struct {
	TotalImprovements int          `json:"totalImprovements"`
	Files             []JSONResult `json:"files"`
}

func (j *JSONExporter) Export(results map[string][]*domain.FunctionData) ([]byte, error) {
	total := 0
	var files []JSONResult

	for file, methods := range results {
		total += len(methods)
		files = append(files, JSONResult{
			FilePath: file,
			Methods:  methods,
		})
	}

	report := JSONReport{
		TotalImprovements: total,
		Files:             files,
	}

	return json.MarshalIndent(report, "", "  ")
}
