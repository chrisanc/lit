package exporter

import (
	"CLI_App/internal/domain"
	"encoding/json"
	"fmt"
)

// OASIS SARIF v2.1.0 Schema Structs
type SarifReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SarifRun `json:"runs"`
}

type SarifRun struct {
	Tool    SarifTool     `json:"tool"`
	Results []SarifResult `json:"results"`
}

type SarifTool struct {
	Driver SarifDriver `json:"driver"`
}

type SarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SarifRule `json:"rules"`
}

type SarifRule struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	ShortDescription SarifMultiformat `json:"shortDescription"`
}

type SarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   SarifMessage    `json:"message"`
	Locations []SarifLocation `json:"locations"`
}

type SarifMessage struct {
	Text string `json:"text"`
}

type SarifLocation struct {
	PhysicalLocation SarifPhysicalLocation `json:"physicalLocation"`
}

type SarifPhysicalLocation struct {
	ArtifactLocation SarifArtifactLocation `json:"artifactLocation"`
	Region           SarifRegion           `json:"region"`
}

type SarifArtifactLocation struct {
	URI string `json:"uri"`
}

type SarifRegion struct {
	StartLine   uint `json:"startLine"`
	StartColumn uint `json:"startColumn"`
}

type SarifMultiformat struct {
	Text string `json:"text"`
}

type SarifExporter struct{}

func (s *SarifExporter) Export(results map[string][]*domain.FunctionData) ([]byte, error) {
	rules := []SarifRule{
		{
			ID:               "LIT001",
			Name:             "HighCyclomaticComplexity",
			ShortDescription: SarifMultiformat{Text: "Function or method has elevated cyclomatic complexity"},
		},
		{
			ID:               "LIT002",
			Name:             "InvalidNamingConvention",
			ShortDescription: SarifMultiformat{Text: "Identifier violates configured naming convention rules"},
		},
		{
			ID:               "LIT003",
			Name:             "ExcessiveParameterCount",
			ShortDescription: SarifMultiformat{Text: "Function accepts more parameters than recommended"},
		},
		{
			ID:               "LIT004",
			Name:             "ExcessiveMethodLength",
			ShortDescription: SarifMultiformat{Text: "Method length exceeds recommended line count"},
		},
	}

	sarifResults := []SarifResult{}

	for filePath, functions := range results {
		for _, fn := range functions {
			startRow := fn.StartPosition.Row + 1
			startCol := fn.StartPosition.Column + 1

			location := SarifLocation{
				PhysicalLocation: SarifPhysicalLocation{
					ArtifactLocation: SarifArtifactLocation{URI: filePath},
					Region:           SarifRegion{StartLine: startRow, StartColumn: startCol},
				},
			}

			if fn.Complexity > 1 {
				sarifResults = append(sarifResults, SarifResult{
					RuleID:    "LIT001",
					Level:     "warning",
					Message:   SarifMessage{Text: fmt.Sprintf("Function '%s' has cyclomatic complexity of %d", fn.Name, fn.Complexity)},
					Locations: []SarifLocation{location},
				})
			}

			if fn.InvalidNames > 0 {
				sarifResults = append(sarifResults, SarifResult{
					RuleID:    "LIT002",
					Level:     "warning",
					Message:   SarifMessage{Text: fmt.Sprintf("Function '%s' contains %d identifier naming convention violations", fn.Name, fn.InvalidNames)},
					Locations: []SarifLocation{location},
				})
			}

			if fn.TotalParams > 4 {
				sarifResults = append(sarifResults, SarifResult{
					RuleID:    "LIT003",
					Level:     "warning",
					Message:   SarifMessage{Text: fmt.Sprintf("Function '%s' accepts %d parameters", fn.Name, fn.TotalParams)},
					Locations: []SarifLocation{location},
				})
			}

			if fn.Size > 50 {
				sarifResults = append(sarifResults, SarifResult{
					RuleID:    "LIT004",
					Level:     "warning",
					Message:   SarifMessage{Text: fmt.Sprintf("Function '%s' spans %d lines of code", fn.Name, fn.Size)},
					Locations: []SarifLocation{location},
				})
			}
		}
	}

	report := SarifReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SarifRun{
			{
				Tool: SarifTool{
					Driver: SarifDriver{
						Name:           "lit",
						Version:        "1.0 release",
						InformationURI: "https://github.com/chrisanc/lit",
						Rules:          rules,
					},
				},
				Results: sarifResults,
			},
		},
	}

	return json.MarshalIndent(report, "", "  ")
}
