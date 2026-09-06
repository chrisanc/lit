package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// python: struct with LanguageInformation embedded
type python struct {
	data types.LanguageData
}

func NewPythonLanguage(pattern string) types.NodeManagement {
	p := &python{
		data: types.LanguageData{
			Language: tree.NewLanguage(pyGrammar.Language()),
		},
	}
	p.data.Queries = buildPythonQuery() + p.GetVarAppearancesQuery(pattern)
	return p
}

// ManageNode - > Function to implement the NodeManagement interface
func (p python) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	switch {
	case captureNames[node.Index] == "variable.name":
		// Manage the variable whenever it's detected
		nodeInfo.UpdateInvalidNames()
		return
	case node.Node.GrammarName() == "boolean_operator" && node.Node.Parent().GrammarName() == "assignment":
		return
	}
	nodeInfo.Complexity++
}

func buildPythonQuery() string {
	return "(function_definition name: (identifier) @function.name " +
		"parameters: (parameters) @function.parameters " +
		"body: (block) @function.body) @function " +
		// Classes
		"(class_definition name: (_) @model.name ) @model" +
		// Keywords
		"[" +
		// If, else-if, else
		"(if_statement condition: (_)) (elif_clause condition: (_)) (else_clause body: (_))" +
		// Loops
		"(for_statement) (while_statement condition: (_) body: (_))" +
		// Operators
		"(boolean_operator left: (_) right: (_))" +
		// Clauses
		"(except_clause value: (_)) (conditional_expression) (case_clause (_))" +
		// List comprehension
		"(list_comprehension body: (_) (for_in_clause left: (_) right: (_))) (if_clause (_))" +
		"] @keyword"
}

func (p python) GetLanguageData() types.LanguageData {
	return p.data
}

func (p python) GetVarAppearancesQuery(pattern string) string {
	return fmt.Sprintf("((identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", pattern, domain.AllowNonNamedVar)
}
