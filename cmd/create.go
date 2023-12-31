package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghdwlsgur/harborctl/internal"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/jedib0t/go-pretty/table"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createRobotInputParams = &internal.CreateRobotInputParams{}
)

var CreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a robot account",
	Long:  `Creating a robot account in harbor`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			panicRed(err)
		}

		name := strings.TrimSpace(args[0])
		duration := viper.GetInt64("create-duration")
		description := viper.GetString("description")

		ctx := rootCmd.Context()
		token := ctx.Value(header("BasicAuth")).(string)

		createRobotInputParams = internal.NewCreateRobotInputParams(
			name,        /* name */
			description, /* description */
			duration,    /* duration */
		)

		createRobotParams, err := createRobotInputParams.CreateRobotParams(ctx)
		if err != nil {
			panicRed(fmt.Errorf("failed to create robot params: %w", err))
		}

		robotCreated, err := utils.NewRobotClient().CreateRobot(
			createRobotParams,                      /* params */
			utils.SetAuthorizationWithToken(token), /* authInfo */
		)
		if err != nil {
			panicRed(fmt.Errorf("failed to create robot - already exists account: %w ", err))
		}

		if robotCreated.IsSuccess() {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			internal.CreateRobotTableOutput(
				t,            /* writer */
				robotCreated, /* robot */
				description /* description */)
			t.Render()

			response, err := internal.SaveSecret(
				robotCreated, /* *robot.CreateRobotCreated */
				token,        /* token */
				description /* description */)
			if err != nil {
				panicRed(fmt.Errorf("failed to store secret: %w", err))
			}

			if response.Code == 200 {
				msg := fmt.Sprintf("Successfully created robot %s\n", robotCreated.GetPayload().Name)
				doneMsg(msg)
			}
		}
	},
}

func init() {
	host, err := os.Hostname()
	if err != nil {
		panicRed(fmt.Errorf("failed to get hostname: %w", err))
	}

	description := fmt.Sprintf("created by harborctl on %s", host)
	CreateCommand.Flags().Int64P("create-duration", "d", 30, "Setting an expiration period for the harbor robot account")
	CreateCommand.Flags().StringP("description", "e", description, "Writing a description for the harbor robot account")

	viper.BindPFlag("create-duration", CreateCommand.Flags().Lookup("create-duration"))
	viper.BindPFlag("description", CreateCommand.Flags().Lookup("description"))

	rootCmd.AddCommand(CreateCommand)
}
