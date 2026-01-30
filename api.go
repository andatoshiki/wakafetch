package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andatoshiki/wakafetch/wakafetch/types"
)

func fetchSummary(apiKey, apiURL string, days int, timeoutSeconds int) (*types.SummaryResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	today := time.Now()
	todayDate := today.Format("2006-01-02")
	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	requestURL := fmt.Sprintf("%s/v1/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
	if strings.HasSuffix(apiURL, "/v1") {
		requestURL = fmt.Sprintf("%s/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
	}
	response, err := fetchApi[types.SummaryResponse](apiKey, requestURL, timeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

func fetchSummaryWithDates(apiKey, apiURL, startDate, endDate string, timeoutSeconds int) (*types.SummaryResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	requestURL := fmt.Sprintf("%s/v1/users/current/summaries?start=%s&end=%s", apiURL, startDate, endDate)
	if strings.HasSuffix(apiURL, "/v1") {
		requestURL = fmt.Sprintf("%s/users/current/summaries?start=%s&end=%s", apiURL, startDate, endDate)
	}
	response, err := fetchApi[types.SummaryResponse](apiKey, requestURL, timeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

func fetchStats(apiKey, apiURL, rangeStr string, timeoutSeconds int) (*types.StatsResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	requestURL := fmt.Sprintf("%s/v1/users/current/stats/%s", apiURL, rangeStr)
	if strings.HasSuffix(apiURL, "/v1") {
		requestURL = fmt.Sprintf("%s/users/current/stats/%s", apiURL, rangeStr)
	}
	// fmt.Println(requestURL)
	response, err := fetchApi[types.StatsResponse](apiKey, requestURL, timeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

func fetchApi[T any](apiKey, requestURL string, timeoutSeconds int) (*T, error) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}
	timeout := time.Duration(timeoutSeconds) * time.Second
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))
	req.Header.Set("Authorization", "Basic "+encodedKey)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		if ne, ok := err.(interface{ Timeout() bool }); ok && ne.Timeout() {
			return nil, fmt.Errorf("Request timed out after %s while contacting server", timeout)
		}
		return nil, fmt.Errorf("Unable to reach server. Check your internet connection")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("Authentication failed (401). Check your API key")
		case http.StatusForbidden:
			return nil, fmt.Errorf("Access forbidden (403). Your API key might not have permission")
		case http.StatusNotFound:
			return nil, fmt.Errorf("Endpoint not found (404). Verify the API URL")
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("Rate limit exceeded (429). Please try again later")
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return nil, fmt.Errorf("Server unavailable (%s). Please try again later", resp.Status)
		default:
			return nil, fmt.Errorf("Api request failed: %s", resp.Status)
		}
	}

	var apiResponse T
	if err = json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("Invalid response from server (failed to decode JSON)")
	}

	return &apiResponse, nil
}
