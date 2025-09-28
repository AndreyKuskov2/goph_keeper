package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Register registers a new user with the server
func (c *CLIClient) Register(login, password string) error {
	user := User{
		Login:    login,
		Password: password,
	}

	reqBody, err := CreateJSONRequest(user)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/register", reqBody, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// Login authenticates a user and returns a JWT token
func (c *CLIClient) Login(login, password string) (string, error) {
	user := UserLoginRequest{
		Login:    login,
		Password: password,
	}

	reqBody, err := CreateJSONRequest(user)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/login", reqBody, false)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.handleErrorResponse(resp)
	}

	// Parse server response structure
	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract token from data field
	var loginResp UserLoginResponse
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &loginResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal login response: %v", err)
	}

	return loginResp.Token, nil
}
