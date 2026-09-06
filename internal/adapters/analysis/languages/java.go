package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type java struct {
	data types.LanguageData
}

func NewJavaLanguage(pattern string) types.NodeManagement {
	j := &java{
		data: types.LanguageData{
			Language: tree.NewLanguage(javaGrammar.Language()),
		},
	}
	j.data.Queries = buildJavaQuery() + j.GetVarAppearancesQuery(pattern)
	return j
}

func (j java) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.UpdateInvalidNames()
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func buildJavaQuery() string {
	return "(method_declaration type: (_) name: (_) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (block) @function.body ) @function " +
		"(constructor_declaration name: (_) @function.name " +
		"parameters: (_) @function.parameters " +
		"body: (_) @function.body ) @function" +
		// Classes, interfaces
		"(class_declaration name: (_) @model.name) @model" +
		"(interface_declaration name: (_) @model.name) @model" +
		// Keywords (+1 complexity)
		"[" +
		// Loops
		"(for_statement) (while_statement) (do_statement) (enhanced_for_statement)" +
		// If, else-if, else
		"(if_statement condition: (_) consequence: (_) alternative: (_)?) (ternary_expression)" +
		// Expressions
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"(switch_block_statement_group) (catch_clause)" +
		"] @keyword"
}

func (j java) GetLanguageData() types.LanguageData {
	return j.data
}

func (j java) GetVarAppearancesQuery(pattern string) string {
	return fmt.Sprintf("((identifier) @variable.name (#not-match? @variable.name \"^%s|%s$\"))", pattern, domain.AllowNonNamedVar)
}
