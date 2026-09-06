package analysis

import (
	"CLI_App/internal/adapters/analysis/types"
	"strings"
)

/*
 * writer.go - > file to write some data on a file using some ast queries.
 */

// FileModifier is an adapter for modifying data into the code script. This way, we can manage the file modification safely
type FileModifier struct {
	management            types.NodeManagement
	activePattern         string
	namingConventionIndex int8
}

func NewFileModifier(management types.NodeManagement, activePattern string, namingConventionIndex int8) FileModifier {
	return FileModifier{
		management:            management,
		activePattern:         activePattern,
		namingConventionIndex: namingConventionIndex,
	}
}

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
// Takes the initial variable name and converts it
// Only converts from one convention to another (safety conditions)
func (f FileModifier) ModifyVariableName(code *[]string) int {
	// get the query, cursor and captures (applying the query to fetch them)
	root := GetAST(code, f.management.GetLanguageData().Language)
	defer root.Close()
	query, cursor, captures := GetCapturesByQueries(f.management.GetLanguageData().Language,
		f.management.GetVarAppearancesQuery(f.activePattern), code, root.RootNode())
	defer query.Close()
	defer cursor.Close()

	localCache := make(map[uint]int)
	totalWrongNames := 0

	// loop through the node captures
	for {
		// get the next match
		match := captures.Next()
		if match == nil {
			break
		}
		copyOf := *match
		// get the node from the captures (just one capture per match)
		node := copyOf.Captures[0].Node
		rowIdx := node.StartPosition().Row
		if int(rowIdx) >= len(*code) {
			continue
		}

		value, ok := localCache[rowIdx]
		if !ok {
			value = 0
		}

		lineLen := len((*code)[rowIdx])
		startCol := int(node.StartPosition().Column) + value
		endCol := int(node.EndPosition().Column) + value

		if startCol < 0 || startCol >= lineLen || endCol < startCol || endCol > lineLen {
			continue
		}

		oldName := strings.Trim((*code)[rowIdx][startCol:endCol], "_")
		if oldName == "" {
			continue
		}
		newName := refactorVarName(GetTokens(oldName), f.namingConventionIndex)
		if oldName == newName || newName == "" {
			continue
		}

		totalWrongNames++
		row := (*code)[rowIdx]
		(*code)[rowIdx] = row[:startCol] + newName + row[endCol:]

		localCache[rowIdx] += len(newName) - int(node.EndPosition().Column-node.StartPosition().Column)
	}

	return totalWrongNames
}

// ---- Writing on files and renaming ----

// GetTokens : get the tokens of a variable name (without underscore or upper chars). Two pointers approach
func GetTokens(line string) []string {
	if len(line) == 0 {
		return nil
	}
	var tokens []string
	var i, j int

	for j < len(line) {
		if line[j] == 95 || (line[j] >= 65 && line[j] <= 90 && j > 0) {
			if i < j {
				tokens = append(tokens, strings.ToLower(line[i:j]))
			}
			for j < len(line) && line[j] == 95 {
				j++
			}
			i = j
			j++
		} else {
			j++
		}
	}

	if i < len(line) {
		tokens = append(tokens, strings.ToLower(line[i:]))
	}

	if len(tokens) > 0 && line[0] >= 65 && line[0] <= 90 {
		tokens[0] = strings.ToUpper(tokens[0][:1]) + tokens[0][1:]
	}

	return tokens
}

// refactorVarName: with the strings split in tokens, returns a []byte of the new line of code.
func refactorVarName(tokens []string, namingConventionIndex int8) string {
	if len(tokens) == 0 {
		return ""
	}
	var newName = tokens[0]

	switch namingConventionIndex {
	case 2:
		newName = ""
		camelCases(&newName, tokens)
	case 1, 3:
		if len(tokens) > 1 {
			camelCases(&newName, tokens[1:])
		}
	case 4:
		for _, token := range tokens[1:] {
			newName += "_" + token
		}
	}
	return newName
}

// function with the logics for the camelCase and CamelCase conversions
func camelCases(target *string, tokens []string) {
	if len(tokens) == 0 {
		return
	}
	for _, token := range tokens {
		if len(token) > 0 {
			*target += strings.ToUpper(token[:1]) + token[1:]
		}
	}
}
