package cli

import (
	"CLI_App/internal/adapters/analysis/languages"
	"CLI_App/internal/adapters/config"
	"CLI_App/internal/domain"
	"CLI_App/internal/service"

	"github.com/spf13/cobra"
)

func Files() *cobra.Command {
	command := &cobra.Command{
		Use:   "scan",
		Short: "Scan repository source files and analyze complexity, metrics, and naming conventions",
		Long:  "Scan traverses repository source files in parallel using Tree-Sitter AST parsing.\nIt evaluates cyclomatic complexity, method length, parameter thresholds, and variable/function naming conventions.\nReports can be output as text, SARIF, JSON, or Markdown.",
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := cmd.Flags().GetBool("loc")
			fix, _ := cmd.Flags().GetBool("fix")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			formatFlag, _ := cmd.Flags().GetString("format")
			outputFlag, _ := cmd.Flags().GetString("output")

			configAdapter := config.NewJSONAdapter()
			cfg := configAdapter.GetConfig()
			scanner := service.NewScannerService(
				languages.NewFileAnalyzer(
					cfg.GetVariableConvention(),
					cfg.GetFunctionConvention(),
					domain.NewFeedback(cfg),
					cfg.GetVariableConventionIndex(),
					cfg.GetFunctionConventionIndex(),
					cfg.IgnoredSymbols,
				),
			)

			switch {
			case loc:
				scanner.ExecuteLOC(cmd.Context())
				scanner.PrintLOCResults()
			case fix || dryRun:
				scanner.FixFile(cmd.Context(), dryRun)
				scanner.PrintFixResults(dryRun)
			default:
				scanner.ScanFiles(cmd.Context())
				if err := scanner.ExportResults(formatFlag, outputFlag); err != nil {
					cmd.PrintErrln("Export error:", err)
				}
			}
		},
	}
	command.Flags().Bool("loc", false, "Analyze repository composition and print lines of code statistics by language")
	command.Flags().Bool("fix", false, "Automatically refactor variable names that violate configured naming conventions")
	command.Flags().Bool("dry-run", false, "Preview variable naming refactorings as unified ANSI git diffs without modifying files on disk")
	command.Flags().StringP("format", "f", "text", "Specify output report format: text, sarif, json, markdown")
	command.Flags().StringP("output", "o", "", "Write report output to a specified file path (defaults to stdout)")

	return command
}
