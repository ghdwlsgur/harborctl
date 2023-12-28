package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configCommand = &cobra.Command{
		Use:       "config",
		Short:     "Configure harborctl",
		Long:      `Configure harborctl`,
		ValidArgs: []string{"login"},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err error
			)

			switch args[0] {
			case "login":
				if _user.Login(); err != nil {
					panicRed(err)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(configCommand)
}
