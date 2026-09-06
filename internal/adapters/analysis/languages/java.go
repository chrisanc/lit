package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type java struct {
	data           types.LanguageData
	ignoredSymbols []string
}

func NewJavaLanguage(varPattern, funcPattern string, ignoredSymbols []string) types.NodeManagement {
	j := &java{
		data: types.LanguageData{
			Language: tree.NewLanguage(javaGrammar.Language()),
		},
		ignoredSymbols: ignoredSymbols,
	}
	j.data.Queries = buildJavaQuery() + j.GetVarAppearancesQuery(varPattern) + j.GetFuncAppearancesQuery(funcPattern)
	return j
}

func (j java) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData, source []byte) {
	captureName := captureNames[node.Index]

	if captureName == "variable.name" || captureName == "function.name" {
		varName := node.Node.Utf8Text(source)
		if !domain.IsBuiltinSymbol("java", varName) && !domain.IsIgnoredSymbol(varName, j.ignoredSymbols) {
			nodeInfo.UpdateInvalidNames()
		}
		return
	}

	alternative := node.Node.ChildByFieldName("alternative")
	parent := node.Node.Parent()
	switch {
	case parent != nil && node.Node.GrammarName() == "binary_expression" && parent.GrammarName() == "variable_declarator":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func buildJavaQuery() string {
	return "(method_declaration type: (_) name: (_) @function.decl " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (block) @function.body ) @function " +
		"(constructor_declaration name: (_) @function.decl " +
		"parameters: (_) @function.parameters " +
		"body: (_) @function.body ) @function" +
		"(class_declaration name: (_) @model.name) @model" +
		"(interface_declaration name: (_) @model.name) @model" +
		"[" +
		"(for_statement) (while_statement) (do_statement) (enhanced_for_statement)" +
		"(if_statement condition: (_) consequence: (_) alternative: (_)?) (ternary_expression)" +
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"(switch_block_statement_group) (catch_clause)" +
		"] @keyword"
}

func (j java) GetLanguageData() types.LanguageData {
	return j.data
}

func (j java) GetVarAppearancesQuery(varPattern string) string {
	return fmt.Sprintf(" (variable_declarator name: (identifier) @variable.name (#not-match? @variable.name \"%s|%s\"))", varPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (formal_parameter name: (identifier) @variable.name (#not-match? @variable.name \"%s|%s\"))", varPattern, domain.AllowNonNamedVar)
}

func (j java) GetFuncAppearancesQuery(funcPattern string) string {
	return fmt.Sprintf(" (method_declaration name: (identifier) @function.name (#not-match? @function.name \"%s|%s\"))", funcPattern, domain.AllowNonNamedVar)
}
