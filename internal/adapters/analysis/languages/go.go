package languages

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type golang struct {
	data           types.LanguageData
	ignoredSymbols []string
}

func NewGolangLanguage(varPattern, funcPattern string, ignoredSymbols []string) types.NodeManagement {
	g := &golang{
		data: types.LanguageData{
			Language: tree.NewLanguage(goGrammar.Language()),
		},
		ignoredSymbols: ignoredSymbols,
	}
	g.data.Queries = buildGolangQuery() + g.GetVarAppearancesQuery(varPattern) + g.GetFuncAppearancesQuery(funcPattern)
	return g
}

func (g golang) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData, source []byte) {
	captureName := captureNames[node.Index]

	if captureName == "variable.name" || captureName == "function.name" {
		varName := node.Node.Utf8Text(source)
		if !domain.IsBuiltinSymbol("go", varName) && !domain.IsIgnoredSymbol(varName, g.ignoredSymbols) {
			nodeInfo.UpdateInvalidNames()
		}
		return
	}

	alternative := node.Node.ChildByFieldName("alternative")
	parent := node.Node.Parent()
	switch {
	case parent != nil && node.Node.GrammarName() == "binary_expression" && parent.GrammarName() == "expression_list":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func buildGolangQuery() string {
	return "(function_declaration name: (_) @function.decl " +
		"parameters: (_) @function.parameters " +
		"body: (_) @function.body ) @function " +
		"(method_declaration name: (_) @function.decl " +
		"parameters: (_) @function.parameters " +
		"body: (_) @function.body ) @function" +
		"(type_declaration (type_spec name: (_) @model.name " +
		"type: ([(struct_type) (interface_type)]))) @model" +
		"[" +
		"(if_statement) (for_statement) (expression_case)" +
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"] @keyword"
}

func (g golang) GetLanguageData() types.LanguageData {
	return g.data
}

func (g golang) GetVarAppearancesQuery(varPattern string) string {
	return fmt.Sprintf(" (short_var_declaration left: (expression_list (identifier) @variable.name) (#not-match? @variable.name \"%s|%s\"))", varPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (var_spec name: (identifier) @variable.name (#not-match? @variable.name \"%s|%s\"))", varPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (parameter_declaration name: (identifier) @variable.name (#not-match? @variable.name \"%s|%s\"))", varPattern, domain.AllowNonNamedVar)
}

func (g golang) GetFuncAppearancesQuery(funcPattern string) string {
	return fmt.Sprintf(" (function_declaration name: (identifier) @function.name (#not-match? @function.name \"%s|%s\"))", funcPattern, domain.AllowNonNamedVar) +
		fmt.Sprintf(" (method_declaration name: (field_identifier) @function.name (#not-match? @function.name \"%s|%s\"))", funcPattern, domain.AllowNonNamedVar)
}
