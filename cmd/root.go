package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ghdwlsgur/harborctl/internal"
	"github.com/spf13/cobra"
)

var (
	_user   = &internal.User{}
	rootCmd = &cobra.Command{
		Use:   "harborctl",
		Short: `it is a cli client tool for managing harbor robot accounts.`,
		Long: `it is a cli client tool for managing harbor robot accounts.
designed for Querypie PE team members, it primarily creates, deletes, 
searches, and updates robot accounts on harbor.chequer.io.
		`,
		Example: `  harborctl create <robot-name> -d <duration> -e <description>
  harborctl delete <robot-name>
  harborctl search
  harborctl update <robot-name> -d <duration>`,
	}

	doneMsg = color.New(color.Bold, color.FgHiGreen).PrintFunc()
)

func panicRed(err error) {
	fmt.Println(color.RedString("[err] %s", err.Error()))
	os.Exit(1)
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		panicRed(fmt.Errorf("failed to execute rootCmd: %w", err))
	}
}

type header string

func initConfig() {
	BasicAuth, err := authentication(_user)
	if err != nil {
		panicRed(err)
	}

	ctx := context.WithValue(context.Background(), header("BasicAuth"), BasicAuth)
	rootCmd.SetContext(ctx)
}

func authentication(u *internal.User) (string, error) {
	if u.Verify() {
		c, err := u.Parsing()
		if err != nil {
			return "", fmt.Errorf("failed to parsing credentials file - authentication: %w", err)
		}

		return c.GetBasicAuth(), nil
	}

	if err := u.Login(); err != nil {
		return "", fmt.Errorf("failed to login - authentication: %w", err)
	}

	return u.GetBasicAuth(), nil
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.InitDefaultVersionFlag()
}
