package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ghdwlsgur/harbor-robot-sdk/pkg/sdk/robot/client/robot"
	"github.com/ghdwlsgur/harborctl/utils"
)

type Request struct {
	Description string `json:"description"`
	Secret      string `json:"secret"`
	CreatedAt   string `json:"creation_time"`
	ExpiresAt   string `json:"expires_at"`
}

type Response struct {
	Code     int    `json:"code"`
	Message  string `json:"msg"`
	Response struct {
		RobotID int    `json:"robot_id"`
		Message string `json:"msg"`
		Data    string `json:"data"`
	} `json:"response"`
}

const (
	baseURL = "http://10.60.10.144"
	port    = "8081"
)

func GetURL(robotId int) string {
	return baseURL + ":" + port + "/api/v1/harborctl/robot/" + strconv.Itoa(robotId)
}

func DeleteSecret(robotID int64, token string) (*Response, error) {
	req, err := http.NewRequest("DELETE", GetURL(int(robotID)), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest - GetSecret: %w", err)
	}

	req.Header.Add("Authorization", "Basic "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do - GetSecret: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll - GetSecret: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal - GetSecret: %w", err)
	}

	return &response, nil
}

func GetSecret(robotID int64, token string) (*Response, error) {
	req, err := http.NewRequest("GET", GetURL(int(robotID)), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest - GetSecret: %w", err)
	}

	req.Header.Add("Authorization", "Basic "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do - GetSecret: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll - GetSecret: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal - GetSecret: %w", err)
	}

	return &response, nil
}

func SaveSecret(
	robot *robot.CreateRobotCreated,
	token,
	description string) (*Response, error) {

	data := Request{
		Description: description,
		Secret:      robot.GetPayload().Secret,
		CreatedAt:   robot.GetPayload().CreationTime.String(),
		ExpiresAt:   utils.ExpiresAtToStringTime(robot.GetPayload().ExpiresAt),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal - RequestStoreSecret: %w", err)
	}

	req, err := http.NewRequest("POST", GetURL(int(robot.GetPayload().ID)), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest - RequestStoreSecret: %w", err)
	}
	req.Header.Add("Authorization", "Basic "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do - RequestStoreSecret: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll - RequestStoreSecret: %w", err)
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal - RequestStoreSecret: %w", err)
	}

	return &response, nil
}
