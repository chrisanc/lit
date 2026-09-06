package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

type python struct {
	data types.LanguageData
}

func NewPythonLanguage(varPattern, funcPattern string) types.NodeManagement {
	p := &python{
		data: types.LanguageData{
			Language: tree.NewLanguage(pyGrammar.Language()),
		},
	}
	p.data.Queries = buildPythonQuery() + p.GetVarAppearancesQuery(varPattern) + p.GetFuncAppearancesQuery(funcPattern)
	return p
}

func (p python) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData, source []byte) {
	captureName := captureNames[node.Index]
	if captureName == "variable.name" || captureName == "function.name" {
		varName := node.Node.Utf8Text(source)
		if !domain.IsBuiltinSymbol("py", varName) {
			nodeInfo.UpdateInvalidNames()
		}
		return
	}
	if node.Node.GrammarName() == "boolean_operator" && node.Node.Parent().GrammarName() == "assignment" {
		return
	}
	nodeInfo.Complexity++
}

func buildPythonQuery() string {
	return "(function_definition name: (identifier) @function.decl " +
		"parameters: (parameters) @function.parameters " +
		"body: (block) @function.body) @function " +
		"(class_definition name: (_) @model.name ) @model " +
		"[" +
		"(if_statement condition: (_)) (elif_clause condition: (_)) (else_clause body: (_))" +
		"(for_statement) (while_statement condition: (_) body: (_))" +
		"(boolean_operator left: (_) right: (_))" +
		"(except_clause value: (_)) (conditional_expression) (case_clause (_))" +
		"(list_comprehension body: (_) (for_in_clause left: (_) right: (_))) (if_clause (_))" +
		"] @keyword"
}

func (p python) GetLanguageData() types.LanguageData {
	return p.data
}

func (p python) GetVarAppearancesQuery(varPattern string) string {
	return fmt.Sprintf(" (assignment left: (identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", varPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (parameters (identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", varPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (for_statement left: (identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", varPattern, domain.AllowNonNamedVar)
}

func (p python) GetFuncAppearancesQuery(funcPattern string) string {
	return fmt.Sprintf(" (function_definition name: (identifier) @function.name (#not-match? @function.name \"^%s|%s$\"))", funcPattern, domain.AllowNonNamedVar)
}
