/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ewosborne/adctl/common"
	"github.com/spf13/cobra"
)

// dhcpCmd represents the dhcp command
var dhcpCmd = &cobra.Command{
	Use:   "dhcp",
	Short: "Manage DHCP server configuration and leases",
	Long:  "Get DHCP status, manage leases, configure DHCP server, and manage static leases.",
}

// dhcpStatusCmd represents the status command
var dhcpStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get DHCP server status and configuration",
	RunE:  dhcpStatusCmdE,
}

// dhcpLeasesCmd represents the leases command
var dhcpLeasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "List active DHCP leases",
	RunE:  dhcpLeasesCmdE,
}

// dhcpCheckCmd represents the check command
var dhcpCheckCmd = &cobra.Command{
	Use:   "check <interface>",
	Short: "Check for active DHCP servers on an interface",
	Long:  "Check for active DHCP servers on the specified network interface. The interface name is required (e.g., eth0, wlan0).",
	Args:  cobra.ExactArgs(1),
	RunE:  dhcpCheckCmdE,
}

// dhcpConfigCmd represents the config command
var dhcpConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure DHCP server",
	RunE:  dhcpConfigCmdE,
}

// dhcpResetCmd represents the reset command
var dhcpResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset DHCP configuration",
	RunE:  dhcpResetCmdE,
}

// dhcpStaticLeaseCmd represents the static-lease command
var dhcpStaticLeaseCmd = &cobra.Command{
	Use:   "static-lease",
	Short: "Manage static DHCP leases",
}

// dhcpStaticLeaseListCmd represents the list command
var dhcpStaticLeaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all static leases",
	RunE:  dhcpStaticLeaseListCmdE,
}

// dhcpStaticLeaseAddCmd represents the add command
var dhcpStaticLeaseAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a static lease",
	RunE:  dhcpStaticLeaseAddCmdE,
}

// dhcpStaticLeaseRemoveCmd represents the remove command
var dhcpStaticLeaseRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a static lease",
	RunE:  dhcpStaticLeaseRemoveCmdE,
}

// dhcpStaticLeaseUpdateCmd represents the update command
var dhcpStaticLeaseUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a static lease",
	RunE:  dhcpStaticLeaseUpdateCmdE,
}

// DHCP status data structures
type DHCPStatus struct {
	Enabled       bool          `json:"enabled"`
	InterfaceName string        `json:"interface_name"`
	V4            V4Config      `json:"v4"`
	V6            V6Config      `json:"v6"`
	Leases        []LeaseDynamic `json:"leases"`
	StaticLeases  []LeaseStatic  `json:"static_leases"`
}

type V4Config struct {
	GatewayIP     string `json:"gateway_ip"`
	SubnetMask    string `json:"subnet_mask"`
	RangeStart    string `json:"range_start"`
	RangeEnd      string `json:"range_end"`
	LeaseDuration uint32 `json:"lease_duration"`
}

type V6Config struct {
	RangeStart    string `json:"range_start"`
	LeaseDuration uint32 `json:"lease_duration"`
}

type LeaseDynamic struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	Expires  string `json:"expires"`
}

type LeaseStatic struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
}

// DHCP check response structures
type DHCPCheckResponse struct {
	V4 DHCPCheckV4 `json:"v4"`
	V6 DHCPCheckV6 `json:"v6"`
}

type DHCPCheckV4 struct {
	OtherServer DHCPCheckOtherServer `json:"other_server"`
	StaticIP    DHCPCheckStaticIP    `json:"static_ip"`
}

type DHCPCheckV6 struct {
	OtherServer DHCPCheckOtherServer `json:"other_server"`
}

type DHCPCheckOtherServer struct {
	Found string `json:"found"`
	Error string `json:"error,omitempty"`
}

type DHCPCheckStaticIP struct {
	Static string `json:"static"`
	IP     string `json:"ip,omitempty"`
}

// Flags for config command
var dhcpConfigEnabled bool
var dhcpConfigInterface string
var dhcpConfigV4Gateway string
var dhcpConfigV4Subnet string
var dhcpConfigV4RangeStart string
var dhcpConfigV4RangeEnd string
var dhcpConfigV4LeaseDuration uint32
var dhcpConfigV6RangeStart string
var dhcpConfigV6LeaseDuration uint32

