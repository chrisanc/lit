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
	ext := filepath.Ext(filePath)
	if len(ext) <= 1 {
		return nil
	}
	activeLanguage := analyzer.getLanguage(ext[1:])
	if activeLanguage == nil {
		return nil
	}

	// Calculate the cyclical complexity and get the functions returned
	functions := analysis.CyclicalComplexity(activeLanguage, code)
	if functions == nil {
		return nil
	}

	messages := analyzer.feedback.GetMessages()
	i := 0
	for i < len(functions) {
		if functions[i] == nil {
			functions = append(functions[:i], functions[i+1:]...)
			continue
		}
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
	ext := filepath.Ext(filePath)
	if len(ext) <= 1 {
		return 0
	}
	activeLanguage := analyzer.getLanguage(ext[1:])
	if activeLanguage == nil {
		return 0
	}
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
