package utils

import (
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client"
	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type HTTPMethod string

const (
	Get    HTTPMethod = "get"
	Post   HTTPMethod = "post"
	Put    HTTPMethod = "put"
	Delete HTTPMethod = "delete"
	Patch  HTTPMethod = "patch"
	List   HTTPMethod = "list"
)

func NewRobotClient() robot.ClientService {
	return client.Default.Robot
}

func SetAuthorizationWithToken(token string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(
		func(req runtime.ClientRequest, reg strfmt.Registry) error {
			return req.SetHeaderParam("Authorization", "Basic "+token)
		})
}