// Flags for static lease commands
var staticLeaseIP string
var staticLeaseMAC string
var staticLeaseHostname string

func init() {
	rootCmd.AddCommand(dhcpCmd)

	dhcpCmd.AddCommand(dhcpStatusCmd)
	dhcpCmd.AddCommand(dhcpLeasesCmd)
	dhcpCmd.AddCommand(dhcpCheckCmd)
	dhcpCmd.AddCommand(dhcpConfigCmd)
	dhcpCmd.AddCommand(dhcpResetCmd)
	dhcpCmd.AddCommand(dhcpStaticLeaseCmd)

	// Config command flags
	dhcpConfigCmd.Flags().BoolVar(&dhcpConfigEnabled, "enabled", false, "Enable/disable DHCP")
	dhcpConfigCmd.Flags().StringVar(&dhcpConfigInterface, "interface", "", "Interface name")
	dhcpConfigCmd.Flags().StringVar(&dhcpConfigV4Gateway, "v4-gateway", "", "IPv4 gateway IP")
	dhcpConfigCmd.Flags().StringVar(&dhcpConfigV4Subnet, "v4-subnet", "", "IPv4 subnet mask")
	dhcpConfigCmd.Flags().StringVar(&dhcpConfigV4RangeStart, "v4-range-start", "", "IPv4 range start")
	dhcpConfigCmd.Flags().StringVar(&dhcpConfigV4RangeEnd, "v4-range-end", "", "IPv4 range end")
	dhcpConfigCmd.Flags().Uint32Var(&dhcpConfigV4LeaseDuration, "v4-lease-duration", 0, "IPv4 lease duration in minutes")
	dhcpConfigCmd.Flags().StringVar(&dhcpConfigV6RangeStart, "v6-range-start", "", "IPv6 range start")
	dhcpConfigCmd.Flags().Uint32Var(&dhcpConfigV6LeaseDuration, "v6-lease-duration", 0, "IPv6 lease duration in minutes")

	// Static lease subcommands
	dhcpStaticLeaseCmd.AddCommand(dhcpStaticLeaseListCmd)
	dhcpStaticLeaseCmd.AddCommand(dhcpStaticLeaseAddCmd)
	dhcpStaticLeaseCmd.AddCommand(dhcpStaticLeaseRemoveCmd)
	dhcpStaticLeaseCmd.AddCommand(dhcpStaticLeaseUpdateCmd)

	// Static lease flags
	dhcpStaticLeaseAddCmd.Flags().StringVar(&staticLeaseIP, "ip", "", "IP address (required)")
	dhcpStaticLeaseAddCmd.Flags().StringVar(&staticLeaseMAC, "mac", "", "MAC address (required)")
	dhcpStaticLeaseAddCmd.Flags().StringVar(&staticLeaseHostname, "hostname", "", "Hostname (required)")
	dhcpStaticLeaseAddCmd.MarkFlagRequired("ip")
	dhcpStaticLeaseAddCmd.MarkFlagRequired("mac")
	dhcpStaticLeaseAddCmd.MarkFlagRequired("hostname")

	dhcpStaticLeaseRemoveCmd.Flags().StringVar(&staticLeaseIP, "ip", "", "IP address")
	dhcpStaticLeaseRemoveCmd.Flags().StringVar(&staticLeaseMAC, "mac", "", "MAC address")

	dhcpStaticLeaseUpdateCmd.Flags().StringVar(&staticLeaseIP, "ip", "", "IP address to update (required)")
	dhcpStaticLeaseUpdateCmd.Flags().StringVar(&staticLeaseMAC, "mac", "", "New MAC address")
	dhcpStaticLeaseUpdateCmd.Flags().StringVar(&staticLeaseHostname, "hostname", "", "New hostname")
	dhcpStaticLeaseUpdateCmd.MarkFlagRequired("ip")
}

