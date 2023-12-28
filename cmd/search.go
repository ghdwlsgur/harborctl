package cmd

import (
	"fmt"
	"os"

	"github.com/ghdwlsgur/captain/internal"
	"github.com/ghdwlsgur/captain/utils"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

var (
	ListRobotInputParams = &internal.ListRobotInputParams{}
)

var SearchCommand = &cobra.Command{
	Use:   "search",
	Short: "search robot",
	Long:  `Sub-command for Search`,
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

					k := fmt.Sprintf("%s [%s]", v.Description, v.Name)
					robotTable[k] = &internal.Robot{
						ID:           v.ID,
						Name:         v.Name,
						Description:  v.Description,
						CreationTime: v.CreationTime.String(),
						ExpiredTime:  v.ExpiresAt,
						Dday:         utils.CountDays(v.ExpiresAt).Validate(),
						Duration:     v.Duration,
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
	},
}

func init() {
	rootCmd.AddCommand(SearchCommand)
}
