package analysis_test

import (
	"CLI_App/internal/adapters/analysis"
	"CLI_App/internal/adapters/analysis/languages"
	"testing"
)

var goCode = []string{
	"package main",
	"",
	"import \"fmt\"",
	"",
	"func CalculateSum(a int, b int) int {",
	"	if a > 0 && b > 0 {",
	"		return a + b",
	"	}",
	"	return 0",
	"}",
	"",
	"type User struct {",
	"	ID int",
	"	Name string",
	"}",
}

var pythonCode = []string{
	"def calculate_sum(a, b):",
	"    if a > 0 and b > 0:",
	"        return a + b",
	"    return 0",
	"",
	"class User:",
	"    def __init__(self, user_id, name):",
	"        self.user_id = user_id",
	"        self.name = name",
}

var jsCode = []string{
	"function calculateSum(a, b) {",
	"    if (a > 0 && b > 0) {",
	"        return a + b;",
	"    }",
	"    return 0;",
	"}",
	"",
	"class User {",
	"    constructor(id, name) {",
	"        this.id = id;",
	"        this.name = name;",
	"    }",
	"}",
}

func BenchmarkASTParser_Go(b *testing.B) {
	lang := languages.NewGolangLanguage("camelCase")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analysis.CyclicalComplexity(lang, &goCode)
	}
}

func BenchmarkASTParser_Python(b *testing.B) {
	lang := languages.NewPythonLanguage("snake_case")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analysis.CyclicalComplexity(lang, &pythonCode)
	}
}

func BenchmarkASTParser_JS(b *testing.B) {
	lang := languages.NewJSLanguage("camelCase")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analysis.CyclicalComplexity(lang, &jsCode)
	}
}
