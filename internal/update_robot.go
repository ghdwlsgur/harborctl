package internal

import (
	"context"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/models"
)

type UpdateRobotInputParams struct {
	RobotID     int64
	Duration    int64
	Name        string
	Description string
}

type RobotUpdater interface {
	GetRobotID() int64
	GetDuration() int64
	GetName() string
	GetDescription() string
	UpdateRobotParams(ctx context.Context) *robot.UpdateRobotParams
}

var _ RobotUpdater = (*UpdateRobotInputParams)(nil)

func NewUpdateRobotInputParams(
	robotID int64,
	duration int64,
	name,
	description string,
) *UpdateRobotInputParams {
	return &UpdateRobotInputParams{
		RobotID:     robotID,
		Duration:    duration,
		Name:        name,
		Description: description,
	}
}

func (r UpdateRobotInputParams) GetRobotID() int64 {
	return r.RobotID
}

func (r UpdateRobotInputParams) GetDuration() int64 {
	return r.Duration
}

func (r UpdateRobotInputParams) GetName() string {
	return r.Name
}

func (r UpdateRobotInputParams) GetDescription() string {
	return r.Description
}

func (r *UpdateRobotInputParams) UpdateRobotParams(
	ctx context.Context) *robot.UpdateRobotParams {

	updateRobotParams := &robot.UpdateRobotParams{
		RobotID: r.GetRobotID(),
		Robot: &models.Robot{
			Name:        r.GetName(),
			Description: r.GetDescription(),
			Disable:     false,
			Duration:    r.GetDuration(),
			Secret:      "random",
			Level:       "system",
			Permissions: []*models.RobotPermission{
				{
					Access: []*models.Access{
						{
							Action:   "list",
							Resource: "repository",
						},
						{
							Action:   "pull",
							Resource: "repository",
						},
					},
					Kind:      "project",
					Namespace: "querypie",
				},
			},
		},
		Context: ctx,
	}

	return updateRobotParams
}
