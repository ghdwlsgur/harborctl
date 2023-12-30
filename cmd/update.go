package cmd

import (
	"fmt"
	"os"

	"github.com/ghdwlsgur/harborctl/internal"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	getRobotInputParams = &internal.GetRobotInputParams{}
)

var UpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Update a robot account's duration",
	Long:  `Updating a robot account's duration in harbor, Only accounts with a duration of 0 or greater are eligible for updates`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cobra.MaximumNArgs(0)(cmd, args); err != nil {
			panicRed(err)
		}

		ctx := rootCmd.Context()
		token := ctx.Value(header("BasicAuth")).(string)

		ListRobotInputParams, err := internal.NewListRobotInputParams(
			ctx,   /* ctx */
			token, /* token */
		)
		if err != nil {
			panicRed(fmt.Errorf("failed to create robot input params: %w", err))
		}

		var (
			robotTable = make(map[string]*internal.Robot)
		)

		spinner := utils.StartSpinner("Retrieving a list of robot accounts with remaining expiration period ")
		for i := 0; i < ListRobotInputParams.GetMaxPage(); i++ {
			success, robotPayload, err := ListRobotInputParams.Payload(ctx)
			if err != nil {
				panicRed(fmt.Errorf("failed to get payload: %w", err))
			}

			if success {
				for _, v := range robotPayload {
					if v.Description == "" {
						v.Description = "undefined"
					}

					Dday := utils.CountDays(v.ExpiresAt).LeftDays
					if utils.Expiration(Dday) != utils.Expired {
						k := fmt.Sprintf("%s [%s]", v.Description, v.Name)
						robotTable[k] = &internal.Robot{
							ID:           v.ID,
							Name:         v.Name,
							Description:  v.Description,
							CreationTime: v.CreationTime.String(),
							ExpiredTime:  v.ExpiresAt,
							Dday:         utils.CountDays(v.ExpiresAt).Validate(),
							Duration:     v.Duration,
							Editable:     v.Editable,
						}
					}
				}
				ListRobotInputParams.NextPage()
			} else {
				panicRed(fmt.Errorf("failed to get payload: %w", err))
			}
		}
		utils.StopSpinner(spinner)

		robots := make([]string, 0, len(robotTable))
		for robotListed := range robotTable {
			robots = append(robots, robotListed)
		}

		msg := "Please select the robot you want to view in detail"
		answer, err := utils.AskPromptOptionList(
			msg,    /* message */
			robots, /* options */
			10 /* size */)
		if err != nil {
			panicRed(fmt.Errorf("failed to ask prompt, option list: %w", err))
		}

		before := table.NewWriter()
		before.SetTitle("Before")
		before.SetOutputMirror(os.Stdout)
		internal.ListRobotTableOutput(before, robotTable[answer])
		before.Render()

		beforeDuration := robotTable[answer].Duration
		afterDuration := viper.GetInt64("update-duration")

		// editable check and update
		if robotTable[answer].Editable {
			updateRobotInputParams := internal.NewUpdateRobotInputParams(
				robotTable[answer].ID,          /* robotID */
				int64(afterDuration),           /* duration */
				robotTable[answer].Name,        /* name */
				robotTable[answer].Description, /* description */
			)

			updateRobotParams, err := updateRobotInputParams.UpdateRobotParams(ctx)
			if err != nil {
				panicRed(fmt.Errorf("failed to create updating robot params: %w", err))
			}
			robotUpdated, err := utils.NewRobotClient().UpdateRobot(
				updateRobotParams,                      /* params */
				utils.SetAuthorizationWithToken(token), /* authInfo */
			)
			if err != nil {
				panicRed(fmt.Errorf("failed to update robot: %w", err))
			}

			if robotUpdated.IsSuccess() {
				getRobotInputParams = internal.NewGetRobotInputParams(robotTable[answer].ID)
				getRobotParams := getRobotInputParams.GetRobotParams(ctx)

				newRobot, err := utils.NewRobotClient().GetRobotByID(
					getRobotParams, /* params */
					utils.SetAuthorizationWithToken(token) /* authInfo */)
				if err != nil {
					panicRed(fmt.Errorf("failed to get robot by id: %w", err))
				}

				after := table.NewWriter()
				after.SetOutputMirror(os.Stdout)
				internal.UpdateRobotTableOutput(after, newRobot.GetPayload())
				after.Render()

				msg := fmt.Sprintf("Successfully updated robot duration %d to %d\n", beforeDuration, afterDuration)
				doneMsg(msg)
			}
		}

	},
}

func init() {
	UpdateCommand.Flags().Int64P("update-duration", "d", 30, "Updating an expiration period for the harbor robot account")

	viper.BindPFlag("update-duration", UpdateCommand.Flags().Lookup("update-duration"))

	rootCmd.AddCommand(UpdateCommand)
}
