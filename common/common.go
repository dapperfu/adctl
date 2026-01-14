package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

// ServerConfig represents a single AdGuard server configuration
type ServerConfig struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Host     string `mapstructure:"host" yaml:"host"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

type CommandArgs struct {
	RequestBody map[string]any
	Method      string
	URL         url.URL
	Server      *ServerConfig // Optional server config, nil means use legacy viper config
}

// GetBaseURL returns the base URL for a server configuration
// If server is nil, uses legacy viper config
func GetBaseURL(server *ServerConfig) (url.URL, error) {
	var ret = url.URL{Scheme: "http"}

	var host string
	var err error
	if server != nil {
		host = server.Host
	} else {
		host, err = getHost()
		if err != nil {
			return ret, err
		}
	}

	if host == "" {
		return ret, fmt.Errorf("host is empty")
	}

	ret.Host = host
	return ret, nil
}

// GetBaseURLLegacy maintains backward compatibility
func GetBaseURLLegacy() (url.URL, error) {
	return GetBaseURL(nil)
}

func getHost() (string, error) {
	ret := viper.GetString("host")
	if ret == "" {
		return "", fmt.Errorf("can't find host (set ADCTL_HOST environment variable or configure in config file)")
	}
	return ret, nil
}

// AbleCommand enables or disables protection on a server
func AbleCommand(server *ServerConfig, state bool, durationString string) error {
	// base url
	baseURL, err := GetBaseURL(server)
	if err != nil {
		return err
	}

	baseURL.Path = "/control/protection"

	// data for post
	var requestBody = make(map[string]any)
	requestBody["enabled"] = state

	var duration uint64
	if len(durationString) > 0 {
		tmp, err := time.ParseDuration(durationString)
		if err != nil {
			return fmt.Errorf("time.ParseDuration: %w", err)
		}
		duration = uint64(tmp.Milliseconds())
	}

	requestBody["duration"] = duration

	// put it all together
	enableQuery := CommandArgs{
		Method:      "POST",
		URL:         baseURL,
		RequestBody: requestBody,
		Server:      server,
	}

	// don't care about body here
	_, err = SendCommand(enableQuery)
	if err != nil {
		return err
	}

	return nil
}

// AbleCommandLegacy maintains backward compatibility
func AbleCommandLegacy(state bool, durationString string) error {
	return AbleCommand(nil, state, durationString)
}

// SendCommand sends a command to a server
func SendCommand(ca CommandArgs) ([]byte, error) {
	var jsonData []byte
	var err error

	// turn params into json.  not sure if I can safely do this to all verbs.
	if ca.Method == "POST" || ca.Method == "PUT" {
		jsonData, err = json.Marshal(ca.RequestBody)
		if err != nil {
			return nil, fmt.Errorf("error marshaling json: %w", err)
		}
	}

	// create the final request
	request, err := http.NewRequest(ca.Method, ca.URL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// set request headers
	request.Header.Set("Content-Type", "application/json")

	// set basic auth - use server config if provided, otherwise use legacy viper
	var username, password string
	if ca.Server != nil {
		username = ca.Server.Username
		password = ca.Server.Password
	} else {
		username = viper.GetString("username")
		if username == "" {
			return nil, fmt.Errorf("can't find username (set ADCTL_USERNAME environment variable or configure in config file)")
		}
		password = viper.GetString("password")
		if password == "" {
			return nil, fmt.Errorf("can't find password (set ADCTL_PASSWORD environment variable or configure in config file)")
		}
	}

	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	request.SetBasicAuth(username, password)

	// connect.  Old implementation let me set timeouts to handle short dns timeouts and
	//   long log fetches.  bother with it here? skipping for now.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error Do'ing request: %w", err)
	}
	defer resp.Body.Close()

	// read response
	// Read response but really I just want to know if there's an error
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code not 200: %v", resp.Status)
	}

	return body, nil
}

// SendCommandToAll sends a command to all servers and collects results
// Returns a map of server name to result (or error)
func SendCommandToAll(servers []ServerConfig, ca CommandArgs) map[string]interface{} {
	results := make(map[string]interface{})

	for _, server := range servers {
		// Create a copy of CommandArgs with this server's config
		serverCA := ca
		serverCA.Server = &server

		// Get base URL for this server
		baseURL, err := GetBaseURL(&server)
		if err != nil {
			results[server.Name] = fmt.Errorf("failed to get base URL: %w", err)
			continue
		}
		serverCA.URL = baseURL
		serverCA.URL.Path = ca.URL.Path // Preserve the path

		// Send command
		body, err := SendCommand(serverCA)
		if err != nil {
			results[server.Name] = err
		} else {
			results[server.Name] = body
		}
	}

	return results
}
