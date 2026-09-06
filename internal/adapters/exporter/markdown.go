package exporter

import (
	"CLI_App/internal/domain"
	"bytes"
	"fmt"
)

type MarkdownExporter struct{}

func (m *MarkdownExporter) Export(results map[string][]*domain.FunctionData) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("# Lit Static Analysis Report\n\n")

	total := 0
	for _, fnList := range results {
		total += len(fnList)
	}

	buf.WriteString(fmt.Sprintf("**Total Potential Improvements Found:** %d\n\n", total))
	buf.WriteString("| File | Symbol | Position | Complexity | Parameters | Lines | Invalid Names |\n")
	buf.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- | :--- |\n")

	for file, methods := range results {
		for _, fn := range methods {
			pos := fmt.Sprintf("%d:%d", fn.StartPosition.Row+1, fn.StartPosition.Column+1)
			buf.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | %d | %d | %d | %d |\n",
				file, fn.Name, pos, fn.Complexity, fn.TotalParams, fn.Size, fn.InvalidNames))
		}
	}

	return buf.Bytes(), nil
}
