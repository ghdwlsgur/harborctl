package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghdwlsgur/captain/internal"
	"github.com/ghdwlsgur/captain/utils"
	"github.com/jedib0t/go-pretty/table"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createRobotInputParams = &internal.CreateRobotInputParams{}
)

var CreateCommand = &cobra.Command{
	Use:   "create",
	Short: "create robot",
	Long:  `Sub-command for Create`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			panicRed(err)
		}

		name := strings.TrimSpace(args[0])
		duration := viper.GetInt64("duration")
		description := viper.GetString("description")

		ctx := rootCmd.Context()
		token := ctx.Value(header("BasicAuth")).(string)

		err := utils.NewClient(utils.Post)
		if err != nil {
			err = fmt.Errorf("failed to create client: %w", err)
			panicRed(err)
		}

		createRobotInputParams = internal.NewCreateRobotInputParams(
			name,        /* name */
			description, /* description */
			duration,    /* duration */
		)

		createRobotParams, err := createRobotInputParams.CreateRobotParams(ctx)
		if err != nil {
			err = fmt.Errorf("failed to create robot params: %w", err)
			panicRed(err)
		}

		robot, err := utils.NewRobotClient().CreateRobot(
			createRobotParams,                      /* params */
			utils.SetAuthorizationWithToken(token), /* authInfoWriter */
		)
		if err != nil {
			err = fmt.Errorf("failed to create robot: %w", err)
			panicRed(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		internal.OutputCreateRobotTable(t, robot)
		t.Render()
	},
}

func init() {
	CreateCommand.Flags().IntP("duration", "d", -1, "duration")
	CreateCommand.Flags().StringP("description", "e", "created by captain cli", "description")

	viper.BindPFlag("duration", CreateCommand.Flags().Lookup("duration"))
	viper.BindPFlag("description", CreateCommand.Flags().Lookup("description"))

	rootCmd.AddCommand(CreateCommand)
}
