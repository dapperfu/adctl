/*
Copyright © 2024 Eric Osborne
No header.
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage AdGuard server configurations",
	Long:  `Add, list, and manage AdGuard Home server configurations.`,
}

// serverAddCmd represents the server add command
var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new AdGuard server configuration",
	Long:  `Interactively prompt for server details and add a new server to the configuration.`,
	RunE:  serverAddCmdE,
}

// serverListCmd represents the server list command
var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured AdGuard servers",
	Long:  `Display all configured AdGuard Home servers (passwords are masked).`,
	RunE:  serverListCmdE,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverAddCmd)
	serverCmd.AddCommand(serverListCmd)
}

func serverAddCmdE(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for server name
	fmt.Print("Server name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read server name: %w", err)
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("server name cannot be empty")
	}
	if name == ReservedServerName {
		return fmt.Errorf("'%s' is a reserved server name", ReservedServerName)
	}

	// Check if server already exists
	if ServerExists(name) {
		return fmt.Errorf("server '%s' already exists", name)
	}

	// Prompt for host
	fmt.Print("Host (host:port, e.g., router.example.com:8080): ")
	host, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read host: %w", err)
	}
	host = strings.TrimSpace(host)

	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Prompt for username
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	username = strings.TrimSpace(username)

	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Prompt for password (masked)
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // New line after password input
	password := string(passwordBytes)

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Create server config
	server := ServerConfig{
		Name:     name,
		Host:     host,
		Username: username,
		Password: password,
	}

	// Add server to config
	if err := AddServer(server); err != nil {
		return fmt.Errorf("failed to add server: %w", err)
	}

	fmt.Printf("Server '%s' added successfully.\n", name)
	return nil
}

func serverListCmdE(cmd *cobra.Command, args []string) error {
	servers, err := GetServers()
	if err != nil {
		return fmt.Errorf("failed to get servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No servers configured.")
		return nil
	}

	// Create a list with masked passwords for display
	type ServerDisplay struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	displayServers := make([]ServerDisplay, len(servers))
	for i, s := range servers {
		displayServers[i] = ServerDisplay{
			Name:     s.Name,
			Host:     s.Host,
			Username: s.Username,
			Password: "***",
		}
	}

	// Output as JSON
	output, err := json.MarshalIndent(displayServers, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal servers: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
