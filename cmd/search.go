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
	ListRobotInputParams = &internal.ListRobotInputParams{}
)

var SearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search robot accounts",
	Long:  `Searching robot accounts in harbor`,
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
			panicRed(fmt.Errorf("failed to create list robot input params: %w", err))
		}

		var (
			robotTable = make(map[string]*internal.Robot)
		)

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
				panicRed(fmt.Errorf("failed to get payload: %w", err))
			}
		}

		robots := make([]string, 0, len(robotTable))
		for robotListed := range robotTable {
			robots = append(robots, robotListed)
		}

		const msg = "Please select the robot you want to view in detail"
		answer, err := utils.AskPromptOptionList(
			msg,    /* message */
			robots, /* options */
			10 /* size */)
		if err != nil {
			panicRed(fmt.Errorf("failed to ask prompt option list: %w", err))
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
