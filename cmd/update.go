package cmd

import (
	"fmt"
	"os"

	"github.com/ghdwlsgur/captain/internal"
	"github.com/ghdwlsgur/captain/utils"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "update robot",
	Long:  `Sub-command for Update`,
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
			err = fmt.Errorf("failed to create list robot input params: %w", err)
			panicRed(err)
		}

		var (
			robotTable = make(map[string]*internal.Robot)
		)

		for i := 0; i < ListRobotInputParams.GetMaxPage(); i++ {
			success, robotPayload, err := ListRobotInputParams.Payload(ctx)
			if err != nil {
				err = fmt.Errorf("failed to get payload: %w", err)
				panicRed(err)
			}

			if success {
				for _, v := range robotPayload {
					if v.Description == "" {
						v.Description = "undefined"
					}

					if v.Duration >= 0 {
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
				break
			}
		}

		robots := make([]string, 0, len(robotTable))
		for robotListed := range robotTable {
			robots = append(robots, robotListed)
		}

		answer, err := utils.AskPromptOptionList("Please select the robot you want to view in detail", robots, 10)
		if err != nil {
			err = fmt.Errorf("failed to ask prompt option list: %w", err)
			panicRed(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		internal.ListRobotTableOutput(t, robotTable[answer])
		t.Render()

		beforeDuration := robotTable[answer].Duration
		afterDuration := viper.GetInt64("duration")

		if robotTable[answer].Editable {

			updateRobotInputParams := internal.NewUpdateRobotInputParams(
				robotTable[answer].ID,          /* robotID */
				int64(afterDuration),           /* duration */
				robotTable[answer].Name,        /* name */
				robotTable[answer].Description, /* description */
			)

			updateRobotParams := updateRobotInputParams.UpdateRobotParams(ctx)
			robotUpdated, err := utils.NewRobotClient().UpdateRobot(
				updateRobotParams,                      /* params */
				utils.SetAuthorizationWithToken(token), /* authInfoWriter */
			)
			if err != nil {
				err = fmt.Errorf("failed to update robot: %w", err)
				panicRed(err)
			}

			if robotUpdated.IsSuccess() {
				msg := fmt.Sprintf("Successfully updated robot duration %d to %d\n", beforeDuration, afterDuration)
				doneMsg(msg)
			}
		}

	},
}

func init() {
	UpdateCommand.Flags().Int64P("duration", "d", 30, "duration")

	viper.BindPFlag("duration", UpdateCommand.Flags().Lookup("duration"))

	rootCmd.AddCommand(UpdateCommand)
}
