/*
Copyright © 2024 Eric Osborne
No header.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ewosborne/adctl/common"
	"github.com/spf13/viper"
)

const ReservedServerName = "all"

// GetConfigDir returns the OS-specific configuration directory
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	var configDir string
	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\adctl
		appData := os.Getenv("APPDATA")
		if appData == "" {
			configDir = filepath.Join(home, "AppData", "Roaming", "adctl")
		} else {
			configDir = filepath.Join(appData, "adctl")
		}
	default:
		// Linux/macOS: ~/.config/adctl
		configDir = filepath.Join(home, ".config", "adctl")
	}

	return configDir, nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "adctl.yaml"), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0755)
}

// GetServers returns all configured servers
func GetServers() ([]common.ServerConfig, error) {
	var servers []common.ServerConfig

	// Check for new multi-server format
	if viper.IsSet("servers") {
		err := viper.UnmarshalKey("servers", &servers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal servers: %w", err)
		}
		return servers, nil
	}

	// Check for legacy single-server format and migrate
	host := viper.GetString("host")
	username := viper.GetString("username")
	password := viper.GetString("password")

	if host != "" && username != "" && password != "" {
		// Legacy config found, create a default server
		servers = []common.ServerConfig{
			{
				Name:     "default",
				Host:     host,
				Username: username,
				Password: password,
			},
		}
		// Optionally migrate to new format (we'll do this on write)
	}

	return servers, nil
}

// GetServer returns a specific server configuration by name
func GetServer(name string) (*common.ServerConfig, error) {
	if name == ReservedServerName {
		return nil, fmt.Errorf("'%s' is a reserved server name", ReservedServerName)
	}

	servers, err := GetServers()
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.Name == name {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("server '%s' not found", name)
}

// AddServer adds a new server to the configuration
func AddServer(server common.ServerConfig) error {
	// Validate server name
	if server.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}
	if server.Name == ReservedServerName {
		return fmt.Errorf("'%s' is a reserved server name", ReservedServerName)
	}

	// Validate required fields
	if server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	if server.Username == "" {
		return fmt.Errorf("server username cannot be empty")
	}
	if server.Password == "" {
		return fmt.Errorf("server password cannot be empty")
	}

	// Check for duplicate names
	servers, err := GetServers()
	if err != nil {
		return err
	}

	for _, s := range servers {
		if s.Name == server.Name {
			return fmt.Errorf("server '%s' already exists", server.Name)
		}
	}

	// Add the new server
	servers = append(servers, server)

	// Save to config file
	return SaveServers(servers)
}

// SaveServers writes the server list to the config file
func SaveServers(servers []common.ServerConfig) error {
	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create a new viper instance for writing
	writeViper := viper.New()
	writeViper.SetConfigType("yaml")
	writeViper.SetConfigFile(configPath)

	// Set the servers
	writeViper.Set("servers", servers)

	// Write the config file
	if err := writeViper.WriteConfig(); err != nil {
		// If file doesn't exist, create it
		if os.IsNotExist(err) {
			return writeViper.SafeWriteConfig()
		}
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Reload the main viper instance
	return viper.ReadInConfig()
}

// ServerExists checks if a server with the given name exists
func ServerExists(name string) bool {
	if name == ReservedServerName {
		return true // "all" always exists
	}
	_, err := GetServer(name)
	return err == nil
}

// GetCurrentServer returns the server config for the current server flag
// Returns nil if "all" is selected (caller should handle multi-server case)
func GetCurrentServer() (*common.ServerConfig, error) {
	if serverFlag == ReservedServerName {
		return nil, nil // nil means "all"
	}
	return GetServer(serverFlag)
}

// GetCurrentServers returns all servers to target based on the server flag
func GetCurrentServers() ([]common.ServerConfig, error) {
	if serverFlag == ReservedServerName {
		return GetServers()
	}
	server, err := GetServer(serverFlag)
	if err != nil {
		return nil, err
	}
	return []common.ServerConfig{*server}, nil
}
