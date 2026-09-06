package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// root: commands entry point. Every command is a subcommand of root
var root = &cobra.Command{
	Use:     "lit",
	Short:   "High-performance AST code analysis and linting CLI for Git repositories",
	Long:    "Lit is a high-performance CLI tool written in Go that analyzes source code across multiple languages using Tree-Sitter AST parsing.\nIt evaluates cyclomatic complexity, method length, parameter counts, and naming conventions, providing colorized diff previews and SARIF/JSON/Markdown reporting.",
	Version: "2.0.0",
}

// Execute function to execute some code
func Execute() {
	// Add commands to the root
	root.AddCommand(Files())
	root.AddCommand(Configuration())

	// Execute the root, registering all the children commands
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// End the execution
	os.Exit(0)
}
