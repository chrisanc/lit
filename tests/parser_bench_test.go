package tests

import (
	"CLI_App/internal/adapters/analysis"
	"CLI_App/internal/adapters/analysis/languages"
	"testing"
)

var sampleGoCode = []string{
	"package sample",
	"import \"fmt\"",
	"type UserStruct struct { Name string }",
	"func CalculateTotal(a, b int) int {",
	"    if a > 0 && b > 0 {",
	"        return a + b",
	"    } else if a == 0 {",
	"        return b",
	"    }",
	"    return 0",
	"}",
}

var samplePyCode = []string{
	"class DataProcessor:",
	"    def __init__(self, data_list):",
	"        self.data_list = data_list",
	"    def process(self):",
	"        for item in self.data_list:",
	"            if item > 10 and item < 50:",
	"                print(item)",
}

var sampleJSCode = []string{
	"class UserSession {",
	"    constructor(user_id) {",
	"        this.user_id = user_id;",
	"    }",
	"    validateUser(role) {",
	"        if (role === 'admin' || role === 'root') {",
	"            return true;",
	"        }",
	"        return false;",
	"    }",
	"}",
}

func BenchmarkASTParser_Go(b *testing.B) {
	lang := languages.NewGolangLanguage("^[a-z]+$")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = analysis.CyclicalComplexity(lang, &sampleGoCode)
	}
}

func BenchmarkASTParser_Python(b *testing.B) {
	lang := languages.NewPythonLanguage("^[a-z]+$")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = analysis.CyclicalComplexity(lang, &samplePyCode)
	}
}

func BenchmarkASTParser_JS(b *testing.B) {
	lang := languages.NewJSLanguage("^[a-z]+$")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = analysis.CyclicalComplexity(lang, &sampleJSCode)
	}
}
