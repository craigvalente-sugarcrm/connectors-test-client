package cmd

import (
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Connectors Test Client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !terminal.IsTerminal(unix.Stdout) {
			logrus.SetFormatter(&logrus.JSONFormatter{})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339Nano,
			})
		}

		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "make output more verbose")
}

// Execute starts the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
