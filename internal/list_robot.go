package internal

import (
	"context"
	"fmt"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/models"
	"github.com/ghdwlsgur/harborctl/utils"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

type Robot struct {
	ID           int64
	Name         string
	Description  string
	CreationTime string
	ExpiredTime  int64
	Dday         string
	Duration     int64
	Editable     bool
}

type ListRobotInputParams struct {
	Page      *int64
	PageSize  *int64
	TotalSize int
	Token     string
}

type RobotLister interface {
	setPage(page *int64)
	setTotalSize(ctx context.Context) (int, error)
	GetPage() *int64
	GetPageSize() *int64
	GetTotalSize() int
	GetToken() string
	GetMaxPage() int
	NextPage()
	Payload(ctx context.Context) (bool, []*models.Robot, error)
}

var _ RobotLister = (*ListRobotInputParams)(nil)

func NewListRobotInputParams(
	ctx context.Context,
	token string) (*ListRobotInputParams, error) {

	var err error
	page := int64(1)
	pageSize := int64(10)
	l := &ListRobotInputParams{
		Page:     &page,
		PageSize: &pageSize,
		Token:    token,
	}

	l.TotalSize, err = l.setTotalSize(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewListRobotInputParams: %w", err)
	}

	return l, nil
}

func (r *ListRobotInputParams) setPage(page *int64) {
	if page == nil {
		page = new(int64)
		*page = 1
		r.Page = page
	}
	r.Page = page
}

func (r *ListRobotInputParams) setTotalSize(ctx context.Context) (int, error) {
	robot, err := utils.NewRobotClient().ListRobot(
		nil, /* params */
		utils.SetAuthorizationWithToken(r.GetToken()), /* authInfo */
	)
	if err != nil {
		return 0, fmt.Errorf("setTotalSize: %w", err)
	}

	return int(robot.XTotalCount), nil
}

func (r ListRobotInputParams) GetPage() *int64 {
	return r.Page
}

func (r ListRobotInputParams) GetPageSize() *int64 {
	return r.PageSize
}

func (r ListRobotInputParams) GetTotalSize() int {
	return r.TotalSize
}

func (r ListRobotInputParams) GetToken() string {
	return r.Token
}

func (r *ListRobotInputParams) GetMaxPage() int {
	return r.GetTotalSize() + int(*r.GetPageSize()) - 1/int(*r.GetPageSize())
}

func (r *ListRobotInputParams) NextPage() {
	if int(*r.GetPage()) < r.GetMaxPage() {
		page := *r.GetPage() + 1
		r.setPage(&page)
	}
}

func (r *ListRobotInputParams) Payload(ctx context.Context) (bool, []*models.Robot, error) {
	params := &robot.ListRobotParams{
		Page:     r.GetPage(),     /* Page */
		PageSize: r.GetPageSize(), /* PageSize */
		Context:  ctx,             /* Context */
	}

	robotListed, err := utils.NewRobotClient().ListRobot(
		params, /* params */
		utils.SetAuthorizationWithToken(r.GetToken()), /* authInfo */
	)
	if err != nil {
		err = fmt.Errorf("Payload: %w", err)
		return false, nil, err
	}

	return robotListed.IsSuccess(), robotListed.GetPayload(), nil
}

func ListRobotTableOutput(
	writer table.Writer,
	robotTable *Robot) (table.Writer, error) {

	creationTime, err := utils.CreationTimeFormatKST(robotTable.CreationTime)
	if err != nil {
		err = fmt.Errorf("utils.CreationTimeFormatKST - ListRobotTableOutput: %w", err)
		return nil, err
	}

	leftDays := utils.CountDays(robotTable.ExpiredTime).Validate()
	writer.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Creation_Time",
		"Expired_Time",
		"D_Day",
		"Duration",
	})
	colorColumn := []string{
		"ID",
		"Name",
		"Description",
		"Creation_Time",
		"Expired_Time",
		"D_Day",
	}
	if leftDays == utils.ExpiredMessage {
		colorConfig := []table.ColumnConfig{}
		for _, v := range colorColumn {
			colorConfig = append(colorConfig, table.ColumnConfig{
				Name:        v,
				Align:       text.AlignLeft,
				AlignHeader: text.AlignLeft,
				AlignFooter: text.AlignLeft,
				Colors:      text.Colors{text.FgHiRed, text.Bold},
			})
		}
		writer.SetColumnConfigs(colorConfig)
	}
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
