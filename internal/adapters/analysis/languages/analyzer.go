package languages

import (
	"CLI_App/internal/adapters/analysis"
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"path/filepath"
)

type FileAnalyzer struct {
	activePattern         string
	feedback              *domain.Feedback
	namingConventionIndex int8
}

func NewFileAnalyzer(activePattern string, feedback *domain.Feedback, namingConventionIndex int8) *FileAnalyzer {
	return &FileAnalyzer{
		activePattern:         activePattern,
		feedback:              feedback,
		namingConventionIndex: namingConventionIndex,
	}
}

// AnalyzeFile analyses the file via DFS (is executed in the scanner)
func (analyzer *FileAnalyzer) AnalyzeFile(filePath string, code *[]string) []*domain.FunctionData {
	// Set the variable to save up the language of the current script
	activeLanguage := analyzer.getLanguage(filepath.Ext(filePath)[1:])

	// Calculate the cyclical complexity and get the functions returned
	functions := analysis.CyclicalComplexity(activeLanguage, code)

	messages := analyzer.feedback.GetMessages()
	i := 0
	for i < len(functions) {
		if functions[i].TotalParams < messages["parameters"][0].MinValue &&
			functions[i].Size < messages["size"][0].MinValue &&
			functions[i].Complexity < messages["complexity"][0].MinValue &&
			functions[i].InvalidNames < 1 &&
			functions[i].Feedback == "" {
			functions = append(functions[:i], functions[i+1:]...)
			continue
		}
		functions[i].SetFunctionFeedback(messages)
		i++
	}

	return functions
}

func (analyzer *FileAnalyzer) FixFile(filePath string, code *[]string) int {
	activeLanguage := analyzer.getLanguage(filepath.Ext(filePath)[1:])
	writer := analysis.NewFileModifier(activeLanguage, analyzer.activePattern, analyzer.namingConventionIndex)
	return writer.ModifyVariableName(code)
}

func (analyzer *FileAnalyzer) getLanguage(ext string) types.NodeManagement {
	// Save the language for the complexity
	switch ext {
	case "js", "jsx":
		return NewJSLanguage(analyzer.activePattern)
	case "go":
		return NewGolangLanguage(analyzer.activePattern)
	case "java":
		return NewJavaLanguage(analyzer.activePattern)
	case "py":
		return NewPythonLanguage(analyzer.activePattern)
	default:
		return nil
	}
}
