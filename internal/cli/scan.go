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
		Short: "Scan the repository files (must include a .gitignore) and retrieves data from them",
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := cmd.Flags().GetBool("loc")
			fix, _ := cmd.Flags().GetBool("fix")

			configAdapter := config.NewJSONAdapter()
			cfg := configAdapter.GetConfig()
			scanner := service.NewScannerService(
				languages.NewFileAnalyzer(
					domain.Conventions[cfg.NamingConventionIndex-1],
					domain.NewFeedback(cfg),
					cfg.NamingConventionIndex,
				),
			)

			switch {
			case loc:
				scanner.ExecuteLOC()
				scanner.PrintLOCResults()
			case fix:
				scanner.FixFile()
				scanner.PrintFixResults()
			default:
				scanner.ScanFiles()
				scanner.PrintScanningResults()
			}
		},
	}
	command.Flags().Bool("loc", false, "Retrieves the languages used with statistics")
	command.Flags().Bool("fix", false, "Fixes up the variables with an invalid naming conventions."+
		"It only one convention to another\nExample: if you have variables snake_case and the active convention is camelCase, it's converted.")

	return command
}
