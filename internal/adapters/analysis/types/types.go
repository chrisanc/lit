package types

import (
	"CLI_App/internal/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
)

// LanguageData registers the language and its compiled query definitions.
type LanguageData struct {
	Language *tree.Language
	Queries  string
}

// NodeManagement defines the methods every language implementation uses.
type NodeManagement interface {
	ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData, source []byte)
	GetLanguageData() LanguageData
	GetVarAppearancesQuery(varPattern string) string
	GetFuncAppearancesQuery(funcPattern string) string
}
