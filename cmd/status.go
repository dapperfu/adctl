/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ewosborne/adctl/common"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check and change adblocking status",
	RunE:  StatusGetCmdE,
}

// statusCmd represents the status command
//
//lint:ignore U1000 not sure why it's unhappy
var statusGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get adblock status",
	RunE:  StatusGetCmdE,
}

type Status struct {
	Protection_enabled           bool
	Protection_disabled_duration uint64
}

type DisableTime struct {
	Duration   string
	HasTimeout bool
}

type ReadableStatus struct {
	Protection_enabled           bool
	Protection_disabled_duration string
}

func StatusGetCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		// Multi-server mode
		return GetStatusAll(servers)
	}

	// Single server mode
	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}
	s, err := GetStatus(server)
	if err != nil {
		return err
	}
	return PrintStatus(s)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func printToggle() error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		// Multi-server mode
		return toggleCommandAll(servers)
	}

	// Single server mode
	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	err = toggleCommand(server)
	if err != nil {
		return err
	}

	status, err := GetStatus(server)
	if err != nil {
		return err
	}
	PrintStatus(status)
	return nil
}

func toggleCommand(server *common.ServerConfig) error {
	status, err := GetStatus(server)
	if err != nil {
		return err
	}

	dTime := DisableTime{HasTimeout: false}
	switch status.Protection_enabled {
	case true:
		_, err = disableCommand(server, dTime)
	case false:
		_, err = enableCommand(server)
	}

	return err
}

func toggleCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string `json:"server"`
		Status Status `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		err := toggleCommand(&server)
		if err != nil {
			result.Error = err.Error()
		} else {
			status, err := GetStatus(&server)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Status = status
			}
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func PrintStatus(status Status) error {
	// status, err := GetStatus()
	// if err != nil {
	// 	return fmt.Errorf("error getting status: %w", err)
	// }

	var readableStatus ReadableStatus
	readableStatus.Protection_enabled = status.Protection_enabled

	if status.Protection_disabled_duration > 0 {
		readableStatus.Protection_disabled_duration = time.Duration(
			status.Protection_disabled_duration * uint64(time.Millisecond)).Truncate(time.Second).String()
	}

	tmp, err := json.MarshalIndent(readableStatus, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(tmp))
	return nil
}

// GetStatus gets status for a specific server (nil means legacy/viper config)
func GetStatus(server *common.ServerConfig) (Status, error) {
	var ret Status

	// build the command, it's specific to status
	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return ret, err
	}
	baseURL.Path = "/control/status"

	statusQuery := common.CommandArgs{
		Method: "GET",
		URL:    baseURL,
		Server: server,
	}

	body, err := common.SendCommand(statusQuery)
	if err != nil {
		return ret, err
	}

	// serialize body into Status and return appropriately
	var s Status
	json.Unmarshal(body, &s)

	return s, nil
}

// GetStatusAll gets status for all servers
func GetStatusAll(servers []common.ServerConfig) error {
	type ServerStatus struct {
		Server string `json:"server"`
		Status Status `json:"status"`
		Error  string `json:"error,omitempty"`
	}

	var results []ServerStatus
	for _, server := range servers {
		status, err := GetStatus(&server)
		result := ServerStatus{
			Server: server.Name,
		}
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Status = status
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
