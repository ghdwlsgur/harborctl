package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	loginCommand = &cobra.Command{
		Use:   "login",
		Short: "logging on harbor",
		Long:  `logging on harbor`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cobra.MaximumNArgs(0)(cmd, args); err != nil {
				panicRed(err)
			}

			_, err := Authentication(_user)
			if err != nil {
				panicRed(fmt.Errorf("failed to login: %v", err))
			}

			msg := fmt.Sprintf("login success: %s\n", _user.GetUsername())
			doneMsg(msg)
		},
	}
)

func init() {
	rootCmd.AddCommand(loginCommand)
}
