package analysis

import (
	"CLI_App/internal/adapters/analysis/types"
	"CLI_App/internal/domain"
	"fmt"
	"os"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
)

func GetAST(code *[]string, language *tree.Language) *tree.Tree {
	// Create a parser for the code
	parser := tree.NewParser()
	defer parser.Close()
	// Set the grammar of the language to parse
	err := parser.SetLanguage(language)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing the language...")
		os.Exit(1)
	}
	// Parse the source code with the configured parser
	treeParser := parser.Parse([]byte(strings.Join(*code, "\n")), nil)

	// Get the root (program node) node from the parser.
	return treeParser
}

func GetCapturesByQueries(language *tree.Language, queries string, code *[]string, root *tree.Node) (*tree.Query,
	*tree.QueryCursor, tree.QueryMatches) {
	// Create a query to extract the data we need
	query, err := tree.NewQuery(language, queries)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing", err)
		os.Exit(1)
	}

	// Create a query cursor for our custom query (just to keep the necessary data)
	cursor := tree.NewQueryCursor()

	// Execute the query and generate the captures
	return query, cursor, cursor.Matches(query, root, []byte(strings.Join(*code, "\n")))
}

// CyclicalComplexity Function that calculates the cyclical complexity of the code. Useful for the user feedback.
func CyclicalComplexity(languageInfo types.NodeManagement, code *[]string) []*domain.FunctionData {
	// Get our ast bases in our code and grammar
	ast := GetAST(code, languageInfo.GetLanguageData().Language)
	defer ast.Close()
	// Get the basics to iterate through the captures and keep the data
	query, cursor, captures := GetCapturesByQueries(languageInfo.GetLanguageData().Language, languageInfo.GetLanguageData().Queries, code, ast.RootNode())
	defer query.Close()
	defer cursor.Close()
	// List used as a stack to get the subfunctions and it's complexity right
	var Stack []*domain.FunctionData
	// Slice to save up the functions we take out the stack (with their final complexity)
	var Functions []*domain.FunctionData
	// Append a default main to the slice (JS, Python)
	Functions = append(Functions, &domain.FunctionData{Name: "Default Main", Complexity: 1, TotalParams: 0})

	// Get the Functions data
	for {
		// Iterator of all the captures of the source code.
		match := captures.Next()
		if match == nil {
			break
		}
		// Get a deep copy of the match (because the memory of it will be overwritten)
		copyOf := *match

		switch {
		// While we iterate through the captures, we save up the copies on a slice of objects.
		case query.CaptureNames()[copyOf.Captures[0].Index] == "function":
			Stack = append(Stack, &domain.FunctionData{Complexity: 1})
			// Add the initial data to the object reference in the stack
			Stack[len(Stack)-1].AddInitialData(
				"Method "+(*code)[copyOf.Captures[1].Node.StartPosition().Row][copyOf.Captures[1].Node.StartPosition().Column:copyOf.Captures[1].Node.EndPosition().Column],
				copyOf.Captures[2].Node.NamedChildCount(),
				copyOf.Captures[3].Node.StartByte(), copyOf.Captures[3].Node.EndByte(),
				copyOf.Captures[3].Node.EndPosition().Row-copyOf.Captures[3].Node.StartPosition().Row,
				domain.Point(copyOf.Captures[3].Node.StartPosition()),
			)

		case query.CaptureNames()[copyOf.Captures[0].Index] == "model":
			Stack = append(Stack, &domain.FunctionData{Complexity: 1})
			Stack[len(Stack)-1].AddInitialData(
				"Model "+(*code)[copyOf.Captures[1].Node.StartPosition().Row][copyOf.Captures[1].Node.StartPosition().Column:copyOf.Captures[1].Node.EndPosition().Column],
				0, copyOf.Captures[0].Node.StartByte(), copyOf.Captures[0].Node.EndByte(),
				copyOf.Captures[0].Node.EndPosition().Row-copyOf.Captures[0].Node.StartPosition().Row,
				domain.Point(copyOf.Captures[0].Node.StartPosition()),
			)

		// If there´s code without a function (JS, Python) before a function definition, we count it as main.
		case len(Stack) <= 0:
			languageInfo.ManageNode(query.CaptureNames(), copyOf.Captures[0], Functions[0])

		// Validate node ranges with the function body. This is the logic for functions inside functions or the main one
		case Stack[len(Stack)-1].IsTargetInRange(copyOf.Captures[0].Node.StartByte(), copyOf.Captures[0].Node.EndByte()):
			// + 1 in complexity in the function.
			languageInfo.ManageNode(query.CaptureNames(), copyOf.Captures[0], Stack[len(Stack)-1])

		// The code only gets here when there's a line out of the scope of a function
		// Remove the most recent element from the stack and add it to the Functions list
		default:
			// Flag to detect if the node is a main function node or not
			isMain := false
			// While the last element on the stack doesn't satisfy the range of the current node, we add it up
			// to the final Functions slice and remove it from the stack to keep going until it finds the right Function.
			for !(Stack[len(Stack)-1].IsTargetInRange(copyOf.Captures[0].Node.StartByte(), copyOf.Captures[0].Node.EndByte())) {
				// If the stack only has the current function, we assign the +1 to the main (default).
				if len(Stack) <= 1 {
					isMain = true
					languageInfo.ManageNode(query.CaptureNames(), copyOf.Captures[0], Functions[0])
					break
				}
				Functions = append(Functions, Stack[len(Stack)-1])
				Stack = Stack[:len(Stack)-1]
			}
			// At the end, in the verified node, we sum +1 to the complexity.
			if !isMain {
				languageInfo.ManageNode(query.CaptureNames(), copyOf.Captures[0], Stack[len(Stack)-1])
			}
		}
	}

	// If there's any function still on the stack, we copy it into the Functions slice.
	return append(Functions, Stack...)
}
