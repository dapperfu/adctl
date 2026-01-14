/*
Copyright © 2024 Eric Osborne
No header.
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugLogger *log.Logger

// var outputFormat string
var enableDebug bool
var serverFlag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "adctl",
	Long: `adctl lets you control AdGuard Home from the CLI. Documentation and source: https://github.com/ewosborne/adctl`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func TestscriptEntryPoint() int {
	Execute()
	return 0
}

func SetVersionInfo(version string) {
	rootCmd.Version = version
}

func init() {
	initConfig()

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "d", os.Getenv("DEBUG") == "true", "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&serverFlag, "server", "s", "all", "Server name to target (use 'all' for all servers)")
	//rootCmd.PersistentFlags().StringVarP(&outputFormat, "output format", "o", "json", "Enable debug mode")

	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)

	// need PreRun because flags aren't parsed until a command is run.
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if enableDebug {
			debugLogger.SetOutput(os.Stderr)
		} else {
			debugLogger.SetOutput(io.Discard)
		}

		// Validate server flag (skip for server command itself)
		if cmd.Name() != "server" && serverFlag != "all" {
			if !ServerExists(serverFlag) {
				fmt.Fprintf(os.Stderr, "Error: server '%s' not found\n", serverFlag)
				os.Exit(1)
			}
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set config file name (without extension)
	viper.SetConfigName("adctl")

	// Set config type priority: yaml, json, toml
	viper.SetConfigType("yaml")

	// Add config search paths - use OS-specific config directory
	configDir, err := GetConfigDir()
	if err == nil {
		viper.AddConfigPath(configDir)
	}

	// Also check legacy locations for backward compatibility
	home, err := os.UserHomeDir()
	if err == nil {
		// ~/.config/ (legacy)
		viper.AddConfigPath(filepath.Join(home, ".config"))
	}
	// Current directory
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("ADCTL")
	viper.AutomaticEnv() // read in environment variables that match

	// Bind environment variables
	viper.BindEnv("host", "ADCTL_HOST")
	viper.BindEnv("username", "ADCTL_USERNAME")
	viper.BindEnv("password", "ADCTL_PASSWORD")

	// If a config file is found, read it in.
	// Note: debugLogger not initialized yet, so we can't log here
	_ = viper.ReadInConfig()
}
