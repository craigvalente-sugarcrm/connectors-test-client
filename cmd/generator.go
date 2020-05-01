package cmd

import (
	"github.com/connectors-test-client/app/common"
	"github.com/connectors-test-client/app/google"
	"github.com/connectors-test-client/app/microsoft"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "gen",
		Short: "Generates test file",
		RunE: func(cmd *cobra.Command, args []string) error {
			streamType := viper.GetString("type")
			settings := &common.GeneratorSettings{
				Account:    viper.GetString("acct"),
				StreamType: streamType,
				TransCount: viper.GetInt("count"),
			}
			switch streamType {
			case "MS":
				microsoft.GenerateTest(settings)
				break
			case "GOOG":
				google.GenerateTest(settings)
				break
			default:
				return errors.New("Invalid stream type")
			}
			return nil
		},
	}

	cmd.Flags().StringP("acct", "a", "all", "user account")
	viper.BindPFlag("acct", cmd.Flags().Lookup("acct"))

	cmd.Flags().StringP("type", "t", "", "stream type")
	viper.BindPFlag("type", cmd.Flags().Lookup("type"))

	cmd.Flags().IntP("count", "c", 0, "action count")
	viper.BindPFlag("count", cmd.Flags().Lookup("count"))

	rootCmd.AddCommand(cmd)
}
