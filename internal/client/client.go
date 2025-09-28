package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CLIClient represents the main client for interacting with the GophKeeper server
type CLIClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewCLIClient creates a new CLI client instance
func NewCLIClient(baseURL string) *CLIClient {
	return &CLIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Token management methods
func (c *CLIClient) SaveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	tokenFile := filepath.Join(homeDir, ".goph_keeper_token")
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

func (c *CLIClient) LoadToken() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	tokenFile := filepath.Join(homeDir, ".goph_keeper_token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return fmt.Errorf("failed to read token file: %v", err)
	}

	// Remove extra whitespace and newlines
	c.token = strings.TrimSpace(string(data))
	return nil
}

func (c *CLIClient) SetToken(token string) {
	c.token = token
}

// Core HTTP methods
func (c *CLIClient) makeRequest(method, path string, body []byte, useAuth bool) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if useAuth {
		if c.token == "" {
			return nil, fmt.Errorf("authentication token required")
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

func (c *CLIClient) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %v", err)
	}

	var errorResp ServerErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("server error (status %d): %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("server error (status %d): %s", resp.StatusCode, errorResp.Error)
}

// Helper method to parse server response data
func (c *CLIClient) parseServerResponseData(resp *http.Response, target interface{}) error {
	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract data from server response
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return nil
}
