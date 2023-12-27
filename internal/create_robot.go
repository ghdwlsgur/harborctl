package internal

import (
	"context"
	"fmt"

	"github.com/ghdwlsgur/captain/utils"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/models"
	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"
)

type CreateRobotInputParams struct {
	Name        string
	Description string
	Duration    int64
}

type RobotCreater interface {
	CreateRobotParams(ctx context.Context) (*robot.CreateRobotParams, error)
}

func NewCreateRobotInputParams(name,
	description string,
	duration int64) *CreateRobotInputParams {

	return &CreateRobotInputParams{
		Name:        name,
		Description: description,
		Duration:    duration,
	}
}

func (r CreateRobotInputParams) GetName() string {
	return r.Name
}

func (r CreateRobotInputParams) GetDescription() string {
	return r.Description
}

func (r CreateRobotInputParams) GetDuration() int64 {
	return r.Duration
}

func (r *CreateRobotInputParams) CreateRobotParams(ctx context.Context) (*robot.CreateRobotParams, error) {
	createRobotParamsWithContext := robot.NewCreateRobotParamsWithContext(ctx)
	createRobotModel := &models.RobotCreate{
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
	}

	if err := createRobotModel.Validate(strfmt.Default); err != nil {
		return nil, fmt.Errorf("failed to validate robot create - CreateRobotParams: %w", err)
	}
	createRobotParamsWithContext.SetRobot(createRobotModel)

	return createRobotParamsWithContext, nil
}

func OutputCreateRobotTable(
	writer table.Writer,
	robot *robot.CreateRobotCreated) (table.Writer, error) {

	days, err := utils.CountDays(robot.GetPayload().ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count days - OutputCreateRobotTable: %w", err)
	}

	creationTime, err := utils.CreationTimeFormatKST(robot.GetPayload().CreationTime.String())
	if err != nil {
		return nil, fmt.Errorf("failed to format creation time - OutputCreateRobotTable: %w", err)
	}

	writer.AppendHeader(table.Row{"Name", "Secret", "Expires_At", "D_Day", "Creation_Time"})
	writer.AppendRow(table.Row{
		robot.GetPayload().Name,                                   /* name */
		robot.GetPayload().Secret,                                 /* secret */
		utils.ExpiresAtToStringTime(robot.GetPayload().ExpiresAt), /* expires_at */
		days,         /* d_day */
		creationTime, /* creation_time */
	})

	return writer, nil
}
