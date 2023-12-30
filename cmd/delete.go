package cmd

import (
	"fmt"
	"os"

	"github.com/ghdwlsgur/harborctl/internal"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var (
	deleteRobotInputParams = &internal.DeleteRobotInputParams{}
)

var DeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a robot account",
	Long:  `Delete a robot account in harbor`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cobra.MaximumNArgs(0)(cmd, args); err != nil {
			panicRed(err)
		}

		ctx := rootCmd.Context()
		token := ctx.Value(header("BasicAuth")).(string)

		ListRobotInputParams, err := internal.NewListRobotInputParams(
			rootCmd.Context(), /* ctx */
			token,             /* token */
		)
		if err != nil {
			panicRed(fmt.Errorf("failed to create listing robot input params: %w", err))
		}

		var (
			robotTable = make(map[string]*internal.Robot)
		)

		spinner := utils.StartSpinner("Retrieving all robot accounts ")
		for i := 0; i < ListRobotInputParams.GetMaxPage(); i++ {
			success, robotPayload, err := ListRobotInputParams.Payload(ctx)
			if err != nil {
				panicRed(fmt.Errorf("failed to get robot payload: %w", err))
			}

			if success {
				for _, v := range robotPayload {
					if v.Description == "" {
						v.Description = "undefined"
					}

					k := fmt.Sprintf("%s [%s]", v.Description, v.Name)
					robotTable[k] = &internal.Robot{
						ID:           v.ID,                                    /* ID */
						Name:         v.Name,                                  /* Name */
						Description:  v.Description,                           /* Description */
						CreationTime: v.CreationTime.String(),                 /* CreationTime */
						ExpiredTime:  v.ExpiresAt,                             /* ExpiredTime */
						Dday:         utils.CountDays(v.ExpiresAt).Validate(), /* Dday */
						Duration:     v.Duration,                              /* Duration */
					}
				}
				ListRobotInputParams.NextPage()
			} else {
				panicRed(fmt.Errorf("failed to get robot payload, success is false: %w", err))
			}
		}
		utils.StopSpinner(spinner)

		robots := make([]string, 0, len(robotTable))
		for robotListed := range robotTable {
			robots = append(robots, robotListed)
		}

		msg := "Please select the robot you want to delete"
		answer, err := utils.AskPromptOptionList(
			msg,    /* message */
			robots, /* options */
			10 /* default */)
		if err != nil {
			panicRed(fmt.Errorf("failed to ask prompt, option list: %w", err))
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		internal.DeleteRobotTableOutput(t, robotTable[answer])
		t.Render()

		msg = "Are you sure you want to delete this robot"
		deleteAnswer, err := utils.AskYesOrNo(msg)
		if err != nil {
			panicRed(fmt.Errorf("failed to ask prompt, yes or no: %w", err))
		}

		if deleteAnswer == "Yes" {
			deleteRobotInputParams = internal.NewDeleteRobotInputParams(robotTable[answer].ID)
			deleteRobotParams := deleteRobotInputParams.DeleteRobotParams(ctx)

			robot, err := utils.NewRobotClient().DeleteRobot(
				deleteRobotParams, /* params */
				utils.SetAuthorizationWithToken(token) /* authInfo */)
			if err != nil {
				panicRed(fmt.Errorf("failed to delete robot: %w", err))
			}

			if robot.IsSuccess() {
				msg := fmt.Sprintf("Successfully deleted robot %s\n", robotTable[answer].Description)
				doneMsg(msg)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(DeleteCommand)
}
