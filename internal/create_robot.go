package internal

import (
	"context"
	"fmt"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/models"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"
)

type CreateRobotInputParams struct {
	Name        string
	Description string
	ExpiresAt   string
	Duration    int64
}

type RobotCreater interface {
	GetName() string
	GetDescription() string
	GetDuration() int64
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

var _ RobotCreater = (*CreateRobotInputParams)(nil)

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
		return nil, fmt.Errorf("failed to validate robot create model - CreateRobotParams: %w", err)
	}
	createRobotParamsWithContext.SetRobot(createRobotModel)

	return createRobotParamsWithContext, nil
}

func CreateRobotTableOutput(
	writer table.Writer,
	robot *robot.CreateRobotCreated,
	description string) (table.Writer, error) {

	count := utils.CountDays(robot.GetPayload().ExpiresAt)
	if count.Err != nil {
		return nil, fmt.Errorf("failed to count days - OutputCreateRobotTable: %w", count.Err)
	}

	creationTime, err := utils.CreationTimeFormatKST(robot.GetPayload().CreationTime.String())
	if err != nil {
		return nil, fmt.Errorf("failed to format creation time - OutputCreateRobotTable: %w", err)
	}

	writer.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Secret",
		"Creation_Time",
		"Expired_time",
		"D_Day",
	})
	writer.AppendRow(table.Row{
		robot.GetPayload().ID,     /* id */
		robot.GetPayload().Name,   /* name */
		description,               /* description */
		robot.GetPayload().Secret, /* secret */
		creationTime,              /* creation_time */
		utils.ExpiresAtToStringTime(robot.GetPayload().ExpiresAt), /* expires_at */
		count.LeftDays, /* d_day */
	})

	return writer, nil
}
