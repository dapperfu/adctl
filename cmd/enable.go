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

var statusEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable ad blocking",
	RunE:  StatusEnableCmdE,
}

func StatusEnableCmdE(cmd *cobra.Command, flags []string) error {
	return printEnable()
}

func init() {
	rootCmd.AddCommand(statusEnableCmd)

}

func printEnable() error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		// Multi-server mode
		return enableCommandAll(servers)
	}

	// Single server mode
	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	status, err := enableCommand(server)
	if err != nil {
		return err
	}
	err = PrintStatus(status)
	if err != nil {
		return err
	}
	return nil
}

func enableCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string `json:"server"`
		Status Status `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		status, err := enableCommand(&server)
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

func enableCommand(server *common.ServerConfig) (Status, error) {
	err := common.AbleCommand(server, true, "")
	if err != nil {
		return Status{}, err
	}

	status, err := GetStatus(server)
	if err != nil {
		return Status{}, err
	}

	return status, nil
}