// dhcpStatusCmdE handles the dhcp status command
func dhcpStatusCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpStatusCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	status, err := getDHCPStatus(server)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(status, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// dhcpLeasesCmdE handles the dhcp leases command
func dhcpLeasesCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpLeasesCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	status, err := getDHCPStatus(server)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(status.Leases, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal leases: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// dhcpCheckCmdE handles the dhcp check command
func dhcpCheckCmdE(cmd *cobra.Command, args []string) error {
	interfaceName := args[0]
	if interfaceName == "" {
		return fmt.Errorf("interface name cannot be empty")
	}

	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpCheckCommandAll(servers, interfaceName)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	result, err := checkDHCP(server, interfaceName)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal check result: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// dhcpConfigCmdE handles the dhcp config command
func dhcpConfigCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpConfigCommandAll(servers, cmd)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	err = setDHCPConfig(server, cmd)
	if err != nil {
		return err
	}

	// Return updated status
	return dhcpStatusCmdE(cmd, args)
}

// dhcpResetCmdE handles the dhcp reset command
func dhcpResetCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpResetCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	err = resetDHCP(server)
	if err != nil {
		return err
	}

	// Return updated status
	return dhcpStatusCmdE(cmd, args)
}

// dhcpStaticLeaseListCmdE handles the static-lease list command
func dhcpStaticLeaseListCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpStaticLeaseListCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	status, err := getDHCPStatus(server)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(status.StaticLeases, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal static leases: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// dhcpStaticLeaseAddCmdE handles the static-lease add command
func dhcpStaticLeaseAddCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpStaticLeaseAddCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	err = addStaticLease(server, staticLeaseIP, staticLeaseMAC, staticLeaseHostname)
	if err != nil {
		return err
	}

	// Return updated list
	return dhcpStaticLeaseListCmdE(cmd, args)
}

// dhcpStaticLeaseRemoveCmdE handles the static-lease remove command
func dhcpStaticLeaseRemoveCmdE(cmd *cobra.Command, args []string) error {
	if staticLeaseIP == "" && staticLeaseMAC == "" {
		return fmt.Errorf("at least one of --ip or --mac is required")
	}

	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpStaticLeaseRemoveCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	err = removeStaticLease(server, staticLeaseIP, staticLeaseMAC)
	if err != nil {
		return err
	}

	// Return updated list
	return dhcpStaticLeaseListCmdE(cmd, args)
}

// dhcpStaticLeaseUpdateCmdE handles the static-lease update command
func dhcpStaticLeaseUpdateCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetCurrentServers()
	if err != nil {
		return err
	}

	if serverFlag == ReservedServerName && len(servers) > 1 {
		return dhcpStaticLeaseUpdateCommandAll(servers)
	}

	var server *common.ServerConfig
	if len(servers) > 0 {
		server = &servers[0]
	}

	err = updateStaticLease(server, staticLeaseIP, staticLeaseMAC, staticLeaseHostname)
	if err != nil {
		return err
	}

	// Return updated list
	return dhcpStaticLeaseListCmdE(cmd, args)
}

// getDHCPStatus gets DHCP status for a server
func getDHCPStatus(server *common.ServerConfig) (DHCPStatus, error) {
	var ret DHCPStatus

	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return ret, err
	}
	baseURL.Path = "/control/dhcp/status"

	statusQuery := common.CommandArgs{
		Method: "GET",
		URL:    baseURL,
		Server: server,
	}

	body, err := common.SendCommand(statusQuery)
	if err != nil {
		return ret, fmt.Errorf("failed to get DHCP status: %w", err)
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return ret, fmt.Errorf("failed to unmarshal DHCP status: %w", err)
	}

	return ret, nil
}

// checkDHCP checks for active DHCP servers on an interface
func checkDHCP(server *common.ServerConfig, interfaceName string) (DHCPCheckResponse, error) {
	var ret DHCPCheckResponse

	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return ret, err
	}
	baseURL.Path = "/control/dhcp/find_active_dhcp"

	requestBody := make(map[string]any)
	requestBody["interface"] = interfaceName

	checkQuery := common.CommandArgs{
		Method:      "POST",
		URL:         baseURL,
		RequestBody: requestBody,
		Server:      server,
	}

	body, err := common.SendCommand(checkQuery)
	if err != nil {
		return ret, fmt.Errorf("failed to check DHCP: %w", err)
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return ret, fmt.Errorf("failed to unmarshal DHCP check result: %w", err)
	}

	return ret, nil
}

