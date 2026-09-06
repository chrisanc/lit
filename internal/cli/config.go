package cli

import (
	"CLI_App/internal/adapters/config"
	"CLI_App/internal/domain"

	"github.com/spf13/cobra"
)

func Configuration() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Interactively configure variable/function naming conventions and threshold alerts",
		Long:  "Config launches an interactive terminal prompt to set naming conventions for variables and functions independently,\nas well as thresholds for parameter count, cyclomatic complexity, and method size in config.json.",
		Run: func(cmd *cobra.Command, args []string) {
			varIdx := GetVariableNamingConvention()
			funcIdx := GetFunctionNamingConvention()
			alerts := GetAlertsConfig()
			jsonAdapter := config.NewJSONAdapter()
			newConfig := &domain.Config{
				NamingConventionIndex:         varIdx,
				VariableNamingConventionIndex: varIdx,
				FunctionNamingConventionIndex: funcIdx,
				Alerts:                        alerts,
			}
			jsonAdapter.SaveConfig(newConfig)
		},
	}
}
