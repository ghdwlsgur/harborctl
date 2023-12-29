package internal

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/jedib0t/go-pretty/table"
)

type DeleteRobotInputParams struct {
	RobotID int64
}

type RobotDeleter interface {
	DeleteRobotParams(ctx context.Context) *robot.DeleteRobotParams
	GetRobotID() int64
}

var _ RobotDeleter = (*DeleteRobotInputParams)(nil)

func NewDeleteRobotInputParams(robotID int64) *DeleteRobotInputParams {
	return &DeleteRobotInputParams{
		RobotID: robotID,
	}
}

func (r DeleteRobotInputParams) GetRobotID() int64 {
	return r.RobotID
}

func (r *DeleteRobotInputParams) DeleteRobotParams(ctx context.Context) *robot.DeleteRobotParams {
	deleteRobotParamsWithContext := robot.NewDeleteRobotParamsWithContext(ctx)
	deleteRobotParamsWithContext.SetRobotID(r.GetRobotID())

	return deleteRobotParamsWithContext
}

func DeleteRobotTableOutput(
	writer table.Writer,
	robotTable *Robot) (table.Writer, error) {

	creationTime, err := utils.CreationTimeFormatKST(robotTable.CreationTime)
	if err != nil {
		err = fmt.Errorf("failed to parse time - CreationTimeFormatKST: %w", err)
		return nil, err
	}

	writer.SetTitle(color.HiRedString("Once you delete this, it cannot be undone"))
	writer.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Creation_Time",
		"Expired_Time",
		"D_Day",
		"Duration",
	})
	writer.AppendRow(table.Row{
		robotTable.ID,
		robotTable.Name,
		robotTable.Description,
		creationTime,
		utils.ExpiresAtToStringTime(robotTable.ExpiredTime),
		utils.CountDays(robotTable.ExpiredTime).Validate(),
		robotTable.Duration,
	})

	return writer, nil
}
