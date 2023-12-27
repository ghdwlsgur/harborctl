package utils

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	basePath = "/api/v2.0/"
)

var (
	ErrInvalidMethod   = errors.New("failed to get path because method is invalid - GetRobotPath")
	ErrRobotIdEmpty    = errors.New("failed to get path because robotId is empty - GetRobotPath")
	ErrRobotIdNotEmpty = errors.New("failed to get path because robotId is not empty - GetRobotPath")
)

func GetRobotPath(method HTTPMethod, robotId ...string) (string, error) {
	var builder strings.Builder

	builder.WriteString(basePath)

	switch method {
	case "post", "get":
		if len(robotId) > 0 {
			return "", ErrRobotIdNotEmpty
		}
		builder.WriteString("robots")
		return builder.String(), nil

	case "delete", "patch", "list", "put":
		if len(robotId) == 0 || robotId[0] == "" {
			return "", ErrRobotIdEmpty
		}
		builder.WriteString("robots/")
		builder.WriteString(robotId[0])
		return builder.String(), nil

	default:
		return "", ErrInvalidMethod
	}
}
