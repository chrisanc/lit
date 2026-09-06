package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// javascript interface with LanguageData embedded
type javascript struct {
	data types.LanguageData
}

func NewJSLanguage(pattern string) types.NodeManagement {
	js := &javascript{
		data: types.LanguageData{
			Language: tree.NewLanguage(jsGrammar.Language()),
		},
	}
	js.data.Queries = buildJSQuery() + js.GetVarAppearancesQuery(pattern)
	return js
}

func (js javascript) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.UpdateInvalidNames()
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator":
		return
	}
	nodeInfo.Complexity++
}

func buildJSQuery() string {
	return "(function_declaration name: (identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ) @function" +
		// Arrow Function
		"(variable_declarator name: (identifier) @function.name " +
		"value: (arrow_function parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body )) @function" +
		// Classes
		"(class_declaration name: (_) @model.name ) @model" +
		// Class Methods
		"(method_definition name: (property_identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ) @function" +
		// Functions body information (keywords)
		"[" +
		// if, else-if, else
		"(if_statement condition: (_) consequence: (_) alternative: (else_clause)?)" +
		"(else_clause (statement_block))" +
		// Statements
		"(for_statement) (for_in_statement) (while_statement) (switch_case) (catch_clause)" +
		// Expressions
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"(ternary_expression)" +
		// ForEach, Map, etc.
		"(call_expression function: (member_expression object: (_) property: (_) @call.name)" +
		"arguments: (arguments (arrow_function)) (#match? @call.name \"^(forEach)$\"))" +
		"] @keyword"
}

func (js javascript) GetLanguageData() types.LanguageData {
	return js.data
}

func (js javascript) GetVarAppearancesQuery(pattern string) string {
	return fmt.Sprintf("([(identifier) @variable.name (property_identifier) @variable.name (array_pattern (identifier) @variable.name)]") +
		fmt.Sprintf("(#not-match? @variable.name \"^%s|%s$\"))", pattern, domain.AllowNonNamedVar) +
		fmt.Sprintf("(variable_declarator name: (object_pattern (shorthand_property_identifier_pattern) @variable.name)") +
		fmt.Sprintf("(#not-match? @variable.name \"%s|%s\"))", pattern, domain.AllowNonNamedVar)
}
