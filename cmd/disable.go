/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ewosborne/adctl/common"
	"github.com/spf13/cobra"
)

func disableCommand(server *common.ServerConfig, dTime DisableTime) (Status, error) {
	var err error

	if dTime.HasTimeout {
		err = common.AbleCommand(server, false, dTime.Duration)
	} else {
		err = common.AbleCommand(server, false, "")
	}

	if err != nil {
		return Status{}, err
	}

	s, err := GetStatus(server)

	return s, err
}

var statusDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable ad blocker. Optional duration in time.Duration format.",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  StatusDisableCmdE,
}

func init() {
	rootCmd.AddCommand(statusDisableCmd)
}

func StatusDisableCmdE(cmd *cobra.Command, args []string) error {

	var dTime = DisableTime{}

	switch len(args) {
	case 0:
		dTime.HasTimeout = false
	case 1:
		dTime.HasTimeout = true
		dTime.Duration = args[0]
	default:
		return fmt.Errorf("only one arg allowed for disable")
	}

	return printDisable(dTime)
}

func printDisable(dTime DisableTime) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		// Multi-server mode
		return disableCommandAll(servers, dTime)
	}

	// Single server mode
	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	status, err := disableCommand(server, dTime)
	if err != nil {
		return err
	}

	PrintStatus(status)
	return nil
}

func disableCommandAll(servers []common.ServerConfig, dTime DisableTime) error {
	type ServerResult struct {
		Server string `json:"server"`
		Status Status `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		status, err := disableCommand(&server, dTime)
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
