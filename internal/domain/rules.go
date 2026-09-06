package domain

import "regexp"

// patterns.go - > File that contains every regex pattern used in the script

// ScanValidScriptPattern for the regex that validates scripts for the scanner.
// Only contains the languages supported for scanning
var ScanValidScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|js|java|jsx|go)$"

// LocValidScriptPattern validates the scripts for the loc.
var LocValidScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|go|java|js|jsx|dart|c|cpp|css|html|ts|md)$"

// NotValidDirPattern for the regex that validates you won´t visit unwanted sites.
// Only contains unwanted directories and files.
var NotValidDirPattern = "^(node_modules|.*\\.exe|target|venv|__pycache__|" +
	"\\.(git|idea|mvn|cmd))$"

// Variable names conventions. For good practices and consistency on the variable names of the project

// Constant - > Defines the regex for the constant variables.
var Constant = "^[A-Z][A-Z0-9]*([_]{1}[A-Z0-9]+)*$"

// LowerCamelCase - > variables with the style 'camelCase'
var LowerCamelCase = "^[a-z]{1}([A-Z0-9]*[a-z0-9]+)*$" + "|" + Constant

// UpperCamelCase - > variables with the style 'CamelCase'
var UpperCamelCase = "^[A-Z]{1}([A-Z0-9]*[a-z0-9]+)*$" + "|" + Constant

// CamelCase - > Variables that can be either 'camelCase' and 'CamelCase' (default)
var CamelCase = "^[A-Za-z]{1}([A-Z0-9]*[a-z0-9]+)*$" + "|" + Constant

// SnakeCase - > Variables with the style 'sneak_case'
var SnakeCase = "^[a-z][a-z0-9]*([_]{1}[a-z0-9]+)*$" + "|" + Constant

// AllowNonNamedVar - > This variable is used interpolated with any of the other ones (used to allow _ vars).
// Required in languages like Go or Python.
// Example: pattern := AllowNonNamedVar + CamelCase
var AllowNonNamedVar = "^[_]{1}$"

var Conventions = []string{LowerCamelCase, UpperCamelCase, CamelCase, SnakeCase}

func RegexMatch(pattern, target string) bool {
	matched, _ := regexp.MatchString(pattern, target)
	return matched
}
