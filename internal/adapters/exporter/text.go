package exporter

import (
	"CLI_App/internal/domain"
	"bytes"
	"fmt"
)

type TextExporter struct{}

func (t *TextExporter) Export(results map[string][]*domain.FunctionData) ([]byte, error) {
	var buf bytes.Buffer

	totalFunctions := 0
	for _, v := range results {
		totalFunctions += len(v)
	}

	buf.WriteString(fmt.Sprintf("Found %d possible improvements\n", totalFunctions))
	for key, value := range results {
		buf.WriteString(fmt.Sprintf("- %s:\n", key))
		for _, item := range value {
			buf.WriteString(fmt.Sprintf(" * %s (at %d:%d)\n", item.Name, item.StartPosition.Row, item.StartPosition.Column))
			buf.WriteString(fmt.Sprintf("   Parameters: %d\n   Total lines of code: %d\n", item.TotalParams, item.Size))
			buf.WriteString(fmt.Sprintf("   Found %d variables/methods/models with the wrong naming convention.\n", item.InvalidNames))
			if item.Feedback != "" {
				buf.WriteString(item.Feedback)
			}
		}
	}

	return buf.Bytes(), nil
}
