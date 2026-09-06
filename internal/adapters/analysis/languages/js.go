package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

type javascript struct {
	data types.LanguageData
}

func NewJSLanguage(varPattern, funcPattern string) types.NodeManagement {
	js := &javascript{
		data: types.LanguageData{
			Language: tree.NewLanguage(jsGrammar.Language()),
		},
	}
	js.data.Queries = buildJSQuery() + js.GetVarAppearancesQuery(varPattern) + js.GetFuncAppearancesQuery(funcPattern)
	return js
}

func (js javascript) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData, source []byte) {
	captureName := captureNames[node.Index]
	if captureName == "variable.name" || captureName == "function.name" {
		varName := node.Node.Utf8Text(source)
		if !domain.IsBuiltinSymbol("js", varName) {
			nodeInfo.UpdateInvalidNames()
		}
		return
	}
	if node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator" {
		return
	}
	nodeInfo.Complexity++
}

func buildJSQuery() string {
	return "(function_declaration name: (identifier) @function.decl " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ) @function" +
		"(variable_declarator name: (identifier) @function.decl " +
		"value: (arrow_function parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body )) @function" +
		"(class_declaration name: (_) @model.name ) @model" +
		"(method_definition name: (property_identifier) @function.decl " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ) @function" +
		"[" +
		"(if_statement condition: (_) consequence: (_) alternative: (else_clause)?)" +
		"(else_clause (statement_block))" +
		"(for_statement) (for_in_statement) (while_statement) (switch_case) (catch_clause)" +
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"(ternary_expression)" +
		"(call_expression function: (member_expression object: (_) property: (_) @call.name)" +
		"arguments: (arguments (arrow_function)) (#match? @call.name \"^(forEach)$\"))" +
		"] @keyword"
}

func (js javascript) GetLanguageData() types.LanguageData {
	return js.data
}

func (js javascript) GetVarAppearancesQuery(varPattern string) string {
	return fmt.Sprintf(" (variable_declarator name: (identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", varPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (formal_parameters (identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", varPattern, domain.AllowNonNamedVar)
}

func (js javascript) GetFuncAppearancesQuery(funcPattern string) string {
	return fmt.Sprintf(" (function_declaration name: (identifier) @function.name (#not-match? @function.name \"^%s|%s$\"))", funcPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (method_definition name: (property_identifier) @function.name (#not-match? @function.name \"^%s|%s$\"))", funcPattern, domain.AllowNonNamedVar)
}
