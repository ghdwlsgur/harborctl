package internal

import (
	"context"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
)

type GetRobotInputParams struct {
	RobotID int64
}

type RobotGetter interface {
	GetRobotParams(ctx context.Context) *robot.GetRobotByIDParams
	GetRobotID() int64
}

var _ RobotGetter = (*GetRobotInputParams)(nil)

// NewGetRobotInputParams is a constructor for GetRobotInputParams
func NewGetRobotInputParams(robotID int64) *GetRobotInputParams {
	return &GetRobotInputParams{
		RobotID: robotID,
	}
}

// GetRobotID returns robotID
func (r GetRobotInputParams) GetRobotID() int64 {
	return r.RobotID
}

// GetRobotParams returns robot.GetRobotByIDParams
func (r *GetRobotInputParams) GetRobotParams(ctx context.Context) *robot.GetRobotByIDParams {
	getRobotParamsWithContext := robot.NewGetRobotByIDParamsWithContext(ctx)
	getRobotParamsWithContext.SetRobotID(r.GetRobotID())

	return getRobotParamsWithContext
}