// setDHCPConfig sets DHCP configuration
func setDHCPConfig(server *common.ServerConfig, cmd *cobra.Command) error {
	// Get current status to preserve existing config
	currentStatus, err := getDHCPStatus(server)
	if err != nil {
		return fmt.Errorf("failed to get current DHCP status: %w", err)
	}

	requestBody := make(map[string]any)

	// Set enabled flag if provided, otherwise preserve current
	if cmd.Flags().Changed("enabled") {
		requestBody["enabled"] = dhcpConfigEnabled
	} else {
		requestBody["enabled"] = currentStatus.Enabled
	}

	// Set interface if provided
	if dhcpConfigInterface != "" {
		requestBody["interface_name"] = dhcpConfigInterface
	} else {
		requestBody["interface_name"] = currentStatus.InterfaceName
	}

	// Build v4 config
	v4Config := make(map[string]any)
	if dhcpConfigV4Gateway != "" {
		v4Config["gateway_ip"] = dhcpConfigV4Gateway
	} else if currentStatus.V4.GatewayIP != "" {
		v4Config["gateway_ip"] = currentStatus.V4.GatewayIP
	}

	if dhcpConfigV4Subnet != "" {
		v4Config["subnet_mask"] = dhcpConfigV4Subnet
	} else if currentStatus.V4.SubnetMask != "" {
		v4Config["subnet_mask"] = currentStatus.V4.SubnetMask
	}

	if dhcpConfigV4RangeStart != "" {
		v4Config["range_start"] = dhcpConfigV4RangeStart
	} else if currentStatus.V4.RangeStart != "" {
		v4Config["range_start"] = currentStatus.V4.RangeStart
	}

	if dhcpConfigV4RangeEnd != "" {
		v4Config["range_end"] = dhcpConfigV4RangeEnd
	} else if currentStatus.V4.RangeEnd != "" {
		v4Config["range_end"] = currentStatus.V4.RangeEnd
	}

	if dhcpConfigV4LeaseDuration > 0 {
		v4Config["lease_duration"] = dhcpConfigV4LeaseDuration
	} else if currentStatus.V4.LeaseDuration > 0 {
		v4Config["lease_duration"] = currentStatus.V4.LeaseDuration
	}

	if len(v4Config) > 0 {
		requestBody["v4"] = v4Config
	} else {
		requestBody["v4"] = currentStatus.V4
	}

	// Build v6 config
	v6Config := make(map[string]any)
	if dhcpConfigV6RangeStart != "" {
		v6Config["range_start"] = dhcpConfigV6RangeStart
	} else if currentStatus.V6.RangeStart != "" {
		v6Config["range_start"] = currentStatus.V6.RangeStart
	}

	if dhcpConfigV6LeaseDuration > 0 {
		v6Config["lease_duration"] = dhcpConfigV6LeaseDuration
	} else if currentStatus.V6.LeaseDuration > 0 {
		v6Config["lease_duration"] = currentStatus.V6.LeaseDuration
	}

	if len(v6Config) > 0 {
		requestBody["v6"] = v6Config
	} else {
		requestBody["v6"] = currentStatus.V6
	}

	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return err
	}
	baseURL.Path = "/control/dhcp/set_config"

	configQuery := common.CommandArgs{
		Method:      "POST",
		URL:         baseURL,
		RequestBody: requestBody,
		Server:      server,
	}

	_, err = common.SendCommand(configQuery)
	if err != nil {
		return fmt.Errorf("failed to set DHCP config: %w", err)
	}

	return nil
}

// resetDHCP resets DHCP configuration
func resetDHCP(server *common.ServerConfig) error {
	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return err
	}
	baseURL.Path = "/control/dhcp/reset"

	resetQuery := common.CommandArgs{
		Method: "POST",
		URL:    baseURL,
		Server: server,
	}

	_, err = common.SendCommand(resetQuery)
	if err != nil {
		return fmt.Errorf("failed to reset DHCP: %w", err)
	}

	return nil
}

// addStaticLease adds a static lease
func addStaticLease(server *common.ServerConfig, ip, mac, hostname string) error {
	requestBody := make(map[string]any)
	requestBody["ip"] = ip
	requestBody["mac"] = mac
	requestBody["hostname"] = hostname

	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return err
	}
	baseURL.Path = "/control/dhcp/add_static_lease"

	addQuery := common.CommandArgs{
		Method:      "POST",
		URL:         baseURL,
		RequestBody: requestBody,
		Server:      server,
	}

	_, err = common.SendCommand(addQuery)
	if err != nil {
		return fmt.Errorf("failed to add static lease: %w", err)
	}

	return nil
}

