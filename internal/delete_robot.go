package internal

import (
	"context"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
)

type DeleteRobotInputParams struct {
	RobotID int64
}

type RobotDeleter interface {
	DeleteRobotParams(ctx context.Context) (*robot.DeleteRobotParams, error)
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

func (r *DeleteRobotInputParams) DeleteRobotParams(ctx context.Context) (*robot.DeleteRobotParams, error) {
	deleteRobotParamsWithContext := robot.NewDeleteRobotParamsWithContext(ctx)
	deleteRobotParamsWithContext.SetRobotID(r.GetRobotID())

	return deleteRobotParamsWithContext, nil
}
