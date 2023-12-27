package cmd

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/ghdwlsgur/captain/utils"
// 	"github.com/spf13/cobra"
// )

// var DeleteCommand = &cobra.Command{
// 	Use:   "delete",
// 	Short: "delete robot",
// 	Long:  `Sub-command for Delete`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
// 			panicRed(err)
// 		}

// 		name := strings.TrimSpace(args[0])

// 		ctx := rootCmd.Context()
// 		token := ctx.Value(header("BasicAuth")).(string)

// 		err := utils.NewClient(utils.Get)
// 		if err != nil {
// 			err = fmt.Errorf("failed to create client: %w", err)
// 			panicRed(err)
// 		}
// 		// err := utils.NewClient(utils.Delete, )

// 	},
// }

// func init() {
// 	rootCmd.AddCommand(DeleteCommand)
// }
