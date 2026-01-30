package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andatoshiki/wakafetch/wakafetch/ui"
)

func parseConfig() (string, string, error) {
	configPath := getConfigPath()

	file, err := os.Open(configPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var apiURL, apiKey string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "api_url") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				apiURL = strings.TrimSpace(parts[1])
			}
		} else if strings.HasPrefix(line, "api_key") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				apiKey = strings.TrimSpace(parts[1])
			}
		}

		if apiURL != "" && apiKey != "" {
			break
		}
	}

	if apiURL == "" {
		return "", "", fmt.Errorf("api_url not found in config")
	}

	if apiKey == "" {
		return "", "", fmt.Errorf("api_key not found in config")
	}

	apiURL = strings.TrimSuffix(apiURL, "/")

	// Normalize API URL for WakaTime vs Wakapi/self-hosted instances.
	// - Official WakaTime: keep as-is (expects /v1/... paths).
	// - Wakapi/self-hosted: ensure we point to the WakaTime-compatible prefix:
	//   <base>/api/compat/wakatime
	if !strings.Contains(apiURL, "wakatime.com") {
		// Treat anything that's not the official WakaTime API as Wakapi-compatible.
		if strings.HasSuffix(apiURL, "/api/compat/wakatime") {
			// Already normalized
		} else if strings.HasSuffix(apiURL, "/api") {
			apiURL = apiURL + "/compat/wakatime"
		} else {
			apiURL = apiURL + "/api/compat/wakatime"
		}
	}

	return apiURL, apiKey, nil
}

func getConfigPath() string {
	var configFile string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		ui.Errorln("Error getting home directory: %v\n", err)
	}
	configFile = filepath.Join(homeDir, ".wakatime.cfg")

	return configFile
}
