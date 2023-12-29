package internal

import (
	"context"
	"fmt"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/models"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
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
	UpdateRobotParams(ctx context.Context) (*robot.UpdateRobotParams, error)
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
	ctx context.Context) (*robot.UpdateRobotParams, error) {

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

	if err := updateRobotParams.Robot.Validate(strfmt.Default); err != nil {
		return nil, fmt.Errorf("failed to validate robot update model - UpdateRobotParams: %w", err)
	}

	return updateRobotParams, nil
}

func UpdateRobotTableOutput(
	writer table.Writer,
	robot *models.Robot) (table.Writer, error) {

	creationTime, err := utils.CreationTimeFormatKST(robot.CreationTime.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse time - CreationTimeFormatKST: %w", err)
	}

	writer.SetTitle("After")
	writer.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Creation_Time",
		"Expired_Time",
		"D_Day",
		"Duration",
	})
	writer.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "Expired_Time",
			Align:       text.AlignLeft,
			AlignHeader: text.AlignLeft,
			AlignFooter: text.AlignLeft,
			Colors:      text.Colors{text.FgHiGreen, text.Bold},
		},
		{
			Name:        "D_Day",
			Align:       text.AlignLeft,
			AlignHeader: text.AlignLeft,
			AlignFooter: text.AlignLeft,
			Colors:      text.Colors{text.FgHiGreen, text.Bold},
		},
		{
			Name:        "Duration",
			Align:       text.AlignRight,
			AlignHeader: text.AlignRight,
			AlignFooter: text.AlignRight,
			Colors:      text.Colors{text.FgHiGreen, text.Bold},
		},
	})
	writer.AppendRow(table.Row{
		robot.ID,
		robot.Name,
		robot.Description,
		creationTime,
		utils.ExpiresAtToStringTime(robot.ExpiresAt),
		utils.CountDays(robot.ExpiresAt).Validate(),
		robot.Duration,
	})

	return writer, nil
}