// removeStaticLease removes a static lease
func removeStaticLease(server *common.ServerConfig, ip, mac string) error {
	requestBody := make(map[string]any)
	if ip != "" {
		requestBody["ip"] = ip
	}
	if mac != "" {
		requestBody["mac"] = mac
	}

	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return err
	}
	baseURL.Path = "/control/dhcp/remove_static_lease"

	removeQuery := common.CommandArgs{
		Method:      "POST",
		URL:         baseURL,
		RequestBody: requestBody,
		Server:      server,
	}

	_, err = common.SendCommand(removeQuery)
	if err != nil {
		return fmt.Errorf("failed to remove static lease: %w", err)
	}

	return nil
}

// updateStaticLease updates a static lease
func updateStaticLease(server *common.ServerConfig, ip, mac, hostname string) error {
	requestBody := make(map[string]any)
	requestBody["ip"] = ip
	if mac != "" {
		requestBody["mac"] = mac
	}
	if hostname != "" {
		requestBody["hostname"] = hostname
	}

	baseURL, err := common.GetBaseURL(server)
	if err != nil {
		return err
	}
	baseURL.Path = "/control/dhcp/update_static_lease"

	updateQuery := common.CommandArgs{
		Method:      "POST",
		URL:         baseURL,
		RequestBody: requestBody,
		Server:      server,
	}

	_, err = common.SendCommand(updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update static lease: %w", err)
	}

	return nil
}

// Multi-server support functions

func dhcpStatusCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string     `json:"server"`
		Result DHCPStatus `json:"result,omitempty"`
		Error  string     `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		status, err := getDHCPStatus(&server)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Result = status
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpLeasesCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string          `json:"server"`
		Result []LeaseDynamic `json:"result,omitempty"`
		Error  string          `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		status, err := getDHCPStatus(&server)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Result = status.Leases
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpCheckCommandAll(servers []common.ServerConfig, interfaceName string) error {
	type ServerResult struct {
		Server string            `json:"server"`
		Result DHCPCheckResponse `json:"result,omitempty"`
		Error  string            `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		checkResult, err := checkDHCP(&server, interfaceName)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Result = checkResult
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpConfigCommandAll(servers []common.ServerConfig, cmd *cobra.Command) error {
	type ServerResult struct {
		Server string     `json:"server"`
		Result DHCPStatus `json:"result,omitempty"`
		Error  string     `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		err := setDHCPConfig(&server, cmd)
		if err != nil {
			result.Error = err.Error()
		} else {
			status, err := getDHCPStatus(&server)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Result = status
			}
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpResetCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string     `json:"server"`
		Result DHCPStatus `json:"result,omitempty"`
		Error  string     `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		err := resetDHCP(&server)
		if err != nil {
			result.Error = err.Error()
		} else {
			status, err := getDHCPStatus(&server)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Result = status
			}
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpStaticLeaseListCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string         `json:"server"`
		Result []LeaseStatic `json:"result,omitempty"`
		Error  string         `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		status, err := getDHCPStatus(&server)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Result = status.StaticLeases
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpStaticLeaseAddCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string         `json:"server"`
		Result []LeaseStatic `json:"result,omitempty"`
		Error  string         `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		err := addStaticLease(&server, staticLeaseIP, staticLeaseMAC, staticLeaseHostname)
		if err != nil {
			result.Error = err.Error()
		} else {
			status, err := getDHCPStatus(&server)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Result = status.StaticLeases
			}
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpStaticLeaseRemoveCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string         `json:"server"`
		Result []LeaseStatic `json:"result,omitempty"`
		Error  string         `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		err := removeStaticLease(&server, staticLeaseIP, staticLeaseMAC)
		if err != nil {
			result.Error = err.Error()
		} else {
			status, err := getDHCPStatus(&server)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Result = status.StaticLeases
			}
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func dhcpStaticLeaseUpdateCommandAll(servers []common.ServerConfig) error {
	type ServerResult struct {
		Server string         `json:"server"`
		Result []LeaseStatic `json:"result,omitempty"`
		Error  string         `json:"error,omitempty"`
	}

	var results []ServerResult
	for _, server := range servers {
		result := ServerResult{Server: server.Name}
		err := updateStaticLease(&server, staticLeaseIP, staticLeaseMAC, staticLeaseHostname)
		if err != nil {
			result.Error = err.Error()
		} else {
			status, err := getDHCPStatus(&server)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Result = status.StaticLeases
			}
		}
		results = append(results, result)
	}

	output, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
