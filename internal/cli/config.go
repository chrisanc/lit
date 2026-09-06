package cli

import (
	"CLI_App/internal/adapters/config"
	"CLI_App/internal/domain"

	"github.com/spf13/cobra"
)

func Configuration() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configure the scan variables.",
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
